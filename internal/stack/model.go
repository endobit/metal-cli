package stack

import (
	"github.com/spf13/cobra"

	pb "github.com/endobit/stack/internal/generated/go/proto/stack/v1"
	"github.com/endobit/table"
)

func newAddModelCmd(client *rpcClient) *cobra.Command {
	var arch string

	cmd := cobra.Command{
		Use:   "model make model",
		Short: "Add a model",
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			create := pb.CreateModelRequest_builder{
				Make: required(args[0]),
				Name: required(args[0]),
			}.Build()
			if _, err := client.stack.CreateModel(client.Context(), create); err != nil {
				return err
			}

			var fields pb.UpdateModelRequest_Fields

			if arch != "" {
				fields.SetArchitecture(pb.Architecture(pb.Architecture_value[arch]))
			}

			update := pb.UpdateModelRequest_builder{
				Make:   required(args[0]),
				Name:   required(args[0]),
				Fields: &fields,
			}.Build()
			_, err := client.stack.UpdateModel(client.Context(), update)

			return err
		},
	}

	cmd.Flags().StringVar(&arch, archFlag, "", "Architecture for the model")

	return &cmd
}

func newListModelCmd(client *rpcClient) *cobra.Command {
	return &cobra.Command{
		Use:   "model make [model_glob]",
		Short: "List one or more models",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(_ *cobra.Command, args []string) error {
			var glob string

			if len(args) == 2 {
				glob = args[1]
			}

			type row struct{ Make, Model, Arch string }
			t := table.New()
			defer t.Flush()

			r := newModelReader(client, args[0], glob)

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
		},
	}
}

func newUpdateModelCmd(client *rpcClient) *cobra.Command {
	var name, arch string

	cmd := cobra.Command{
		Use:   "model make model",
		Short: "Set a model's properties",
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			var pbarch *pb.Architecture

			if arch != "" {
				a := pb.Architecture(pb.Architecture_value[arch])
				pbarch = &a

			}
			req := pb.UpdateModelRequest_builder{
				Make: required(args[0]),
				Name: required(args[1]),
				Fields: pb.UpdateModelRequest_Fields_builder{
					Name:         optional(name),
					Architecture: pbarch,
				}.Build(),
			}.Build()

			_, err := client.stack.UpdateModel(client.Context(), req)

			return err
		},
	}

	cmd.Flags().StringVar(&name, modelFlag, "", "New name for the model")
	cmd.Flags().StringVar(&arch, archFlag, "", "Architecture for the model")

	return &cmd
}

func newRemoveModelCmd(client *rpcClient) *cobra.Command {
	return &cobra.Command{
		Use:   "model make model_glob",
		Short: "Remove one or more models",
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			req := pb.DeleteModelsRequest_builder{
				Make: required(args[0]),
				Glob: required(args[1]),
			}.Build()

			_, err := client.stack.DeleteModels(client.Context(), req)
			if err != nil {
				return err
			}

			return nil
		},
	}
}
