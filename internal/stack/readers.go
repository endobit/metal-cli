package stack

import (
	pb "github.com/endobit/stack/internal/generated/go/proto/stack/v1"
	"github.com/endobit/stack/internal/stack/stream"
)

//
// Global
//

// newGlobalAttrReader creates a new stream reader for global attributes.
func newGlobalAttrReader(client *rpcClient, attr string) *stream.Reader[pb.ReadGlobalAttrsResponse] {
	var req pb.ReadGlobalAttrsRequest

	if attr != "" {
		req.SetGlob(attr)
	}

	return stream.NewReader(client.stack.ReadGlobalAttrs(client.Context(), &req))
}

//
// Makes and Models
//

// newMakeReader creates a new stream reader for makes.
func newMakeReader(client *rpcClient, mk string) *stream.Reader[pb.ReadMakesResponse] {
	var req pb.ReadMakesRequest

	if mk != "" {
		req.SetGlob(mk)
	}

	return stream.NewReader(client.stack.ReadMakes(client.Context(), &req))
}

// newModelReader creates a new stream reader for models.
func newModelReader(client *rpcClient, mk, model string) *stream.Reader[pb.ReadModelsResponse] {
	var req pb.ReadModelsRequest

	if mk != "" {
		req.SetMake(mk)
	}
	if model != "" {
		req.SetGlob(model)
	}

	return stream.NewReader(client.stack.ReadModels(client.Context(), &req))
}

//
// Zones
//

// newZoneReader creates a new stream reader for zones.
func newZoneReader(client *rpcClient, zone string) *stream.Reader[pb.ReadZonesResponse] {
	var req pb.ReadZonesRequest

	if zone != "" {
		req.SetGlob(zone)
	}

	return stream.NewReader(client.stack.ReadZones(client.Context(), &req))
}

// newZoneAttrReader creates a new stream reader for zone attributes.
func newZoneAttrReader(client *rpcClient, zone, attr string) *stream.Reader[pb.ReadZoneAttrsResponse] {
	var req pb.ReadZoneAttrsRequest

	if zone != "" {
		req.SetZone(zone)
	}
	if attr != "" {
		req.SetGlob(attr)
	}

	return stream.NewReader(client.stack.ReadZoneAttrs(client.Context(), &req))
}

//
// Appliances
//

// newApplianceReader creates a new stream reader for appliances.
func newApplianceReader(client *rpcClient, zone, appliance string) *stream.Reader[pb.ReadAppliancesResponse] {
	var req pb.ReadAppliancesRequest

	if zone != "" {
		req.SetZone(zone)
	}
	if appliance != "" {
		req.SetGlob(appliance)
	}

	return stream.NewReader(client.stack.ReadAppliances(client.Context(), &req))
}

// newApplianceAttrReader creates a new stream reader for appliance attributes.
func newApplianceAttrReader(client *rpcClient, zone, appliance, attr string) *stream.Reader[pb.ReadApplianceAttrsResponse] {
	var req pb.ReadApplianceAttrsRequest

	if zone != "" {
		req.SetZone(zone)
	}
	if appliance != "" {
		req.SetAppliance(appliance)
	}
	if attr != "" {
		req.SetGlob(attr)
	}

	return stream.NewReader(client.stack.ReadApplianceAttrs(client.Context(), &req))
}

//
// Environments
//

// newEnvironmentReader creates a new stream reader for environments.
func newEnvironmentReader(client *rpcClient, zone, environment string) *stream.Reader[pb.ReadEnvironmentsResponse] {
	var req pb.ReadEnvironmentsRequest

	if zone != "" {
		req.SetZone(zone)
	}
	if environment != "" {
		req.SetGlob(environment)
	}

	return stream.NewReader(client.stack.ReadEnvironments(client.Context(), &req))
}

// newEnvironmentAttrReader creates a new stream reader for environment attributes.
func newEnvironmentAttrReader(client *rpcClient, zone, environment, attr string) *stream.Reader[pb.ReadEnvironmentAttrsResponse] {
	var req pb.ReadEnvironmentAttrsRequest

	if zone != "" {
		req.SetZone(zone)
	}
	if environment != "" {
		req.SetEnvironment(environment)
	}
	if attr != "" {
		req.SetGlob(attr)
	}

	return stream.NewReader(client.stack.ReadEnvironmentAttrs(client.Context(), &req))
}

