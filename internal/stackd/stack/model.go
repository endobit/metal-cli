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

//// Model

//
// CREATE
//

func (s Service) CreateModel(ctx context.Context, in *pb.CreateModelRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasMake():
		return nil, errMissingMake
	case !in.HasName():
		return nil, errMissingModel
	}

	p := db.CreateModelParams{
		Make: in.GetMake(),
		Name: in.GetName(),
	}

	if err := s.createModel(ctx, p); err != nil {
		return nil, dberr(err)
	}

	return new(emptypb.Empty), nil
}

func (s Service) createModel(ctx context.Context, params db.CreateModelParams) error {
	err := s.createMake(ctx, db.CreateMakeParams{Name: params.Make}) // create make if it doesn't exist
	if err != nil && status.Code(err) != codes.AlreadyExists {
		return dberr(err)
	}

	_, err = s.db.ReadModel(ctx, db.ReadModelParams{
		Make:  params.Make,
		Model: params.Name,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return s.db.CreateModel(ctx, params)
		}

		return err
	}

	return status.Errorf(codes.AlreadyExists, "make %q model %q exists", params.Make, params.Name)
}

//
// READ
//

func (s Service) ReadModels(in *pb.ReadModelsRequest, out grpc.ServerStreamingServer[pb.ReadModelsResponse]) error {
	var (
		rows []struct {
			ID           int64
			Make         string
			Name         string
			Architecture *string
		}
		err error
	)

	ctx := context.Background()

	switch {
	case in.HasMake() && in.HasName():
		var res db.ReadModelRow

		res, err = s.db.ReadModel(ctx, db.ReadModelParams{
			Make:  in.GetMake(),
			Model: in.GetName(),
		})
		rows = append(rows, res)

	case in.HasMake() && in.HasGlob():
		var res []db.ReadModelsByGlobRow

		res, err = s.db.ReadModelsByGlob(ctx, db.ReadModelsByGlobParams{
			Make: in.GetMake(),
			Glob: in.GetGlob(),
		})
		for i := range res {
			rows = append(rows, res[i])
		}

	case in.HasMake():
		var res []db.ReadModelsByMakeRow

		res, err = s.db.ReadModelsByMake(ctx, db.ReadModelsByMakeParams{
			Make: in.GetMake(),
		})
		for i := range res {
			rows = append(rows, res[i])
		}

	default:
		var res []db.ReadModelsRow

		res, err = s.db.ReadModels(ctx, db.ReadModelsParams{})
		for i := range res {
			rows = append(rows, res[i])
		}
	}

	if err != nil {
		return dberr(err)
	}

	for i := range rows {
		resp := pb.ReadModelsResponse_builder{
			Make: &rows[i].Make,
			Name: &rows[i].Name,
		}.Build()
		if rows[i].Architecture != nil {
			resp.SetArchitecture(pb.Architecture(pb.Architecture_value[*rows[i].Architecture]))
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

func (s Service) UpdateModel(ctx context.Context, in *pb.UpdateModelRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasMake():
		return nil, errMissingMake
	case !in.HasName():
		return nil, errMissingModel
	}

	fields := in.GetFields()

	if fields.HasName() {
		err := s.db.UpdateModelName(ctx, db.UpdateModelNameParams{
			Make:  in.GetMake(),
			Model: in.GetName(),
			Name:  fields.GetName(),
		})
		if err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasArchitecture() {
		err := s.db.UpdateModelArchitecture(ctx, db.UpdateModelArchitectureParams{
			Model:        in.GetName(),
			Architecture: proto.String(fields.GetArchitecture().String()),
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

func (s Service) DeleteModels(ctx context.Context, in *pb.DeleteModelsRequest) (*emptypb.Empty, error) {
	if !in.HasMake() {
		return nil, errMissingMake
	}

	var ids []int64

	switch {
	case in.HasName():
		row, err := s.db.ReadModel(ctx, db.ReadModelParams{
			Make:  in.GetMake(),
			Model: in.GetName(),
		})
		if err != nil {
			return nil, dberr(err)
		}
		ids = append(ids, row.ID)
	case in.HasGlob():
		rows, err := s.db.ReadModelsByGlob(ctx, db.ReadModelsByGlobParams{
			Make: in.GetMake(),
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
		if err := s.db.DeleteModel(ctx, db.DeleteModelParams{ID: id}); err != nil {
			return nil, dberr(err)
		}
	}

	return new(emptypb.Empty), nil
}

//// Attributes

//
// CREATE
//

func (s *Service) CreateModelAttr(ctx context.Context, in *pb.CreateModelAttrRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasMake():
		return nil, errMissingMake
	case !in.HasModel():
		return nil, errMissingModel
	case !in.HasName():
		return nil, errMissingAttribute
	}

	p := db.CreateModelAttributeParams{
		Make:  in.GetMake(),
		Model: in.GetModel(),
		Name:  in.GetName(),
	}

	if err := s.db.CreateModelAttribute(ctx, p); err != nil {
		return nil, dberr(err)
	}

	return new(emptypb.Empty), nil
}

//
// READ
//

func (s *Service) ReadModelAttrs(in *pb.ReadModelAttrsRequest, out grpc.ServerStreamingServer[pb.ReadModelAttrsResponse]) error {
	var (
		rows []struct {
			ID          int64
			Model       string
			Make        string
			Name        string
			Value       *string
			IsProtected int64
		}
		err error
	)

	ctx := context.Background()

	switch {
	case in.HasMake() && in.HasModel() && in.HasName():
		var res db.ReadModelAttributeRow

		res, err = s.db.ReadModelAttribute(ctx, db.ReadModelAttributeParams{
			Make:  in.GetMake(),
			Model: in.GetModel(),
			Attr:  in.GetName(),
		})
		rows = append(rows, res)

	case in.HasMake() && in.HasModel() && in.HasGlob():
		var res []db.ReadModelAttributesByGlobRow

		res, err = s.db.ReadModelAttributesByGlob(ctx, db.ReadModelAttributesByGlobParams{
			Make:  in.GetMake(),
			Model: in.GetModel(),
			Glob:  in.GetGlob(),
		})
		for i := range res {
			rows = append(rows, res[i])
		}

	case in.HasMake() && in.HasModel():
		var res []db.ReadModelAttributesByMakeModelRow

		res, err = s.db.ReadModelAttributesByMakeModel(ctx, db.ReadModelAttributesByMakeModelParams{
			Make:  in.GetMake(),
			Model: in.GetModel(),
		})
		for i := range res {
			rows = append(rows, res[i])
		}

	case in.HasMake():
		var res []db.ReadModelAttributesByMakeRow

		res, err = s.db.ReadModelAttributesByMake(ctx, db.ReadModelAttributesByMakeParams{
			Make: in.GetMake(),
		})
		for i := range res {
			rows = append(rows, res[i])
		}

	default:
		var res []db.ReadModelAttributesRow

		res, err = s.db.ReadModelAttributes(ctx, db.ReadModelAttributesParams{})
		for i := range res {
			rows = append(rows, res[i])
		}
	}

	if err != nil {
		return dberr(err)
	}

	for i := range rows {
		resp := pb.ReadModelAttrsResponse_builder{
			Make:      &rows[i].Make,
			Model:     &rows[i].Model,
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

func (s *Service) UpdateModelAttr(ctx context.Context, in *pb.UpdateModelAttrRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasMake():
		return nil, errMissingMake
	case !in.HasModel():
		return nil, errMissingModel
	case !in.HasName():
		return nil, errMissingAttribute
	}

	fields := in.GetFields()

	if fields.HasName() {
		p := db.UpdateModelAttributeNameParams{
			Make:  in.GetMake(),
			Model: in.GetModel(),
			Attr:  in.GetName(),
			Name:  fields.GetName(),
		}
		if err := s.db.UpdateModelAttributeName(ctx, p); err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasValue() {
		p := db.UpdateModelAttributeValueParams{
			Model: in.GetModel(),
			Attr:  in.GetName(),
			Value: proto.String(fields.GetValue()),
		}
		if err := s.db.UpdateModelAttributeValue(ctx, p); err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasProtected() {
		p := db.UpdateModelAttributeProtectionParams{
			Model: in.GetModel(),
			Attr:  in.GetName(),
		}
		if fields.GetProtected() {
			p.IsProtected = 1
		}
		if err := s.db.UpdateModelAttributeProtection(ctx, p); err != nil {
			return nil, dberr(err)
		}
	}

	return new(emptypb.Empty), nil
}

//
// Delete
//

func (s *Service) DeleteModelAttrs(ctx context.Context, in *pb.DeleteModelAttrsRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasMake():
		return nil, errMissingMake
	case !in.HasModel():
		return nil, errMissingModel
	}

	var ids []int64

	switch {
	case in.HasName():
		row, err := s.db.ReadModelAttribute(ctx, db.ReadModelAttributeParams{
			Model: in.GetModel(),
			Attr:  in.GetName(),
		})
		if err != nil {
			return nil, dberr(err)
		}
		ids = append(ids, row.ID)
	case in.HasGlob():
		rows, err := s.db.ReadModelAttributesByGlob(ctx, db.ReadModelAttributesByGlobParams{
			Model: in.GetModel(),
			Glob:  in.GetGlob(),
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
