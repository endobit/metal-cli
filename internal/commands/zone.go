package commands

import (
	"github.com/spf13/cobra"

	"endobit.io/metal"
	"endobit.io/metal-cli/internal/flags"
	pb "endobit.io/metal/gen/go/proto/metal/v1"
	"endobit.io/table"
)

type Zone struct {
	Client       *metal.Client
	renameFlag   flags.Rename
	timeZoneFlag flags.TimeZone
	templateFlag flags.Template
}

type ZoneAttr struct {
	Client     *metal.Client
	renameFlag flags.Rename
	zoneFlag   flags.Zone
}

func (z *Zone) New(verb Verb) *cobra.Command {
	var cmd cobra.Command

	switch verb {
	case Add:
		cmd = cobra.Command{
			Use:   zone + " name",
			Short: "Add a " + zone,
			Args:  cobra.ExactArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				if err := z.create(args[0]); err != nil {
					return err
				}

				return z.update(args[0])
			},
		}
	case Set:
		cmd = cobra.Command{
			Use:   zone + " name",
			Short: "Set a " + zone + "'s properties",
			Args:  cobra.ExactArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				return z.update(args[0])
			},
		}
	case List:
		cmd = cobra.Command{
			Use:   zone + " [glob]",
			Short: "List one or more " + zone + "s",
			Args:  cobra.MaximumNArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				var glob string

				if len(args) > 0 {
					glob = args[0]
				}
				return z.list(glob)
			},
		}
	case Remove:
		cmd = cobra.Command{
			Use:   zone + " glob",
			Short: "Remove one or more " + zone + "s",
			Args:  cobra.ExactArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				return z.remove(args[0])
			},
		}
	}

	if verb == Add || verb == Set {
		z.timeZoneFlag.Add(cmd.Flags(), zone)
	}

	if verb == Set {
		z.renameFlag.Add(cmd.Flags(), zone)
	}

	if verb == Report {
		z.templateFlag.Add(cmd.Flags(), zone)
	}

	return &cmd
}

func (z *Zone) create(zone string) error {
	var req pb.CreateZoneRequest

	req.SetName(zone)
	_, err := z.Client.Metal.CreateZone(z.Client.Context(), &req)

	return err
}

func (z *Zone) list(glob string) error {
	type row struct{ Zone, TimeZone string }
	t := table.New()
	defer t.Flush()

	r := z.Client.NewZoneReader(glob)

	for resp, err := range r.Responses() {
		if err != nil {
			return err
		}

		_ = t.Write(row{
			Zone:     resp.GetName(),
			TimeZone: resp.GetTimeZone(),
		})
	}

	return nil
}

func (z *Zone) update(zone string) error {
	req := pb.UpdateZoneRequest_builder{
		Name: &zone,
		Fields: pb.UpdateZoneRequest_Fields_builder{
			Name:     z.renameFlag.Ptr(),
			TimeZone: z.timeZoneFlag.Ptr(),
		}.Build(),
	}.Build()

	_, err := z.Client.Metal.UpdateZone(z.Client.Context(), req)

	return err
}

func (z *Zone) remove(glob string) error {
	var req pb.DeleteZonesRequest

	req.SetGlob(glob)
	_, err := z.Client.Metal.DeleteZones(z.Client.Context(), &req)

	return err
}

func (a *ZoneAttr) New(verb Verb) *cobra.Command {
	var cmd cobra.Command

	switch verb {
	case Add:
		cmd = cobra.Command{
			Use:   attribute + " name",
			Short: "Add an " + attribute + " to a " + zone,
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
			Short: "Set a " + zone + " " + attribute + "'s properties",
			Args:  cobra.ExactArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				return a.update(args[0])
			},
		}
	case List:
		cmd = cobra.Command{
			Use:   attribute + " [glob]",
			Short: "List one or more " + zone + " " + attribute + "s",
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
			Short: "Remove one or more " + zone + " " + attribute + "s",
			Args:  cobra.ExactArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				return a.remove(args[0])
			},
		}
	}

	a.zoneFlag.Add(cmd.Flags(), zone)

	if verb == Add || verb == Set || verb == Remove {
		a.zoneFlag.Required(cmd.Flags())
	}

	if verb == Set {
		a.renameFlag.Add(cmd.Flags(), zone)
	}

	return &cmd
}

func (a *ZoneAttr) create(attr string) error {
	req := pb.CreateZoneAttrRequest_builder{
		Zone: a.zoneFlag.Ptr(),
		Name: &attr,
	}.Build()

	_, err := a.Client.Metal.CreateZoneAttr(a.Client.Context(), req)

	return err
}

func (a *ZoneAttr) list(glob string) error {
	type row struct{ Zone, Attr, Value string }
	t := table.New()
	defer t.Flush()

	r := a.Client.NewZoneAttrReader(a.zoneFlag.Val(), glob)

	for resp, err := range r.Responses() {
		if err != nil {
			return err
		}

		_ = t.Write(row{
			Zone:  resp.GetZone(),
			Attr:  resp.GetName(),
			Value: resp.GetValue(),
		})
	}

	return nil
}

func (a *ZoneAttr) update(attr string) error {
	req := pb.UpdateZoneAttrRequest_builder{
		Zone: a.zoneFlag.Ptr(),
		Name: &attr,
		Fields: pb.UpdateZoneAttrRequest_Fields_builder{
			Name: a.renameFlag.Ptr(),
		}.Build(),
	}.Build()

	_, err := a.Client.Metal.UpdateZoneAttr(a.Client.Context(), req)

	return err
}

func (a *ZoneAttr) remove(glob string) error {
	req := pb.DeleteZoneAttrsRequest_builder{
		Zone: a.zoneFlag.Ptr(),
		Glob: &glob,
	}.Build()

	_, err := a.Client.Metal.DeleteZoneAttrs(a.Client.Context(), req)

	return err
}
