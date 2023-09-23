package main

import "context"

type UserUsecase interface {
	FetchByID(ctx context.Context, userID string) (User, error)
	Update(ctx context.Context, u User) err
}

type UserInteractor struct {
	ur UserRepository
}

func NewUserUsecase(ur UserRepository) {
	return &UserInteractor{
		ur: ur,
	}
}

type UserRepository interface {
	FetchByID(ctx context.Context, userID string) (User, error)
	Update(ctx context.Context, u User) (rowsAffected int64, error)
}