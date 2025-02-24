// Package main implements the stack CLI.
package main

import (
	"crypto/tls"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"endobit.io/metal"
	"endobit.io/metal-cli/internal/commands"
	authpb "endobit.io/metal/gen/go/proto/auth/v1"
	metalpb "endobit.io/metal/gen/go/proto/metal/v1"
	"endobit.io/metal/logging"
)

var version string

func main() {
	cmd := newRootCmd()
	cmd.Version = version

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	var (
		username, password, metalServer string
		rpc                             metal.Client
		logOpts                         *logging.Options
	)

	cmd := cobra.Command{
		Use:   "stack",
		Short: "Stack Client",
		Long:  "Stack Command Line Client",
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			logger, err := logOpts.NewLogger()
			if err != nil {
				return err
			}

			creds := credentials.NewTLS(&tls.Config{
				InsecureSkipVerify: true, //nolint:gosec
				MinVersion:         tls.VersionTLS12,
			})

			conn, err := grpc.NewClient(metalServer, grpc.WithTransportCredentials(creds))
			if err != nil {
				return err
			}

			rpc = metal.Client{
				Logger: logger,
				Metal:  metalpb.NewMetalServiceClient(conn),
				Auth:   authpb.NewAuthServiceClient(conn),
			}

			return rpc.Authorize(username, password)
		},
	}

	logOpts = logging.NewOptions(cmd.PersistentFlags())

	cmd.PersistentFlags().StringVar(&username, "username", "admin", "username for authentication")
	cmd.PersistentFlags().StringVar(&password, "password", "admin", "password for authentication")
	cmd.PersistentFlags().StringVar(&metalServer, "metal-server", "localhost:"+strconv.Itoa(metal.DefaultPort),
		"address of the metal server")

	root := commands.Root{Client: &rpc}

	cmd.AddCommand(
		root.New(commands.Add),
		root.New(commands.Dump),
		root.New(commands.List),
		root.New(commands.Load),
		root.New(commands.Remove),
		root.New(commands.Report),
		root.New(commands.Set))

	return &cmd
}
