// Package stackd implements the Stack Daemon.
package stackd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3" // sqlite3 driver

	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/endobit/stack/internal/cert"
	"github.com/endobit/stack/internal/generated/go/db"
	pb "github.com/endobit/stack/internal/generated/go/proto/stack/v1"
	"github.com/endobit/stack/internal/logging"
	"github.com/endobit/stack/internal/stackd/auth"
	"github.com/endobit/stack/internal/stackd/stack"
	stackdb "github.com/endobit/stack/sql"
)

// DefaultPort is the default port to listen on. It can be overridden with the
// --port flag.
const DefaultPort = 8080

// NewRootCmd creates the top level command object.
func NewRootCmd() *cobra.Command {
	var (
		dbDir         string
		useTailscale  bool
		port          int
		tokenTTL      time.Duration
		certDir       string
		debug         bool
		adminPassword string
	)
	cmd := cobra.Command{
		Use:   "stackd",
		Short: "Stack Server",
		RunE: func(cmd *cobra.Command, _ []string) error {
			logger := setupLogging(debug)
			logger.Info("Starting Stack Server", "version", cmd.Version, "port", port)

			dbPath := filepath.Join(dbDir, "stack.db")
			if err := setupDatabase(logger, dbPath); err != nil {
				return fmt.Errorf("failed to setup database: %w", err)
			}
			db, err := sql.Open("sqlite3", dbPath)
			if err != nil {
				return fmt.Errorf("failed to open database: %w", err)
			}

			listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
			if err != nil {
				return fmt.Errorf("failed to listen: %w", err)
			}

			// Create a self-signed certificate if one does not exist.

			certPath := filepath.Join(certDir, "cert.pem")
			keyPath := filepath.Join(certDir, "key.pem")
			if err := setupCertificates(certPath, keyPath); err != nil {
				return err
			}

			cert, err := tls.LoadX509KeyPair(certPath, keyPath)
			if err != nil {
				return fmt.Errorf("failed to load certificate: %w", err)
			}

			ca := x509.NewCertPool()

			caBytes, err := os.ReadFile(certPath)
			if err != nil {
				return err
			}

			if !ca.AppendCertsFromPEM(caBytes) {
				return fmt.Errorf("failed to append CA certificate")
			}

			tlsConfig := tls.Config{
				Certificates: []tls.Certificate{cert},
				ClientAuth:   tls.NoClientCert,
				ClientCAs:    ca,
				MinVersion:   tls.VersionTLS12,
			}

			if adminPassword != "" {
				if err := setupAdminPassword(db, "admin", adminPassword); err != nil {
					return err
				}
			}

			authService := auth.NewService(
				auth.WithLogger(logger),
				auth.WithTTL(tokenTTL),
				auth.WithUser("admin", "admin"), // TODO: keep track of users in the database
				auth.WithDB(db),
			)

			stackService := stack.NewService(
				stack.WithLogger(logger),
				stack.WithDB(db),
			)

			server := grpc.NewServer(
				grpc.Creds(credentials.NewTLS(&tlsConfig)),
				grpc.ChainUnaryInterceptor(
					unaryLoggingInterceptor, // log the call
					authService.UnaryInterceptor(pb.AuthService_Login_FullMethodName), // check for a valid token and authorization
				),
				grpc.ChainStreamInterceptor(
					streamLoggingInterceptor,        // log the call
					authService.StreamInterceptor(), // check for a valid token and authorization
				),
			)

			pb.RegisterAuthServiceServer(server, authService)
			pb.RegisterStackServiceServer(server, stackService)

			if err := server.Serve(listener); err != nil {
				return fmt.Errorf("failed to serve: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&dbDir, "dbpath", ".", "Database directory")
	cmd.Flags().BoolVar(&useTailscale, "tailscale", false, "Get certificate from tailscale")
	cmd.Flags().IntVar(&port, "port", DefaultPort, "port to listen on")
	cmd.Flags().DurationVar(&tokenTTL, "token-ttl", 5*time.Minute, "token time to live")
	cmd.Flags().StringVar(&certDir, "cert-path", ".", "path to cert.pem and key.pem files")
	cmd.Flags().BoolVar(&debug, "debug", false, "enable debug logging")
	cmd.Flags().StringVar(&adminPassword, "admin-password", "admin", "admin password")

	return &cmd
}

func setupAdminPassword(database *sql.DB, username, password string) error {
	ctx := context.Background()

	q := db.New(database)

	err := q.CreateUser(ctx, db.CreateUserParams{
		Name: username,
	})
	if err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}

	err = q.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
		User:         username,
		PasswordHash: password,
	})
	if err != nil {
		return fmt.Errorf("failed set admin password: %w", err)
	}

	user, err := q.ReadUser(ctx, db.ReadUserParams{User: username})
	if err != nil {
		return fmt.Errorf("failed to read admin user: %w", err)
	}

	err = q.CreateRole(ctx, db.CreateRoleParams{Name: "admin"})
	if err != nil {
		return fmt.Errorf("failed to create admin role: %w", err)
	}

	role, err := q.ReadRole(ctx, db.ReadRoleParams{Role: "admin"})
	if err != nil {
		return fmt.Errorf("failed to read admin role: %w", err)
	}

	return q.AssignRoleToUser(ctx, db.AssignRoleToUserParams{
		UserID: user.ID,
		RoleID: role.ID,
	})
}

func setupLogging(debug bool) *slog.Logger {
	logLevel := slog.LevelInfo
	if debug {
		logLevel = slog.LevelDebug
	}

	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel,
	})

	logger := slog.New(handler)
	if debug {
		logger.Debug("Debug logging enabled")
	}

	return logger
}

func setupCertificates(certPath, keyPath string) error {
	var missing int

	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		missing++
	}
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		missing++
	}
	if missing == 1 {
		return fmt.Errorf("missing cert or key file")
	}
	if missing == 0 {
		return nil
	}

	certFile, err := os.Create(certPath)
	if err != nil {
		return fmt.Errorf("failed to create cert file: %w", err)
	}
	defer certFile.Close()

	keyFile, err := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("failed to create key file: %w", err)
	}
	defer keyFile.Close()

	options := cert.NewOptions()
	return options.Create(certFile, keyFile)
}

func setupDatabase(logger *slog.Logger, dbPath string) error {
	store, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	defer store.Close()

	goose.SetLogger(logging.Legacy{Logger: logger})
	goose.SetBaseFS(stackdb.Migrations)

	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	if err := goose.Up(store, "migrations"); err != nil {
		return err
	}

	return nil
}
