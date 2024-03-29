package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/learn/api-shop/internal/repo"
	"go.uber.org/dig"
)

type (
	Checkout struct {
		Items       []string `json:"items"`
		TotalAmount float64  `json:"total_amount"`
	}

	CheckoutUsecase interface {
		Checkout(ctx context.Context, form []repo.OrderDetail) (res Checkout, err error)
	}

	CheckoutUsecaseImpl struct {
		dig.In
		OrderRepo   repo.OrderRepository
		ProductRepo repo.ProductRepository
		PromoRepo   repo.PromoRepository
	}
)

func NewCheckoutUsecase(impl CheckoutUsecaseImpl) CheckoutUsecase {
	return &impl
}

func (c *CheckoutUsecaseImpl) Checkout(ctx context.Context, form []repo.OrderDetail) (res Checkout, err error) {
	tx, err := c.OrderRepo.BeginTx()
	if err != nil {
		log.Printf("error while do BeginTx %+v", err)
		return res, err
	}

	defer c.OrderRepo.RollbackTx(tx)

	orderID, err := c.createOrder(tx, ctx, &res)
	if err != nil {
		return res, err
	}

	for i, v := range form {
		err := c.processOrderItem(ctx, tx, &form[i], v, orderID, &res)
		if err != nil {
			return res, err
		}
	}

	err = c.OrderRepo.CreateOrderDetails(tx, ctx, form)
	if err != nil {
		log.Printf("error while do CreateOrderDetails %+v", err)
		return res, err
	}

	c.OrderRepo.CommitTx(tx)

	return res, nil
}

func (c *CheckoutUsecaseImpl) createOrder(tx *sqlx.Tx, ctx context.Context, res *Checkout) (int64, error) {
	orderID, err := c.OrderRepo.CreateOrder(tx, ctx, repo.Order{
		Date:  time.Now(),
		Total: res.TotalAmount,
	})
	if err != nil {
		log.Printf("error while do CreateOrder %+v", err)
		return 0, err
	}
	return orderID, nil
}

func (c *CheckoutUsecaseImpl) processOrderItem(ctx context.Context, tx *sqlx.Tx, item *repo.OrderDetail, v repo.OrderDetail, orderID int64, res *Checkout) error {
	promo, err := c.PromoRepo.GetPromoByProductID(ctx, v.ProductID)
	if err != nil {
		log.Printf("error while do GetPromoByProductID %+v", err)
		return err
	}

	productDetail, err := c.ProductRepo.GetProductByProductID(ctx, v.ProductID)
	if err != nil {
		log.Printf("error while do GetProductByProductID %+v", err)
		return err
	}

	if productDetail.Qty < v.Qty {
		return fmt.Errorf("the product %s qty is not enough to fulfill the request", productDetail.Name)
	}

	err = c.calculatePriceAndRewards(ctx, item, v, &productDetail, &promo, orderID, res)
	if err != nil {
		return err
	}
	err = c.updateProductQty(ctx, tx, v)
	if err != nil {
		return err
	}

	res.TotalAmount += item.Price
	return nil
}

func (c *CheckoutUsecaseImpl) calculatePriceAndRewards(ctx context.Context, item *repo.OrderDetail, v repo.OrderDetail, productDetail *repo.Product, promo *repo.Promo, orderID int64, res *Checkout) error {
	item.Price = float64(v.Qty) * productDetail.Price
	item.OrderID = orderID
	for i := 0; i < int(v.Qty); i++ {
		res.Items = append(res.Items, productDetail.Name)
	}

	var promotion Promotion

	if v.Qty >= promo.MinQty {
		switch promo.PromoType {
		case "product":
			promotion = c.calculateProductPromo(item, v, productDetail, promo, res)
		case "discount":
			promotion = &DiscountPromo{}
		}

		item.PromoID = promo.PromoID
	}

	if promotion != nil {
		err := promotion.ApplyPromotion(item, v, productDetail, promo, res)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *CheckoutUsecaseImpl) calculateProductPromo(item *repo.OrderDetail, v repo.OrderDetail, productDetail *repo.Product, promo *repo.Promo, res *Checkout) Promotion {
	if v.ProductID == int64(promo.Reward) {
		return &ProductPromoDiscount{}
	}

	return &ProductPromoFree{
		ProductRepo: c.ProductRepo,
	}
}

func (c *CheckoutUsecaseImpl) updateProductQty(ctx context.Context, tx *sqlx.Tx, v repo.OrderDetail) error {
	err := c.ProductRepo.UpdateProductQtyByProductID(tx, ctx, repo.Product{
		ProductID: v.ProductID,
		Qty:       v.Qty,
	})
	if err != nil {
		log.Printf("error while do UpdateProductQtyByProductID %+v", err)
		return err
	}
	return nil
}
