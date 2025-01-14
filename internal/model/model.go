package model

import (
	"github.com/spf13/cobra"

	"endobit.io/metal"
	"endobit.io/metal-cli/internal/cli"
	pb "endobit.io/metal/gen/go/proto/metal/v1"
	"endobit.io/table"
)

const object = "model"

type options struct {
	Arch   *string
	Rename *string
}

func NewCmd(verb cli.Verb, rpc *metal.Client) *cobra.Command {
	var (
		cmd  cobra.Command
		opts options
	)

	switch verb {
	case cli.Add:
		cmd = cobra.Command{
			Use:   object + "make name",
			Short: "Add a " + object,
			Args:  cobra.ExactArgs(2),
			RunE: func(_ *cobra.Command, args []string) error {
				if err := opts.create(rpc, args[0], args[1]); err != nil {
					return err
				}

				return opts.update(rpc, args[0], args[1])
			},
		}
	case cli.Set:
		cmd = cobra.Command{
			Use:   object + "make name",
			Short: "Set a " + object + "'s properties",
			Args:  cobra.ExactArgs(2),
			RunE: func(_ *cobra.Command, args []string) error {
				return opts.update(rpc, args[0], args[1])
			},
		}
	case cli.List:
		cmd = cobra.Command{
			Use:   object + "make [glob]",
			Short: "List one or more " + object + "s",
			Args:  cobra.RangeArgs(1, 2),
			RunE: func(_ *cobra.Command, args []string) error {
				var glob string

				if len(args) > 1 {
					glob = args[1]
				}

				return opts.list(rpc, args[0], glob)
			},
		}
	case cli.Remove:
		cmd = cobra.Command{
			Use:   object + "make glob",
			Short: "Remove one or more " + object + "s",
			Args:  cobra.ExactArgs(2),
			RunE: func(_ *cobra.Command, args []string) error {
				return opts.remove(rpc, args[0], args[1])
			},
		}
	}

	if verb == cli.Add || verb == cli.Set {
		opts.Arch = cmd.Flags().String("arch", "", "architecture for the "+object)
	}

	if verb == cli.Set {
		opts.Rename = cli.AddRenameFlag(&cmd, object)
	}

	return &cmd
}

func (options) create(rpc *metal.Client, vendor, model string) error {
	req := pb.CreateModelRequest_builder{
		Make: &vendor,
		Name: &model,
	}.Build()

	_, err := rpc.Metal.CreateModel(rpc.Context(), req)
	return err
}

func (options) list(rpc *metal.Client, vendor, glob string) error {
	type row struct{ Make, Model, Arch string }
	t := table.New()
	defer t.Flush()

	r := rpc.NewModelReader(vendor, glob)

	for resp, err := range r.Responses() {
		if err != nil {
			return err
		}

		_ = t.Write(row{
			Make:  resp.GetMake(),
			Model: resp.GetName(),
			Arch:  resp.GetArchitecture().String(),
		})
	}

	return nil
}

func (o options) update(rpc *metal.Client, vendor, model string) error {
	var pbarch *pb.Architecture

	if o.Arch != nil {
		a := pb.Architecture(pb.Architecture_value[*o.Arch])
		pbarch = &a
	}

	req := pb.UpdateModelRequest_builder{
		Make: &vendor,
		Name: &model,
		Fields: pb.UpdateModelRequest_Fields_builder{
			Name:         o.Rename,
			Architecture: pbarch,
		}.Build(),
	}.Build()

	_, err := rpc.Metal.UpdateModel(rpc.Context(), req)

	return err
}

func (o options) remove(rpc *metal.Client, vendor, glob string) error {
	req := pb.DeleteModelsRequest_builder{
		Make: &vendor,
		Glob: &glob,
	}.Build()

	_, err := rpc.Metal.DeleteModels(rpc.Context(), req)
	return err
}
