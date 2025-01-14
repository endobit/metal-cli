package stack

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"

	pb "github.com/endobit/stack/internal/generated/go/proto/stack/v1"
)

func newDumpCmd(client *rpcClient) *cobra.Command {
	return &cobra.Command{
		Use:   Dump,
		Short: "Dump stack schema",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {

			d := dumper{Client: client}
			schema, err := d.Dump()
			if err != nil {
				return err
			}

			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(schema)
		},
	}
}

type dumper struct {
	Client *rpcClient
}

func (d dumper) Dump() (*pb.Schema, error) {
	attrs, err := d.GlobalAttrs()
	if err != nil {
		return nil, err
	}

	makes, err := d.Makes()
	if err != nil {
		return nil, err
	}

	models, err := d.Models()
	if err != nil {
		return nil, err
	}

	zones, err := d.Zones()
	if err != nil {
		return nil, err
	}

	return &pb.Schema{
		Attributes: attrs,
		Makes:      makes,
		Models:     models,
		Zones:      zones,
	}, nil
}

func (d dumper) GlobalAttrs() (map[string]string, error) {
	attrs := make(map[string]string)

	r := newGlobalAttrReader(d.Client, "")
	for ga, err := range r.Responses() {
		if err != nil {
			return nil, err
		}

		attrs[ga.GetName()] = ga.GetValue()
	}

	return attrs, nil
}

func (d dumper) Makes() ([]*pb.Make, error) {
	var makes []*pb.Make

	r := newMakeReader(d.Client, "")
	for m, err := range r.Responses() {
		if err != nil {
			return nil, err
		}

		mk := pb.Make{
			Name: ptr(m.GetName()),
		}

		makes = append(makes, &mk)
	}

	return makes, nil
}

func (d dumper) Models() ([]*pb.Model, error) {
	var models []*pb.Model

	r := newModelReader(d.Client, "", "")
	for m, err := range r.Responses() {
		if err != nil {
			return nil, err
		}

		model := pb.Model{
			Name: ptr(m.GetName()),
			Make: ptr(m.GetMake()),
		}

		if m.HasArchitecture() {
			model.Architecture = ptr(m.GetArchitecture())
		}

		models = append(models, &model)
	}

	return models, nil
}

func (d dumper) Zones() ([]*pb.Zone, error) {
	var zones []*pb.Zone

	r := newZoneReader(d.Client, "")
	for z, err := range r.Responses() {
		if err != nil {
			return nil, err
		}

		attrs, err := d.ZoneAttrs(z.GetName())
		if err != nil {
			return nil, err
		}

		networks, err := d.Networks(z.GetName())
		if err != nil {
			return nil, err
		}

		appliances, err := d.Appliances(z.GetName())
		if err != nil {
			return nil, err
		}

		environments, err := d.Environments(z.GetName())
		if err != nil {
			return nil, err
		}

		hosts, err := d.Hosts(z.GetName(), "")
		if err != nil {
			return nil, err
		}

		clusters, err := d.Clusters(z.GetName())
		if err != nil {
			return nil, err
		}

		zone := pb.Zone{
			Attributes:   attrs,
			Name:         ptr(z.GetName()),
			Networks:     networks,
			Appliances:   appliances,
			Environments: environments,
			Hosts:        hosts,
			Clusters:     clusters,
		}

		if z.HasTimeZone() {
			zone.TimeZone = ptr(z.GetTimeZone())
		}

		zones = append(zones, &zone)
	}

	return zones, nil
}

func (d dumper) ZoneAttrs(zone string) (map[string]string, error) {
	attrs := make(map[string]string)

	r := newZoneAttrReader(d.Client, zone, "")
	for za, err := range r.Responses() {
		if err != nil {
			return nil, err
		}

		attrs[za.GetName()] = za.GetValue()
	}

	return attrs, nil
}

func (d dumper) Networks(zone string) ([]*pb.Network, error) {
	var networks []*pb.Network

	r := newNetworkReader(d.Client, zone, "")
	for n, err := range r.Responses() {
		if err != nil {
			return nil, err
		}

		network := pb.Network{
			Name:    ptr(n.GetName()),
			Address: ptr(n.GetAddress()),
			Gateway: ptr(n.GetGateway()),
			Pxe:     ptr(n.GetPxe()),
			Mtu:     ptr(n.GetMtu()),
		}

		networks = append(networks, &network)
	}

	return networks, nil
}

func (d dumper) Appliances(zone string) ([]*pb.Appliance, error) {
	var appliances []*pb.Appliance

	r := newApplianceReader(d.Client, zone, "")
	for a, err := range r.Responses() {
		if err != nil {
			return nil, err
		}

		attrs, err := d.ApplianceAttrs(zone, a.GetName())
		if err != nil {
			return nil, err
		}

		appliance := pb.Appliance{
			Name:       ptr(a.GetName()),
			Attributes: attrs,
		}

		appliances = append(appliances, &appliance)
	}

	return appliances, nil
}

