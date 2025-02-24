package commands

import (
	"github.com/spf13/cobra"

	"endobit.io/metal"
	"endobit.io/table"

	"endobit.io/metal-cli/internal/flags"
	pb "endobit.io/metal/gen/go/proto/metal/v1"
)

type Appliance struct {
	Client     *metal.Client
	renameFlag flags.Rename
	zoneFlag   flags.Zone
}

type ApplianceAttr struct {
	Client        *metal.Client
	renameFlag    flags.Rename
	zoneFlag      flags.Zone
	applianceFlag flags.Appliance
}

func (a *Appliance) New(verb Verb) *cobra.Command {
	var cmd cobra.Command

	switch verb {
	case Add:
		cmd = cobra.Command{
			Use:   appliance + " name",
			Short: "Add an " + appliance + " to a zone",
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
			Use:   appliance + " name",
			Short: "Set an " + appliance + "'s properties",
			Args:  cobra.ExactArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				return a.update(args[0])
			},
		}
	case List:
		cmd = cobra.Command{
			Use:   appliance + " [glob]",
			Short: "List one or more " + appliance + "s",
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
			Use:   appliance + " glob",
			Short: "Remove one or more " + appliance + "s",
			Args:  cobra.ExactArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				return a.remove(args[0])
			},
		}
	}

	a.zoneFlag.Add(cmd.Flags(), appliance)

	if verb == Add || verb == Set || verb == Remove {
		a.zoneFlag.Required(cmd.Flags())
	}

	if verb == Set {
		a.renameFlag.Add(cmd.Flags(), appliance)
	}

	attr := ApplianceAttr{Client: a.Client}
	cmd.AddCommand(attr.New(verb))

	return &cmd
}

func (a *Appliance) create(appliance string) error {
	req := pb.CreateApplianceRequest_builder{
		Zone: a.zoneFlag.Ptr(),
		Name: &appliance,
	}.Build()

	_, err := a.Client.Metal.CreateAppliance(a.Client.Context(), req)

	return err
}

func (a *Appliance) list(glob string) error {
	type row struct{ Zone, Appliance string }
	t := table.New()
	defer t.Flush()

	r := a.Client.NewApplianceReader(a.zoneFlag.Val(), glob)

	for resp, err := range r.Responses() {
		if err != nil {
			return err
		}

		_ = t.Write(row{
			Zone:      resp.GetZone(),
			Appliance: resp.GetName(),
		})
	}

	return nil
}

func (a *Appliance) update(appliance string) error {
	req := pb.UpdateApplianceRequest_builder{
		Zone: a.zoneFlag.Ptr(),
		Name: &appliance,
		Fields: pb.UpdateApplianceRequest_Fields_builder{
			Name: a.renameFlag.Ptr(),
		}.Build(),
	}.Build()

	_, err := a.Client.Metal.UpdateAppliance(a.Client.Context(), req)

	return err
}

func (a *Appliance) remove(glob string) error {
	req := pb.DeleteAppliancesRequest_builder{
		Zone: a.zoneFlag.Ptr(),
		Glob: &glob,
	}.Build()

	_, err := a.Client.Metal.DeleteAppliances(a.Client.Context(), req)

	return err
}

func (a *ApplianceAttr) New(verb Verb) *cobra.Command {
	var cmd cobra.Command

	switch verb {
	case Add:
		cmd = cobra.Command{
			Use:   attribute + " name",
			Short: "Add an " + attribute + " to an " + appliance,
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
			Short: "Set an " + appliance + " " + attribute + "'s properties",
			Args:  cobra.ExactArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				return a.update(args[0])
			},
		}
	case List:
		cmd = cobra.Command{
			Use:   attribute + " [glob]",
			Short: "List one or more " + appliance + " " + attribute + "s",
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
			Short: "Remove one or more " + appliance + " " + attribute + "s",
			Args:  cobra.ExactArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				return a.remove(args[0])
			},
		}
	}

	a.zoneFlag.Add(cmd.Flags(), appliance)
	a.applianceFlag.Add(cmd.Flags(), attribute)

	if verb == Add || verb == Set || verb == Remove {
		a.zoneFlag.Required(cmd.Flags())
		a.applianceFlag.Required(cmd.Flags())
	}

	if verb == Set {
		a.renameFlag.Add(cmd.Flags(), appliance)
	}

	return &cmd
}

func (a *ApplianceAttr) create(attr string) error {
	req := pb.CreateApplianceAttrRequest_builder{
		Zone:      a.zoneFlag.Ptr(),
		Appliance: a.applianceFlag.Ptr(),
		Name:      &attr,
	}.Build()

	_, err := a.Client.Metal.CreateApplianceAttr(a.Client.Context(), req)

	return err
}

func (a *ApplianceAttr) list(glob string) error {
	type row struct{ Zone, Appliance, Attr, Value string }
	t := table.New()
	defer t.Flush()

	r := a.Client.NewApplianceAttrReader(a.zoneFlag.Val(), a.applianceFlag.Val(), glob)

	for resp, err := range r.Responses() {
		if err != nil {
			return err
		}

		_ = t.Write(row{
			Zone:      resp.GetZone(),
			Appliance: resp.GetAppliance(),
			Attr:      resp.GetName(),
			Value:     resp.GetValue(),
		})
	}

	return nil
}

func (a *ApplianceAttr) update(attr string) error {
	req := pb.UpdateApplianceAttrRequest_builder{
		Zone:      a.zoneFlag.Ptr(),
		Appliance: a.applianceFlag.Ptr(),
		Name:      &attr,
		Fields: pb.UpdateApplianceAttrRequest_Fields_builder{
			Name: a.renameFlag.Ptr(),
		}.Build(),
	}.Build()

	_, err := a.Client.Metal.UpdateApplianceAttr(a.Client.Context(), req)

	return err
}

func (a *ApplianceAttr) remove(glob string) error {
	req := pb.DeleteApplianceAttrsRequest_builder{
		Zone:      a.zoneFlag.Ptr(),
		Appliance: a.applianceFlag.Ptr(),
		Glob:      &glob,
	}.Build()

	_, err := a.Client.Metal.DeleteApplianceAttrs(a.Client.Context(), req)

	return err
}
