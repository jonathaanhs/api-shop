package service

import (
	"context"
	"log"

	"github.com/learn/api-shop/internal/repo"
)

type Promotion interface {
	ApplyPromotion(item *repo.OrderDetail, v repo.OrderDetail, productDetail *repo.Product, promo *repo.Promo, res *Checkout) error
}

type ProductPromoDiscount struct {
}

func (p *ProductPromoDiscount) ApplyPromotion(item *repo.OrderDetail, v repo.OrderDetail, productDetail *repo.Product, promo *repo.Promo, res *Checkout) error {
	tmpQty := v.Qty - 1
	item.Price = float64(tmpQty) * productDetail.Price
	return nil
}

type ProductPromoFree struct {
	ProductRepo repo.ProductRepository
}

func (p *ProductPromoFree) ApplyPromotion(item *repo.OrderDetail, v repo.OrderDetail, productDetail *repo.Product, promo *repo.Promo, res *Checkout) error {
	productRewardDetail, err := p.ProductRepo.GetProductByProductID(context.Background(), int64(promo.Reward))
	if err != nil {
		log.Printf("error while do GetProductByProductID %+v", err)
		return err
	}
	res.Items = append(res.Items, productRewardDetail.Name)
	return nil
}

type DiscountPromo struct {
}

func (p *DiscountPromo) ApplyPromotion(item *repo.OrderDetail, v repo.OrderDetail, productDetail *repo.Product, promo *repo.Promo, res *Checkout) error {
	item.Price = (productDetail.Price * float64(v.Qty)) - ((productDetail.Price * float64(v.Qty)) * (promo.Reward / 100))
	return nil
}