func (d dumper) ApplianceAttrs(zone, appliance string) (map[string]string, error) {
	attrs := make(map[string]string)

	r := newApplianceAttrReader(d.Client, zone, appliance, "")
	for aa, err := range r.Responses() {
		if err != nil {
			return nil, err
		}

		attrs[aa.GetName()] = aa.GetValue()
	}

	return attrs, nil
}

func (d dumper) Environments(zone string) ([]*pb.Environment, error) {
	var environments []*pb.Environment

	r := newEnvironmentReader(d.Client, zone, "")
	for e, err := range r.Responses() {
		if err != nil {
			return nil, err
		}

		attrs, err := d.EnvironmentAttrs(zone, e.GetName())
		if err != nil {
			return nil, err
		}

		environment := pb.Environment{
			Name:       ptr(e.GetName()),
			Attributes: attrs,
		}

		environments = append(environments, &environment)
	}

	return environments, nil
}

func (d dumper) EnvironmentAttrs(zone, environment string) (map[string]string, error) {
	attrs := make(map[string]string)

	r := newEnvironmentAttrReader(d.Client, zone, environment, "")
	for ea, err := range r.Responses() {
		if err != nil {
			return nil, err
		}

		attrs[ea.GetName()] = ea.GetValue()
	}

	return attrs, nil
}

func (d dumper) Hosts(zone, cluster string) ([]*pb.Host, error) {
	var hosts []*pb.Host

	r := newHostReader(d.Client, zone, cluster, "")
	for h, err := range r.Responses() {
		if err != nil {
			return nil, err
		}

		if cluster != "" && h.HasCluster() { // just want stand-alone hosts
			continue
		}

		attrs, err := d.HostAttrs(zone, cluster, h.GetName())
		if err != nil {
			return nil, err
		}

		ifaces, err := d.HostInterfaces(zone, cluster, h.GetName(), "")
		if err != nil {
			return nil, err
		}

		host := pb.Host{
			Name:       ptr(h.GetName()),
			Attributes: attrs,
			Interfaces: ifaces,
		}

		if h.HasMake() {
			host.Make = ptr(h.GetMake())
		}
		if h.HasModel() {
			host.Model = ptr(h.GetModel())
		}
		if h.HasEnvironment() {
			host.Environment = ptr(h.GetEnvironment())
		}
		if h.HasAppliance() {
			host.Appliance = ptr(h.GetAppliance())
		}
		if h.HasLocation() {
			host.Location = ptr(h.GetLocation())
		}
		if h.HasRack() {
			host.Rack = ptr(h.GetRack())
		}
		if h.HasRank() {
			host.Rank = ptr(h.GetRank())
		}
		if h.HasSlot() {
			host.Slot = ptr(h.GetSlot())
		}
		if h.HasType() {
			host.Type = ptr(h.GetType())
		}

		hosts = append(hosts, &host)
	}

	return hosts, nil
}

func (d dumper) HostAttrs(zone, cluster, host string) (map[string]string, error) {
	attrs := make(map[string]string)

	r := newHostAttrReader(d.Client, zone, cluster, host, "")
	for ha, err := range r.Responses() {
		if err != nil {
			return nil, err
		}

		attrs[ha.GetName()] = ha.GetValue()
	}

	return attrs, nil
}

func (d dumper) HostInterfaces(zone, cluster, host, iface string) ([]*pb.Host_Interface, error) {
	var ifaces []*pb.Host_Interface

	r := newHostInterfaceReader(d.Client, zone, cluster, host, iface)
	for i, err := range r.Responses() {
		if err != nil {
			return nil, err
		}

		iface := pb.Host_Interface{
			Name:       ptr(i.GetName()),
			Network:    ptr(i.GetNetwork()),
			Ip:         ptr(i.GetIp()),
			Mac:        ptr(i.GetMac()),
			Dhcp:       ptr(i.GetDhcp()),
			Pxe:        ptr(i.GetPxe()),
			Management: ptr(i.GetManagement()),
		}

		if i.HasType() {
			iface.Type = ptr(i.GetType())
		}

		if i.HasBondMode() {
			iface.BondMode = ptr(i.GetBondMode())
		}

		ifaces = append(ifaces, &iface)
	}

	return ifaces, nil
}

func (d dumper) Clusters(zone string) ([]*pb.Cluster, error) {
	var clusters []*pb.Cluster

	r := newClusterReader(d.Client, zone, "")
	for c, err := range r.Responses() {
		if err != nil {
			return nil, err
		}

		attrs, err := d.ClusterAttrs(zone, c.GetName())
		if err != nil {
			return nil, err
		}

		hosts, err := d.Hosts(zone, c.GetName())
		if err != nil {
			return nil, err
		}

		cluster := pb.Cluster{
			Attributes: attrs,
			Hosts:      hosts,
			Name:       ptr(c.GetName()),
		}

		clusters = append(clusters, &cluster)
	}

	return clusters, nil
}

func (d dumper) ClusterAttrs(zone, cluster string) (map[string]string, error) {
	attrs := make(map[string]string)

	r := newClusterAttrReader(d.Client, zone, cluster, "")
	for ca, err := range r.Responses() {
		if err != nil {
			return nil, err
		}

		attrs[ca.GetName()] = ca.GetValue()
	}

	return attrs, nil
}
