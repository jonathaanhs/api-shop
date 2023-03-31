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

func TestGetPromoByProductID(t *testing.T) {
	testCases := []struct {
		name          string
		productID     int64
		expectedPromo repo.Promo
		expectedErr   error
		mockFunc      func(mock sqlmock.Sqlmock)
	}{
		{
			name:      "success",
			productID: 1,
			expectedPromo: repo.Promo{
				PromoID:   1,
				ProductID: 1,
				PromoType: "type1",
				Reward:    1.23,
				MinQty:    1,
			},
			expectedErr: nil,
			mockFunc: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"promo_id", "product_id", "promo_type", "reward", "min_qty"}).
					AddRow(1, 1, "type1", 1.23, 1)
				mock.ExpectQuery("select promo_id, product_id, promo_type, reward, min_qty from promos where product_id = \\$1").
					WithArgs(1).WillReturnRows(rows)
			},
		},
		{
			name:          "database error",
			productID:     1,
			expectedPromo: repo.Promo{},
			expectedErr:   errors.New("database error"),
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("select promo_id, product_id, promo_type, reward, min_qty from promos where product_id = \\$1").
					WithArgs(1).WillReturnError(errors.New("database error"))
			},
		},
		{
			name:          "error scanning promo rows",
			productID:     1,
			expectedPromo: repo.Promo{},
			expectedErr:   errors.New("sql: Scan error on column index 3, name \"reward\": converting driver.Value type string (\"not a float\") to a float64: invalid syntax"),
			mockFunc: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"promo_id", "product_id", "promo_type", "reward", "min_qty"}).
					AddRow(1, 1, "discount", "not a float", 2)
				mock.ExpectQuery("select promo_id, product_id, promo_type, reward, min_qty from promos where product_id = \\$1").
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

			repoImpl := repo.PromoRepoImpl{DB: sqlx.NewDb(db, "sqlmock")}
			repo := repo.NewPromoRepository(repoImpl)

			promo, err := repo.GetPromoByProductID(context.Background(), tc.productID)
			if err != nil {
				if err.Error() != tc.expectedErr.Error() {
					t.Errorf("expected error '%s', but got '%s'", tc.expectedErr, err)
				}
			} else {
				assert.Equal(t, tc.expectedPromo, promo)
			}
		})
	}
}

func TestPromoRepoImpl_GetAllPromo(t *testing.T) {
	testCases := []struct {
		name          string
		expectedPromo []repo.Promo
		expectedErr   error
		mockFunc      func(mock sqlmock.Sqlmock)
	}{
		{
			name:          "successfully get all promos",
			expectedPromo: []repo.Promo{{PromoID: 1, ProductID: 1, PromoType: "discount", Reward: 10.0, MinQty: 2}, {PromoID: 2, ProductID: 2, PromoType: "free gift", Reward: 0.0, MinQty: 5}},
			expectedErr:   nil,
			mockFunc: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"promo_id", "product_id", "promo_type", "reward", "min_qty"}).
					AddRow(1, 1, "discount", 10.0, 2).
					AddRow(2, 2, "free gift", 0.0, 5)
				mock.ExpectQuery("select promo_id, product_id, promo_type, reward, min_qty from promos order by promo_id asc").
					WillReturnRows(rows)
			},
		},
		{
			name:          "database error",
			expectedPromo: []repo.Promo{},
			expectedErr:   errors.New("database error"),
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("select promo_id, product_id, promo_type, reward, min_qty from promos order by promo_id asc").
					WillReturnError(errors.New("database error"))
			},
		},
		{
			name:          "error scanning promo rows",
			expectedPromo: []repo.Promo{},
			expectedErr:   errors.New("sql: Scan error on column index 3, name \"reward\": converting driver.Value type string (\"not a float\") to a float64: invalid syntax"),
			mockFunc: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"promo_id", "product_id", "promo_type", "reward", "min_qty"}).
					AddRow(1, 1, "discount", "not a float", 2)
				mock.ExpectQuery("select promo_id, product_id, promo_type, reward, min_qty from promos order by promo_id asc").
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

			repoImpl := repo.PromoRepoImpl{DB: sqlx.NewDb(db, "sqlmock")}
			repo := repo.NewPromoRepository(repoImpl)

			promo, err := repo.GetAllPromo(context.Background())
			if err != nil {
				if err.Error() != tc.expectedErr.Error() {
					t.Errorf("expected error '%s', but got '%s'", tc.expectedErr, err)
				}
			} else {
				assert.Equal(t, tc.expectedPromo, promo)
			}
		})
	}
}
