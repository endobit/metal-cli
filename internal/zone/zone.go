package zone

import (
	"errors"

	"github.com/spf13/cobra"

	"endobit.io/metal"
	"endobit.io/metal-cli/internal/cli"
	pb "endobit.io/metal/gen/go/proto/metal/v1"
	"endobit.io/table"
)

const object = "zone"

type options struct {
	Rename   *string
	TimeZone *string
	Template *string
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
			Short: "Add a " + object,
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
			Short: "Set a " + object + "'s properties",
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

	case cli.Report:
		cmd = cobra.Command{
			Use:   object + " name",
			Short: "Evaluate a template for a " + object,
			Args:  cobra.ExactArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				return opts.report(rpc, args[0])
			},
		}
	}

	if verb == cli.Add || verb == cli.Set {
		opts.TimeZone = cmd.Flags().String("timezone", "", "time zone for the "+object)
	}

	if verb == cli.Set {
		opts.Rename = cli.AddRenameFlag(&cmd, object)
	}

	if verb == cli.Report {
		opts.Template = cmd.Flags().String("template", "", "template for the "+object)
	}

	return &cmd
}

func (options) create(rpc *metal.Client, zone string) error {
	var req pb.CreateZoneRequest

	req.SetName(zone)
	_, err := rpc.Metal.CreateZone(rpc.Context(), &req)

	return err
}

func (options) list(rpc *metal.Client, glob string) error {
	type row struct{ Name, TimeZone string }
	t := table.New()
	defer t.Flush()

	r := rpc.NewZoneReader(glob)

	for resp, err := range r.Responses() {
		if err != nil {
			return err
		}

		_ = t.Write(row{
			Name:     resp.GetName(),
			TimeZone: resp.GetTimeZone(),
		})
	}

	return nil
}

func (o options) update(rpc *metal.Client, zone string) error {
	req := pb.UpdateZoneRequest_builder{
		Name: &zone,
		Fields: pb.UpdateZoneRequest_Fields_builder{
			Name:     o.Rename,
			TimeZone: o.TimeZone,
		}.Build(),
	}.Build()

	_, err := rpc.Metal.UpdateZone(rpc.Context(), req)

	return err
}

func (options) remove(rpc *metal.Client, glob string) error {
	var req pb.DeleteZonesRequest

	req.SetGlob(glob)
	_, err := rpc.Metal.DeleteZones(rpc.Context(), &req)

	return err
}

func (o options) report(rpc *metal.Client, zone string) error {
	var req pb.ReadSchemaRequest

	if zone != "" {
		req.SetZone(zone)
	}

	doc, err := rpc.Metal.ReadSchema(rpc.Context(), &req)
	if err != nil {
		return err
	}

	switch len(doc.GetSchema().Zones) {
	case 0:
		return errors.New("no zones found")
	case 1:
		return errors.New("more than one zone found")
	}

	return nil
}
