package cli

// import (
// 	"fmt"

// 	pb "endobit.io/metal/gen/go/proto/metal/v1"
// )

// type Dumper struct {
// 	*RPC
// 	Filter Filter
// }

// type Filter struct {
// 	Zone    string
// 	Cluster string
// 	Host    string
// }

// func (d Dumper) Dump() (schema *pb.Schema, err error) {
// 	defer func() {
// 		if err != nil {
// 			err = fmt.Errorf("dump failed: %w", err)
// 		}
// 	}()

// 	attrs, err := d.GlobalAttrs()
// 	if err != nil {
// 		return nil, err
// 	}

// 	makes, err := d.Makes()
// 	if err != nil {
// 		return nil, err
// 	}

// 	models, err := d.Models()
// 	if err != nil {
// 		return nil, err
// 	}

// 	zones, err := d.Zones()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &pb.Schema{
// 		Attributes: attrs,
// 		Makes:      makes,
// 		Zones:      zones,
// 	}, nil
// }

// func (d Dumper) GlobalAttrs() (map[string]string, error) {
// 	attrs := make(map[string]string)

// 	r := d.NewGlobalAttrReader("")
// 	for ga, err := range r.Responses() {
// 		if err != nil {
// 			return nil, err
// 		}

// 		attrs[ga.GetName()] = ga.GetValue()
// 	}

// 	return attrs, nil
// }

// func (d Dumper) Makes() (makes []*pb.Make, err error) {
// 	var name string

// 	defer func() {
// 		if err != nil {
// 			err = fmt.Errorf("make %q: %w", name, err)
// 		}
// 	}()

// 	r := d.NewMakeReader("")
// 	for m, err := range r.Responses() {
// 		name = m.GetName()

// 		if err != nil {
// 			return nil, err
// 		}

// 		mk := pb.Make{
// 			Name: Ptr(m.GetName()),
// 		}

// 		makes = append(makes, &mk)
// 	}

// 	return makes, nil
// }

// func (d Dumper) Models() (models []*pb.Model, err error) {
// 	var name string

// 	defer func() {
// 		if err != nil {
// 			err = fmt.Errorf("model %q: %w", name, err)
// 		}
// 	}()

// 	r := d.NewModelReader("", "")
// 	for m, err := range r.Responses() {
// 		name = m.GetName()

// 		if err != nil {
// 			return nil, err
// 		}

// 		model := pb.Model{
// 			Name: Ptr(m.GetName()),
// 		}

// 		if m.HasArchitecture() {
// 			model.Architecture = Ptr(m.GetArchitecture())
// 		}

// 		models = append(models, &model)
// 	}

// 	return models, nil
// }

// func (d Dumper) Zones() (zones []*pb.Zone, err error) {
// 	var name string

// 	defer func() {
// 		if err != nil {
// 			err = fmt.Errorf("zone %q: %w", name, err)
// 		}
// 	}()

// 	r := d.NewZoneReader(d.Filter.Zone)
// 	for z, err := range r.Responses() {
// 		name = z.GetName()

// 		if err != nil {
// 			return nil, err
// 		}

// 		attrs, err := d.ZoneAttrs(z.GetName())
// 		if err != nil {
// 			return nil, err
// 		}

// 		networks, err := d.Networks(z.GetName())
// 		if err != nil {
// 			return nil, err
// 		}

// 		appliances, err := d.Appliances(z.GetName())
// 		if err != nil {
// 			return nil, err
// 		}

// 		environments, err := d.Environments(z.GetName())
// 		if err != nil {
// 			return nil, err
// 		}

// 		hosts, err := d.Hosts(z.GetName(), "")
// 		if err != nil {
// 			return nil, err
// 		}

// 		clusters, err := d.Clusters(z.GetName())
// 		if err != nil {
// 			return nil, err
// 		}

// 		zone := pb.Zone{
// 			Attributes:   attrs,
// 			Name:         Ptr(z.GetName()),
// 			Networks:     networks,
// 			Appliances:   appliances,
// 			Environments: environments,
// 			Hosts:        hosts,
// 			Clusters:     clusters,
// 		}

// 		if z.HasTimeZone() {
// 			zone.TimeZone = Ptr(z.GetTimeZone())
// 		}

// 		zones = append(zones, &zone)
// 	}

// 	return zones, nil
// }

// func (d Dumper) ZoneAttrs(zone string) (map[string]string, error) {
// 	attrs := make(map[string]string)

