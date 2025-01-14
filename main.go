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
	"endobit.io/metal-cli/internal/appliance"
	"endobit.io/metal-cli/internal/cli"
	"endobit.io/metal-cli/internal/model"
	"endobit.io/metal-cli/internal/zone"
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

	cmd.AddCommand(
		newVerbCmd(cli.Add, &rpc),
		newVerbCmd(cli.Dump, &rpc),
		newVerbCmd(cli.List, &rpc),
		newVerbCmd(cli.Load, &rpc),
		newVerbCmd(cli.Remove, &rpc),
		newVerbCmd(cli.Report, &rpc),
		newVerbCmd(cli.Set, &rpc))

	return &cmd
}

var (
	jsonFlag bool
)

func newVerbCmd(verb cli.Verb, rpc *metal.Client) *cobra.Command {
	var cmd cobra.Command

	switch verb {
	case cli.Add:
		cmd = cobra.Command{
			Use:     "add",
			Aliases: []string{"create"},
			Short:   "Add objects",
		}

		cmd.AddCommand(
			appliance.NewCmd(verb, rpc),
			model.NewCmd(verb, rpc),
			zone.NewCmd(verb, rpc))

	case cli.Dump:
		cmd = cobra.Command{
			Use:   "dump",
			Short: "Dump stack schema",
			Args:  cobra.NoArgs,
			RunE: func(_ *cobra.Command, _ []string) error {
				return dump(rpc)
			},
		}

		cmd.PersistentFlags().BoolVar(&jsonFlag, "json", false, "output JSON instead of YAML")

	case cli.Set:
		cmd = cobra.Command{
			Use:     "set",
			Aliases: []string{"update"},
			Short:   "Set object properties",
		}

		cmd.AddCommand(
			appliance.NewCmd(verb, rpc),
			model.NewCmd(verb, rpc),
			zone.NewCmd(verb, rpc))

	case cli.List:
		cmd = cobra.Command{
			Use:     "list",
			Aliases: []string{"ls"},
			Short:   "List objects",
			Long:    "List is for humans.",
		}

		cmd.AddCommand(
			appliance.NewCmd(verb, rpc),
			model.NewCmd(verb, rpc),
			zone.NewCmd(verb, rpc))

	case cli.Load:
		cmd = cobra.Command{
			Use:     "load filename",
			Aliases: []string{"ld"},
			Args:    cobra.ExactArgs(1),
			Short:   "Load objects",
			RunE: func(_ *cobra.Command, args []string) error {
				return load(rpc, args[0])
			},
		}

	case cli.Report:
		cmd = cobra.Command{
			Use:   "report",
			Short: "Report objects",
			Long:  "Report is for computers.",
		}

	case cli.Remove:
		cmd = cobra.Command{
			Use:     "remove",
			Aliases: []string{"del", "rm"},
			Short:   "Remove objects",
		}

		cmd.AddCommand(
			appliance.NewCmd(verb, rpc),
			model.NewCmd(verb, rpc),
			zone.NewCmd(verb, rpc))
	}

	return &cmd
}
