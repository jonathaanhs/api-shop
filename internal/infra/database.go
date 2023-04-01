package infra

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"

	_ "github.com/lib/pq"
)

type (
	Databases struct {
		dig.Out
		Pg *sqlx.DB
	}

	DatabaseCfgs struct {
		dig.In
		Pg *DatabaseCfg
	}

	DatabaseCfg struct {
		DBName string `envconfig:"DBNAME" required:"true" default:"dbname"`
		DBUser string `envconfig:"DBUSER" required:"true" default:"dbuser"`
		DBPass string `envconfig:"DBPASS" required:"true" default:"dbpass"`
		Host   string `envconfig:"HOST" required:"true" default:"localhost"`
		Port   string `envconfig:"PORT" required:"true" default:"9999"`

		MaxOpenConns    int           `envconfig:"MAX_OPEN_CONNS" default:"20" required:"true"`
		MaxIdleConns    int           `envconfig:"MAX_IDLE_CONNS" default:"5" required:"true"`
		ConnMaxLifetime time.Duration `envconfig:"CONN_MAX_LIFETIME" default:"15m" required:"true"`
	}
)

func NewDatabases(cfgs DatabaseCfgs) Databases {
	return Databases{
		Pg: openPostgres(cfgs.Pg),
	}
}

func openPostgres(p *DatabaseCfg) *sqlx.DB {
	conn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		p.DBUser, p.DBPass, p.Host, p.Port, p.DBName,
	)
	db, err := sqlx.Open("postgres", conn)
	if err != nil {
		logrus.Fatalf("postgres: %s", err.Error())
	}

	db.SetConnMaxLifetime(p.ConnMaxLifetime)
	db.SetMaxIdleConns(p.MaxIdleConns)
	db.SetMaxOpenConns(p.MaxOpenConns)

	if err = db.Ping(); err != nil {
		logrus.Fatalf("postgres: %s", err.Error())
	}

	return db
}
