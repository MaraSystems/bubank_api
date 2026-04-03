package users

import (
	db "github.com/MaraSystems/bubank_api/db/sqlc"
	"github.com/MaraSystems/bubank_api/models"
	"github.com/MaraSystems/bubank_api/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func UserToHTTP(user db.User) models.UserResponse {
	return models.UserResponse{
		Username:  user.Username,
		Email:     user.Email,
		FullName:  user.FullName,
		CreatedAt: user.CreatedAt.Time,
	}
}

func UserToGRPC(user db.User) *pb.User {
	return &pb.User{
		Username:  user.Username,
		Email:     user.Email,
		FullName:  user.FullName,
		CreatedAt: timestamppb.New(user.CreatedAt.Time),
	}
}
