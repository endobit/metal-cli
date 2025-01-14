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

//// Cluster

//
// CREATE
//

func (s Service) CreateCluster(ctx context.Context, in *pb.CreateClusterRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasName():
		return nil, errMissingCluster
	}

	p := db.CreateClusterParams{
		Zone: in.GetZone(),
		Name: in.GetName(),
	}

	if err := s.createCluster(ctx, p); err != nil {
		return nil, dberr(err)
	}

	return new(emptypb.Empty), nil
}

func (s Service) createCluster(ctx context.Context, params db.CreateClusterParams) error {
	if _, err := s.db.ReadCluster(ctx, db.ReadClusterParams(params)); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return s.db.CreateCluster(ctx, params)
		}

		return err
	}

	return status.Errorf(codes.AlreadyExists, "zone %q cluster %q exists", params.Zone, params.Name)
}

//
// READ
//

func (s Service) ReadClusters(in *pb.ReadClustersRequest, out grpc.ServerStreamingServer[pb.ReadClustersResponse]) error {
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
		var res db.ReadClusterRow

		res, err = s.db.ReadCluster(ctx, db.ReadClusterParams{
			Zone: in.GetZone(),
			Name: in.GetName(),
		})
		rows = append(rows, res)

	case in.HasZone() && in.HasGlob():
		var res []db.ReadClustersByGlobRow

		res, err = s.db.ReadClustersByGlob(ctx, db.ReadClustersByGlobParams{
			Zone: in.GetZone(),
			Glob: in.GetGlob(),
		})
		for i := range res {
			rows = append(rows, res[i])
		}

	case in.HasZone():
		var res []db.ReadClustersByZoneRow

		res, err = s.db.ReadClustersByZone(ctx, db.ReadClustersByZoneParams{
			Zone: in.GetZone(),
		})
		for i := range res {
			rows = append(rows, res[i])
		}

	default:
		var res []db.ReadClustersRow

		res, err = s.db.ReadClusters(ctx, db.ReadClustersParams{})
		for i := range res {
			rows = append(rows, res[i])
		}
	}

	if err != nil {
		return dberr(err)
	}

	for i := range rows {
		resp := pb.ReadClustersResponse_builder{
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

func (s Service) UpdateCluster(ctx context.Context, in *pb.UpdateClusterRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasName():
		return nil, errMissingCluster
	}

	fields := in.GetFields()

	if fields.HasName() {
		err := s.db.UpdateClusterName(ctx, db.UpdateClusterNameParams{
			Zone:    in.GetZone(),
			Cluster: in.GetName(),
			Name:    fields.GetName(),
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

func (s Service) DeleteClusters(ctx context.Context, in *pb.DeleteClustersRequest) (*emptypb.Empty, error) {
	if !in.HasZone() {
		return nil, errMissingZone
	}

	var ids []int64

	switch {
	case in.HasName():
		row, err := s.db.ReadCluster(ctx, db.ReadClusterParams{
			Zone: in.GetZone(),
			Name: in.GetName(),
		})
		if err != nil {
			return nil, dberr(err)
		}
		ids = append(ids, row.ID)
	case in.HasGlob():
		rows, err := s.db.ReadClustersByGlob(ctx, db.ReadClustersByGlobParams{
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
		if err := s.db.DeleteCluster(ctx, db.DeleteClusterParams{ID: id}); err != nil {
			return nil, dberr(err)
		}
	}

	return new(emptypb.Empty), nil
}

//// Cluster Attributes

//
// CREATE
//

func (s *Service) CreateClusterAttr(ctx context.Context, in *pb.CreateClusterAttrRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasCluster():
		return nil, errMissingCluster
	case !in.HasName():
		return nil, errMissingAttribute
	}

	p := db.CreateClusterAttributeParams{
		Zone:    in.GetZone(),
		Cluster: in.GetCluster(),
		Name:    in.GetName(),
	}

	if err := s.db.CreateClusterAttribute(ctx, p); err != nil {
		return nil, dberr(err)
	}

	return new(emptypb.Empty), nil
}

//
// READ
//

func (s *Service) ReadClusterAttrs(in *pb.ReadClusterAttrsRequest, out grpc.ServerStreamingServer[pb.ReadClusterAttrsResponse]) error {
	var (
		rows []struct {
			ID          int64
			Cluster     string
			Name        string
			Value       *string
			IsProtected int64
		}
		err error
	)

	ctx := context.Background()

	switch {
	case in.HasZone() && in.HasCluster() && in.HasName():
		var res db.ReadClusterAttributeRow

		res, err = s.db.ReadClusterAttribute(ctx, db.ReadClusterAttributeParams{
			Zone:    in.GetZone(),
			Cluster: in.GetCluster(),
			Attr:    in.GetName(),
		})
		rows = append(rows, res)

	case in.HasZone() && in.HasCluster() && in.HasGlob():
		var res []db.ReadClusterAttributesByGlobRow

		res, err = s.db.ReadClusterAttributesByGlob(ctx, db.ReadClusterAttributesByGlobParams{
			Zone:    in.GetZone(),
			Cluster: in.GetCluster(),
			Glob:    in.GetGlob(),
		})
		for i := range res {
			rows = append(rows, res[i])
		}

	case in.HasZone() && in.HasCluster():
		var res []db.ReadClusterAttributesByClusterRow

		res, err = s.db.ReadClusterAttributesByCluster(ctx, db.ReadClusterAttributesByClusterParams{
			Zone:    in.GetZone(),
			Cluster: in.GetCluster(),
		})
		for i := range res {
			rows = append(rows, res[i])
		}

	case in.HasZone():
		var res []db.ReadClusterAttributesByZoneRow

		res, err = s.db.ReadClusterAttributesByZone(ctx, db.ReadClusterAttributesByZoneParams{
			Zone: in.GetZone(),
		})
		for i := range res {
			rows = append(rows, res[i])
		}

	default:
		var res []db.ReadClusterAttributesRow

		res, err = s.db.ReadClusterAttributes(ctx, db.ReadClusterAttributesParams{})
		for i := range res {
			rows = append(rows, res[i])
		}
	}

	if err != nil {
		return dberr(err)
	}

	for i := range rows {
		resp := pb.ReadClusterAttrsResponse_builder{
			Cluster:   &rows[i].Cluster,
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

func (s *Service) UpdateClusterAttr(ctx context.Context, in *pb.UpdateClusterAttrRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasCluster():
		return nil, errMissingCluster
	case !in.HasName():
		return nil, errMissingAttribute
	}

	fields := in.GetFields()

	if fields.HasName() {
		p := db.UpdateClusterAttributeNameParams{
			Zone:    in.GetZone(),
			Cluster: in.GetCluster(),
			Attr:    in.GetName(),
			Name:    fields.GetName(),
		}
		if err := s.db.UpdateClusterAttributeName(ctx, p); err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasValue() {
		p := db.UpdateClusterAttributeValueParams{
			Zone:    in.GetZone(),
			Cluster: in.GetCluster(),
			Attr:    in.GetName(),
			Value:   proto.String(fields.GetValue()),
		}
		if err := s.db.UpdateClusterAttributeValue(ctx, p); err != nil {
			return nil, dberr(err)
		}
	}

	if fields.HasProtected() {
		p := db.UpdateClusterAttributeProtectionParams{
			Zone:    in.GetZone(),
			Cluster: in.GetCluster(),
			Attr:    in.GetName(),
		}
		if fields.GetProtected() {
			p.IsProtected = 1
		}
		if err := s.db.UpdateClusterAttributeProtection(ctx, p); err != nil {
			return nil, dberr(err)
		}
	}

	return new(emptypb.Empty), nil
}

//
// Delete
//

func (s *Service) DeleteClusterAttrs(ctx context.Context, in *pb.DeleteClusterAttrsRequest) (*emptypb.Empty, error) {
	switch {
	case !in.HasZone():
		return nil, errMissingZone
	case !in.HasCluster():
		return nil, errMissingCluster
	}

	var ids []int64

	switch {
	case in.HasName():
		row, err := s.db.ReadClusterAttribute(ctx, db.ReadClusterAttributeParams{
			Zone:    in.GetZone(),
			Cluster: in.GetCluster(),
			Attr:    in.GetName(),
		})
		if err != nil {
			return nil, dberr(err)
		}
		ids = append(ids, row.ID)
	case in.HasGlob():
		rows, err := s.db.ReadClusterAttributesByGlob(ctx, db.ReadClusterAttributesByGlobParams{
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
		if err := s.db.DeleteAttribute(ctx, db.DeleteAttributeParams{ID: id}); err != nil {
			return nil, dberr(err)
		}
	}

	return new(emptypb.Empty), nil
}
