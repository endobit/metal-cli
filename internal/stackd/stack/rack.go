package stack

import (
	"context"
	"database/sql"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/endobit/stack/internal/generated/go/db"
	pb "github.com/endobit/stack/internal/generated/go/proto/stack/v1"
)

//// Rack

//
// CREATE
//

func (s Service) CreateRack(ctx context.Context, in *pb.CreateRackRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasName():
		return nil, errMissingRack
	}

	p := db.CreateRackParams{
		Zone: in.GetZone(),
		Name: in.GetName(),
	}

	if err := s.createRack(ctx, p); err != nil {
		return nil, dberr(err)
	}

	return new(emptypb.Empty), nil
}

func (s Service) createRack(ctx context.Context, params db.CreateRackParams) error {
	_, err := s.db.ReadRack(ctx, db.ReadRackParams{
		Zone: params.Zone,
		Name: params.Name,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}

	return s.db.CreateRack(ctx, params)
}

//
// READ
//

func (s Service) ReadRacks(in *pb.ReadRacksRequest, out grpc.ServerStreamingServer[pb.ReadRacksResponse]) error {
	var (
		rows []struct {
			ID   int64
			Name string
			Zone string
		}
		err error
	)

	ctx := context.Background()

	switch {
	case in.HasZone() && in.HasName():
		var res db.ReadRackRow

		res, err = s.db.ReadRack(ctx, db.ReadRackParams{
			Zone: in.GetZone(),
			Name: in.GetName(),
		})
		rows = append(rows, res)

	case in.HasZone() && in.HasGlob():
		var res []db.ReadRacksByGlobRow

		res, err = s.db.ReadRacksByGlob(ctx, db.ReadRacksByGlobParams{
			Zone: in.GetZone(),
			Glob: in.GetGlob(),
		})
		for i := range res {
			rows = append(rows, res[i])
		}

	case in.HasZone():
		var res []db.ReadRacksByZoneRow

		res, err = s.db.ReadRacksByZone(ctx, db.ReadRacksByZoneParams{
			Zone: in.GetZone(),
		})
		for i := range res {
			rows = append(rows, res[i])
		}

	default:
		var res []db.ReadRacksRow

		res, err = s.db.ReadRacks(ctx, db.ReadRacksParams{})
		for i := range res {
			rows = append(rows, res[i])
		}
	}

	if err != nil {
		return dberr(err)
	}

	for i := range rows {
		resp := pb.ReadRacksResponse_builder{
			Zone: &rows[i].Zone,
			Name: &rows[i].Name,
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

func (s Service) UpdateRack(ctx context.Context, in *pb.UpdateRackRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasName():
		return nil, errMissingRack
	}

	fields := in.GetFields()

	if fields.HasName() {
		err := s.db.UpdateRackName(ctx, db.UpdateRackNameParams{
			Zone: in.GetZone(),
			Rack: in.GetName(),
			Name: fields.GetName(),
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

func (s Service) DeleteRacks(ctx context.Context, in *pb.DeleteRacksRequest) (*emptypb.Empty, error) {
	if !in.HasZone() {
		return nil, errMissingZone
	}

	var ids []int64

	switch {
	case in.HasName():
		row, err := s.db.ReadRack(ctx, db.ReadRackParams{
			Zone: in.GetZone(),
			Name: in.GetName(),
		})
		if err != nil {
			return nil, dberr(err)
		}
		ids = append(ids, row.ID)
	case in.HasGlob():
		rows, err := s.db.ReadRacksByGlob(ctx, db.ReadRacksByGlobParams{
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
		if err := s.db.DeleteRack(ctx, db.DeleteRackParams{ID: id}); err != nil {
			return nil, dberr(err)
		}
	}

	return new(emptypb.Empty), nil
}

//// Rack Attributes

//
// CREATE
//

func (s *Service) CreateRackAttr(ctx context.Context, in *pb.CreateRackAttrRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasRack():
		return nil, errMissingRack
	case !in.HasName():
		return nil, errMissingAttribute
	}

	p := db.CreateRackAttributeParams{
		Zone: in.GetZone(),
		Rack: in.GetRack(),
		Name: in.GetName(),
	}

	if err := s.db.CreateRackAttribute(ctx, p); err != nil {
		return nil, dberr(err)
	}

	return new(emptypb.Empty), nil
}

//
// READ
//

func (s *Service) ReadRackAttrs(in *pb.ReadRackAttrsRequest, out grpc.ServerStreamingServer[pb.ReadRackAttrsResponse]) error {
	var (
		rows []struct {
			ID          int64
			Rack        string
			Name        string
			Value       *string
			IsProtected int64
		}
		err error
	)

	ctx := context.Background()

	switch {
	case in.HasZone() && in.HasRack() && in.HasName():
		var res db.ReadRackAttributeRow

		res, err = s.db.ReadRackAttribute(ctx, db.ReadRackAttributeParams{
			Zone: in.GetZone(),
			Rack: in.GetRack(),
			Attr: in.GetName(),
		})
		rows = append(rows, res)

	case in.HasZone() && in.HasRack() && in.HasGlob():
		var res []db.ReadRackAttributesByGlobRow

		res, err = s.db.ReadRackAttributesByGlob(ctx, db.ReadRackAttributesByGlobParams{
			Zone: in.GetZone(),
			Rack: in.GetRack(),
			Glob: in.GetGlob(),
		})
		for i := range res {
			rows = append(rows, res[i])
		}

	case in.HasZone() && in.HasRack():
		var res []db.ReadRackAttributesByRackRow

		res, err = s.db.ReadRackAttributesByRack(ctx, db.ReadRackAttributesByRackParams{
			Zone: in.GetZone(),
			Rack: in.GetRack(),
		})
		for i := range res {
			rows = append(rows, res[i])
		}

	case in.HasZone():
		var res []db.ReadRackAttributesByZoneRow

		res, err = s.db.ReadRackAttributesByZone(ctx, db.ReadRackAttributesByZoneParams{
			Zone: in.GetZone(),
		})
		for i := range res {
			rows = append(rows, res[i])
		}

	default:
		var res []db.ReadRackAttributesRow

		res, err = s.db.ReadRackAttributes(ctx, db.ReadRackAttributesParams{})
		for i := range res {
			rows = append(rows, res[i])
		}
	}

	if err != nil {
		return dberr(err)
	}

	for i := range rows {
		resp := pb.ReadRackAttrsResponse_builder{
			Rack:      &rows[i].Rack,
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

func (s *Service) UpdateRackAttr(ctx context.Context, in *pb.UpdateRackAttrRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasRack():
		return nil, errMissingRack
	case !in.HasName():
		return nil, errMissingAttribute
	}

	fields := in.GetFields()

	if fields.HasName() {
		p := db.UpdateRackAttributeNameParams{
			Zone: in.GetZone(),
			Rack: in.GetRack(),
			Attr: in.GetName(),
			Name: fields.GetName(),
		}
		if err := s.db.UpdateRackAttributeName(ctx, p); err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasValue() {
		p := db.UpdateRackAttributeValueParams{
			Zone:  in.GetZone(),
			Rack:  in.GetRack(),
			Attr:  in.GetName(),
			Value: proto.String(fields.GetValue()),
		}
		if err := s.db.UpdateRackAttributeValue(ctx, p); err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasProtected() {
		p := db.UpdateRackAttributeProtectionParams{
			Zone: in.GetZone(),
			Rack: in.GetRack(),
			Attr: in.GetName(),
		}
		if fields.GetProtected() {
			p.IsProtected = 1
		}
		if err := s.db.UpdateRackAttributeProtection(ctx, p); err != nil {
			return nil, dberr(err)
		}
	}

	return new(emptypb.Empty), nil
}

//
// Delete
//

func (s *Service) DeleteRackAttrs(ctx context.Context, in *pb.DeleteRackAttrsRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasRack():
		return nil, errMissingRack
	}

	var ids []int64

	switch {
	case in.HasName():
		row, err := s.db.ReadRackAttribute(ctx, db.ReadRackAttributeParams{
			Zone: in.GetZone(),
			Rack: in.GetRack(),
			Attr: in.GetName(),
		})
		if err != nil {
			return nil, dberr(err)
		}
		ids = append(ids, row.ID)
	case in.HasGlob():
		rows, err := s.db.ReadRackAttributesByGlob(ctx, db.ReadRackAttributesByGlobParams{
			Zone: in.GetZone(),
			Rack: in.GetRack(),
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
