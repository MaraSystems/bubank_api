package gapi

import (
	"context"
	"errors"

	db "github.com/MaraSystems/bubank_api/db/sqlc"
	"github.com/MaraSystems/bubank_api/domains/users"
	pb "github.com/MaraSystems/bubank_api/pb"
	"github.com/MaraSystems/bubank_api/utils"
	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	hashedPassword, err := utils.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to hash password")
	}

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		Email:          req.GetEmail(),
		FullName:       req.GetFullName(),
		HashedPassword: hashedPassword,
	}

	user, err := s.Store.CreateUser(ctx, arg)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			switch pgError.Code {
			case "23505":
				return nil, status.Error(codes.AlreadyExists, "username or email is already in use")
			}
		}

		return nil, status.Error(codes.Internal, "failed to save user")
	}

	res := &pb.RegisterResponse{
		User: users.UserToGRPC(user),
	}
	return res, nil
}
