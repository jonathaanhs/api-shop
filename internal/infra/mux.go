package infra

import (
	"net/http"
	"time"
)

type (
	MuxCfg struct {
		Address      string        `envconfig:"ADDRESS" default:":8089" required:"true"`
		ReadTimeout  time.Duration `envconfig:"READ_TIMEOUT" default:"5s"`
		WriteTimeout time.Duration `envconfig:"WRITE_TIMEOUT" default:"10s"`
		Debug        bool          `envconfig:"DEBUG" default:"true"`
	}
)

func NewMux() *http.ServeMux {
	return http.NewServeMux()
}
