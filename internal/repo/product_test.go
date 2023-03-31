package repo_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/learn/api-shop/internal/repo"
	"github.com/stretchr/testify/assert"
)

func TestProductRepoImpl_GetProductByProductID(t *testing.T) {
	testCases := []struct {
		name         string
		productID    int64
		expectedResp repo.Product
		expectedErr  error
		mockFunc     func(mock sqlmock.Sqlmock)
	}{
		{
			name:      "success",
			productID: 1,
			expectedResp: repo.Product{
				ProductID: 1,
				Sku:       "abc",
				Name:      "sepatu",
				Price:     2.2,
				Qty:       10,
			},
			expectedErr: nil,
			mockFunc: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"product_id", "sku", "name", "price", "qty"}).
					AddRow(1, "abc", "sepatu", 2.2, 10)
				mock.ExpectQuery("select product_id, sku, name, price, qty from products where product_id = \\$1").
					WillReturnRows(rows)
			},
		},
		{
			name:         "database error",
			productID:    1,
			expectedResp: repo.Product{},
			expectedErr:  errors.New("database error"),
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("select product_id, sku, name, price, qty from products where product_id = \\$1").
					WithArgs(1).WillReturnError(errors.New("database error"))
			},
		},
		{
			name:         "error scanning product rows",
			productID:    1,
			expectedResp: repo.Product{},
			expectedErr:  errors.New("sql: Scan error on column index 3, name \"price\": converting driver.Value type string (\"not a float\") to a float64: invalid syntax"),
			mockFunc: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"product_id", "sku", "name", "price", "qty"}).
					AddRow(1, "abc", "sepatu", "not a float", 10)
				mock.ExpectQuery("select product_id, sku, name, price, qty from products where product_id = \\$1").
					WithArgs(1).WillReturnRows(rows).WillReturnError(nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			mockFunc := tc.mockFunc
			if mockFunc == nil {
				mockFunc = func(mock sqlmock.Sqlmock) {}
			}
			mockFunc(mock)

			repoImpl := repo.ProductRepoImpl{DB: sqlx.NewDb(db, "sqlmock")}
			repo := repo.NewProductRepository(repoImpl)

			resp, err := repo.GetProductByProductID(context.Background(), tc.productID)
			if err != nil {
				if err.Error() != tc.expectedErr.Error() {
					t.Errorf("expected error '%s', but got '%s'", tc.expectedErr, err)
				}
			} else {
				assert.Equal(t, tc.expectedResp, resp)
			}
		})
	}
}

func TestProductRepoImpl_GetAllProduct(t *testing.T) {
	testCases := []struct {
		name         string
		expectedResp []repo.Product
		expectedErr  error
		mockFunc     func(mock sqlmock.Sqlmock)
	}{
		{
			name: "successfully get all products",
			expectedResp: []repo.Product{
				{
					ProductID: 1,
					Sku:       "abc",
					Name:      "sepatu",
					Price:     2.2,
					Qty:       10,
				},
				{
					ProductID: 2,
					Sku:       "cda",
					Name:      "jam",
					Price:     20,
					Qty:       20,
				},
			},
			expectedErr: nil,
			mockFunc: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"product_id", "sku", "name", "price", "qty"}).
					AddRow(1, "abc", "sepatu", 2.2, 10).
					AddRow(2, "cda", "jam", 20, 20)
				mock.ExpectQuery("select product_id, sku, name, price, qty from products order by product_id asc").
					WillReturnRows(rows)
			},
		},
		{
			name:         "database error",
			expectedResp: []repo.Product{},
			expectedErr:  errors.New("database error"),
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("select product_id, sku, name, price, qty from products order by product_id asc").
					WillReturnError(errors.New("database error"))
			},
		},
		{
			name:         "error scanning product rows",
			expectedResp: []repo.Product{},
			expectedErr:  errors.New("sql: Scan error on column index 3, name \"price\": converting driver.Value type string (\"not a float\") to a float64: invalid syntax"),
			mockFunc: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"product_id", "sku", "name", "price", "qty"}).
					AddRow(1, "abc", "sepatu", "not a float", 10)
				mock.ExpectQuery("select product_id, sku, name, price, qty from products order by product_id asc").
					WillReturnRows(rows).WillReturnError(nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			mockFunc := tc.mockFunc
			if mockFunc == nil {
				mockFunc = func(mock sqlmock.Sqlmock) {}
			}
			mockFunc(mock)

			repoImpl := repo.ProductRepoImpl{DB: sqlx.NewDb(db, "sqlmock")}
			repo := repo.NewProductRepository(repoImpl)

			resp, err := repo.GetAllProduct(context.Background())
			if err != nil {
				if err.Error() != tc.expectedErr.Error() {
					t.Errorf("expected error '%s', but got '%s'", tc.expectedErr, err)
				}
			} else {
				assert.Equal(t, tc.expectedResp, resp)
			}
		})
	}
}

func TestProductRepoImpl_UpdateProductQtyByProductID(t *testing.T) {
	type args struct {
		form repo.Product
	}

	testCases := []struct {
		name           string
		mockSQL        func(sqlmock.Sqlmock)
		args           args
		expectedResult error
	}{
		{
			name: "success",
			mockSQL: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE products SET qty = qty - \\$1 WHERE product_id = \\$2").
					WithArgs(10, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			args: args{
				form: repo.Product{
					ProductID: 1,
					Qty:       10,
				},
			},
			expectedResult: nil,
		},
		{
			name: "db error",
			mockSQL: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE products SET qty = qty - \\$1 WHERE product_id = \\$2").
					WithArgs(10, 1).
					WillReturnError(errors.New("db error"))
			},
			args: args{
				form: repo.Product{
					ProductID: 1,
					Qty:       10,
				},
			},
			expectedResult: errors.New("db error"),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create mock: %v", err)
			}
			defer db.Close()

			sqlxDB := sqlx.NewDb(db, "sqlmock")

			tt.mockSQL(mock)

			repo := repo.NewProductRepository(repo.ProductRepoImpl{
				DB: sqlxDB,
			})

			tx, err := sqlxDB.Beginx()
			if err != nil {
				t.Fatalf("failed to begin transaction: %v", err)
			}

			err = repo.UpdateProductQtyByProductID(tx, context.Background(), tt.args.form)
			if !assert.Equal(t, tt.expectedResult, err) {
				t.Errorf("Unexpected error result. Got %v, want %v", err, tt.expectedResult)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled SQL expectations: %v", err)
			}
		})
	}
}
