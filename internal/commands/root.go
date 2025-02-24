package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
	"gopkg.in/yaml.v2"

	"endobit.io/metal"

	"endobit.io/metal-cli/internal/flags"
	pb "endobit.io/metal/gen/go/proto/metal/v1"
)

type Root struct {
	Client      *metal.Client
	jsonFlag    flags.JSON
	zoneFlag    flags.Zone
	clusterFlag flags.Cluster
	hostFlag    flags.Host
}

func (r *Root) New(verb Verb) *cobra.Command {
	var cmd cobra.Command

	appliance := Appliance{Client: r.Client}
	environment := Environment{Client: r.Client}
	cluster := Cluster{Client: r.Client}
	rack := Rack{Client: r.Client}
	model := Model{Client: r.Client}
	zone := Zone{Client: r.Client}

	switch verb {
	case Add:
		cmd = cobra.Command{
			Use:     "add",
			Aliases: []string{"create"},
			Short:   "Add objects",
		}

		cmd.AddCommand(
			appliance.New(verb),
			cluster.New(verb),
			environment.New(verb),
			model.New(verb),
			rack.New(verb),
			zone.New(verb))

	case Dump:
		cmd = cobra.Command{
			Use:   "dump",
			Short: "Dump stack schema",
			Args:  cobra.NoArgs,
			RunE: func(_ *cobra.Command, _ []string) error {
				return r.dump()
			},
		}

		r.jsonFlag.Add(cmd.Flags(), "schema")
		r.zoneFlag.Add(cmd.Flags(), "schema")
		r.clusterFlag.Add(cmd.Flags(), "schema")
		r.hostFlag.Add(cmd.Flags(), "schema")

	case Set:
		cmd = cobra.Command{
			Use:     "set",
			Aliases: []string{"update"},
			Short:   "Set object properties",
		}

		cmd.AddCommand(
			appliance.New(verb),
			cluster.New(verb),
			environment.New(verb),
			model.New(verb),
			rack.New(verb),
			zone.New(verb))

	case List:
		cmd = cobra.Command{
			Use:     "list",
			Aliases: []string{"ls"},
			Short:   "List objects",
			Long:    "List is for humans.",
		}

		cmd.AddCommand(
			appliance.New(verb),
			cluster.New(verb),
			environment.New(verb),
			model.New(verb),
			rack.New(verb),
			zone.New(verb))

	case Load:
		cmd = cobra.Command{
			Use:     "load filename",
			Aliases: []string{"ld"},
			Args:    cobra.ExactArgs(1),
			Short:   "Load objects",
			RunE: func(_ *cobra.Command, args []string) error {
				return r.load(args[0])
			},
		}

	case Report:
		cmd = cobra.Command{
			Use:   "report",
			Short: "Report objects",
			Long:  "Report is for computers.",
		}

	case Remove:
		cmd = cobra.Command{
			Use:     "remove",
			Aliases: []string{"del", "rm"},
			Short:   "Remove objects",
		}

		cmd.AddCommand(
			appliance.New(verb),
			cluster.New(verb),
			environment.New(verb),
			model.New(verb),
			rack.New(verb),
			zone.New(verb))
	}

	return &cmd
}

func (r *Root) dump() error {
	var req pb.ReadSchemaRequest

	if r.zoneFlag.Val() != "" {
		req.SetZone(r.zoneFlag.Val())
	}
	if r.clusterFlag.Val() != "" {
		req.SetCluster(r.clusterFlag.Val())
	}
	if r.hostFlag.Val() != "" {
		req.SetHost(r.hostFlag.Val())
	}

	resp, err := r.Client.Metal.ReadSchema(r.Client.Context(), &req)
	if err != nil {
		return err
	}

	doc := resp.GetSchema()

	if !r.jsonFlag.Val() { // parse json as yaml and re-marshal
		var obj map[string]interface{}

		b, err := protojson.MarshalOptions{
			UseProtoNames: true,
		}.Marshal(doc)
		if err != nil {
			return err
		}

		if err := yaml.Unmarshal(b, &obj); err != nil {
			return err
		}

		return yaml.NewEncoder(os.Stdout).Encode(obj)
	}

	b, err := protojson.MarshalOptions{
		Multiline:     true,
		UseProtoNames: true,
	}.Marshal(doc)
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	return nil
}

func (r *Root) load(filename string) error {
	fin, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fin.Close()

	data, err := io.ReadAll(fin)
	if err != nil {
		return err
	}

	var doc pb.Schema

	switch filepath.Ext(filename) {
	case ".json":
		if err := protojson.Unmarshal(data, &doc); err != nil {
			return err
		}
	case ".yaml", ".yml":
		var jsonMap map[string]interface{}

		if err := yaml.Unmarshal(data, &jsonMap); err != nil {
			return err
		}

		jsonData, err := json.Marshal(jsonMap)
		if err != nil {
			return err
		}

		if err := protojson.Unmarshal(jsonData, &doc); err != nil {
			return err
		}

	default:
		return errors.New("unknown file type")
	}

	req := pb.CreateSchemaRequest_builder{
		Schema: &doc,
	}.Build()

	_, err = r.Client.Metal.CreateSchema(r.Client.Context(), req)

	return err
}
