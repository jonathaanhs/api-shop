package infra

import (
	"fmt"
	"net/http"

	"github.com/kelseyhightower/envconfig"
	"go.uber.org/dig"
)

func LoadPgDatabaseCfg() (*DatabaseCfg, error) {
	var cfg DatabaseCfg
	prefix := "PG"
	if err := envconfig.Process(prefix, &cfg); err != nil {
		return nil, fmt.Errorf("%s: %w", prefix, err)
	}

	return &cfg, nil
}

func LoadMuxCfg() (*MuxCfg, error) {
	var cfg MuxCfg
	prefix := "APP"
	if err := envconfig.Process(prefix, &cfg); err != nil {
		return nil, fmt.Errorf("%s: %w", prefix, err)
	}
	return &cfg, nil
}

func LoadHttpServer(p struct {
	dig.In
	Cfg *MuxCfg
	M   *http.ServeMux
}) *http.Server {
	return &http.Server{
		Addr:         p.Cfg.Address,
		ReadTimeout:  p.Cfg.ReadTimeout,
		WriteTimeout: p.Cfg.WriteTimeout,
		Handler:      p.M,
	}
}
