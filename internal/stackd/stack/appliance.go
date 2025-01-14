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

//// Appliance

//
// CREATE
//

func (s Service) CreateAppliance(ctx context.Context, in *pb.CreateApplianceRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasName():
		return nil, errMissingAppliance
	}

	p := db.CreateApplianceParams{
		Zone: in.GetZone(),
		Name: in.GetName(),
	}

	if err := s.createAppliance(ctx, p); err != nil {
		return nil, dberr(err)
	}

	return new(emptypb.Empty), nil
}

func (s Service) createAppliance(ctx context.Context, params db.CreateApplianceParams) error {
	if _, err := s.db.ReadAppliance(ctx, db.ReadApplianceParams(params)); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return s.db.CreateAppliance(ctx, params)
		}

		return err
	}

	return status.Errorf(codes.AlreadyExists, "zone %q appliance %q exists", params.Zone, params.Name)
}

//
// READ
//

func (s Service) ReadAppliances(in *pb.ReadAppliancesRequest, out grpc.ServerStreamingServer[pb.ReadAppliancesResponse]) error {
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
		var res db.ReadApplianceRow

		res, err = s.db.ReadAppliance(ctx, db.ReadApplianceParams{
			Zone: in.GetZone(),
			Name: in.GetName(),
		})
		rows = append(rows, res)

	case in.HasZone() && in.HasGlob():
		var res []db.ReadAppliancesByGlobRow

		res, err = s.db.ReadAppliancesByGlob(ctx, db.ReadAppliancesByGlobParams{
			Zone: in.GetZone(),
			Glob: in.GetGlob(),
		})
		for i := range res {
			rows = append(rows, res[i])
		}

	case in.HasZone():
		var res []db.ReadAppliancesByZoneRow

		res, err = s.db.ReadAppliancesByZone(ctx, db.ReadAppliancesByZoneParams{
			Zone: in.GetZone(),
		})
		for i := range res {
			rows = append(rows, res[i])
		}

	default:
		var res []db.ReadAppliancesRow

		res, err = s.db.ReadAppliances(ctx, db.ReadAppliancesParams{})
		for i := range res {
			rows = append(rows, res[i])
		}
	}

	if err != nil {
		return dberr(err)
	}

	for i := range rows {
		resp := pb.ReadAppliancesResponse_builder{
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

func (s Service) UpdateAppliance(ctx context.Context, in *pb.UpdateApplianceRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasName():
		return nil, errMissingAppliance
	}

	fields := in.GetFields()

	if fields.HasName() {
		err := s.db.UpdateApplianceName(ctx, db.UpdateApplianceNameParams{
			Zone:      in.GetZone(),
			Appliance: in.GetName(),
			Name:      fields.GetName(),
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

func (s Service) DeleteAppliances(ctx context.Context, in *pb.DeleteAppliancesRequest) (*emptypb.Empty, error) {
	if !in.HasZone() {
		return nil, errMissingZone
	}

	var ids []int64

	switch {
	case in.HasName():
		row, err := s.db.ReadAppliance(ctx, db.ReadApplianceParams{
			Zone: in.GetZone(),
			Name: in.GetName(),
		})
		if err != nil {
			return nil, dberr(err)
		}
		ids = append(ids, row.ID)
	case in.HasGlob():
		rows, err := s.db.ReadAppliancesByGlob(ctx, db.ReadAppliancesByGlobParams{
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
		if err := s.db.DeleteAppliance(ctx, db.DeleteApplianceParams{ID: id}); err != nil {
			return nil, dberr(err)
		}
	}

	return new(emptypb.Empty), nil
}

//// Appliance Attributes

//
// CREATE
//

func (s *Service) CreateApplianceAttr(ctx context.Context, in *pb.CreateApplianceAttrRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasAppliance():
		return nil, errMissingAppliance
	case !in.HasName():
		return nil, errMissingAttribute
	}

	p := db.CreateApplianceAttributeParams{
		Zone:      in.GetZone(),
		Appliance: in.GetAppliance(),
		Name:      in.GetName(),
	}

	if err := s.db.CreateApplianceAttribute(ctx, p); err != nil {
		return nil, dberr(err)
	}

	return new(emptypb.Empty), nil
}

//
// READ
//

func (s *Service) ReadApplianceAttrs(in *pb.ReadApplianceAttrsRequest, out grpc.ServerStreamingServer[pb.ReadApplianceAttrsResponse]) error {
	var (
		rows []struct {
			ID          int64
			Appliance   string
			Name        string
			Value       *string
			IsProtected int64
		}
		err error
	)

	ctx := context.Background()

	switch {
	case in.HasZone() && in.HasAppliance() && in.HasName():
		var res db.ReadApplianceAttributeRow

		res, err = s.db.ReadApplianceAttribute(ctx, db.ReadApplianceAttributeParams{
			Zone:      in.GetZone(),
			Appliance: in.GetAppliance(),
			Attr:      in.GetName(),
		})
		rows = append(rows, res)

	case in.HasZone() && in.HasAppliance() && in.HasGlob():
		var res []db.ReadApplianceAttributesByGlobRow

		res, err = s.db.ReadApplianceAttributesByGlob(ctx, db.ReadApplianceAttributesByGlobParams{
			Zone:      in.GetZone(),
			Appliance: in.GetAppliance(),
			Glob:      in.GetGlob(),
		})
		for i := range res {
			rows = append(rows, res[i])
		}

	case in.HasZone() && in.HasAppliance():
		var res []db.ReadApplianceAttributesByApplianceRow

		res, err = s.db.ReadApplianceAttributesByAppliance(ctx, db.ReadApplianceAttributesByApplianceParams{
			Zone:      in.GetZone(),
			Appliance: in.GetAppliance(),
		})
		for i := range res {
			rows = append(rows, res[i])
		}

	case in.HasZone():
		var res []db.ReadApplianceAttributesByZoneRow

		res, err = s.db.ReadApplianceAttributesByZone(ctx, db.ReadApplianceAttributesByZoneParams{
			Zone: in.GetZone(),
		})
		for i := range res {
			rows = append(rows, res[i])
		}

	default:
		var res []db.ReadApplianceAttributesRow

		res, err = s.db.ReadApplianceAttributes(ctx, db.ReadApplianceAttributesParams{})
		for i := range res {
			rows = append(rows, res[i])
		}
	}

	if err != nil {
		return dberr(err)
	}

	for i := range rows {
		resp := pb.ReadApplianceAttrsResponse_builder{
			Appliance: &rows[i].Appliance,
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

func (s *Service) UpdateApplianceAttr(ctx context.Context, in *pb.UpdateApplianceAttrRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasAppliance():
		return nil, errMissingAppliance
	case !in.HasName():
		return nil, errMissingAttribute
	}

	fields := in.GetFields()

	if fields.HasName() {
		p := db.UpdateApplianceAttributeNameParams{
			Zone:      in.GetZone(),
			Appliance: in.GetAppliance(),
			Attr:      in.GetName(),
			Name:      fields.GetName(),
		}
		if err := s.db.UpdateApplianceAttributeName(ctx, p); err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasValue() {
		p := db.UpdateApplianceAttributeValueParams{
			Zone:      in.GetZone(),
			Appliance: in.GetAppliance(),
			Attr:      in.GetName(),
			Value:     proto.String(fields.GetValue()),
		}
		if err := s.db.UpdateApplianceAttributeValue(ctx, p); err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasProtected() {
		p := db.UpdateApplianceAttributeProtectionParams{
			Zone:      in.GetZone(),
			Appliance: in.GetAppliance(),
			Attr:      in.GetName(),
		}
		if fields.GetProtected() {
			p.IsProtected = 1
		}
		if err := s.db.UpdateApplianceAttributeProtection(ctx, p); err != nil {
			return nil, dberr(err)
		}
	}

	return new(emptypb.Empty), nil
}

//
// Delete
//

func (s *Service) DeleteApplianceAttrs(ctx context.Context, in *pb.DeleteApplianceAttrsRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasAppliance():
		return nil, errMissingAppliance
	}

	var ids []int64

	switch {
	case in.HasName():
		row, err := s.db.ReadApplianceAttribute(ctx, db.ReadApplianceAttributeParams{
			Zone:      in.GetZone(),
			Appliance: in.GetAppliance(),
			Attr:      in.GetName(),
		})
		if err != nil {
			return nil, dberr(err)
		}
		ids = append(ids, row.ID)
	case in.HasGlob():
		rows, err := s.db.ReadApplianceAttributesByGlob(ctx, db.ReadApplianceAttributesByGlobParams{
			Zone:      in.GetZone(),
			Appliance: in.GetAppliance(),
			Glob:      in.GetGlob(),
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