//
// Networks
//

// newNetworkReader creates a new stream reader for networks.
func newNetworkReader(client *rpcClient, zone, network string) *stream.Reader[pb.ReadNetworksResponse] {
	var req pb.ReadNetworksRequest

	if zone != "" {
		req.SetZone(zone)
	}
	if network != "" {
		req.SetGlob(network)
	}

	return stream.NewReader(client.stack.ReadNetworks(client.Context(), &req))
}

//
// Racks
//

// newRackReader creates a new stream reader for racks.
func newRackReader(client *rpcClient, zone, rack string) *stream.Reader[pb.ReadRacksResponse] {
	var req pb.ReadRacksRequest

	if zone != "" {
		req.SetZone(zone)
	}
	if rack != "" {
		req.SetGlob(rack)
	}

	return stream.NewReader(client.stack.ReadRacks(client.Context(), &req))
}

// newRackAttrReader creates a new stream reader for rack attributes.
func newRackAttrReader(client *rpcClient, zone, rack, attr string) *stream.Reader[pb.ReadRackAttrsResponse] {
	var req pb.ReadRackAttrsRequest

	if zone != "" {
		req.SetZone(zone)
	}
	if rack != "" {
		req.SetRack(rack)
	}
	if attr != "" {
		req.SetGlob(attr)
	}

	return stream.NewReader(client.stack.ReadRackAttrs(client.Context(), &req))
}

//
// Clusters
//

// newClusterReader creates a new stream reader for clusters.
func newClusterReader(client *rpcClient, zone, cluster string) *stream.Reader[pb.ReadClustersResponse] {
	var req pb.ReadClustersRequest

	if zone != "" {
		req.SetZone(zone)
	}
	if cluster != "" {
		req.SetGlob(cluster)
	}

	return stream.NewReader(client.stack.ReadClusters(client.Context(), &req))
}

// newClusterAttrReader creates a new stream reader for cluster attributes.
func newClusterAttrReader(client *rpcClient, zone, cluster, attr string) *stream.Reader[pb.ReadClusterAttrsResponse] {
	var req pb.ReadClusterAttrsRequest

	if zone != "" {
		req.SetZone(zone)
	}
	if cluster != "" {
		req.SetCluster(cluster)
	}
	if attr != "" {
		req.SetGlob(attr)
	}

	return stream.NewReader(client.stack.ReadClusterAttrs(client.Context(), &req))
}

//
// Hosts
//

// newHostReader creates a new stream reader for hosts.
func newHostReader(client *rpcClient, zone, cluster, host string) *stream.Reader[pb.ReadHostsResponse] {
	var req pb.ReadHostsRequest

	if zone != "" {
		req.SetZone(zone)
	}
	if cluster != "" {
		req.SetCluster(cluster)
	}
	if host != "" {
		req.SetGlob(host)
	}

	return stream.NewReader(client.stack.ReadHosts(client.Context(), &req))
}

// newHostAttrReader creates a new stream reader for host attributes.
func newHostAttrReader(client *rpcClient, zone, cluster, host, attr string) *stream.Reader[pb.ReadHostAttrsResponse] {
	var req pb.ReadHostAttrsRequest

	if zone != "" {
		req.SetZone(zone)
	}
	if cluster != "" {
		req.SetCluster(cluster)
	}
	if host != "" {
		req.SetHost(host)
	}
	if attr != "" {
		req.SetGlob(attr)
	}

	return stream.NewReader(client.stack.ReadHostAttrs(client.Context(), &req))
}

// newHostInterfaceReader creates a new stream reader for host interfaces.
func newHostInterfaceReader(client *rpcClient, zone, cluster, host, iface string) *stream.Reader[pb.ReadHostInterfacesResponse] {
	var req pb.ReadHostInterfacesRequest

	if zone != "" {
		req.SetZone(zone)
	}
	if cluster != "" {
		req.SetCluster(cluster)
	}
	if host != "" {
		req.SetHost(host)
	}
	if iface != "" {
		req.SetGlob(iface)
	}

	return stream.NewReader(client.stack.ReadHostInterfaces(client.Context(), &req))
}
