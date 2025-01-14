package cli

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"log/slog"

// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/status"

// 	pb "endobit.io/metal/gen/go/proto/metal/v1"
// )

// type Loader struct {
// 	*RPC
// }

// func (l Loader) Load(ctx context.Context, doc *pb.Schema) (err error) {
// 	defer func() {
// 		if err != nil {
// 			err = fmt.Errorf("load failed: %w", err)
// 		}
// 	}()

// 	err = l.GlobalAttrs(ctx, doc.GetAttributes())
// 	if err != nil {
// 		return err
// 	}

// 	err = l.Makes(ctx, doc.GetMakes())
// 	if err != nil {
// 		return err
// 	}

// 	return l.Zones(ctx, doc.GetZones())
// }

// func (l Loader) GlobalAttrs(ctx context.Context, attrs map[string]string) error {
// 	for key, value := range attrs {
// 		if key == "" {
// 			return errors.New("missing attribute name")
// 		}

// 		if err := l.GlobalAttr(ctx, key, value); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (l Loader) GlobalAttr(ctx context.Context, key, value string) (err error) {
// 	var fields pb.UpdateGlobalAttrRequest_Fields

// 	defer func() {
// 		if err != nil {
// 			err = fmt.Errorf("attribute %q: %w", key, err)
// 		}
// 	}()

// 	create := pb.CreateGlobalAttrRequest_builder{
// 		Name: Ptr(key),
// 	}.Build()

// 	if _, err = l.RPC.Stack.CreateGlobalAttr(ctx, create); err != nil {
// 		if status.Code(err) != codes.AlreadyExists {
// 			return err
// 		}

// 		l.debug("created attr", slog.String("name", key))
// 	}

// 	if value != "" {
// 		fields.SetValue(value)
// 	}

// 	update := pb.UpdateGlobalAttrRequest_builder{
// 		Name:   Ptr(key),
// 		Fields: &fields,
// 	}.Build()

// 	if _, err = l.RPC.Stack.UpdateGlobalAttr(ctx, update); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (l Loader) Makes(ctx context.Context, makes []*pb.Make) error {
// 	for _, m := range makes {
// 		if m.Name == nil {
// 			return errors.New("missing make name")
// 		}

// 		if err := l.Make(ctx, m); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (l Loader) Make(ctx context.Context, mk *pb.Make) error {
// 	_, err := l.RPC.Stack.CreateMake(ctx, pb.CreateMakeRequest_builder{
// 		Name: mk.Name,
// 	}.Build())
// 	if err != nil && status.Code(err) != codes.AlreadyExists {
// 		return fmt.Errorf("make %q: %w", mk.GetName(), err)
// 	}

// 	return nil
// }

// func (l Loader) Models(ctx context.Context, models []*pb.Model) error {
// 	for _, m := range models {
// 		if m.Name == nil {
// 			return errors.New("missing model name")
// 		}

// 		if err := l.Model(ctx, m); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (l Loader) Model(ctx context.Context, model *pb.Model) (err error) {
// 	var fields pb.UpdateModelRequest_Fields

// 	defer func() {
// 		if err != nil {
// 			err = fmt.Errorf("model %q: %w", model.GetName(), err)
// 		}
// 	}()

// 	if model.Make == nil {
// 		return errors.New("missing make")
// 	}

// 	_, err = l.RPC.Stack.CreateModel(ctx, pb.CreateModelRequest_builder{
// 		Make: Ptr(model.GetMake()),
// 		Name: model.Name,
// 	}.Build())
// 	if err != nil && status.Code(err) != codes.AlreadyExists {
// 		return err
// 	}

// 	if model.Architecture != nil {
// 		fields.SetArchitecture(model.GetArchitecture())
// 	}

// 	update := pb.UpdateModelRequest_builder{
// 		Make:   Ptr(model.GetMake()),
// 		Name:   model.Name,
// 		Fields: &fields,
// 	}.Build()

