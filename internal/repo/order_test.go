package repo_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/learn/api-shop/internal/repo"
	"github.com/stretchr/testify/assert"
)

func TestOrderRepoImpl_CreateOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	tests := []struct {
		name        string
		args        repo.Order
		wantOrderID int64
		wantErr     bool
	}{
		{
			name: "successful insert",
			args: repo.Order{
				Date:  time.Now(),
				Total: 1000.0,
			},
			wantOrderID: 1,
			wantErr:     false,
		},
		{
			name: "insert error",
			args: repo.Order{
				Date:  time.Now(),
				Total: 1000.0,
			},
			wantOrderID: 0,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// set up expectations
			mock.ExpectBegin()
			if tt.wantErr {
				mock.ExpectQuery("insert into orders").WillReturnError(errors.New("insert error"))
			} else {
				mock.ExpectQuery("insert into orders").WithArgs(tt.args.Date, tt.args.Total).WillReturnRows(sqlmock.NewRows([]string{"order_id"}).AddRow(tt.wantOrderID))
			}

			tx, err := sqlxDB.Beginx()
			if err != nil {
				t.Fatalf("failed to begin transaction: %v", err)
			}

			repoImpl := repo.NewOrderRepository(repo.OrderRepoImpl{DB: sqlxDB})
			gotOrderID, err := repoImpl.CreateOrder(tx, context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderRepoImpl.CreateOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOrderID != tt.wantOrderID {
				t.Errorf("OrderRepoImpl.CreateOrder() = %v, want %v", gotOrderID, tt.wantOrderID)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestOrderRepoImpl_CreateOrderDetails(t *testing.T) {
	orderDetail := repo.OrderDetail{
		OrderID:   1,
		ProductID: 1,
		PromoID:   0,
		Price:     10.0,
		Qty:       1,
	}

	tests := []struct {
		name    string
		form    []repo.OrderDetail
		mockFn  func(mock sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "should create order details successfully",
			form: []repo.OrderDetail{orderDetail},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectPrepare("insert into order_details").
					ExpectExec().
					WithArgs(orderDetail.OrderID, orderDetail.ProductID, orderDetail.PromoID, orderDetail.Price, orderDetail.Qty).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "should return error when prepare fails",
			form: []repo.OrderDetail{orderDetail},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectPrepare("insert into order_details").WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
		{
			name: "should return error when exec fails",
			form: []repo.OrderDetail{orderDetail},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectPrepare("insert into order_details").ExpectExec().WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("error creating mock DB: %v", err)
			}
			defer db.Close()

			sqlxDB := sqlx.NewDb(db, "sqlmock")

			orderRepo := repo.NewOrderRepository(repo.OrderRepoImpl{DB: sqlxDB})

			tt.mockFn(mock)

			tx, err := sqlxDB.Beginx()
			if err != nil {
				t.Fatalf("failed to begin transaction: %v", err)
			}

			err = orderRepo.CreateOrderDetails(tx, context.Background(), tt.form)

			assert.Equal(t, tt.wantErr, err != nil, "error does not match the expectation")
			assert.NoError(t, mock.ExpectationsWereMet(), "all expectations were not met")
		})
	}
}

func TestOrderRepoImpl_BeginTx(t *testing.T) {
	type mockExpectations struct {
		query    string
		args     []driver.Value
		rowCount int
		err      error
	}

	tests := []struct {
		name          string
		mockExpect    mockExpectations
		expectedError error
	}{
		{
			name: "successful begin transaction",
			mockExpect: mockExpectations{
				query: "BEGIN",
				args:  []driver.Value{},
				err:   nil,
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()

			mock.ExpectBegin()
			mock.ExpectExec(tt.mockExpect.query).WithArgs(tt.mockExpect.args...).WillReturnResult(sqlmock.NewResult(0, int64(tt.mockExpect.rowCount))).WillReturnError(tt.mockExpect.err)

			orderRepo := repo.NewOrderRepository(repo.OrderRepoImpl{DB: sqlx.NewDb(db, "sqlmock")})

			tx, err := orderRepo.BeginTx()
			assert.Equal(t, tt.expectedError, err)
			if err == nil {
				tx.Rollback()
			}
		})
	}
}
