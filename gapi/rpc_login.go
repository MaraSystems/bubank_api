package gapi

import (
	"context"
	"fmt"

	"github.com/MaraSystems/bubank_api/domains/users"
	pb "github.com/MaraSystems/bubank_api/pb"
	"github.com/MaraSystems/bubank_api/utils"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := s.Store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, status.Error(codes.NotFound, fmt.Sprintf("user not found: %s", req.GetUsername()))
		}

		return nil, status.Error(codes.Internal, "failed to get user")
	}

	err = utils.VerifyPassword(user.HashedPassword, req.Password)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "password is incorrect")
	}

	token, err := s.TokenMaker.Create(req.Username, s.Config.AccessDuration)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create token")
	}

	res := &pb.LoginResponse{
		Token: token,
		User:  users.UserToGRPC(user),
	}
	return res, nil
}
