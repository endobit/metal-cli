package stack

import (
	"github.com/spf13/cobra"

	pb "github.com/endobit/stack/internal/generated/go/proto/stack/v1"
	"github.com/endobit/table"
)

func newAddZoneCmd(client *rpcClient) *cobra.Command {
	var timeZone string

	cmd := cobra.Command{
		Use:   "zone name",
		Short: "Add a zone",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			req := pb.CreateZoneRequest_builder{
				Name:     required(args[0]),
				TimeZone: optional(timeZone),
			}.Build()

			_, err := client.stack.CreateZone(client.Context(), req)

			return err
		},
	}

	cmd.Flags().StringVar(&timeZone, timeZoneFlag, "", "Time zone for the zone")

	return &cmd
}

func newListZoneCmd(client *rpcClient) *cobra.Command {
	return &cobra.Command{
		Use:   "zone [glob]",
		Short: "List one or more zones",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			var glob string

			if len(args) != 0 {
				glob = args[0]
			}

			type row struct{ Name, TimeZone string }
			t := table.New()
			defer t.Flush()

			r := newZoneReader(client, glob)

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
		},
	}
}

func newUpdateZoneCmd(client *rpcClient) *cobra.Command {
	var name, timeZone string

	cmd := cobra.Command{
		Use:   "zone name",
		Short: "Set a zone's properties",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			req := pb.UpdateZoneRequest_builder{
				Name: required(args[0]),
				Fields: pb.UpdateZoneRequest_Fields_builder{
					Name:     optional(name),
					TimeZone: optional(timeZone),
				}.Build(),
			}.Build()

			_, err := client.stack.UpdateZone(client.Context(), req)

			return err
		},
	}

	cmd.Flags().StringVar(&name, nameFlag, "", "New name for the zone")
	cmd.Flags().StringVar(&timeZone, timeZoneFlag, "", "New time zone for the zone")

	return &cmd
}

func newRemoveZoneCmd(client *rpcClient) *cobra.Command {
	return &cobra.Command{
		Use:   "zone glob",
		Short: "Remove one or more zones",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			req := pb.DeleteZonesRequest_builder{
				Glob: required(args[0]),
			}.Build()

			_, err := client.stack.DeleteZones(client.Context(), req)
			if err != nil {
				return err
			}

			return nil
		},
	}
}
