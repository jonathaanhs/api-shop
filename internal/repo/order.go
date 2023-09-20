//go:generate mockery --dir=$PROJECT_DIR/internal/repo  --name=OrderRepository --filename=$GOFILE --output=$PROJECT_DIR/internal/generated/mock --outpkg=mock
package repo

import (
	"context"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/learn/api-shop/pkg/sqlkit"
	"go.uber.org/dig"
)

type (
	Order struct {
		OrderID int64     `json:"order_id" db:"order_id"`
		Date    time.Time `json:"date" db:"date"`
		Total   float64   `json:"total" db:"total"`
	}

	OrderDetail struct {
		OrderDetailID int64   `json:"order_detail_id" db:"order_detail_id"`
		OrderID       int64   `json:"order_id" db:"order_id"`
		ProductID     int64   `json:"product_id" db:"product_id"`
		PromoID       int64   `json:"promo_id" db:"promo_id"`
		Price         float64 `json:"price" db:"price"`
		Qty           int64   `json:"qty" db:"qty"`
	}

	OrderRepository interface {
		CreateOrder(tx *sqlx.Tx, ctx context.Context, form Order) (orderID int64, err error)
		CreateOrderDetails(tx *sqlx.Tx, ctx context.Context, form []OrderDetail) (err error)
		BeginTx() (tx *sqlx.Tx, err error)
		RollbackTx(tx *sqlx.Tx) (err error)
		CommitTx(tx *sqlx.Tx) (err error)
	}

	OrderRepoImpl struct {
		dig.In
		*sqlx.DB
	}
)

func NewOrderRepository(impl OrderRepoImpl) OrderRepository {
	return &impl
}

func (r *OrderRepoImpl) CreateOrder(tx *sqlx.Tx, ctx context.Context, form Order) (orderID int64, err error) {
	err = tx.QueryRowxContext(ctx, "insert into orders(date, total) values($1, $2) RETURNING order_id", form.Date, form.Total).Scan(&orderID)
	if err != nil {
		return orderID, err
	}

	return orderID, nil
}

func (r *OrderRepoImpl) CreateOrderDetails(tx *sqlx.Tx, ctx context.Context, form []OrderDetail) (err error) {
	sqlInsert := "insert into order_details(order_id, product_id, promo_id, price, qty) values"
	rowSQL := "(?, ?, ?, ?, ?)"

	vals := []interface{}{}
	var inserts []string

	for _, val := range form {
		vals = append(vals, val.OrderID, val.ProductID, val.PromoID, val.Price, val.Qty)
		inserts = append(inserts, rowSQL)
	}

	sqlInsert = sqlInsert + strings.Join(inserts, ",")
	sqlInsert = sqlkit.ReplaceSQL(sqlInsert, "?")

	stmt, err := tx.PrepareContext(ctx, sqlInsert)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, vals...)
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderRepoImpl) BeginTx() (tx *sqlx.Tx, err error) {
	return r.DB.Beginx()
}

func (r *OrderRepoImpl) RollbackTx(tx *sqlx.Tx) (err error) {
	return tx.Rollback()
}

func (r *OrderRepoImpl) CommitTx(tx *sqlx.Tx) (err error) {
	return tx.Commit()
}
