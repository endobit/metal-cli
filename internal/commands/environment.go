package commands

import (
	"github.com/spf13/cobra"

	"endobit.io/metal"
	"endobit.io/table"

	"endobit.io/metal-cli/internal/flags"
	pb "endobit.io/metal/gen/go/proto/metal/v1"
)

type Environment struct {
	Client     *metal.Client
	renameFlag flags.Rename
	zoneFlag   flags.Zone
}

type EnvironmentAttr struct {
	Client          *metal.Client
	renameFlag      flags.Rename
	zoneFlag        flags.Zone
	environmentFlag flags.Environment
}

func (e *Environment) New(verb Verb) *cobra.Command {
	var cmd cobra.Command

	switch verb {
	case Add:
		cmd = cobra.Command{
			Use:   environment + " name",
			Short: "Add an " + environment + " to a zone",
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
			Use:   environment + " name",
			Short: "Set an " + environment + "'s properties",
			Args:  cobra.ExactArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				return e.update(args[0])
			},
		}
	case List:
		cmd = cobra.Command{
			Use:   environment + " [glob]",
			Short: "List one or more " + environment + "s",
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
			Use:   environment + " glob",
			Short: "Remove one or more " + environment + "s",
			Args:  cobra.ExactArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				return e.remove(args[0])
			},
		}
	}

	e.zoneFlag.Add(cmd.Flags(), environment)

	if verb == Add || verb == Set || verb == Remove {
		e.zoneFlag.Required(cmd.Flags())
	}

	if verb == Set {
		e.renameFlag.Add(cmd.Flags(), environment)
	}

	attr := EnvironmentAttr{Client: e.Client}
	cmd.AddCommand(attr.New(verb))

	return &cmd
}

func (e *Environment) create(environment string) error {
	req := pb.CreateEnvironmentRequest_builder{
		Zone: e.zoneFlag.Ptr(),
		Name: &environment,
	}.Build()

	_, err := e.Client.Metal.CreateEnvironment(e.Client.Context(), req)

	return err
}

func (e *Environment) list(glob string) error {
	type row struct{ Zone, Environment string }
	t := table.New()
	defer t.Flush()

	r := e.Client.NewEnvironmentReader(e.zoneFlag.Val(), glob)

	for resp, err := range r.Responses() {
		if err != nil {
			return err
		}

		_ = t.Write(row{
			Zone:        resp.GetZone(),
			Environment: resp.GetName(),
		})
	}

	return nil
}

func (e *Environment) update(environment string) error {
	req := pb.UpdateEnvironmentRequest_builder{
		Zone: e.zoneFlag.Ptr(),
		Name: &environment,
		Fields: pb.UpdateEnvironmentRequest_Fields_builder{
			Name: e.renameFlag.Ptr(),
		}.Build(),
	}.Build()

	_, err := e.Client.Metal.UpdateEnvironment(e.Client.Context(), req)

	return err
}

func (e *Environment) remove(glob string) error {
	req := pb.DeleteEnvironmentsRequest_builder{
		Zone: e.zoneFlag.Ptr(),
		Glob: &glob,
	}.Build()

	_, err := e.Client.Metal.DeleteEnvironments(e.Client.Context(), req)

	return err
}

func (e *EnvironmentAttr) New(verb Verb) *cobra.Command {
	var cmd cobra.Command

	switch verb {
	case Add:
		cmd = cobra.Command{
			Use:   attribute + " name",
			Short: "Add an " + attribute + " to an " + environment,
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
			Use:   attribute + " name",
			Short: "Set an " + environment + " " + attribute + "'s properties",
			Args:  cobra.ExactArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				return e.update(args[0])
			},
		}
	case List:
		cmd = cobra.Command{
			Use:   attribute + " [glob]",
			Short: "List one or more " + environment + " " + attribute + "s",
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
			Use:   attribute + " glob",
			Short: "Remove one or more " + environment + " " + attribute + "s",
			Args:  cobra.ExactArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				return e.remove(args[0])
			},
		}
	}

	e.zoneFlag.Add(cmd.Flags(), environment)
	e.environmentFlag.Add(cmd.Flags(), attribute)

	if verb == Add || verb == Set || verb == Remove {
		e.zoneFlag.Required(cmd.Flags())
		e.environmentFlag.Required(cmd.Flags())
	}

	if verb == Set {
		e.renameFlag.Add(cmd.Flags(), environment)
	}

	return &cmd
}

func (e *EnvironmentAttr) create(attr string) error {
	req := pb.CreateEnvironmentAttrRequest_builder{
		Zone:        e.zoneFlag.Ptr(),
		Environment: e.environmentFlag.Ptr(),
		Name:        &attr,
	}.Build()

	_, err := e.Client.Metal.CreateEnvironmentAttr(e.Client.Context(), req)

	return err
}

func (e *EnvironmentAttr) list(glob string) error {
	type row struct{ Zone, Environment, Attr, Value string }
	t := table.New()
	defer t.Flush()

	r := e.Client.NewEnvironmentAttrReader(e.zoneFlag.Val(), e.environmentFlag.Val(), glob)

	for resp, err := range r.Responses() {
		if err != nil {
			return err
		}

		_ = t.Write(row{
			Zone:        resp.GetZone(),
			Environment: resp.GetEnvironment(),
			Attr:        resp.GetName(),
			Value:       resp.GetValue(),
		})
	}

	return nil
}

func (e *EnvironmentAttr) update(attr string) error {
	req := pb.UpdateEnvironmentAttrRequest_builder{
		Zone:        e.zoneFlag.Ptr(),
		Environment: e.environmentFlag.Ptr(),
		Name:        &attr,
		Fields: pb.UpdateEnvironmentAttrRequest_Fields_builder{
			Name: e.renameFlag.Ptr(),
		}.Build(),
	}.Build()

	_, err := e.Client.Metal.UpdateEnvironmentAttr(e.Client.Context(), req)

	return err
}

func (e *EnvironmentAttr) remove(glob string) error {
	req := pb.DeleteEnvironmentAttrsRequest_builder{
		Zone:        e.zoneFlag.Ptr(),
		Environment: e.environmentFlag.Ptr(),
		Glob:        &glob,
	}.Build()

	_, err := e.Client.Metal.DeleteEnvironmentAttrs(e.Client.Context(), req)

	return err
}
