package database

import (
	"context"
	"errors"
	pgx "github.com/jackc/pgx/v4"
	"time"
)

const (
	orderCreated = iota
	orderProcess
	orderCanceled
)

type dataBase struct {
	conn *pgx.Conn
}

func New(ctx context.Context) (*dataBase, error) {

	connection, err := pgx.Connect(ctx, "postgres://postgres:postgres@database:5432/master")
	if err != nil {
		return nil, err
	}
	err = connection.Ping(ctx)
	if err != nil {
		return nil, err
	}
	return &dataBase{conn: connection}, nil
}

func (d *dataBase) CreateUser(ctx context.Context, amount int64) (int64, error) {
	var clientId int64

	if amount < 0 {
		return 0, errors.New("invalid data")
	}
	var query = "INSERT INTO accounts (balance) VALUES ($1) returning client_id"
	err := d.conn.QueryRow(
		ctx,
		query,
		amount,
	).Scan(&clientId)
	if err != nil {
		return 0, err
	}
	return clientId, nil
}

func (d *dataBase) Refill(ctx context.Context, clientId, amount int64) (bool, error) {
	if amount < 0 {
		return false, errors.New("negative balance")
	}
	opts := pgx.TxOptions{
		IsoLevel: "serializable",
	}
	err := d.conn.BeginTxFunc(ctx, opts,
		func(tx pgx.Tx) error {
			var query = "UPDATE accounts SET balance = balance + $1 WHERE client_id = $2"
			_, err := tx.Exec(ctx, query, amount, clientId)
			if err != nil {
				return err
			}
			return nil
		},
	)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (d *dataBase) GetBalance(ctx context.Context, userId int64) (int64, error) {
	var balance int64
	var query = "SELECT balance FROM accounts WHERE client_id = $1"
	row := d.conn.QueryRow(ctx, query, userId)
	err := row.Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

func (d *dataBase) Withdrawal(ctx context.Context, clientId, serviceId, orderId, amount int64) (bool, error) {
	opts := pgx.TxOptions{
		IsoLevel: "serializable",
	}
	err := d.conn.BeginTxFunc(ctx, opts,
		func(tx pgx.Tx) error {
			var balance int64
			var queryGetBalance = "SELECT balance FROM accounts WHERE client_id = $1"

			row := tx.QueryRow(context.Background(), queryGetBalance, clientId)
			err := row.Scan(&balance)
			if err != nil {
				return err
			}
			if balance < amount || balance < 0 {
				return errors.New("underfunded account")
			}

			var queryWithdrawal = "UPDATE accounts SET balance = balance - $1 WHERE client_id = $2"
			ct, err := tx.Exec(ctx, queryWithdrawal, amount, clientId)
			if err != nil {
				return err
			}
			if ct.RowsAffected() == 0 {
				return errors.New("not updated")
			}

			var queryOrder = "INSERT INTO orders (client_id, order_id, service_id, amount, create_at, status) VALUES($1, $2, $3, $4, $5, $6)"
			_, err = tx.Exec(ctx, queryOrder, clientId, orderId, serviceId, amount, time.Now(), orderCreated)
			if err != nil {
				return err
			}
			return nil
		},
	)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (d *dataBase) ProcessWithdrawal(ctx context.Context, clientId, serviceId, orderId, amount int64) (bool, error) {
	opts := pgx.TxOptions{
		IsoLevel: "serializable",
	}
	err := d.conn.BeginTxFunc(ctx, opts,
		func(tx pgx.Tx) error {
			var query = "UPDATE orders SET(done_at, status) = ($1, $2) WHERE client_id = $3 and order_id = $4 and service_id = $5 and amount = $6"
			ct, err := tx.Exec(ctx, query, time.Now(), orderProcess, clientId, orderId, serviceId, amount)
			if err != nil {
				return err
			}
			if ct.RowsAffected() == 0 {
				return errors.New("not updated")
			}
			return nil
		},
	)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (d *dataBase) CancelWithdrawal(ctx context.Context, clientId, serviceId, orderId, amount int64) (bool, error) {
	opts := pgx.TxOptions{
		IsoLevel: "serializable",
	}
	err := d.conn.BeginTxFunc(ctx, opts,
		func(tx pgx.Tx) error {
			var query = "UPDATE orders SET(done_at, status) = ($1, $2) WHERE order_id = $3 and service_id = $4"
			ct, err := tx.Exec(ctx, query, time.Now(), orderCanceled, orderId, serviceId)
			if err != nil {
				return err
			}
			if ct.RowsAffected() == 0 {
				return errors.New("not updated")
			}
			var queryBackAmount = "UPDATE accounts SET balance = balance + $1 WHERE client_id = $2"
			ct, err = tx.Exec(ctx, queryBackAmount, amount, clientId)
			if err != nil {
				return err
			}
			if ct.RowsAffected() == 0 {
				return errors.New("not updated")
			}
			return nil
		},
	)
	if err != nil {
		return false, err
	}
	return true, nil
}