// 	if _, err = l.RPC.Stack.UpdateModel(ctx, update); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (l Loader) Zones(ctx context.Context, zones []*pb.Zone) error {
// 	for _, z := range zones {
// 		if z.Name == nil {
// 			return errors.New("missing zone name")
// 		}

// 		if err := l.Zone(ctx, z); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (l Loader) Zone(ctx context.Context, zone *pb.Zone) (err error) {
// 	var fields pb.UpdateZoneRequest_Fields

// 	defer func() {
// 		if err != nil {
// 			err = fmt.Errorf("zone %q: %w", zone.GetName(), err)
// 		}
// 	}()

// 	_, err = l.RPC.Stack.CreateZone(ctx, pb.CreateZoneRequest_builder{
// 		Name: zone.Name,
// 	}.Build())
// 	if err != nil && status.Code(err) != codes.AlreadyExists {
// 		return err
// 	}

// 	if zone.TimeZone != nil {
// 		fields.SetTimeZone(zone.GetTimeZone())
// 	}

// 	update := pb.UpdateZoneRequest_builder{
// 		Name:   zone.Name,
// 		Fields: &fields,
// 	}.Build()

// 	_, err = l.RPC.Stack.UpdateZone(ctx, update)
// 	if err != nil {
// 		return err
// 	}

// 	err = l.ZoneAttrs(ctx, zone.GetName(), zone.GetAttributes())
// 	if err != nil {
// 		return err
// 	}

// 	err = l.Networks(ctx, zone.GetName(), zone.GetNetworks())
// 	if err != nil {
// 		return err
// 	}

// 	err = l.Appliances(ctx, zone.GetName(), zone.GetAppliances())
// 	if err != nil {
// 		return err
// 	}

// 	err = l.Environments(ctx, zone.GetName(), zone.GetEnvironments())
// 	if err != nil {
// 		return err
// 	}

// 	err = l.Hosts(ctx, zone.GetName(), "", zone.GetHosts())
// 	if err != nil {
// 		return err
// 	}

