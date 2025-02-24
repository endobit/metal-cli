package commands

import (
	"github.com/spf13/cobra"

	"endobit.io/metal"
	"endobit.io/table"

	"endobit.io/metal-cli/internal/flags"
	pb "endobit.io/metal/gen/go/proto/metal/v1"
)

type Rack struct {
	Client     *metal.Client
	renameFlag flags.Rename
	zoneFlag   flags.Zone
}

type RackAttr struct {
	Client     *metal.Client
	renameFlag flags.Rename
	zoneFlag   flags.Zone
	rackFlag   flags.Rack
}

func (e *Rack) New(verb Verb) *cobra.Command {
	var cmd cobra.Command

	switch verb {
	case Add:
		cmd = cobra.Command{
			Use:   rack + " name",
			Short: "Add a " + rack + " to a zone",
			Args:  cobra.ExactArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				if err := e.create(args[0]); err != nil {
					return err
				}

				return e.update(args[0])
			},
		}
	case Set:
		cmd = cobra.Command{
			Use:   rack + " name",
			Short: "Set a " + rack + "'s properties",
			Args:  cobra.ExactArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				return e.update(args[0])
			},
		}
	case List:
		cmd = cobra.Command{
			Use:   rack + " [glob]",
			Short: "List one or more " + rack + "s",
			Args:  cobra.MaximumNArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				var glob string

				if len(args) > 0 {
					glob = args[0]
				}
				return e.list(glob)
			},
		}
	case Remove:
		cmd = cobra.Command{
			Use:   rack + " glob",
			Short: "Remove one or more " + rack + "s",
			Args:  cobra.ExactArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				return e.remove(args[0])
			},
		}
	}

	e.zoneFlag.Add(cmd.Flags(), rack)

	if verb == Add || verb == Set || verb == Remove {
		e.zoneFlag.Required(cmd.Flags())
	}

	if verb == Set {
		e.renameFlag.Add(cmd.Flags(), rack)
	}

	attr := RackAttr{Client: e.Client}
	cmd.AddCommand(attr.New(verb))

	return &cmd
}

func (e *Rack) create(rack string) error {
	req := pb.CreateRackRequest_builder{
		Zone: e.zoneFlag.Ptr(),
		Name: &rack,
	}.Build()

	_, err := e.Client.Metal.CreateRack(e.Client.Context(), req)

	return err
}

func (e *Rack) list(glob string) error {
	type row struct{ Zone, Rack string }
	t := table.New()
	defer t.Flush()

	r := e.Client.NewRackReader(e.zoneFlag.Val(), glob)

	for resp, err := range r.Responses() {
		if err != nil {
			return err
		}

		_ = t.Write(row{
			Zone: resp.GetZone(),
			Rack: resp.GetName(),
		})
	}

	return nil
}

func (e *Rack) update(rack string) error {
	req := pb.UpdateRackRequest_builder{
		Zone: e.zoneFlag.Ptr(),
		Name: &rack,
		Fields: pb.UpdateRackRequest_Fields_builder{
			Name: e.renameFlag.Ptr(),
		}.Build(),
	}.Build()

	_, err := e.Client.Metal.UpdateRack(e.Client.Context(), req)

	return err
}

func (e *Rack) remove(glob string) error {
	req := pb.DeleteRacksRequest_builder{
		Zone: e.zoneFlag.Ptr(),
		Glob: &glob,
	}.Build()

	_, err := e.Client.Metal.DeleteRacks(e.Client.Context(), req)

	return err
}

func (a *RackAttr) New(verb Verb) *cobra.Command {
	var cmd cobra.Command

	switch verb {
	case Add:
		cmd = cobra.Command{
			Use:   attribute + " name",
			Short: "Add an " + attribute + " to a " + rack,
			Args:  cobra.ExactArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				if err := a.create(args[0]); err != nil {
					return err
				}

				return a.update(args[0])
			},
		}
	case Set:
		cmd = cobra.Command{
			Use:   attribute + " name",
			Short: "Set a " + rack + " " + attribute + "'s properties",
			Args:  cobra.ExactArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				return a.update(args[0])
			},
		}
	case List:
		cmd = cobra.Command{
			Use:   attribute + " [glob]",
			Short: "List one or more " + rack + " " + attribute + "s",
			Args:  cobra.MaximumNArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				var glob string

				if len(args) > 0 {
					glob = args[0]
				}
				return a.list(glob)
			},
		}
	case Remove:
		cmd = cobra.Command{
			Use:   attribute + " glob",
			Short: "Remove one or more " + rack + " " + attribute + "s",
			Args:  cobra.ExactArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				return a.remove(args[0])
			},
		}
	}

	a.zoneFlag.Add(cmd.Flags(), rack)
	a.rackFlag.Add(cmd.Flags(), attribute)

	if verb == Add || verb == Set || verb == Remove {
		a.zoneFlag.Required(cmd.Flags())
		a.rackFlag.Required(cmd.Flags())
	}

	if verb == Set {
		a.renameFlag.Add(cmd.Flags(), rack)
	}

	return &cmd
}

func (a *RackAttr) create(attr string) error {
	req := pb.CreateRackAttrRequest_builder{
		Zone: a.zoneFlag.Ptr(),
		Rack: a.rackFlag.Ptr(),
		Name: &attr,
	}.Build()

	_, err := a.Client.Metal.CreateRackAttr(a.Client.Context(), req)

	return err
}

func (a *RackAttr) list(glob string) error {
	type row struct{ Zone, Rack, Attr, Value string }
	t := table.New()
	defer t.Flush()

	r := a.Client.NewRackAttrReader(a.zoneFlag.Val(), a.rackFlag.Val(), glob)

	for resp, err := range r.Responses() {
		if err != nil {
			return err
		}

		_ = t.Write(row{
			Zone:  resp.GetZone(),
			Rack:  resp.GetRack(),
			Attr:  resp.GetName(),
			Value: resp.GetValue(),
		})
	}

	return nil
}

func (a *RackAttr) update(attr string) error {
	req := pb.UpdateRackAttrRequest_builder{
		Zone: a.zoneFlag.Ptr(),
		Rack: a.rackFlag.Ptr(),
		Name: &attr,
		Fields: pb.UpdateRackAttrRequest_Fields_builder{
			Name: a.renameFlag.Ptr(),
		}.Build(),
	}.Build()

	_, err := a.Client.Metal.UpdateRackAttr(a.Client.Context(), req)

	return err
}

func (a *RackAttr) remove(glob string) error {
	req := pb.DeleteRackAttrsRequest_builder{
		Zone: a.zoneFlag.Ptr(),
		Rack: a.rackFlag.Ptr(),
		Glob: &glob,
	}.Build()

	_, err := a.Client.Metal.DeleteRackAttrs(a.Client.Context(), req)

	return err
}
