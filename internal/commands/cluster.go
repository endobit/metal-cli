package commands

import (
	"github.com/spf13/cobra"

	"endobit.io/metal"
	"endobit.io/metal-cli/internal/flags"
	pb "endobit.io/metal/gen/go/proto/metal/v1"
	"endobit.io/table"
)

type Cluster struct {
	Client     *metal.Client
	renameFlag flags.Rename
	zoneFlag   flags.Zone
}

type ClusterAttr struct {
	Client      *metal.Client
	renameFlag  flags.Rename
	zoneFlag    flags.Zone
	clusterFlag flags.Cluster
}

func (c *Cluster) New(verb Verb) *cobra.Command {
	var cmd cobra.Command

	switch verb {
	case Add:
		cmd = cobra.Command{
			Use:   cluster + " name",
			Short: "Add a " + cluster,
			Args:  cobra.ExactArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				if err := c.create(args[0]); err != nil {
					return err
				}

				return c.update(args[0])
			},
		}
	case Set:
		cmd = cobra.Command{
			Use:   cluster + " name",
			Short: "Set a " + cluster + "'s properties",
			Args:  cobra.ExactArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				return c.update(args[0])
			},
		}
	case List:
		cmd = cobra.Command{
			Use:   cluster + " [glob]",
			Short: "List one or more " + cluster + "s",
			Args:  cobra.MaximumNArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				var glob string

				if len(args) > 0 {
					glob = args[0]
				}
				return c.list(glob)
			},
		}
	case Remove:
		cmd = cobra.Command{
			Use:   cluster + " glob",
			Short: "Remove one or more " + cluster + "s",
			Args:  cobra.ExactArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				return c.remove(args[0])
			},
		}
	}

	if verb == Add || verb == Set || verb == Remove {
		c.zoneFlag.Add(cmd.PersistentFlags(), cluster)
		c.zoneFlag.Required(cmd.PersistentFlags())
	}

	if verb == Set {
		c.renameFlag.Add(cmd.Flags(), cluster)
	}

	return &cmd
}

func (c *Cluster) create(cluster string) error {
	var req pb.CreateClusterRequest

	req.SetName(cluster)
	_, err := c.Client.Metal.CreateCluster(c.Client.Context(), &req)

	return err
}

func (c *Cluster) list(glob string) error {
	type row struct{ Zone, Cluster string }
	t := table.New()
	defer t.Flush()

	r := c.Client.NewClusterReader(c.zoneFlag.Val(), glob)

	for resp, err := range r.Responses() {
		if err != nil {
			return err
		}

		_ = t.Write(row{
			Zone:    resp.GetZone(),
			Cluster: resp.GetName(),
		})
	}

	return nil
}

func (c *Cluster) update(cluster string) error {
	req := pb.UpdateClusterRequest_builder{
		Zone: c.zoneFlag.Ptr(),
		Name: &cluster,
		Fields: pb.UpdateClusterRequest_Fields_builder{
			Name: c.renameFlag.Ptr(),
		}.Build(),
	}.Build()

	_, err := c.Client.Metal.UpdateCluster(c.Client.Context(), req)

	return err
}

func (c *Cluster) remove(glob string) error {
	req := pb.DeleteClustersRequest_builder{
		Zone: c.zoneFlag.Ptr(),
		Glob: &glob,
	}.Build()

	_, err := c.Client.Metal.DeleteClusters(c.Client.Context(), req)

	return err
}

func (c *ClusterAttr) New(verb Verb) *cobra.Command {
	var cmd cobra.Command

	switch verb {
	case Add:
		cmd = cobra.Command{
			Use:   attribute + " name",
			Short: "Add an " + attribute + " to a " + cluster,
			Args:  cobra.ExactArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				if err := c.create(args[0]); err != nil {
					return err
				}

				return c.update(args[0])
			},
		}
	case Set:
		cmd = cobra.Command{
			Use:   attribute + " name",
			Short: "Set a " + cluster + " " + attribute + "'s properties",
			Args:  cobra.ExactArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				return c.update(args[0])
			},
		}
	case List:
		cmd = cobra.Command{
			Use:   attribute + " [glob]",
			Short: "List one or more " + cluster + " " + attribute + "s",
			Args:  cobra.MaximumNArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				var glob string

				if len(args) > 0 {
					glob = args[0]
				}
				return c.list(glob)
			},
		}
	case Remove:
		cmd = cobra.Command{
			Use:   attribute + " glob",
			Short: "Remove one or more " + cluster + " " + attribute + "s",
			Args:  cobra.ExactArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				return c.remove(args[0])
			},
		}
	}

	c.zoneFlag.Add(cmd.Flags(), cluster)
	c.clusterFlag.Add(cmd.Flags(), attribute)

	if verb == Add || verb == Set || verb == Remove {
		c.zoneFlag.Required(cmd.Flags())
		c.clusterFlag.Required(cmd.Flags())
	}

	if verb == Set {
		c.renameFlag.Add(cmd.Flags(), cluster)
	}

	return &cmd
}

func (c *ClusterAttr) create(attr string) error {
	req := pb.CreateClusterAttrRequest_builder{
		Zone:    c.zoneFlag.Ptr(),
		Cluster: c.clusterFlag.Ptr(),
		Name:    &attr,
	}.Build()

	_, err := c.Client.Metal.CreateClusterAttr(c.Client.Context(), req)

	return err
}

func (c *ClusterAttr) list(glob string) error {
	type row struct{ Zone, Cluster, Attr, Value string }
	t := table.New()
	defer t.Flush()

	r := c.Client.NewClusterAttrReader(c.zoneFlag.Val(), c.clusterFlag.Val(), glob)

	for resp, err := range r.Responses() {
		if err != nil {
			return err
		}

		_ = t.Write(row{
			Zone:    resp.GetZone(),
			Cluster: resp.GetCluster(),
			Attr:    resp.GetName(),
			Value:   resp.GetValue(),
		})
	}

	return nil
}

func (c *ClusterAttr) update(attr string) error {
	req := pb.UpdateClusterAttrRequest_builder{
		Zone:    c.zoneFlag.Ptr(),
		Cluster: c.clusterFlag.Ptr(),
		Name:    &attr,
		Fields: pb.UpdateClusterAttrRequest_Fields_builder{
			Name: c.renameFlag.Ptr(),
		}.Build(),
	}.Build()

	_, err := c.Client.Metal.UpdateClusterAttr(c.Client.Context(), req)

	return err
}

func (c *ClusterAttr) remove(glob string) error {
	req := pb.DeleteClusterAttrsRequest_builder{
		Zone:    c.zoneFlag.Ptr(),
		Cluster: c.clusterFlag.Ptr(),
		Glob:    &glob,
	}.Build()

	_, err := c.Client.Metal.DeleteClusterAttrs(c.Client.Context(), req)

	return err
}
