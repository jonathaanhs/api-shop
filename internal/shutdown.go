package internal

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/dig"
)

func Shutdown(p struct {
	dig.In
	Pg  *sqlx.DB
	Srv *http.Server
}) error {
	log.Printf("Shutdown at %s\n", time.Now().String())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := p.Pg.Close(); err != nil {
		return err
	}

	if err := p.Srv.Shutdown(ctx); err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}

	log.Println("Server exiting")
	return nil
}