// 	err = l.Clusters(ctx, zone.GetName(), zone.GetClusters())
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (l Loader) ZoneAttrs(ctx context.Context, zone string, attrs map[string]string) error {
// 	for key, value := range attrs {
// 		if key == "" {
// 			return errors.New("missing attribute name")
// 		}

// 		if err := l.ZoneAttr(ctx, zone, key, value); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (l Loader) ZoneAttr(ctx context.Context, zone, key, value string) (err error) {
// 	var fields pb.UpdateZoneAttrRequest_Fields

// 	defer func() {
// 		if err != nil {
// 			err = fmt.Errorf("attribute %q: %w", key, err)
// 		}
// 	}()

// 	create := pb.CreateZoneAttrRequest_builder{
// 		Zone: Ptr(zone),
// 		Name: Ptr(key),
// 	}.Build()

// 	if _, err = l.RPC.Stack.CreateZoneAttr(ctx, create); err != nil {
// 		if status.Code(err) != codes.AlreadyExists {
// 			return err
// 		}

// 		l.debug("created attr", slog.String("zone", zone), slog.String("name", key))
// 	}

// 	if value != "" {
// 		fields.SetValue(value)
// 	}

// 	update := pb.UpdateZoneAttrRequest_builder{
// 		Zone:   Ptr(zone),
// 		Name:   Ptr(key),
// 		Fields: &fields,
// 	}.Build()

// 	if _, err = l.RPC.Stack.UpdateZoneAttr(ctx, update); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (l Loader) Networks(ctx context.Context, zone string, networks []*pb.Network) error {
// 	for _, n := range networks {
// 		if n.Name == nil {
// 			return errors.New("missing network name")
// 		}

// 		if err := l.Network(ctx, zone, n); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (l Loader) Network(ctx context.Context, zone string, network *pb.Network) (err error) {
// 	var fields pb.UpdateNetworkRequest_Fields

// 	defer func() {
// 		if err != nil {
// 			err = fmt.Errorf("network %q: %w", network.GetName(), err)
// 		}
// 	}()

// 	create := pb.CreateNetworkRequest_builder{
// 		Zone: Ptr(zone),
// 		Name: network.Name,
// 	}.Build()

// 	if _, err = l.RPC.Stack.CreateNetwork(ctx, create); err != nil {
// 		if status.Code(err) != codes.AlreadyExists {
// 			return err
// 		}

// 		l.debug("created network", slog.String("zone", zone), slog.String("name", network.GetName()))
// 	}

// 	if network.Address != nil {
// 		fields.SetAddress(network.GetAddress())
// 	}
// 	if network.Gateway != nil {
// 		fields.SetGateway(network.GetGateway())
// 	}
// 	if network.Pxe != nil {
// 		fields.SetPxe(network.GetPxe())
// 	}
// 	if network.Mtu != nil {
// 		fields.SetMtu(network.GetMtu())
// 	}

// 	update := pb.UpdateNetworkRequest_builder{
// 		Zone:   Ptr(zone),
// 		Name:   network.Name,
// 		Fields: &fields,
// 	}.Build()

// 	if _, err = l.RPC.Stack.UpdateNetwork(ctx, update); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (l Loader) Appliances(ctx context.Context, zone string, appliances []*pb.Appliance) error {
// 	for _, a := range appliances {
// 		if a.Name == nil {
// 			return errors.New("missing appliance name")
// 		}

// 		if err := l.Appliance(ctx, zone, a); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (l Loader) Appliance(ctx context.Context, zone string, appliance *pb.Appliance) (err error) {
// 	defer func() {
// 		if err != nil {
// 			err = fmt.Errorf("appliance %q: %w", appliance.GetName(), err)
// 		}
// 	}()

// 	create := pb.CreateApplianceRequest_builder{
// 		Zone: Ptr(zone),
// 		Name: appliance.Name,
// 	}.Build()

// 	if _, err = l.RPC.Stack.CreateAppliance(ctx, create); err != nil {
// 		if status.Code(err) != codes.AlreadyExists {
// 			return err
// 		}

// 		l.debug("created appliance", slog.String("zone", zone), slog.String("name", appliance.GetName()))
// 	}

// 	err = l.ApplianceAttrs(ctx, zone, appliance.GetName(), appliance.GetAttributes())
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (l Loader) ApplianceAttrs(ctx context.Context, zone, appliance string, attrs map[string]string) error {
// 	for key, value := range attrs {
// 		if key == "" {
// 			return errors.New("missing attribute name")
// 		}

// 		if err := l.ApplianceAttr(ctx, zone, appliance, key, value); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (l Loader) ApplianceAttr(ctx context.Context, zone, appliance, key, value string) (err error) {
// 	var fields pb.UpdateApplianceAttrRequest_Fields

// 	defer func() {
// 		if err != nil {
// 			err = fmt.Errorf("attribute %q: %w", key, err)
// 		}
// 	}()

// 	create := pb.CreateApplianceAttrRequest_builder{
// 		Zone:      Ptr(zone),
// 		Appliance: Ptr(appliance),
// 		Name:      Ptr(key),
// 	}.Build()

// 	if _, err = l.RPC.Stack.CreateApplianceAttr(ctx, create); err != nil {
// 		if status.Code(err) != codes.AlreadyExists {
// 			return err
// 		}

// 		l.debug("created attr",
// 			slog.String("zone", zone), slog.String("appliance", appliance), slog.String("name", key))
// 	}

// 	if value != "" {
// 		fields.SetValue(value)
// 	}

// 	update := pb.UpdateApplianceAttrRequest_builder{
// 		Zone:      Ptr(zone),
// 		Appliance: Ptr(appliance),
// 		Name:      Ptr(key),
// 		Fields:    &fields,
// 	}.Build()

// 	if _, err = l.RPC.Stack.UpdateApplianceAttr(ctx, update); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (l Loader) Environments(ctx context.Context, zone string, environments []*pb.Environment) error {
// 	for _, e := range environments {
// 		if e.Name == nil {
// 			return errors.New("missing environment name")
// 		}

// 		if err := l.Environment(ctx, zone, e); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (l Loader) Environment(ctx context.Context, zone string, environment *pb.Environment) (err error) {
// 	defer func() {
// 		if err != nil {
// 			err = fmt.Errorf("environment %q: %w", environment.GetName(), err)
// 		}
// 	}()

// 	create := pb.CreateEnvironmentRequest_builder{
// 		Zone: Ptr(zone),
// 		Name: environment.Name,
// 	}.Build()

// 	if _, err = l.RPC.Stack.CreateEnvironment(ctx, create); err != nil {
// 		if status.Code(err) != codes.AlreadyExists {
// 			return err
// 		}

// 		l.debug("created environment", slog.String("zone", zone), slog.String("name", environment.GetName()))
// 	}

// 	err = l.EnvironmentAttrs(ctx, zone, environment.GetName(), environment.GetAttributes())
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (l Loader) EnvironmentAttrs(ctx context.Context, zone, environment string, attrs map[string]string) error {
// 	for key, value := range attrs {
// 		if key == "" {
// 			return errors.New("missing attribute name")
// 		}

// 		if err := l.EnvironmentAttr(ctx, zone, environment, key, value); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (l Loader) EnvironmentAttr(ctx context.Context, zone, environment, key, value string) (err error) {
// 	var fields pb.UpdateEnvironmentAttrRequest_Fields

// 	defer func() {
// 		if err != nil {
// 			err = fmt.Errorf("attribute %q: %w", key, err)
// 		}
// 	}()

// 	create := pb.CreateEnvironmentAttrRequest_builder{
// 		Zone:        Ptr(zone),
// 		Environment: Ptr(environment),
// 		Name:        Ptr(key),
// 	}.Build()

// 	if _, err = l.RPC.Stack.CreateEnvironmentAttr(ctx, create); err != nil {
// 		if status.Code(err) != codes.AlreadyExists {
// 			return err
// 		}

// 		l.debug("created attr",
// 			slog.String("zone", zone), slog.String("environment", environment), slog.String("name", key))
// 	}

// 	if value != "" {
// 		fields.SetValue(value)
// 	}

// 	update := pb.UpdateEnvironmentAttrRequest_builder{
// 		Zone:        Ptr(zone),
// 		Environment: Ptr(environment),
// 		Name:        Ptr(key),
// 		Fields:      &fields,
// 	}.Build()

// 	if _, err = l.RPC.Stack.UpdateEnvironmentAttr(ctx, update); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (l Loader) Hosts(ctx context.Context, zone, cluster string, host []*pb.Host) error {
// 	for _, h := range host {
// 		if h.Name == nil {
// 			return errors.New("missing host name")
// 		}

// 		if err := l.Host(ctx, zone, cluster, h); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (l Loader) Host(ctx context.Context, zone, cluster string, host *pb.Host) (err error) {
// 	var fields pb.UpdateHostRequest_Fields

// 	defer func() {
// 		if err != nil {
// 			err = fmt.Errorf("host %q: %w", host.GetName(), err)
// 		}
// 	}()

// 	create := pb.CreateHostRequest_builder{
// 		Zone: Ptr(zone),
// 		Name: host.Name,
// 	}.Build()

// 	if cluster != "" {
// 		create.SetCluster(cluster)
// 	}

// 	if _, err = l.RPC.Stack.CreateHost(ctx, create); err != nil {
// 		if status.Code(err) != codes.AlreadyExists {
// 			return err
// 		}

// 		l.debug("created host", slog.String("zone", zone), slog.String("cluster", cluster),
// 			slog.String("name", host.GetName()))
// 	}

// 	if host.Make != nil {
// 		fields.SetMake(host.GetMake())
// 	}
// 	if host.Model != nil {
// 		fields.SetModel(host.GetModel())
// 	}
// 	if host.Environment != nil {
// 		fields.SetEnvironment(host.GetEnvironment())
// 	}
// 	if host.Appliance != nil {
// 		fields.SetAppliance(host.GetAppliance())
// 	}
// 	if host.Location != nil {
// 		fields.SetLocation(host.GetLocation())
// 	}
// 	if host.Rack != nil {
// 		fields.SetRack(host.GetRack())
// 	}
// 	if host.Rank != nil {
// 		fields.SetRank(host.GetRank())
// 	}
// 	if host.Slot != nil {
// 		fields.SetSlot(host.GetSlot())
// 	}
// 	if host.Type != nil {
// 		fields.SetType(host.GetType())
// 	}

// 	update := pb.UpdateHostRequest_builder{
// 		Zone:   Ptr(zone),
// 		Name:   host.Name,
// 		Fields: &fields,
// 	}.Build()

// 	if cluster != "" {
// 		update.SetCluster(cluster)
// 	}

// 	if _, err = l.RPC.Stack.UpdateHost(ctx, update); err != nil {
// 		return err
// 	}

// 	err = l.HostAttrs(ctx, zone, cluster, host.GetName(), host.GetAttributes())
// 	if err != nil {
// 		return err
// 	}

// 	err = l.HostInterfaces(ctx, zone, cluster, host.GetName(), host.GetInterfaces())
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (l Loader) HostAttrs(ctx context.Context, zone, cluster, host string, attrs map[string]string) error {
// 	for key, value := range attrs {
// 		if key == "" {
// 			return errors.New("missing attribute name")
// 		}

// 		if err := l.HostAttr(ctx, zone, cluster, host, key, value); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (l Loader) HostAttr(ctx context.Context, zone, cluster, host, key, value string) (err error) {
// 	var fields pb.UpdateHostAttrRequest_Fields

// 	defer func() {
// 		if err != nil {
// 			err = fmt.Errorf("attribute %q: %w", key, err)
// 		}
// 	}()

// 	create := pb.CreateHostAttrRequest_builder{
// 		Zone: Ptr(zone),
// 		Host: Ptr(host),
// 		Name: Ptr(key),
// 	}.Build()

// 	if cluster != "" {
// 		create.SetCluster(cluster)
// 	}

// 	if _, err = l.RPC.Stack.CreateHostAttr(ctx, create); err != nil {
// 		if status.Code(err) != codes.AlreadyExists {
// 			return err
// 		}

// 		l.debug("created attr",
// 			slog.String("zone", zone), slog.String("cluster", cluster),
// 			slog.String("host", host), slog.String("name", key))
// 	}

// 	if value != "" {
// 		fields.SetValue(value)
// 	}

// 	update := pb.UpdateHostAttrRequest_builder{
// 		Zone:   Ptr(zone),
// 		Host:   Ptr(host),
// 		Name:   Ptr(key),
// 		Fields: &fields,
// 	}.Build()

// 	if cluster != "" {
// 		update.SetCluster(cluster)
// 	}

// 	if _, err = l.RPC.Stack.UpdateHostAttr(ctx, update); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (l Loader) HostInterfaces(ctx context.Context, zone, cluster, host string, ifaces []*pb.Host_Interface) error {
// 	for _, i := range ifaces {
// 		if i.Name == nil {
// 			return errors.New("missing interface name")
// 		}

// 		if err := l.HostInterface(ctx, zone, cluster, host, i); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (l Loader) HostInterface(ctx context.Context, zone, cluster, host string, iface *pb.Host_Interface) (err error) {
// 	var fields pb.UpdateHostInterfaceRequest_Fields

// 	defer func() {
// 		if err != nil {
// 			err = fmt.Errorf("interface %q: %w", iface.GetName(), err)
// 		}
// 	}()

// 	create := pb.CreateHostInterfaceRequest_builder{
// 		Zone: Ptr(zone),
// 		Host: Ptr(host),
// 		Name: iface.Name,
// 	}.Build()

// 	if cluster != "" {
// 		create.SetCluster(cluster)
// 	}

// 	if _, err = l.RPC.Stack.CreateHostInterface(ctx, create); err != nil {
// 		if status.Code(err) != codes.AlreadyExists {
// 			return err
// 		}

// 		l.debug("created interface", slog.String("zone", zone), slog.String("cluster", cluster),
// 			slog.String("host", host), slog.String("name", iface.GetName()))
// 	}

// 	if iface.Network != nil {
// 		fields.SetNetwork(iface.GetNetwork())
// 	}
// 	if iface.Ip != nil {
// 		fields.SetIp(iface.GetIp())
// 	}
// 	if iface.Mac != nil {
// 		fields.SetMac(iface.GetMac())
// 	}
// 	if iface.Dhcp != nil {
// 		fields.SetDhcp(iface.GetDhcp())
// 	}
// 	if iface.Pxe != nil {
// 		fields.SetPxe(iface.GetPxe())
// 	}
// 	if iface.Management != nil {
// 		fields.SetManagement(iface.GetManagement())
// 	}

// 	update := pb.UpdateHostInterfaceRequest_builder{
// 		Zone:   Ptr(zone),
// 		Host:   Ptr(host),
// 		Name:   iface.Name,
// 		Fields: &fields,
// 	}.Build()

// 	if cluster != "" {
// 		update.SetCluster(cluster)
// 	}

// 	if _, err = l.RPC.Stack.UpdateHostInterface(ctx, update); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (l Loader) Clusters(ctx context.Context, zone string, clusters []*pb.Cluster) error {
// 	for _, c := range clusters {
// 		if c.Name == nil {
// 			return errors.New("missing cluster name")
// 		}

// 		if err := l.Cluster(ctx, zone, c); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (l Loader) Cluster(ctx context.Context, zone string, cluster *pb.Cluster) (err error) {
// 	defer func() {
// 		if err != nil {
// 			err = fmt.Errorf("cluster %q: %w", cluster.GetName(), err)
// 		}
// 	}()

// 	create := pb.CreateClusterRequest_builder{
// 		Zone: Ptr(zone),
// 		Name: cluster.Name,
// 	}.Build()

// 	if _, err = l.RPC.Stack.CreateCluster(ctx, create); err != nil {
// 		if status.Code(err) != codes.AlreadyExists {
// 			return err
// 		}

// 		l.debug("created cluster", slog.String("zone", zone), slog.String("name", cluster.GetName()))
// 	}

// 	err = l.ClusterAttrs(ctx, zone, cluster.GetName(), cluster.GetAttributes())
// 	if err != nil {
// 		return err
// 	}

// 	err = l.Hosts(ctx, zone, cluster.GetName(), cluster.GetHosts())
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (l Loader) ClusterAttrs(ctx context.Context, zone, cluster string, attrs map[string]string) error {
// 	for key, value := range attrs {
// 		if key == "" {
// 			return errors.New("missing attribute name")
// 		}

// 		if err := l.ClusterAttr(ctx, zone, cluster, key, value); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (l Loader) ClusterAttr(ctx context.Context, zone, cluster, key, value string) (err error) {
// 	var fields pb.UpdateClusterAttrRequest_Fields

// 	defer func() {
// 		if err != nil {
// 			err = fmt.Errorf("attribute %q: %w", key, err)
// 		}
// 	}()

// 	create := pb.CreateClusterAttrRequest_builder{
// 		Zone:    Ptr(zone),
// 		Cluster: Ptr(cluster),
// 		Name:    Ptr(key),
// 	}.Build()

// 	if _, err = l.RPC.Stack.CreateClusterAttr(ctx, create); err != nil {
// 		if status.Code(err) != codes.AlreadyExists {
// 			return err
// 		}

// 		l.debug("created attr",
// 			slog.String("zone", zone), slog.String("cluster", cluster), slog.String("name", key))
// 	}

// 	if value != "" {
// 		fields.SetValue(value)
// 	}

// 	update := pb.UpdateClusterAttrRequest_builder{
// 		Zone:    Ptr(zone),
// 		Cluster: Ptr(cluster),
// 		Name:    Ptr(key),
// 		Fields:  &fields,
// 	}.Build()

// 	if _, err := l.RPC.Stack.UpdateClusterAttr(ctx, update); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (l Loader) debug(msg string, args ...any) { l.RPC.Logger.Debug(msg, args...) }
