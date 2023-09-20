package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	mockRepo "github.com/learn/api-shop/internal/generated/mock"
	"github.com/learn/api-shop/internal/repo"
	"github.com/learn/api-shop/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCheckout(t *testing.T) {
	tests := []struct {
		name          string
		orderDetails  []repo.OrderDetail
		mockSetupFunc func(orderRepo *mockRepo.OrderRepository, productRepo *mockRepo.ProductRepository, promoRepo *mockRepo.PromoRepository)
		expectedResp  service.Checkout
		wantErr       bool
	}{
		{
			name: "Buying more than 3 Alexa Speakers will have a 10% discount on all Alexa speakers",
			orderDetails: []repo.OrderDetail{
				{
					ProductID: 3,
					Qty:       3,
				},
			},
			mockSetupFunc: func(orderRepo *mockRepo.OrderRepository, productRepo *mockRepo.ProductRepository, promoRepo *mockRepo.PromoRepository) {
				db, mocksql, err := sqlmock.New()
				assert.NoError(t, err)
				sqlxDb := sqlx.NewDb(db, "sqlmock")
				mocksql.ExpectBegin()
				tx, err := sqlxDb.Beginx()
				if err != nil {
					t.Fatalf("failed to begin transaction: %v", err)
				}

				orderRepo.On("BeginTx").Return(tx, nil)
				orderRepo.On("CreateOrder", tx, mock.Anything, mock.Anything).Return(int64(1), nil)
				productRepo.On("GetProductByProductID", mock.Anything, mock.Anything).Return(repo.Product{
					ProductID: 3,
					Sku:       "A304SD",
					Name:      "Alexa Speaker",
					Price:     109.500,
					Qty:       10,
				}, nil)
				productRepo.On("UpdateProductQtyByProductID", mock.Anything, mock.Anything, mock.Anything).Return(nil)

				promoRepo.On("GetPromoByProductID", mock.Anything, mock.Anything).Return(repo.Promo{
					PromoID:   3,
					PromoType: "discount",
					Reward:    10,
					MinQty:    3,
				}, nil)

				orderRepo.On("CreateOrderDetails", tx, mock.Anything, mock.Anything).Return(nil)
				mocksql.ExpectCommit()
				tx.Commit()

			},
			expectedResp: service.Checkout{
				Items:       []string{"Alexa Speaker", "Alexa Speaker", "Alexa Speaker"},
				TotalAmount: 295.65,
			},
			wantErr: false,
		},
		{
			name: "Buy 3 Google Homes for the price of 2",
			orderDetails: []repo.OrderDetail{
				{
					ProductID: 1,
					Qty:       3,
				},
			},
			mockSetupFunc: func(orderRepo *mockRepo.OrderRepository, productRepo *mockRepo.ProductRepository, promoRepo *mockRepo.PromoRepository) {
				db, mocksql, err := sqlmock.New()
				assert.NoError(t, err)
				sqlxDb := sqlx.NewDb(db, "sqlmock")
				mocksql.ExpectBegin()
				tx, err := sqlxDb.Beginx()
				if err != nil {
					t.Fatalf("failed to begin transaction: %v", err)
				}

				orderRepo.On("BeginTx").Return(tx, nil)
				orderRepo.On("CreateOrder", tx, mock.Anything, mock.Anything).Return(int64(1), nil)
				productRepo.On("GetProductByProductID", mock.Anything, mock.Anything).Return(repo.Product{
					ProductID: 1,
					Sku:       "120P90",
					Name:      "Google Home",
					Price:     49.990,
					Qty:       10,
				}, nil)
				productRepo.On("UpdateProductQtyByProductID", mock.Anything, mock.Anything, mock.Anything).Return(nil)

				promoRepo.On("GetPromoByProductID", mock.Anything, mock.Anything).Return(repo.Promo{
					PromoID:   1,
					PromoType: "product",
					Reward:    1,
					MinQty:    3,
				}, nil)
				orderRepo.On("CreateOrderDetails", tx, mock.Anything, mock.Anything).Return(nil)
				mocksql.ExpectCommit()
				tx.Commit()

			},
			expectedResp: service.Checkout{
				Items:       []string{"Google Home", "Google Home", "Google Home"},
				TotalAmount: 99.98,
			},
			wantErr: false,
		},
		{
			name: "Each sale of a MacBook Pro comes with a free Raspberry Pi B",
			orderDetails: []repo.OrderDetail{
				{
					ProductID: 2,
					Qty:       1,
				},
			},
			mockSetupFunc: func(orderRepo *mockRepo.OrderRepository, productRepo *mockRepo.ProductRepository, promoRepo *mockRepo.PromoRepository) {
				db, mocksql, err := sqlmock.New()
				assert.NoError(t, err)
				sqlxDb := sqlx.NewDb(db, "sqlmock")
				mocksql.ExpectBegin()
				tx, err := sqlxDb.Beginx()
				if err != nil {
					t.Fatalf("failed to begin transaction: %v", err)
				}

				orderRepo.On("BeginTx").Return(tx, nil)
				orderRepo.On("CreateOrder", tx, mock.Anything, mock.Anything).Return(int64(1), nil)
				productRepo.On("GetProductByProductID", context.TODO(), int64(2)).Return(repo.Product{
					ProductID: 2,
					Sku:       "43N23P",
					Name:      "MacBook Pro",
					Price:     5399.990,
					Qty:       5,
				}, nil)
				productRepo.On("UpdateProductQtyByProductID", mock.Anything, mock.Anything, mock.Anything).Return(nil)

				promoRepo.On("GetPromoByProductID", mock.Anything, mock.Anything).Return(repo.Promo{
					PromoID:   2,
					PromoType: "product",
					Reward:    4,
					MinQty:    1,
				}, nil)

				productRepo.On("GetProductByProductID", context.TODO(), int64(4)).Return(repo.Product{
					ProductID: 4,
					Sku:       "234234",
					Name:      "Raspberry Pi B",
					Price:     30.000,
					Qty:       2,
				}, nil)
				orderRepo.On("CreateOrderDetails", tx, mock.Anything, mock.Anything).Return(nil)
				mocksql.ExpectCommit()
				tx.Commit()

			},
			expectedResp: service.Checkout{
				Items:       []string{"MacBook Pro", "Raspberry Pi B"},
				TotalAmount: 5399.99,
			},
			wantErr: false,
		},
		{
			name: "the product qty is not enough to fulfill the request",
			orderDetails: []repo.OrderDetail{
				{
					ProductID: 2,
					Qty:       2,
				},
			},
			mockSetupFunc: func(orderRepo *mockRepo.OrderRepository, productRepo *mockRepo.ProductRepository, promoRepo *mockRepo.PromoRepository) {
				db, mocksql, err := sqlmock.New()
				assert.NoError(t, err)
				sqlxDb := sqlx.NewDb(db, "sqlmock")
				mocksql.ExpectBegin()
				tx, err := sqlxDb.Beginx()
				if err != nil {
					t.Fatalf("failed to begin transaction: %v", err)
				}

				orderRepo.On("BeginTx").Return(tx, nil)
				orderRepo.On("CreateOrder", tx, mock.Anything, mock.Anything).Return(int64(1), nil)
				promoRepo.On("GetPromoByProductID", mock.Anything, mock.Anything).Return(repo.Promo{
					PromoID:   2,
					PromoType: "discount",
					Reward:    4,
					MinQty:    1,
				}, nil)
				productRepo.On("GetProductByProductID", context.TODO(), int64(2)).Return(repo.Product{
					ProductID: 2,
					Sku:       "43N23P",
					Name:      "MacBook Pro",
					Price:     5399.990,
					Qty:       1,
				}, nil)
				mocksql.ExpectCommit()
				tx.Commit()
			},
			expectedResp: service.Checkout{},
			wantErr:      true,
		},
		{
			name: "error while begin transaction",
			orderDetails: []repo.OrderDetail{
				{
					ProductID: 2,
					Qty:       1,
				},
			},
			mockSetupFunc: func(orderRepo *mockRepo.OrderRepository, productRepo *mockRepo.ProductRepository, promoRepo *mockRepo.PromoRepository) {
				db, mocksql, err := sqlmock.New()
				assert.NoError(t, err)
				sqlxDb := sqlx.NewDb(db, "sqlmock")
				mocksql.ExpectBegin()
				tx, err := sqlxDb.Beginx()
				if err != nil {
					t.Fatalf("failed to begin transaction: %v", err)
				}

				orderRepo.On("BeginTx").Return(tx, errors.New("error"))
			},
			expectedResp: service.Checkout{},
			wantErr:      true,
		},
		{
			name: "error while create order",
			orderDetails: []repo.OrderDetail{
				{
					ProductID: 2,
					Qty:       1,
				},
			},
			mockSetupFunc: func(orderRepo *mockRepo.OrderRepository, productRepo *mockRepo.ProductRepository, promoRepo *mockRepo.PromoRepository) {
				db, mocksql, err := sqlmock.New()
				assert.NoError(t, err)
				sqlxDb := sqlx.NewDb(db, "sqlmock")
				mocksql.ExpectBegin()
				tx, err := sqlxDb.Beginx()
				if err != nil {
					t.Fatalf("failed to begin transaction: %v", err)
				}

				orderRepo.On("BeginTx").Return(tx, nil)
				orderRepo.On("CreateOrder", tx, mock.Anything, mock.Anything).Return(int64(1), errors.New("error"))

			},
			expectedResp: service.Checkout{},
			wantErr:      true,
		},
		{
			name: "error while create order detail",
			orderDetails: []repo.OrderDetail{
				{
					ProductID: 2,
					Qty:       1,
				},
			},
			mockSetupFunc: func(orderRepo *mockRepo.OrderRepository, productRepo *mockRepo.ProductRepository, promoRepo *mockRepo.PromoRepository) {
				db, mocksql, err := sqlmock.New()
				assert.NoError(t, err)
				sqlxDb := sqlx.NewDb(db, "sqlmock")
				mocksql.ExpectBegin()
				tx, err := sqlxDb.Beginx()
				if err != nil {
					t.Fatalf("failed to begin transaction: %v", err)
				}

				orderRepo.On("BeginTx").Return(tx, nil)
				orderRepo.On("CreateOrder", tx, mock.Anything, mock.Anything).Return(int64(1), nil)
				productRepo.On("GetProductByProductID", context.TODO(), int64(2)).Return(repo.Product{
					ProductID: 2,
					Sku:       "43N23P",
					Name:      "MacBook Pro",
					Price:     5399.990,
					Qty:       5,
				}, nil)
				productRepo.On("UpdateProductQtyByProductID", mock.Anything, mock.Anything, mock.Anything).Return(nil)

				promoRepo.On("GetPromoByProductID", mock.Anything, mock.Anything).Return(repo.Promo{
					PromoID:   2,
					PromoType: "product",
					Reward:    4,
					MinQty:    1,
				}, nil)

				productRepo.On("GetProductByProductID", context.TODO(), int64(4)).Return(repo.Product{
					ProductID: 4,
					Sku:       "234234",
					Name:      "Raspberry Pi B",
					Price:     30.000,
					Qty:       2,
				}, nil)
				orderRepo.On("CreateOrderDetails", tx, mock.Anything, mock.Anything).Return(errors.New("error"))
			},
			expectedResp: service.Checkout{
				Items:       []string{"MacBook Pro", "Raspberry Pi B"},
				TotalAmount: 5399.99,
			},
			wantErr: true,
		},
		{
			name: "error while get promo",
			orderDetails: []repo.OrderDetail{
				{
					ProductID: 2,
					Qty:       1,
				},
			},
			mockSetupFunc: func(orderRepo *mockRepo.OrderRepository, productRepo *mockRepo.ProductRepository, promoRepo *mockRepo.PromoRepository) {
				db, mocksql, err := sqlmock.New()
				assert.NoError(t, err)
				sqlxDb := sqlx.NewDb(db, "sqlmock")
				mocksql.ExpectBegin()
				tx, err := sqlxDb.Beginx()
				if err != nil {
					t.Fatalf("failed to begin transaction: %v", err)
				}

				orderRepo.On("BeginTx").Return(tx, nil)
				orderRepo.On("CreateOrder", tx, mock.Anything, mock.Anything).Return(int64(1), nil)
				promoRepo.On("GetPromoByProductID", mock.Anything, mock.Anything).Return(repo.Promo{
					PromoID:   2,
					PromoType: "product",
					Reward:    4,
					MinQty:    1,
				}, errors.New("error"))
			},
			expectedResp: service.Checkout{},
			wantErr:      true,
		},
		{
			name: "error while GetProductByProductID",
			orderDetails: []repo.OrderDetail{
				{
					ProductID: 2,
					Qty:       2,
				},
			},
			mockSetupFunc: func(orderRepo *mockRepo.OrderRepository, productRepo *mockRepo.ProductRepository, promoRepo *mockRepo.PromoRepository) {
				db, mocksql, err := sqlmock.New()
				assert.NoError(t, err)
				sqlxDb := sqlx.NewDb(db, "sqlmock")
				mocksql.ExpectBegin()
				tx, err := sqlxDb.Beginx()
				if err != nil {
					t.Fatalf("failed to begin transaction: %v", err)
				}

				orderRepo.On("BeginTx").Return(tx, nil)
				orderRepo.On("CreateOrder", tx, mock.Anything, mock.Anything).Return(int64(1), nil)
				promoRepo.On("GetPromoByProductID", mock.Anything, mock.Anything).Return(repo.Promo{
					PromoID:   2,
					PromoType: "discount",
					Reward:    4,
					MinQty:    1,
				}, nil)
				productRepo.On("GetProductByProductID", context.TODO(), int64(2)).Return(repo.Product{
					ProductID: 2,
					Sku:       "43N23P",
					Name:      "MacBook Pro",
					Price:     5399.990,
					Qty:       1,
				}, errors.New("error"))

			},
			expectedResp: service.Checkout{},
			wantErr:      true,
		},
		{
			name: "error while update product qty",
			orderDetails: []repo.OrderDetail{
				{
					ProductID: 2,
					Qty:       1,
				},
			},
			mockSetupFunc: func(orderRepo *mockRepo.OrderRepository, productRepo *mockRepo.ProductRepository, promoRepo *mockRepo.PromoRepository) {
				db, mocksql, err := sqlmock.New()
				assert.NoError(t, err)
				sqlxDb := sqlx.NewDb(db, "sqlmock")
				mocksql.ExpectBegin()
				tx, err := sqlxDb.Beginx()
				if err != nil {
					t.Fatalf("failed to begin transaction: %v", err)
				}

				orderRepo.On("BeginTx").Return(tx, nil)
				orderRepo.On("CreateOrder", tx, mock.Anything, mock.Anything).Return(int64(1), nil)
				promoRepo.On("GetPromoByProductID", mock.Anything, mock.Anything).Return(repo.Promo{
					PromoID:   2,
					PromoType: "product",
					Reward:    4,
					MinQty:    1,
				}, nil)
				productRepo.On("GetProductByProductID", context.TODO(), int64(2)).Return(repo.Product{
					ProductID: 2,
					Sku:       "43N23P",
					Name:      "MacBook Pro",
					Price:     5399.990,
					Qty:       5,
				}, nil)
				productRepo.On("GetProductByProductID", context.TODO(), int64(4)).Return(repo.Product{
					ProductID: 4,
					Sku:       "234234",
					Name:      "Raspberry Pi B",
					Price:     30.000,
					Qty:       2,
				}, nil)
				productRepo.On("UpdateProductQtyByProductID", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error"))

			},
			expectedResp: service.Checkout{
				Items:       []string{"MacBook Pro", "Raspberry Pi B"},
				TotalAmount: 5399.99,
			},
			wantErr: true,
		},
		{
			name: "get free product detail",
			orderDetails: []repo.OrderDetail{
				{
					ProductID: 2,
					Qty:       1,
				},
			},
			mockSetupFunc: func(orderRepo *mockRepo.OrderRepository, productRepo *mockRepo.ProductRepository, promoRepo *mockRepo.PromoRepository) {
				db, mocksql, err := sqlmock.New()
				assert.NoError(t, err)
				sqlxDb := sqlx.NewDb(db, "sqlmock")
				mocksql.ExpectBegin()
				tx, err := sqlxDb.Beginx()
				if err != nil {
					t.Fatalf("failed to begin transaction: %v", err)
				}

				orderRepo.On("BeginTx").Return(tx, nil)
				orderRepo.On("CreateOrder", tx, mock.Anything, mock.Anything).Return(int64(1), nil)
				promoRepo.On("GetPromoByProductID", mock.Anything, mock.Anything).Return(repo.Promo{
					PromoID:   2,
					PromoType: "product",
					Reward:    4,
					MinQty:    1,
				}, nil)
				productRepo.On("GetProductByProductID", context.TODO(), int64(2)).Return(repo.Product{
					ProductID: 2,
					Sku:       "43N23P",
					Name:      "MacBook Pro",
					Price:     5399.990,
					Qty:       5,
				}, nil)
				productRepo.On("GetProductByProductID", context.TODO(), int64(4)).Return(repo.Product{
					ProductID: 4,
					Sku:       "234234",
					Name:      "Raspberry Pi B",
					Price:     30.000,
					Qty:       2,
				}, errors.New("error"))

			},
			expectedResp: service.Checkout{
				Items:       []string{"MacBook Pro", "Raspberry Pi B"},
				TotalAmount: 5399.99,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orderRepo := new(mockRepo.OrderRepository)
			productRepo := new(mockRepo.ProductRepository)
			promoRepo := new(mockRepo.PromoRepository)

			if tt.mockSetupFunc != nil {
				tt.mockSetupFunc(orderRepo, productRepo, promoRepo)
			}

			checkoutUsecase := service.NewCheckoutUsecase(service.CheckoutUsecaseImpl{
				OrderRepo:   orderRepo,
				ProductRepo: productRepo,
				PromoRepo:   promoRepo,
			})

			res, err := checkoutUsecase.Checkout(context.Background(), tt.orderDetails)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResp, res)
			}

			orderRepo.AssertExpectations(t)
			productRepo.AssertExpectations(t)
			promoRepo.AssertExpectations(t)
		})
	}
}
