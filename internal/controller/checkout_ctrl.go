package controller

import (
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/learn/api-shop/internal/repo"
	"github.com/learn/api-shop/internal/service"
	"go.uber.org/dig"
)

type (
	CheckoutCntrlImpl struct {
		dig.In
		CheckoutSvc service.CheckoutUsecase
	}
)

func NewCheckoutHandler(mux *http.ServeMux, checkoutSvc service.CheckoutUsecase) {
	hc := &CheckoutCntrlImpl{
		CheckoutSvc: checkoutSvc,
	}

	schema, err := CreateCheckoutSchema(hc)
	if err != nil {
		panic(err)
	}

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	mux.Handle("/graphql", hc.AdaptHTTPHandler(h))
}

func (cc CheckoutCntrlImpl) AdaptHTTPHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	}
}

func CreateCheckoutSchema(handler *CheckoutCntrlImpl) (graphql.Schema, error) {
	// Define the Checkout type
	checkoutType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Checkout",
		Fields: graphql.Fields{
			"items": &graphql.Field{
				Type: graphql.NewList(graphql.String),
			},
			"total_amount": &graphql.Field{
				Type: graphql.Float,
			},
		},
	})

	inputItemType := graphql.NewInputObject(
		graphql.InputObjectConfig{
			Name: "ItemInput",
			Fields: graphql.InputObjectConfigFieldMap{
				"product_id": &graphql.InputObjectFieldConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"qty": &graphql.InputObjectFieldConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
		},
	)

	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"checkout": &graphql.Field{
				Type: checkoutType,
				Args: graphql.FieldConfigArgument{
					"items": &graphql.ArgumentConfig{
						Type: graphql.NewList(inputItemType),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					items := p.Args["items"].([]interface{})

					orderDetails := make([]repo.OrderDetail, len(items))
					for i, item := range items {
						itemMap := item.(map[string]interface{})
						orderDetails[i] = repo.OrderDetail{
							ProductID: int64(itemMap["product_id"].(int)),
							Qty:       int64(itemMap["qty"].(int)),
						}
					}

					ctx := p.Context

					checkoutResult, err := handler.CheckoutSvc.Checkout(ctx, orderDetails)
					if err != nil {
						return nil, err
					}

					return checkoutResult, nil
				},
			},
		},
	})

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"ping": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return "pong", nil
				},
			},
		},
	})

	schemaConfig := graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	}
	schema, err := graphql.NewSchema(schemaConfig)

	if err != nil {
		return schema, fmt.Errorf("failed to create new schema, error: %v", err)
	}

	return schema, nil
}