// 	r := d.NewZoneAttrReader(zone, "")
// 	for za, err := range r.Responses() {
// 		if err != nil {
// 			return nil, err
// 		}

// 		attrs[za.GetName()] = za.GetValue()
// 	}

// 	return attrs, nil
// }

// func (d Dumper) Networks(zone string) (networks []*pb.Network, err error) {
// 	var name string

// 	defer func() {
// 		if err != nil {
// 			err = fmt.Errorf("network %q: %w", name, err)
// 		}
// 	}()

// 	r := d.NewNetworkReader(zone, "")
// 	for n, err := range r.Responses() {
// 		name = n.GetName()

// 		if err != nil {
// 			return nil, err
// 		}

// 		network := pb.Network{
// 			Name:    Ptr(n.GetName()),
// 			Address: Ptr(n.GetAddress()),
// 			Gateway: Ptr(n.GetGateway()),
// 			Pxe:     Ptr(n.GetPxe()),
// 			Mtu:     Ptr(n.GetMtu()),
// 		}

// 		networks = append(networks, &network)
// 	}

// 	return networks, nil
// }

// func (d Dumper) Appliances(zone string) (appliances []*pb.Appliance, err error) {
// 	var name string

// 	defer func() {
// 		if err != nil {
// 			err = fmt.Errorf("appliance %q: %w", name, err)
// 		}
// 	}()

// 	r := d.NewApplianceReader(zone, "")
// 	for a, err := range r.Responses() {
// 		name = a.GetName()

// 		if err != nil {
// 			return nil, err
// 		}

// 		attrs, err := d.ApplianceAttrs(zone, a.GetName())
// 		if err != nil {
// 			return nil, err
// 		}

// 		appliance := pb.Appliance{
// 			Name:       Ptr(a.GetName()),
// 			Attributes: attrs,
// 		}

// 		appliances = append(appliances, &appliance)
// 	}

// 	return appliances, nil
// }

// func (d Dumper) ApplianceAttrs(zone, appliance string) (map[string]string, error) {
// 	attrs := make(map[string]string)

// 	r := d.NewApplianceAttrReader(zone, appliance, "")
// 	for aa, err := range r.Responses() {
// 		if err != nil {
// 			return nil, err
// 		}

// 		attrs[aa.GetName()] = aa.GetValue()
// 	}

// 	return attrs, nil
// }

// func (d Dumper) Environments(zone string) (environments []*pb.Environment, err error) {
// 	var name string

// 	defer func() {
// 		if err != nil {
// 			err = fmt.Errorf("environment %q: %w", name, err)
// 		}
// 	}()

// 	r := d.NewEnvironmentReader(zone, "")
// 	for e, err := range r.Responses() {
// 		name = e.GetName()

// 		if err != nil {
// 			return nil, err
// 		}

// 		attrs, err := d.EnvironmentAttrs(zone, e.GetName())
// 		if err != nil {
// 			return nil, err
// 		}

// 		environment := pb.Environment{
// 			Name:       Ptr(e.GetName()),
// 			Attributes: attrs,
// 		}

// 		environments = append(environments, &environment)
// 	}

// 	return environments, nil
// }

// func (d Dumper) EnvironmentAttrs(zone, environment string) (map[string]string, error) {
// 	attrs := make(map[string]string)

// 	r := d.NewEnvironmentAttrReader(zone, environment, "")
// 	for ea, err := range r.Responses() {
// 		if err != nil {
// 			return nil, err
// 		}

// 		attrs[ea.GetName()] = ea.GetValue()
// 	}

// 	return attrs, nil
// }

// func (d Dumper) Hosts(zone, cluster string) (hosts []*pb.Host, err error) {
// 	var name string

// 	defer func() {
// 		if err != nil {
// 			err = fmt.Errorf("host %q: %w", name, err)
// 		}
// 	}()

// 	r := d.NewHostReader(zone, cluster, "")
// 	for h, err := range r.Responses() {
// 		name = h.GetName()

// 		if err != nil {
// 			return nil, err
// 		}

// 		if cluster != "" && h.HasCluster() { // just want stand-alone hosts
// 			continue
// 		}

// 		attrs, err := d.HostAttrs(zone, cluster, h.GetName())
// 		if err != nil {
// 			return nil, err
// 		}

// 		ifaces, err := d.HostInterfaces(zone, cluster, h.GetName(), "")
// 		if err != nil {
// 			return nil, err
// 		}

