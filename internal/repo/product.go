package repo

import (
	"context"

	"github.com/jmoiron/sqlx"
	"go.uber.org/dig"
)

type (
	Product struct {
		ProductID int64   `json:"product_id" db:"product_id"`
		Sku       string  `json:"sku" db:"sku"`
		Name      string  `json:"name" db:"name"`
		Price     float64 `json:"price" db:"price"`
		Qty       int64   `json:"qty" db:"qty"`
	}

	ProductRepository interface {
		GetProductByProductID(ctx context.Context, id int64) (res Product, err error)
		GetAllProduct(ctx context.Context) (res []Product, err error)
		UpdateProductQtyByProductID(tx *sqlx.Tx, ctx context.Context, form Product) (err error)
	}

	ProductRepoImpl struct {
		dig.In
		*sqlx.DB
	}
)

func NewProductRepository(impl ProductRepoImpl) ProductRepository {
	return &impl
}

func (r *ProductRepoImpl) GetProductByProductID(ctx context.Context, productID int64) (res Product, err error) {
	rows, err := r.DB.QueryxContext(ctx, "select product_id, sku, name, price, qty from products where product_id = $1", productID)
	if err != nil {
		return res, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.StructScan(&res)
		if err != nil {
			return res, err
		}
	}

	return res, nil
}

func (r *ProductRepoImpl) GetAllProduct(ctx context.Context) (res []Product, err error) {
	rows, err := r.DB.QueryxContext(ctx, "select product_id, sku, name, price, qty from products order by product_id asc")
	if err != nil {
		return res, err
	}

	defer rows.Close()

	for rows.Next() {
		tmp := Product{}
		err = rows.StructScan(&tmp)
		if err != nil {
			return res, err
		}

		res = append(res, tmp)
	}

	return res, nil
}

func (r *ProductRepoImpl) UpdateProductQtyByProductID(tx *sqlx.Tx, ctx context.Context, form Product) (err error) {
	_, err = tx.ExecContext(ctx, "UPDATE products SET qty = qty - $1 WHERE product_id = $2", form.Qty, form.ProductID)
	if err != nil {
		return err
	}

	return nil
}
