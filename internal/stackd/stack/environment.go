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

//// Environment

//
// CREATE
//

func (s Service) CreateEnvironment(ctx context.Context, in *pb.CreateEnvironmentRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasName():
		return nil, errMissingEnvironment
	}

	p := db.CreateEnvironmentParams{
		Zone: in.GetZone(),
		Name: in.GetName(),
	}

	if err := s.createEnvironment(ctx, p); err != nil {
		return nil, dberr(err)
	}

	return new(emptypb.Empty), nil
}

func (s Service) createEnvironment(ctx context.Context, params db.CreateEnvironmentParams) error {
	if _, err := s.db.ReadEnvironment(ctx, db.ReadEnvironmentParams(params)); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return s.db.CreateEnvironment(ctx, params)
		}

		return err
	}

	return status.Errorf(codes.AlreadyExists, "zone %q environment %q exists", params.Zone, params.Name)
}

//
// READ
//

func (s Service) ReadEnvironments(in *pb.ReadEnvironmentsRequest, out grpc.ServerStreamingServer[pb.ReadEnvironmentsResponse]) error {
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
		var res db.ReadEnvironmentRow

		res, err = s.db.ReadEnvironment(ctx, db.ReadEnvironmentParams{
			Zone: in.GetZone(),
			Name: in.GetName(),
		})
		rows = append(rows, res)

	case in.HasZone() && in.HasGlob():
		var res []db.ReadEnvironmentsByGlobRow

		res, err = s.db.ReadEnvironmentsByGlob(ctx, db.ReadEnvironmentsByGlobParams{
			Zone: in.GetZone(),
			Glob: in.GetGlob(),
		})
		for i := range res {
			rows = append(rows, res[i])
		}

	case in.HasZone():
		var res []db.ReadEnvironmentsByZoneRow

		res, err = s.db.ReadEnvironmentsByZone(ctx, db.ReadEnvironmentsByZoneParams{
			Zone: in.GetZone(),
		})
		for i := range res {
			rows = append(rows, res[i])
		}

	default:
		var res []db.ReadEnvironmentsRow

		res, err = s.db.ReadEnvironments(ctx, db.ReadEnvironmentsParams{})
		for i := range res {
			rows = append(rows, res[i])
		}
	}

	if err != nil {
		return dberr(err)
	}

	for i := range rows {
		resp := pb.ReadEnvironmentsResponse_builder{
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

func (s Service) UpdateEnvironment(ctx context.Context, in *pb.UpdateEnvironmentRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasName():
		return nil, errMissingEnvironment
	}

	fields := in.GetFields()

	if fields.HasName() {
		err := s.db.UpdateEnvironmentName(ctx, db.UpdateEnvironmentNameParams{
			Zone:        in.GetZone(),
			Environment: in.GetName(),
			Name:        fields.GetName(),
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

func (s Service) DeleteEnvironments(ctx context.Context, in *pb.DeleteEnvironmentsRequest) (*emptypb.Empty, error) {
	if !in.HasZone() {
		return nil, errMissingZone
	}

	var ids []int64

	switch {
	case in.HasName():
		row, err := s.db.ReadEnvironment(ctx, db.ReadEnvironmentParams{
			Zone: in.GetZone(),
			Name: in.GetName(),
		})
		if err != nil {
			return nil, dberr(err)
		}
		ids = append(ids, row.ID)
	case in.HasGlob():
		rows, err := s.db.ReadEnvironmentsByGlob(ctx, db.ReadEnvironmentsByGlobParams{
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
		if err := s.db.DeleteEnvironment(ctx, db.DeleteEnvironmentParams{ID: id}); err != nil {
			return nil, dberr(err)
		}
	}

	return new(emptypb.Empty), nil
}

//// Environment Attributes

//
// CREATE
//

func (s *Service) CreateEnvironmentAttr(ctx context.Context, in *pb.CreateEnvironmentAttrRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasEnvironment():
		return nil, errMissingEnvironment
	case !in.HasName():
		return nil, errMissingAttribute
	}

	p := db.CreateEnvironmentAttributeParams{
		Zone:        in.GetZone(),
		Environment: in.GetEnvironment(),
		Name:        in.GetName(),
	}

	if err := s.db.CreateEnvironmentAttribute(ctx, p); err != nil {
		return nil, dberr(err)
	}

	return new(emptypb.Empty), nil
}

//
// READ
//

func (s *Service) ReadEnvironmentAttrs(in *pb.ReadEnvironmentAttrsRequest, out grpc.ServerStreamingServer[pb.ReadEnvironmentAttrsResponse]) error {
	var (
		rows []struct {
			ID          int64
			Environment string
			Name        string
			Value       *string
			IsProtected int64
		}
		err error
	)

	ctx := context.Background()

	switch {
	case in.HasZone() && in.HasEnvironment() && in.HasName():
		var res db.ReadEnvironmentAttributeRow

		res, err = s.db.ReadEnvironmentAttribute(ctx, db.ReadEnvironmentAttributeParams{
			Zone:        in.GetZone(),
			Environment: in.GetEnvironment(),
			Attr:        in.GetName(),
		})
		rows = append(rows, res)

	case in.HasZone() && in.HasEnvironment() && in.HasGlob():
		var res []db.ReadEnvironmentAttributesByGlobRow

		res, err = s.db.ReadEnvironmentAttributesByGlob(ctx, db.ReadEnvironmentAttributesByGlobParams{
			Zone:        in.GetZone(),
			Environment: in.GetEnvironment(),
			Glob:        in.GetGlob(),
		})
		for i := range res {
			rows = append(rows, res[i])
		}

	case in.HasZone() && in.HasEnvironment():
		var res []db.ReadEnvironmentAttributesByEnvironmentRow

		res, err = s.db.ReadEnvironmentAttributesByEnvironment(ctx, db.ReadEnvironmentAttributesByEnvironmentParams{
			Zone:        in.GetZone(),
			Environment: in.GetEnvironment(),
		})
		for i := range res {
			rows = append(rows, res[i])
		}

	case in.HasZone():
		var res []db.ReadEnvironmentAttributesByZoneRow

		res, err = s.db.ReadEnvironmentAttributesByZone(ctx, db.ReadEnvironmentAttributesByZoneParams{
			Zone: in.GetZone(),
		})
		for i := range res {
			rows = append(rows, res[i])
		}

	default:
		var res []db.ReadEnvironmentAttributesRow

		res, err = s.db.ReadEnvironmentAttributes(ctx, db.ReadEnvironmentAttributesParams{})
		for i := range res {
			rows = append(rows, res[i])
		}
	}

	if err != nil {
		return dberr(err)
	}

	for i := range rows {
		resp := pb.ReadEnvironmentAttrsResponse_builder{
			Environment: &rows[i].Environment,
			Name:        &rows[i].Name,
			Value:       rows[i].Value,
			Protected:   proto.Bool(rows[i].IsProtected == 1),
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

func (s *Service) UpdateEnvironmentAttr(ctx context.Context, in *pb.UpdateEnvironmentAttrRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasEnvironment():
		return nil, errMissingEnvironment
	case !in.HasName():
		return nil, errMissingAttribute
	}

	fields := in.GetFields()

	if fields.HasName() {
		p := db.UpdateEnvironmentAttributeNameParams{
			Zone:        in.GetZone(),
			Environment: in.GetEnvironment(),
			Attr:        in.GetName(),
			Name:        fields.GetName(),
		}
		if err := s.db.UpdateEnvironmentAttributeName(ctx, p); err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasValue() {
		p := db.UpdateEnvironmentAttributeValueParams{
			Zone:        in.GetZone(),
			Environment: in.GetEnvironment(),
			Attr:        in.GetName(),
			Value:       proto.String(fields.GetValue()),
		}
		if err := s.db.UpdateEnvironmentAttributeValue(ctx, p); err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasProtected() {
		p := db.UpdateEnvironmentAttributeProtectionParams{
			Zone:        in.GetZone(),
			Environment: in.GetEnvironment(),
			Attr:        in.GetName(),
		}
		if fields.GetProtected() {
			p.IsProtected = 1
		}
		if err := s.db.UpdateEnvironmentAttributeProtection(ctx, p); err != nil {
			return nil, dberr(err)
		}
	}

	return new(emptypb.Empty), nil
}

//
// Delete
//

func (s *Service) DeleteEnvironmentAttrs(ctx context.Context, in *pb.DeleteEnvironmentAttrsRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasEnvironment():
		return nil, errMissingEnvironment
	}

	var ids []int64

	switch {
	case in.HasName():
		row, err := s.db.ReadEnvironmentAttribute(ctx, db.ReadEnvironmentAttributeParams{
			Zone:        in.GetZone(),
			Environment: in.GetEnvironment(),
			Attr:        in.GetName(),
		})
		if err != nil {
			return nil, dberr(err)
		}
		ids = append(ids, row.ID)
	case in.HasGlob():
		rows, err := s.db.ReadEnvironmentAttributesByGlob(ctx, db.ReadEnvironmentAttributesByGlobParams{
			Zone:        in.GetZone(),
			Environment: in.GetEnvironment(),
			Glob:        in.GetGlob(),
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
