package commands

import (
	"github.com/spf13/cobra"

	"endobit.io/metal"
	"endobit.io/metal-cli/internal/flags"
	pb "endobit.io/metal/gen/go/proto/metal/v1"
	"endobit.io/table"
)

type Model struct {
	Client     *metal.Client
	makeFlag   flags.Make
	archFlag   flags.Arch
	renameFlag flags.Rename
}

type ModelAttr struct {
	Client     *metal.Client
	makeFlag   flags.Make
	renameFlag flags.Rename
	modelFlag  flags.Model
}

func (m *Model) New(verb Verb) *cobra.Command {
	var cmd cobra.Command

	switch verb {
	case Add:
		cmd = cobra.Command{
			Use:   model + "make name",
			Short: "Add a " + model,
			Args:  cobra.ExactArgs(2),
			RunE: func(_ *cobra.Command, args []string) error {
				if err := m.create(args[0], args[1]); err != nil {
					return err
				}

				return m.update(args[0], args[1])
			},
		}
	case Set:
		cmd = cobra.Command{
			Use:   model + "make name",
			Short: "Set a " + model + "'s properties",
			Args:  cobra.ExactArgs(2),
			RunE: func(_ *cobra.Command, args []string) error {
				return m.update(args[0], args[1])
			},
		}
	case List:
		cmd = cobra.Command{
			Use:   model + "make [glob]",
			Short: "List one or more " + model + "s",
			Args:  cobra.RangeArgs(1, 2),
			RunE: func(_ *cobra.Command, args []string) error {
				var glob string

				if len(args) > 1 {
					glob = args[1]
				}

				return m.list(args[0], glob)
			},
		}
	case Remove:
		cmd = cobra.Command{
			Use:   model + "make glob",
			Short: "Remove one or more " + model + "s",
			Args:  cobra.ExactArgs(2),
			RunE: func(_ *cobra.Command, args []string) error {
				return m.remove(args[0], args[1])
			},
		}
	}

	if verb == Add || verb == Set {
		m.makeFlag.Add(cmd.Flags(), model)
		m.archFlag.Add(cmd.Flags(), model)
	}

	if verb == Set {
		m.renameFlag.Add(cmd.Flags(), model)
	}

	return &cmd
}

func (m *Model) create(vendor, model string) error {
	req := pb.CreateModelRequest_builder{
		Make: &vendor,
		Name: &model,
	}.Build()

	_, err := m.Client.Metal.CreateModel(m.Client.Context(), req)
	return err
}

func (m *Model) list(vendor, glob string) error {
	type row struct{ Make, Model, Arch string }
	t := table.New()
	defer t.Flush()

	r := m.Client.NewModelReader(vendor, glob)

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

func (m *Model) update(vendor, model string) error {
	var pbarch *pb.Architecture

	if m.archFlag.Val() != "" {
		a := pb.Architecture(pb.Architecture_value[m.archFlag.Val()])
		pbarch = &a
	}

	req := pb.UpdateModelRequest_builder{
		Make: &vendor,
		Name: &model,
		Fields: pb.UpdateModelRequest_Fields_builder{
			Name:         m.renameFlag.Ptr(),
			Architecture: pbarch,
		}.Build(),
	}.Build()

	_, err := m.Client.Metal.UpdateModel(m.Client.Context(), req)

	return err
}

func (m *Model) remove(vendor, glob string) error {
	req := pb.DeleteModelsRequest_builder{
		Make: &vendor,
		Glob: &glob,
	}.Build()

	_, err := m.Client.Metal.DeleteModels(m.Client.Context(), req)
	return err
}

func (a *ModelAttr) New(verb Verb) *cobra.Command {
	var cmd cobra.Command

	switch verb {
	case Add:
		cmd = cobra.Command{
			Use:   attribute + " name",
			Short: "Add an " + attribute + " to a " + model,
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
			Short: "Set a " + model + " " + attribute + "'s properties",
			Args:  cobra.ExactArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				return a.update(args[0])
			},
		}
	case List:
		cmd = cobra.Command{
			Use:   attribute + " [glob]",
			Short: "List one or more " + model + " " + attribute + "s",
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
			Short: "Remove one or more " + model + " " + attribute + "s",
			Args:  cobra.ExactArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				return a.remove(args[0])
			},
		}
	}

	a.makeFlag.Add(cmd.Flags(), model)
	a.modelFlag.Add(cmd.Flags(), attribute)

	if verb == Add || verb == Set || verb == Remove {
		a.makeFlag.Required(cmd.Flags())
		a.modelFlag.Required(cmd.Flags())
	}

	if verb == Set {
		a.renameFlag.Add(cmd.Flags(), model)
	}

	return &cmd
}

func (a *ModelAttr) create(attr string) error {
	req := pb.CreateModelAttrRequest_builder{
		Model: a.modelFlag.Ptr(),
		Name:  &attr,
	}.Build()

	_, err := a.Client.Metal.CreateModelAttr(a.Client.Context(), req)

	return err
}

func (a *ModelAttr) list(glob string) error {
	type row struct{ Make, Model, Attr, Value string }
	t := table.New()
	defer t.Flush()

	r := a.Client.NewModelAttrReader(a.modelFlag.Val(), glob)

	for resp, err := range r.Responses() {
		if err != nil {
			return err
		}

		_ = t.Write(row{
			Make:  resp.GetMake(),
			Model: resp.GetModel(),
			Attr:  resp.GetName(),
			Value: resp.GetValue(),
		})
	}

	return nil
}

func (a *ModelAttr) update(attr string) error {
	req := pb.UpdateModelAttrRequest_builder{
		Model: a.modelFlag.Ptr(),
		Name:  &attr,
		Fields: pb.UpdateModelAttrRequest_Fields_builder{
			Name: a.renameFlag.Ptr(),
		}.Build(),
	}.Build()

	_, err := a.Client.Metal.UpdateModelAttr(a.Client.Context(), req)

	return err
}

func (a *ModelAttr) remove(glob string) error {
	req := pb.DeleteModelAttrsRequest_builder{
		Model: a.modelFlag.Ptr(),
		Glob:  &glob,
	}.Build()

	_, err := a.Client.Metal.DeleteModelAttrs(a.Client.Context(), req)

	return err
}
