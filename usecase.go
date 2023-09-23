package main

import (
	"context"
	"errors"
)

type UserUsecase interface {
	FetchByID(ctx context.Context, tenantID, userID string) (*User, error)
	Update(ctx context.Context, tenantID string, u User) error
}

func NewUserUsecase(tx Transaction, ur UserRepository) UserUsecase {
	return &userInteractor{
		tx: tx,
		ur: ur,
	}
}

type userInteractor struct {
	tx Transaction
	ur UserRepository
}

func (i *userInteractor) FetchByID(ctx context.Context, tenantID, userID string) (*User, error) {
	user, err := i.ur.FetchByID(ctx, tenantID, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (i *userInteractor) Update(ctx context.Context, tenantID string, u User) error {
	i.tx.DoInTx(ctx, tenantID, func(ctx context.Context) (any, error) {
		rowsAffected, err := i.ur.Update(ctx, tenantID, u)
		if err != nil {
			return nil, err
		}
		if rowsAffected == 0 {
			return nil, errors.New("no rows affected")
		}
		return nil, nil
	})
	return nil
}

type UserRepository interface {
	FetchByID(ctx context.Context, tenantID, userID string) (*User, error)
	Update(ctx context.Context, tenantID string, u User) (rowsAffected int64, err error)
}