// 		host := pb.Host{
// 			Name:       Ptr(h.GetName()),
// 			Attributes: attrs,
// 			Interfaces: ifaces,
// 		}

// 		if h.HasMake() {
// 			host.Make = Ptr(h.GetMake())
// 		}
// 		if h.HasModel() {
// 			host.Model = Ptr(h.GetModel())
// 		}
// 		if h.HasEnvironment() {
// 			host.Environment = Ptr(h.GetEnvironment())
// 		}
// 		if h.HasAppliance() {
// 			host.Appliance = Ptr(h.GetAppliance())
// 		}
// 		if h.HasLocation() {
// 			host.Location = Ptr(h.GetLocation())
// 		}
// 		if h.HasRack() {
// 			host.Rack = Ptr(h.GetRack())
// 		}
// 		if h.HasRank() {
// 			host.Rank = Ptr(h.GetRank())
// 		}
// 		if h.HasSlot() {
// 			host.Slot = Ptr(h.GetSlot())
// 		}
// 		if h.HasType() {
// 			host.Type = Ptr(h.GetType())
// 		}

// 		hosts = append(hosts, &host)
// 	}

// 	return hosts, nil
// }

// func (d Dumper) HostAttrs(zone, cluster, host string) (map[string]string, error) {
// 	attrs := make(map[string]string)

// 	r := d.NewHostAttrReader(zone, cluster, host, "")
// 	for ha, err := range r.Responses() {
// 		if err != nil {
// 			return nil, err
// 		}

// 		attrs[ha.GetName()] = ha.GetValue()
// 	}

// 	return attrs, nil
// }

// func (d Dumper) HostInterfaces(zone, cluster, host, iface string) (ifaces []*pb.Host_Interface, err error) {
// 	var name string

// 	defer func() {
// 		if err != nil {
// 			err = fmt.Errorf("interface %q: %w", name, err)
// 		}
// 	}()

// 	r := d.NewHostInterfaceReader(zone, cluster, host, iface)
// 	for i, err := range r.Responses() {
// 		name = i.GetName()

// 		if err != nil {
// 			return nil, err
// 		}

// 		iface := pb.Host_Interface{Name: Ptr(i.GetName())}

// 		if i.HasIp() {
// 			iface.Ip = Ptr(i.GetIp())
// 		}
// 		if i.HasMac() {
// 			iface.Mac = Ptr(i.GetMac())
// 		}
// 		if i.HasDhcp() {
// 			iface.Dhcp = Ptr(i.GetDhcp())
// 		}
// 		if i.HasPxe() {
// 			iface.Pxe = Ptr(i.GetPxe())
// 		}
// 		if i.HasManagement() {
// 			iface.Management = Ptr(i.GetManagement())
// 		}
// 		if i.HasType() {
// 			iface.Type = Ptr(i.GetType())
// 		}
// 		if i.HasBondMode() {
// 			iface.BondMode = Ptr(i.GetBondMode())
// 		}

// 		ifaces = append(ifaces, &iface)
// 	}

// 	return ifaces, nil
// }

// func (d Dumper) Clusters(zone string) (clusters []*pb.Cluster, err error) {
// 	var name string

// 	defer func() {
// 		if err != nil {
// 			err = fmt.Errorf("cluster %q: %w", name, err)
// 		}
// 	}()

// 	r := d.NewClusterReader(zone, "")
// 	for c, err := range r.Responses() {
// 		name = c.GetName()

// 		if err != nil {
// 			return nil, err
// 		}

// 		attrs, err := d.ClusterAttrs(zone, c.GetName())
// 		if err != nil {
// 			return nil, err
// 		}

// 		hosts, err := d.Hosts(zone, c.GetName())
// 		if err != nil {
// 			return nil, err
// 		}

// 		cluster := pb.Cluster{
// 			Attributes: attrs,
// 			Hosts:      hosts,
// 			Name:       Ptr(c.GetName()),
// 		}

// 		clusters = append(clusters, &cluster)
// 	}

// 	return clusters, nil
// }

// func (d Dumper) ClusterAttrs(zone, cluster string) (map[string]string, error) {
// 	attrs := make(map[string]string)

// 	r := d.NewClusterAttrReader(zone, cluster, "")
// 	for ca, err := range r.Responses() {
// 		if err != nil {
// 			return nil, err
// 		}

// 		attrs[ca.GetName()] = ca.GetValue()
// 	}

// 	return attrs, nil
// }
