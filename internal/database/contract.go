package database

import "context"

//go:generate mockery --all --output=./mocks --case=underscore --exported --with-expecter

type DbManager interface {
	CreateUser(ctx context.Context, amount int64) (int64, error)
	Refill(ctx context.Context, userId, amount int64) (bool, error)
	GetBalance(ctx context.Context, userId int64) (int64, error)
	Withdrawal(ctx context.Context, clientId, serviceId, orderId, Amount int64) (bool, error)
	ProcessWithdrawal(ctx context.Context, clientId, serviceId, orderId, Amount int64) (bool, error)
	CancelWithdrawal(ctx context.Context, clientId, serviceId, orderId, amount int64) (bool, error)
}
