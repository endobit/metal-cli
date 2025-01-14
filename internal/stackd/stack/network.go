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

//// Network

//
// CREATE
//

func (s Service) CreateNetwork(ctx context.Context, in *pb.CreateNetworkRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasName():
		return nil, errMissingNetwork
	}

	p := db.CreateNetworkParams{
		Zone: in.GetZone(),
		Name: in.GetName(),
	}

	if err := s.db.CreateNetwork(ctx, p); err != nil {
		return nil, dberr(err)
	}

	return new(emptypb.Empty), nil
}

//
// READ
//

func (s Service) ReadNetworks(in *pb.ReadNetworksRequest, out grpc.ServerStreamingServer[pb.ReadNetworksResponse]) error {
	var (
		rows []struct {
			ID      int64
			Name    string
			Address *string
			Gateway *string
			IsPXE   int64
			MTU     int64
			Zone    string
		}
		err error
	)

	ctx := context.Background()

	switch {
	case in.HasZone() && in.HasName():
		var res db.ReadNetworkRow

		res, err = s.db.ReadNetwork(ctx, db.ReadNetworkParams{
			Zone: in.GetZone(),
			Name: in.GetName(),
		})
		rows = append(rows, res)

	case in.HasZone() && in.HasGlob():
		var res []db.ReadNetworksByGlobRow

		res, err = s.db.ReadNetworksByGlob(ctx, db.ReadNetworksByGlobParams{
			Zone: in.GetZone(),
			Glob: in.GetGlob(),
		})
		for i := range res {
			rows = append(rows, res[i])
		}

	case in.HasZone():
		var res []db.ReadNetworksByZoneRow

		res, err = s.db.ReadNetworksByZone(ctx, db.ReadNetworksByZoneParams{
			Zone: in.GetZone(),
		})
		for i := range res {
			rows = append(rows, res[i])
		}

	default:
		var res []db.ReadNetworksRow

		res, err = s.db.ReadNetworks(ctx, db.ReadNetworksParams{})
		for i := range res {
			rows = append(rows, res[i])
		}
	}

	if err != nil {
		return dberr(err)
	}

	for i := range rows {
		mtu := rows[i].MTU
		if 0 < mtu || mtu > math.MaxUint32 {
			return fmt.Errorf("mtu %v is not valid", mtu)
		}

		resp := pb.ReadNetworksResponse_builder{
			Zone:    &rows[i].Zone,
			Name:    &rows[i].Name,
			Address: rows[i].Address,
			Gateway: rows[i].Gateway,
			Pxe:     proto.Bool(rows[i].IsPXE == 1),
			Mtu:     proto.Uint32(uint32(mtu)),
		}.Build()

		if err := out.Send(resp); err != nil {
			return err
		}
	}

	return nil
}

//
// UPDATE
//

func (s Service) UpdateNetwork(ctx context.Context, in *pb.UpdateNetworkRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasName():
		return nil, errMissingNetwork
	}

	fields := in.GetFields()

	if fields.HasName() {
		err := s.db.UpdateNetworkName(ctx, db.UpdateNetworkNameParams{
			Zone:    in.GetZone(),
			Network: in.GetName(),
			Name:    fields.GetName(),
		})
		if err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasAddress() {
		err := s.db.UpdateNetworkAddress(ctx, db.UpdateNetworkAddressParams{
			Zone:    in.GetZone(),
			Network: in.GetName(),
			Address: proto.String(fields.GetAddress()),
		})
		if err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasGateway() {
		err := s.db.UpdateNetworkGateway(ctx, db.UpdateNetworkGatewayParams{
			Zone:    in.GetZone(),
			Network: in.GetName(),
			Gateway: proto.String(fields.GetGateway()),
		})
		if err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasPxe() {
		p := db.UpdateNetworkPXEParams{
			Zone:    in.GetZone(),
			Network: in.GetName(),
		}
		if fields.GetPxe() {
			p.IsPXE = 1
		}
		if err := s.db.UpdateNetworkPXE(ctx, p); err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasMtu() {
		err := s.db.UpdateNetworkMTU(ctx, db.UpdateNetworkMTUParams{
			Zone:    in.GetZone(),
			Network: in.GetName(),
			MTU:     int64(fields.GetMtu()),
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

func (s Service) DeleteNetworks(ctx context.Context, in *pb.DeleteNetworksRequest) (*emptypb.Empty, error) {
	if !in.HasZone() {
		return nil, errMissingZone
	}

	var ids []int64

	switch {
	case in.HasName():
		row, err := s.db.ReadNetwork(ctx, db.ReadNetworkParams{
			Zone: in.GetZone(),
			Name: in.GetName(),
		})
		if err != nil {
			return nil, dberr(err)
		}
		ids = append(ids, row.ID)
	case in.HasGlob():
		rows, err := s.db.ReadNetworksByGlob(ctx, db.ReadNetworksByGlobParams{
			Zone: in.GetZone(),
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
		if err := s.db.DeleteNetwork(ctx, db.DeleteNetworkParams{ID: id}); err != nil {
			return nil, dberr(err)
		}
	}

	return new(emptypb.Empty), nil
}
