package stack

import (
	"context"
	"crypto/tls"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"

	pb "github.com/endobit/stack/internal/generated/go/proto/stack/v1"
	"github.com/endobit/stack/internal/stackd"
	"github.com/endobit/stack/internal/stackd/auth"
)

const (
	Add    = "add"
	Dump   = "dump"
	List   = "list"
	Load   = "load"
	Remove = "remove"
	Report = "report"
	Update = "update"
)

const (
	archFlag     = "arch"
	makeFlag     = "make"
	modelFlag    = "model"
	nameFlag     = "name"
	timeZoneFlag = "time-zone"
)

type rpcClient struct {
	logger *slog.Logger
	stack  pb.StackServiceClient
	auth   pb.AuthServiceClient
	token  string
}

func (c *rpcClient) Context() context.Context {
	ctx := context.Background()
	if c.token != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, auth.AuthorizationMetaData, c.token)
	}

	return ctx
}

func (c *rpcClient) Authorize() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := pb.LoginRequest_builder{
		Username: proto.String("admin"),
		Password: proto.String("admin"),
	}.Build()

	resp, err := c.auth.Login(ctx, req)
	if err != nil {
		return err
	}

	c.token = resp.GetToken()

	return nil
}

// NewRootCmd creates the top level command object and initializes the cli.
func NewRootCmd() *cobra.Command {
	var (
		hostname string
		port     int
		debug    bool
		client   rpcClient
	)

	cmd := cobra.Command{
		Use:   "cast",
		Short: "Stack Client",
		Long:  "Stack Command Line Client",
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			logLevel := slog.LevelInfo
			if debug {
				logLevel = slog.LevelDebug
			}

			handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
				AddSource: true,
				Level:     logLevel,
			})
			logger := slog.New(handler)

			creds := credentials.NewTLS(&tls.Config{
				InsecureSkipVerify: true, //nolint:gosec
				MinVersion:         tls.VersionTLS12,
			})

			conn, err := grpc.NewClient(hostname+":"+strconv.Itoa(port),
				grpc.WithTransportCredentials(creds))
			if err != nil {
				return err
			}

			client = rpcClient{
				logger: logger,
				stack:  pb.NewStackServiceClient(conn),
				auth:   pb.NewAuthServiceClient(conn),
			}

			return client.Authorize()
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			return nil
		},
	}

	cmd.PersistentFlags().IntVar(&port, "port", stackd.DefaultPort, "service port")
	cmd.PersistentFlags().StringVar(&hostname, "host", "localhost", "service hostname")

	cmd.AddCommand(
		newAddCmd(&client),
		newDumpCmd(&client),
		newListCmd(&client),
		newLoadCmd(&client),
		newRemoveCmd(&client),
		newReportCmd(&client),
		newUpdateCmd(&client),
	)

	return &cmd
}

func newAddCmd(client *rpcClient) *cobra.Command {
	cmd := cobra.Command{
		Use:   Add,
		Short: "Add objects",
	}

	cmd.AddCommand(
		newAddZoneCmd(client),
	)

	return &cmd
}

func newListCmd(client *rpcClient) *cobra.Command {
	cmd := cobra.Command{
		Use:     List,
		Aliases: []string{"ls"},
		Short:   "List objects",
		Long:    "List is for humans.",
	}

	cmd.AddCommand(
		newListZoneCmd(client),
	)

	return &cmd
}

func newRemoveCmd(client *rpcClient) *cobra.Command {
	cmd := cobra.Command{
		Use:   Remove,
		Short: "Remove objects",
	}

	cmd.AddCommand(
		newRemoveZoneCmd(client),
	)

	return &cmd
}

func newReportCmd(_ *rpcClient) *cobra.Command {
	return &cobra.Command{
		Use:   Report,
		Short: "Report objects",
		Long:  "Report is for computers.",
	}
}

func newUpdateCmd(client *rpcClient) *cobra.Command {
	cmd := cobra.Command{
		Use:   Update,
		Short: "update objects",
	}

	cmd.AddCommand(
		newUpdateZoneCmd(client),
	)

	return &cmd
}

func required[T comparable](v T) *T {
	var zero T

	if v == zero {
		panic("required value is zero")
	}

	return &v
}

func optional[T comparable](v T) *T {
	var zero T

	if v != zero {
		return &v
	}

	return nil
}
