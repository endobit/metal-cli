package stack

import (
	"context"
	"database/sql"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/endobit/stack/internal/generated/go/db"
	pb "github.com/endobit/stack/internal/generated/go/proto/stack/v1"
)

//// Make

//
// CREATE
//

func (s Service) CreateMake(ctx context.Context, in *pb.CreateMakeRequest) (*emptypb.Empty, error) {
	if !in.HasName() {
		return nil, errMissingMake
	}

	p := db.CreateMakeParams{
		Name: in.GetName(),
	}

	if err := s.createMake(ctx, p); err != nil {
		return nil, dberr(err)
	}

	return new(emptypb.Empty), nil
}

func (s Service) createMake(ctx context.Context, params db.CreateMakeParams) error {
	_, err := s.db.ReadMake(ctx, db.ReadMakeParams{
		Make: params.Name,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return s.db.CreateMake(ctx, params)
		}

		return err
	}

	return status.Errorf(codes.AlreadyExists, "make %q exists", params.Name)
}

//
// READ
//

func (s Service) ReadMakes(in *pb.ReadMakesRequest, out grpc.ServerStreamingServer[pb.ReadMakesResponse]) error {
	var (
		rows []db.Make
		err  error
	)

	ctx := context.Background()

	switch {
	case in.HasName():
		var res db.Make

		res, err = s.db.ReadMake(ctx, db.ReadMakeParams{Make: in.GetName()})
		rows = append(rows, res)

	case in.HasGlob():
		rows, err = s.db.ReadMakesByGlob(ctx, db.ReadMakesByGlobParams{Glob: in.GetGlob()})

	default:
		rows, err = s.db.ReadMakes(ctx, db.ReadMakesParams{})
	}

	if err != nil {
		return dberr(err)
	}

	for i := range rows {
		resp := pb.ReadMakesResponse_builder{
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

func (s Service) UpdateMake(ctx context.Context, in *pb.UpdateMakeRequest) (*emptypb.Empty, error) {
	fields := in.GetFields()

	if fields.HasName() {
		err := s.db.UpdateMakeName(ctx, db.UpdateMakeNameParams{
			Make: in.GetName(),
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

func (s Service) DeleteMakes(ctx context.Context, in *pb.DeleteMakesRequest) (*emptypb.Empty, error) {
	var ids []int64

	switch {
	case in.HasName():
		row, err := s.db.ReadMake(ctx, db.ReadMakeParams{Make: in.GetName()})
		if err != nil {
			return nil, dberr(err)
		}
		ids = append(ids, row.ID)
	case in.HasGlob():
		rows, err := s.db.ReadMakesByGlob(ctx, db.ReadMakesByGlobParams{Glob: in.GetGlob()})
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
		if err := s.db.DeleteMake(ctx, db.DeleteMakeParams{ID: id}); err != nil {
			return nil, dberr(err)
		}
	}

	return new(emptypb.Empty), nil
}
