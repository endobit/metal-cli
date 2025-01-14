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

//// Attributes

//
// CREATE
//

func (s Service) CreateGlobalAttr(ctx context.Context, in *pb.CreateGlobalAttrRequest) (*emptypb.Empty, error) {
	if !in.HasName() {
		return nil, errMissingAttribute
	}

	p := db.CreateGlobalAttributeParams{
		Name: in.GetName(),
	}

	if err := s.createGlobalAttr(ctx, p); err != nil {
		return nil, dberr(err)
	}

	return new(emptypb.Empty), nil
}

func (s Service) createGlobalAttr(ctx context.Context, params db.CreateGlobalAttributeParams) error {
	_, err := s.db.ReadGlobalAttribute(ctx, db.ReadGlobalAttributeParams{
		Attr: params.Name,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return s.db.CreateGlobalAttribute(ctx, params)
		}

		return err
	}

	return status.Errorf(codes.AlreadyExists, "attribute %q exists", params.Name)
}

//
// READ
//

func (s Service) ReadGlobalAttrs(in *pb.ReadGlobalAttrsRequest, out grpc.ServerStreamingServer[pb.ReadGlobalAttrsResponse]) error {
	var (
		rows []struct {
			ID          int64
			Name        string
			Value       *string
			IsProtected int64
		}
		err error
	)

	ctx := context.Background()

	switch {
	case in.HasName():
		var res db.ReadGlobalAttributeRow

		res, err = s.db.ReadGlobalAttribute(ctx, db.ReadGlobalAttributeParams{
			Attr: in.GetName(),
		})
		rows = append(rows, res)

	case in.HasGlob():
		var res []db.ReadGlobalAttributesByGlobRow

		res, err = s.db.ReadGlobalAttributesByGlob(ctx, db.ReadGlobalAttributesByGlobParams{
			Glob: in.GetGlob(),
		})
		for i := range res {
			rows = append(rows, res[i])
		}

	default:
		var res []db.ReadGlobalAttributesRow

		res, err = s.db.ReadGlobalAttributes(ctx, db.ReadGlobalAttributesParams{})
		for i := range res {
			rows = append(rows, res[i])
		}
	}

	if err != nil {
		return dberr(err)
	}

	for i := range rows {
		resp := pb.ReadGlobalAttrsResponse_builder{
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

func (s Service) UpdateGlobalAttr(ctx context.Context, in *pb.UpdateGlobalAttrRequest) (*emptypb.Empty, error) {
	if !in.HasName() {
		return nil, errMissingAttribute
	}

	fields := in.GetFields()

	if fields.HasName() {
		p := db.UpdateGlobalAttributeNameParams{
			Attr: in.GetName(),
			Name: fields.GetName(),
		}
		if err := s.db.UpdateGlobalAttributeName(ctx, p); err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasValue() {
		p := db.UpdateGlobalAttributeValueParams{
			Attr:  in.GetName(),
			Value: proto.String(fields.GetValue()),
		}
		if err := s.db.UpdateGlobalAttributeValue(ctx, p); err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasProtected() {
		p := db.UpdateGlobalAttributeProtectionParams{
			Attr: in.GetName(),
		}
		if fields.GetProtected() {
			p.IsProtected = 1
		}
		if err := s.db.UpdateGlobalAttributeProtection(ctx, p); err != nil {
			return nil, dberr(err)
		}
	}

	return new(emptypb.Empty), nil
}

//
// Delete
//

func (s Service) DeleteGlobalAttrs(ctx context.Context, in *pb.DeleteGlobalAttrsRequest) (*emptypb.Empty, error) {
	var ids []int64

	switch {
	case in.HasName():
		row, err := s.db.ReadGlobalAttribute(ctx, db.ReadGlobalAttributeParams{
			Attr: in.GetName(),
		})
		if err != nil {
			return nil, dberr(err)
		}
		ids = append(ids, row.ID)
	case in.HasGlob():
		rows, err := s.db.ReadGlobalAttributesByGlob(ctx, db.ReadGlobalAttributesByGlobParams{
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
