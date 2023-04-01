package internal

import (
	"log"
	"net/http"

	"github.com/learn/api-shop/internal/infra"
	"go.uber.org/dig"
)

func Start(p struct {
	dig.In
	Cfg *infra.MuxCfg
	Srv *http.Server
}) (err error) {
	log.Println("Server Start ", p.Cfg.Address)
	return p.Srv.ListenAndServe()
}
