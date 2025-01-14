package stack

import (
	"context"
	"fmt"
	"math"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/endobit/stack/internal/generated/go/db"
	pb "github.com/endobit/stack/internal/generated/go/proto/stack/v1"
)

//// Host

//
// CREATE
//

func (s Service) CreateHost(ctx context.Context, in *pb.CreateHostRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasName():
		return nil, errMissingHost
	}

	p := db.CreateHostParams{
		Zone: in.GetZone(),
		Name: in.GetName(),
	}

	if in.HasCluster() {
		p.Cluster = in.GetCluster()
		p.Zone = "" // get zone from cluster
	}

	if err := s.db.CreateHost(ctx, p); err != nil {
		return nil, dberr(err)
	}

	return new(emptypb.Empty), nil
}

//
// READ
//

func (s Service) ReadHosts(in *pb.ReadHostsRequest, out grpc.ServerStreamingServer[pb.ReadHostsResponse]) error {
	var (
		rows []struct {
			ID          int64
			Name        string
			Make        *string
			Model       *string
			Environment *string
			Appliance   *string
			Location    *string
			Rack        *string
			Rank        *int64
			Slot        *int64
			Zone        *string
			Cluster     *string
		}
		err error
	)

	ctx := context.Background()

	switch {
	case in.HasZone() && in.HasName():
		var res db.ReadHostRow

		res, err = s.db.ReadHost(ctx, db.ReadHostParams{
			Zone:    in.GetZone(),
			Cluster: in.GetCluster(), // "" means standalone host (same for queries below)
			Name:    in.GetName(),
		})
		rows = append(rows, res)

	case in.HasZone() && in.HasGlob():
		var res []db.ReadHostsByGlobRow

		res, err = s.db.ReadHostsByGlob(ctx, db.ReadHostsByGlobParams{
			Zone:    in.GetZone(),
			Cluster: in.GetCluster(),
			Glob:    in.GetGlob(),
		})
		for i := range res {
			rows = append(rows, res[i])
		}

	case in.HasCluster() && in.HasZone():
		var res []db.ReadHostsByClusterRow

		res, err = s.db.ReadHostsByCluster(ctx, db.ReadHostsByClusterParams{
			Zone:    in.GetZone(),
			Cluster: in.GetCluster(),
		})
		for i := range res {
			rows = append(rows, res[i])
		}

	case in.HasZone():
		var res []db.ReadHostsByZoneRow

		res, err = s.db.ReadHostsByZone(ctx, db.ReadHostsByZoneParams{
			Zone: in.GetZone(),
		})
		for i := range res {
			rows = append(rows, res[i])
		}

	default:
		var res []db.ReadHostsRow

		res, err = s.db.ReadHosts(ctx, db.ReadHostsParams{})
		for i := range res {
			rows = append(rows, res[i])
		}
	}

	if err != nil {
		return dberr(err)
	}

	for i := range rows {
		resp := pb.ReadHostsResponse_builder{
			Zone:        rows[i].Zone,
			Cluster:     rows[i].Cluster,
			Name:        &rows[i].Name,
			Make:        rows[i].Make,
			Model:       rows[i].Model,
			Environment: rows[i].Environment,
			Appliance:   rows[i].Appliance,
			Rack:        rows[i].Rack,
		}.Build()

		if rows[i].Location != nil {
			resp.SetLocation(*rows[i].Location)
		}

		if rows[i].Rack != nil {
			resp.SetRack(*rows[i].Rack)
		}

		if rows[i].Rank != nil {
			n := *rows[i].Rank
			if 0 < n || n > math.MaxUint32 {
				return fmt.Errorf("rank %v is not valid", n)
			}

			resp.SetRank(uint32(n))
		}

		if rows[i].Slot != nil {
			n := *rows[i].Slot
			if 0 < n || n > math.MaxUint32 {
				return fmt.Errorf("slot %v is not valid", n)
			}

			resp.SetSlot(uint32(n))
		}

		if err := out.Send(resp); err != nil {
			return err
		}
	}

	return nil
}

//
// UPDATE
//

