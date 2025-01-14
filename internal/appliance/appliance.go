package appliance

import (
	"github.com/spf13/cobra"

	"endobit.io/metal"
	"endobit.io/table"

	"endobit.io/metal-cli/internal/cli"
	pb "endobit.io/metal/gen/go/proto/metal/v1"
)

const object = "appliance"

type options struct {
	Rename *string
	Zone   *string
}

func NewCmd(verb cli.Verb, rpc *metal.Client) *cobra.Command {
	var (
		cmd  cobra.Command
		opts options
	)

	switch verb {
	case cli.Add:
		cmd = cobra.Command{
			Use:   object + " name",
			Short: "Add an " + object + " to a zone",
			Args:  cobra.ExactArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				if err := opts.create(rpc, args[0]); err != nil {
					return err
				}

				return opts.update(rpc, args[0])
			},
		}
	case cli.Set:
		cmd = cobra.Command{
			Use:   object + " name",
			Short: "Set an " + object + "'s properties",
			Args:  cobra.ExactArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				return opts.update(rpc, args[0])
			},
		}
	case cli.List:
		cmd = cobra.Command{
			Use:   object + " [glob]",
			Short: "List one or more " + object + "s",
			Args:  cobra.MaximumNArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				var glob string

				if len(args) > 0 {
					glob = args[0]
				}
				return opts.list(rpc, glob)
			},
		}
	case cli.Remove:
		cmd = cobra.Command{
			Use:   object + " glob",
			Short: "Remove one or more " + object + "s",
			Args:  cobra.ExactArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				return opts.remove(rpc, args[0])
			},
		}
	}

	opts.Zone = cli.AddZoneFlag(&cmd, object)

	if verb == cli.Add || verb == cli.Set || verb == cli.Remove {
		if err := cmd.MarkFlagRequired("zone"); err != nil {
			panic(err)
		}
	}

	if verb == cli.Set {
		opts.Rename = cli.AddRenameFlag(&cmd, object)
	}

	return &cmd
}

func (o options) create(rpc *metal.Client, appliance string) error {
	req := pb.CreateApplianceRequest_builder{
		Zone: o.Zone,
		Name: &appliance,
	}.Build()

	_, err := rpc.Metal.CreateAppliance(rpc.Context(), req)

	return err
}

func (o options) list(rpc *metal.Client, glob string) error {
	type row struct{ Zone, Name string }
	t := table.New()
	defer t.Flush()

	r := rpc.NewApplianceReader(cli.Val(o.Zone), glob)

	for resp, err := range r.Responses() {
		if err != nil {
			return err
		}

		_ = t.Write(row{
			Zone: resp.GetZone(),
			Name: resp.GetName(),
		})
	}

	return nil
}

func (o options) update(rpc *metal.Client, appliance string) error {
	req := pb.UpdateApplianceRequest_builder{
		Zone: o.Zone,
		Name: &appliance,
		Fields: pb.UpdateApplianceRequest_Fields_builder{
			Name: o.Rename,
		}.Build(),
	}.Build()

	_, err := rpc.Metal.UpdateAppliance(rpc.Context(), req)

	return err
}

func (o options) remove(rpc *metal.Client, glob string) error {
	req := pb.DeleteAppliancesRequest_builder{
		Zone: o.Zone,
		Glob: &glob,
	}.Build()

	_, err := rpc.Metal.DeleteAppliances(rpc.Context(), req)

	return err
}
