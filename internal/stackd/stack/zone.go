package stack

import (
	"context"
	"database/sql"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/endobit/stack/internal/generated/go/db"
	pb "github.com/endobit/stack/internal/generated/go/proto/stack/v1"
)

//// Zone

//
// CREATE
//

func (s Service) CreateZone(ctx context.Context, in *pb.CreateZoneRequest) (*emptypb.Empty, error) {
	if !in.HasName() {
		return nil, errMissingZone
	}

	p := db.CreateZoneParams{
		Name: in.GetName(),
	}

	if err := s.createZone(ctx, p); err != nil {
		return nil, dberr(err)
	}

	return new(emptypb.Empty), nil
}

func (s Service) createZone(ctx context.Context, params db.CreateZoneParams) error {
	_, err := s.db.ReadZone(ctx, db.ReadZoneParams{
		Zone: params.Name,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return s.db.CreateZone(ctx, params)
		}

		return err
	}

	return status.Errorf(codes.AlreadyExists, "zone %q exists", params.Name)
}

//
// READ
//

func (s Service) ReadZones(in *pb.ReadZonesRequest, out grpc.ServerStreamingServer[pb.ReadZonesResponse]) error {
	var (
		rows []db.Zone
		err  error
	)

	ctx := context.Background()

	switch {
	case in.HasName():
		var res db.Zone

		res, err = s.db.ReadZone(ctx, db.ReadZoneParams{Zone: in.GetName()})
		rows = append(rows, res)

	case in.HasGlob():
		rows, err = s.db.ReadZonesByGlob(ctx, db.ReadZonesByGlobParams{Glob: in.GetGlob()})

	default:
		rows, err = s.db.ReadZones(ctx, db.ReadZonesParams{})
	}

	if err != nil {
		return dberr(err)
	}

	for i := range rows {
		resp := pb.ReadZonesResponse_builder{
			Name:     &rows[i].Name,
			TimeZone: rows[i].TimeZone,
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

func (s Service) UpdateZone(ctx context.Context, in *pb.UpdateZoneRequest) (*emptypb.Empty, error) {
	fields := in.GetFields()

	if fields.HasName() {
		err := s.db.UpdateZoneName(ctx, db.UpdateZoneNameParams{
			Zone: in.GetName(),
			Name: fields.GetName(),
		})
		if err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasTimeZone() {
		err := s.db.UpdateZoneTimeZone(ctx, db.UpdateZoneTimeZoneParams{
			Zone:     in.GetName(),
			TimeZone: proto.String(fields.GetTimeZone()),
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

func (s Service) DeleteZones(ctx context.Context, in *pb.DeleteZonesRequest) (*emptypb.Empty, error) {
	var ids []int64

	switch {
	case in.HasName():
		row, err := s.db.ReadZone(ctx, db.ReadZoneParams{Zone: in.GetName()})
		if err != nil {
			return nil, dberr(err)
		}
		ids = append(ids, row.ID)
	case in.HasGlob():
		rows, err := s.db.ReadZonesByGlob(ctx, db.ReadZonesByGlobParams{Glob: in.GetGlob()})
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
		if err := s.db.DeleteZone(ctx, db.DeleteZoneParams{ID: id}); err != nil {
			return nil, dberr(err)
		}
	}

	return new(emptypb.Empty), nil
}

//// Zone Attributes

//
// CREATE
//

func (s *Service) CreateZoneAttr(ctx context.Context, in *pb.CreateZoneAttrRequest) (*emptypb.Empty, error) {
	p := db.CreateZoneAttributeParams{
		Zone: in.GetZone(),
		Name: in.GetName(),
	}

	if err := s.db.CreateZoneAttribute(ctx, p); err != nil {
		return nil, dberr(err)
	}

	return new(emptypb.Empty), nil
}

//
// READ
//

func (s *Service) ReadZoneAttrs(in *pb.ReadZoneAttrsRequest, out grpc.ServerStreamingServer[pb.ReadZoneAttrsResponse]) error {
	var (
		rows []struct {
			ID          int64
			Zone        string
			Name        string
			Value       *string
			IsProtected int64
		}
		err error
	)

	ctx := context.Background()

	switch {
	case in.HasZone() && in.HasName():
		var res db.ReadZoneAttributeRow

		res, err = s.db.ReadZoneAttribute(ctx, db.ReadZoneAttributeParams{
			Zone: in.GetZone(),
			Attr: in.GetName(),
		})
		rows = append(rows, res)

	case in.HasZone() && in.HasGlob():
		var res []db.ReadZoneAttributesByGlobRow

		res, err = s.db.ReadZoneAttributesByGlob(ctx, db.ReadZoneAttributesByGlobParams{
			Zone: in.GetZone(),
			Glob: in.GetGlob(),
		})
		for i := range res {
			rows = append(rows, res[i])
		}

	case in.HasZone():
		var res []db.ReadZoneAttributesByZoneRow

		res, err = s.db.ReadZoneAttributesByZone(ctx, db.ReadZoneAttributesByZoneParams{
			Zone: in.GetZone(),
		})
		for i := range res {
			rows = append(rows, res[i])
		}

	default:
		var res []db.ReadZoneAttributesRow

		res, err = s.db.ReadZoneAttributes(ctx, db.ReadZoneAttributesParams{})
		for i := range res {
			rows = append(rows, res[i])
		}
	}

	if err != nil {
		return dberr(err)
	}

	for i := range rows {
		resp := pb.ReadZoneAttrsResponse_builder{
			Zone:      &rows[i].Zone,
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

func (s *Service) UpdateZoneAttr(ctx context.Context, in *pb.UpdateZoneAttrRequest) (*emptypb.Empty, error) {
	fields := in.GetFields()

	if fields.HasName() {
		p := db.UpdateZoneAttributeNameParams{
			Zone: in.GetZone(),
			Attr: in.GetName(),
			Name: fields.GetName(),
		}
		if err := s.db.UpdateZoneAttributeName(ctx, p); err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasValue() {
		p := db.UpdateZoneAttributeValueParams{
			Zone:  in.GetZone(),
			Attr:  in.GetName(),
			Value: proto.String(fields.GetValue()),
		}
		if err := s.db.UpdateZoneAttributeValue(ctx, p); err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasProtected() {
		p := db.UpdateZoneAttributeProtectionParams{
			Zone: in.GetZone(),
			Attr: in.GetName(),
		}
		if fields.GetProtected() {
			p.IsProtected = 1
		}
		if err := s.db.UpdateZoneAttributeProtection(ctx, p); err != nil {
			return nil, dberr(err)
		}
	}

	return new(emptypb.Empty), nil
}

//
// Delete
//

func (s *Service) DeleteZoneAttrs(ctx context.Context, in *pb.DeleteZoneAttrsRequest) (*emptypb.Empty, error) {
	var ids []int64

	switch {
	case in.HasName():
		row, err := s.db.ReadZoneAttribute(ctx, db.ReadZoneAttributeParams{
			Zone: in.GetZone(),
			Attr: in.GetName(),
		})
		if err != nil {
			return nil, dberr(err)
		}
		ids = append(ids, row.ID)
	case in.HasGlob():
		rows, err := s.db.ReadZoneAttributesByGlob(ctx, db.ReadZoneAttributesByGlobParams{
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
		if err := s.db.DeleteAttribute(ctx, db.DeleteAttributeParams{ID: id}); err != nil {
			return nil, dberr(err)
		}
	}

	return new(emptypb.Empty), nil
}