func (s Service) UpdateHost(ctx context.Context, in *pb.UpdateHostRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasName():
		return nil, errMissingHost
	}

	fields := in.GetFields()

	fmt.Printf("update host %#v\n", fields)

	if fields.HasMake() {
		if !fields.HasModel() {
			return nil, errMissingModel
		}
	}

	if fields.HasName() {
		err := s.db.UpdateHostName(ctx, db.UpdateHostNameParams{
			Zone:    in.GetZone(),
			Cluster: in.GetCluster(),
			Host:    in.GetName(),
			Name:    fields.GetName(),
		})
		if err != nil {
			return nil, dberr(err)
		}
	}

	switch {
	case fields.HasMake() && fields.HasModel():
		_ = s.createModel(ctx, db.CreateModelParams{ // try to create if doesn't exists
			Make: fields.GetMake(),
			Name: fields.GetModel(),
		})

		err := s.db.UpdateHostModel(ctx, db.UpdateHostModelParams{
			Zone:    in.GetZone(),
			Cluster: in.GetCluster(),
			Host:    in.GetName(),
			Make:    fields.GetMake(),
			Model:   fields.GetModel(),
		})
		if err != nil {
			return nil, dberr(err)
		}
	case !fields.HasMake():
		return nil, errMissingMake
	case !fields.HasModel():
		return nil, errMissingModel
	}

	if fields.HasEnvironment() {
		_ = s.createEnvironment(ctx, db.CreateEnvironmentParams{ // try to create if doesn't exists
			Zone: in.GetZone(),
			Name: fields.GetEnvironment(),
		})

		err := s.db.UpdateHostEnvironment(ctx, db.UpdateHostEnvironmentParams{
			Zone:        in.GetZone(),
			Cluster:     in.GetCluster(),
			Host:        in.GetName(),
			Environment: fields.GetEnvironment(),
		})
		if err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasAppliance() {
		_ = s.createAppliance(ctx, db.CreateApplianceParams{ // try to create if doesn't exists
			Zone: in.GetZone(),
			Name: fields.GetAppliance(),
		})

		err := s.db.UpdateHostAppliance(ctx, db.UpdateHostApplianceParams{
			Zone:      in.GetZone(),
			Cluster:   in.GetCluster(),
			Host:      in.GetName(),
			Appliance: fields.GetAppliance(),
		})
		if err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasLocation() {
		err := s.db.UpdateHostLocation(ctx, db.UpdateHostLocationParams{
			Zone:     in.GetZone(),
			Cluster:  in.GetCluster(),
			Host:     in.GetName(),
			Location: proto.String(fields.GetLocation()),
		})
		if err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasLocation() {
		err := s.db.UpdateHostLocation(ctx, db.UpdateHostLocationParams{
			Zone:     in.GetZone(),
			Cluster:  in.GetCluster(),
			Host:     in.GetName(),
			Location: proto.String(fields.GetLocation()),
		})
		if err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasRack() {
		_ = s.createRack(ctx, db.CreateRackParams{ // try to create if doesn't exists
			Zone: in.GetZone(),
			Name: fields.GetRack(),
		})

		err := s.db.UpdateHostRack(ctx, db.UpdateHostRackParams{
			Zone:    in.GetZone(),
			Cluster: in.GetCluster(),
			Host:    in.GetName(),
			Rack:    fields.GetRack(),
		})
		if err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasRank() {
		err := s.db.UpdateHostRank(ctx, db.UpdateHostRankParams{
			Zone:    in.GetZone(),
			Cluster: in.GetCluster(),
			Host:    in.GetName(),
			Rank:    proto.Int64(int64(fields.GetRank())),
		})
		if err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasSlot() {
		err := s.db.UpdateHostSlot(ctx, db.UpdateHostSlotParams{
			Zone:    in.GetZone(),
			Cluster: in.GetCluster(),
			Host:    in.GetName(),
			Slot:    proto.Int64(int64(fields.GetSlot())),
		})
		if err != nil {
			return nil, dberr(err)
		}
	}

	return new(emptypb.Empty), nil
}

//
// DELETE
//

func (s Service) DeleteHosts(ctx context.Context, in *pb.DeleteHostsRequest) (*emptypb.Empty, error) {
	if !in.HasZone() {
		return nil, errMissingZone
	}

	var ids []int64

	switch {
	case in.HasName():
		row, err := s.db.ReadHost(ctx, db.ReadHostParams{
			Zone:    in.GetZone(),
			Cluster: in.GetCluster(),
			Name:    in.GetName(),
		})
		if err != nil {
			return nil, dberr(err)
		}
		ids = append(ids, row.ID)
	case in.HasGlob():
		rows, err := s.db.ReadHostsByGlob(ctx, db.ReadHostsByGlobParams{
			Zone:    in.GetZone(),
			Cluster: in.GetCluster(),
			Glob:    in.GetGlob(),
		})
		if err != nil {
			return nil, dberr(err)
		}
		for i := range rows {
			ids = append(ids, rows[i].ID)
		}
	default:
		return nil, errNameOrGlob
	}

	for _, id := range ids {
		if err := s.db.DeleteHost(ctx, db.DeleteHostParams{ID: id}); err != nil {
			return nil, dberr(err)
		}
	}

	return new(emptypb.Empty), nil
}

//// Host Attributes

//
// CREATE
//

func (s *Service) CreateHostAttr(ctx context.Context, in *pb.CreateHostAttrRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasHost():
		return nil, errMissingHost
	case !in.HasName():
		return nil, errMissingAttribute
	}

	p := db.CreateHostAttributeParams{
		Zone: in.GetZone(),
		Host: in.GetHost(),
		Name: in.GetName(),
	}

	if err := s.db.CreateHostAttribute(ctx, p); err != nil {
		return nil, dberr(err)
	}

	return new(emptypb.Empty), nil
}

//
// READ
//

func (s *Service) ReadHostAttrs(in *pb.ReadHostAttrsRequest, out grpc.ServerStreamingServer[pb.ReadHostAttrsResponse]) error {
	var (
		rows []struct {
			ID          int64
			Name        string
			Value       *string
			IsProtected int64
			Host        string
			Zone        *string
			Cluster     *string
		}
		err error
	)

	ctx := context.Background()

	switch {
	case in.HasZone() && in.HasHost() && in.HasName():
		var res db.ReadHostAttributeRow

		res, err = s.db.ReadHostAttribute(ctx, db.ReadHostAttributeParams{
			Zone:    in.GetZone(),
			Cluster: in.GetCluster(), // "" means standalone (same for queries below)
			Host:    in.GetHost(),
			Attr:    in.GetName(),
		})
		rows = append(rows, res)

	case in.HasZone() && in.HasHost() && in.HasGlob():
		var res []db.ReadHostAttributesByGlobRow

		res, err = s.db.ReadHostAttributesByGlob(ctx, db.ReadHostAttributesByGlobParams{
			Zone:    in.GetZone(),
			Cluster: in.GetCluster(),
			Host:    in.GetHost(),
			Glob:    in.GetGlob(),
		})
		for i := range res {
			rows = append(rows, res[i])
		}

	case in.HasZone() && in.HasHost():
		var res []db.ReadHostAttributesByHostRow

		res, err = s.db.ReadHostAttributesByHost(ctx, db.ReadHostAttributesByHostParams{
			Zone:    in.GetZone(),
			Cluster: in.GetCluster(),
			Host:    in.GetHost(),
		})
		for i := range res {
			rows = append(rows, res[i])
		}

	case in.HasZone() && in.HasGlob():
		var res []db.ReadHostAttributesByClusterRow

		res, err = s.db.ReadHostAttributesByCluster(ctx, db.ReadHostAttributesByClusterParams{
			Zone:    in.GetZone(),
			Cluster: in.GetCluster(),
		})
		for i := range res {
			rows = append(rows, res[i])
		}

	case in.HasZone():
		var res []db.ReadHostAttributesByZoneRow

		res, err = s.db.ReadHostAttributesByZone(ctx, db.ReadHostAttributesByZoneParams{
			Zone: in.GetZone(),
		})
		for i := range res {
			rows = append(rows, res[i])
		}

	default:
		var res []db.ReadHostAttributesRow

		res, err = s.db.ReadHostAttributes(ctx, db.ReadHostAttributesParams{})
		for i := range res {
			rows = append(rows, res[i])
		}
	}

	if err != nil {
		return dberr(err)
	}

	for i := range rows {
		resp := pb.ReadHostAttrsResponse_builder{
			Zone:      rows[i].Zone,
			Cluster:   rows[i].Cluster,
			Host:      &rows[i].Host,
			Name:      &rows[i].Name,
			Value:     rows[i].Value,
			Protected: proto.Bool(rows[i].IsProtected == 1),
		}.Build()

		if err := out.Send(resp); err != nil {
			return err
		}
	}

	return nil
}

//
// Update
//

func (s *Service) UpdateHostAttr(ctx context.Context, in *pb.UpdateHostAttrRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasHost():
		return nil, errMissingHost
	case !in.HasName():
		return nil, errMissingAttribute
	}

	fields := in.GetFields()

	if fields.HasName() {
		p := db.UpdateHostAttributeNameParams{
			Zone: in.GetZone(),
			Host: in.GetHost(),
			Attr: in.GetName(),
			Name: fields.GetName(),
		}
		if err := s.db.UpdateHostAttributeName(ctx, p); err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasValue() {
		p := db.UpdateHostAttributeValueParams{
			Zone:  in.GetZone(),
			Host:  in.GetHost(),
			Attr:  in.GetName(),
			Value: proto.String(fields.GetValue()),
		}
		if err := s.db.UpdateHostAttributeValue(ctx, p); err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasProtected() {
		p := db.UpdateHostAttributeProtectionParams{
			Zone: in.GetZone(),
			Host: in.GetHost(),
			Attr: in.GetName(),
		}
		if fields.GetProtected() {
			p.IsProtected = 1
		}
		if err := s.db.UpdateHostAttributeProtection(ctx, p); err != nil {
			return nil, dberr(err)
		}
	}

	return new(emptypb.Empty), nil
}

//
// Delete
//

func (s *Service) DeleteHostAttrs(ctx context.Context, in *pb.DeleteHostAttrsRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasHost():
		return nil, errMissingHost
	}

	var ids []int64

	switch {
	case in.HasName():
		row, err := s.db.ReadHostAttribute(ctx, db.ReadHostAttributeParams{
			Zone: in.GetZone(),
			Host: in.GetHost(),
			Attr: in.GetName(),
		})
		if err != nil {
			return nil, dberr(err)
		}
		ids = append(ids, row.ID)
	case in.HasGlob():
		rows, err := s.db.ReadHostAttributesByGlob(ctx, db.ReadHostAttributesByGlobParams{
			Zone: in.GetZone(),
			Host: in.GetHost(),
			Glob: in.GetGlob(),
		})
		if err != nil {
			return nil, dberr(err)
		}
		for i := range rows {
			ids = append(ids, rows[i].ID)
		}
	default:
		return nil, errNameOrGlob
	}

	for _, id := range ids {
		if err := s.db.DeleteAttribute(ctx, db.DeleteAttributeParams{ID: id}); err != nil {
			return nil, dberr(err)
		}
	}

	return new(emptypb.Empty), nil
}

//// Host Interfaces

//
// CREATE
//

func (s Service) CreateHostInterface(ctx context.Context, in *pb.CreateHostInterfaceRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasName():
		return nil, errMissingHost
	}

	p := db.CreateHostInterfaceParams{
		Zone: in.GetZone(),
		Name: in.GetName(),
	}

	if in.HasCluster() {
		p.Cluster = in.GetCluster()
		p.Zone = "" // get zone from cluster
	}

	if err := s.db.CreateHostInterface(ctx, p); err != nil {
		return nil, dberr(err)
	}

	return new(emptypb.Empty), nil
}

//
// READ
//

func (s *Service) ReadHostInterfaces(in *pb.ReadHostInterfacesRequest, out grpc.ServerStreamingServer[pb.ReadHostInterfacesResponse]) error {
	var (
		rows []struct {
			ID              int64
			Name            string
			IP              *string
			MAC             *string
			Netmask         *string
			IsDHCP          int64
			IsPXE           int64
			IsManagement    int64
			Type            *string
			BondMode        *string
			MasterInterface string
			Network         string
			Host            string
			Zone            *string
			Cluster         *string
		}
		err error
	)

	ctx := context.Background()

	switch {
	case in.HasZone() && in.HasHost() && in.HasName():
		var res db.ReadHostInterfaceRow

		res, err = s.db.ReadHostInterface(ctx, db.ReadHostInterfaceParams{
			Zone:    in.GetZone(),
			Cluster: in.GetCluster(), // "" means standalone host (same for queries below)
			Host:    in.GetHost(),
			Name:    in.GetName(),
		})
		rows = append(rows, res)

	case in.HasZone() && in.HasHost() && in.HasGlob(): // clustered, standalone
		var res []db.ReadHostInterfacesByGlobRow

		res, err = s.db.ReadHostInterfacesByGlob(ctx, db.ReadHostInterfacesByGlobParams{
			Zone:    in.GetZone(),
			Cluster: in.GetCluster(),
			Host:    in.GetHost(),
			Glob:    in.GetGlob(),
		})
		for i := range res {
			rows = append(rows, res[i])
		}

	case in.HasZone() && in.HasHost():
		var res []db.ReadHostInterfacesByHostRow

		res, err = s.db.ReadHostInterfacesByHost(ctx, db.ReadHostInterfacesByHostParams{
			Zone:    in.GetZone(),
			Cluster: in.GetCluster(),
			Host:    in.GetHost(),
		})
		for i := range res {
			rows = append(rows, res[i])
		}

	case in.HasZone() && in.HasGlob():
		var res []db.ReadHostInterfacesByClusterRow

		res, err = s.db.ReadHostInterfacesByCluster(ctx, db.ReadHostInterfacesByClusterParams{
			Zone:    in.GetZone(),
			Cluster: in.GetCluster(),
		})
		for i := range res {
			rows = append(rows, res[i])
		}

	case in.HasZone():
		var res []db.ReadHostInterfacesByZoneRow

		res, err = s.db.ReadHostInterfacesByZone(ctx, db.ReadHostInterfacesByZoneParams{
			Zone: in.GetZone(),
		})
		for i := range res {
			rows = append(rows, res[i])
		}

	default:
		var res []db.ReadHostInterfacesRow

		res, err = s.db.ReadHostInterfaces(ctx, db.ReadHostInterfacesParams{})
		for i := range res {
			rows = append(rows, res[i])
		}
	}

	if err != nil {
		return dberr(err)
	}

	for i := range rows {
		resp := pb.ReadHostInterfacesResponse_builder{
			Zone:       rows[i].Zone,
			Cluster:    rows[i].Cluster,
			Host:       &rows[i].Host,
			Name:       &rows[i].Name,
			Ip:         rows[i].IP,
			Mac:        rows[i].MAC,
			Netmask:    rows[i].Netmask,
			Dhcp:       proto.Bool(rows[i].IsDHCP == 1),
			Pxe:        proto.Bool(rows[i].IsPXE == 1),
			Management: proto.Bool(rows[i].IsManagement == 1),
			Master:     &rows[i].MasterInterface,
			Network:    &rows[i].Network,
		}.Build()

		if rows[i].Type != nil {
			resp.SetType(pb.InterfaceType(pb.InterfaceType_value[*rows[i].Type]))
		}

		if rows[i].BondMode != nil {
			resp.SetBondMode(pb.BondMode(pb.BondMode_value[*rows[i].BondMode]))
		}

		if err := out.Send(resp); err != nil {
			return err
		}
	}

	return nil
}

//
// Update
//

func (s *Service) UpdateHostInterface(ctx context.Context, in *pb.UpdateHostInterfaceRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasHost():
		return nil, errMissingHost
	case !in.HasName():
		return nil, errMissingInterface
	}

	fields := in.GetFields()

	if fields.HasName() {
		p := db.UpdateHostInterfaceNameParams{
			Zone:      in.GetZone(),
			Host:      in.GetHost(),
			Interface: in.GetName(),
			Name:      fields.GetName(),
		}
		if err := s.db.UpdateHostInterfaceName(ctx, p); err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasIp() {
		p := db.UpdateHostInterfaceIPParams{
			Zone:      in.GetZone(),
			Host:      in.GetHost(),
			Interface: in.GetName(),
			IP:        proto.String(fields.GetIp()),
		}
		if err := s.db.UpdateHostInterfaceIP(ctx, p); err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasMac() {
		p := db.UpdateHostInterfaceMACParams{
			Zone:      in.GetZone(),
			Host:      in.GetHost(),
			Interface: in.GetName(),
			MAC:       proto.String(fields.GetMac()),
		}
		if err := s.db.UpdateHostInterfaceMAC(ctx, p); err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasNetmask() {
		p := db.UpdateHostInterfaceNetmaskParams{
			Zone:      in.GetZone(),
			Host:      in.GetHost(),
			Interface: in.GetName(),
			Mask:      proto.String(fields.GetNetmask()),
		}
		if err := s.db.UpdateHostInterfaceNetmask(ctx, p); err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasDhcp() {
		p := db.UpdateHostInterfaceDHCPParams{
			Zone:      in.GetZone(),
			Host:      in.GetHost(),
			Interface: in.GetName(),
		}
		if fields.GetDhcp() {
			p.DHCP = 1
		}
		if err := s.db.UpdateHostInterfaceDHCP(ctx, p); err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasPxe() {
		p := db.UpdateHostInterfacePXEParams{
			Zone:      in.GetZone(),
			Host:      in.GetHost(),
			Interface: in.GetName(),
		}
		if fields.GetPxe() {
			p.PXE = 1
		}
		if err := s.db.UpdateHostInterfacePXE(ctx, p); err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasManagement() {
		p := db.UpdateHostInterfaceManagementParams{
			Zone:      in.GetZone(),
			Host:      in.GetHost(),
			Interface: in.GetName(),
		}
		if fields.GetManagement() {
			p.Management = 1
		}
		if err := s.db.UpdateHostInterfaceManagement(ctx, p); err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasType() {
		p := db.UpdateHostInterfaceTypeParams{
			Zone:      in.GetZone(),
			Host:      in.GetHost(),
			Interface: in.GetName(),
			Type:      proto.String(fields.GetType().String()),
		}
		if err := s.db.UpdateHostInterfaceType(ctx, p); err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasBondMode() {
		p := db.UpdateHostInterfaceBondModeParams{
			Zone:      in.GetZone(),
			Host:      in.GetHost(),
			Interface: in.GetName(),
			Bond:      proto.String(fields.GetBondMode().String()),
		}
		if err := s.db.UpdateHostInterfaceBondMode(ctx, p); err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasMaster() {
		p := db.UpdateHostInterfaceMasterParams{
			Zone:      in.GetZone(),
			Host:      in.GetHost(),
			Interface: in.GetName(),
			Master:    fields.GetMaster(),
		}
		if err := s.db.UpdateHostInterfaceMaster(ctx, p); err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasNetwork() {
		p := db.UpdateHostInterfaceNetworkParams{
			Zone:      in.GetZone(),
			Host:      in.GetHost(),
			Interface: in.GetName(),
			Network:   fields.GetNetwork(),
		}
		if err := s.db.UpdateHostInterfaceNetwork(ctx, p); err != nil {
			return nil, dberr(err)
		}
	}

	return new(emptypb.Empty), nil
}

//
// Delete
//

func (s *Service) DeleteHostInterfaces(ctx context.Context, in *pb.DeleteHostInterfacesRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasHost():
		return nil, errMissingHost
	}

	var ids []int64

	switch {
	case in.HasName():
		row, err := s.db.ReadHostInterface(ctx, db.ReadHostInterfaceParams{
			Zone:    in.GetZone(),
			Cluster: in.GetCluster(),
			Host:    in.GetHost(),
			Name:    in.GetName(),
		})
		if err != nil {
			return nil, dberr(err)
		}
		ids = append(ids, row.ID)
	case in.HasGlob():
		rows, err := s.db.ReadHostInterfacesByGlob(ctx, db.ReadHostInterfacesByGlobParams{
			Zone:    in.GetZone(),
			Cluster: in.GetCluster(),
			Host:    in.GetHost(),
			Glob:    in.GetGlob(),
		})
		if err != nil {
			return nil, dberr(err)
		}
		for i := range rows {
			ids = append(ids, rows[i].ID)
		}
	default:
		return nil, errNameOrGlob
	}

	for _, id := range ids {
		if err := s.db.DeleteHostInterface(ctx, db.DeleteHostInterfaceParams{ID: id}); err != nil {
			return nil, dberr(err)
		}
	}

	return new(emptypb.Empty), nil
}
