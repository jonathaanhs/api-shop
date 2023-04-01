package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/learn/api-shop/internal"
	"github.com/learn/api-shop/internal/controller"
	"github.com/learn/api-shop/internal/infra"
	"github.com/learn/api-shop/internal/repo"
	"github.com/learn/api-shop/internal/service"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatal(err.Error())
	}

	container := dig.New()

	container.Provide(infra.LoadPgDatabaseCfg)
	container.Provide(infra.LoadMuxCfg)
	container.Provide(infra.LoadHttpServer)
	container.Provide(infra.NewDatabases)
	container.Provide(infra.NewMux)
	container.Provide(repo.NewOrderRepository)
	container.Provide(repo.NewProductRepository)
	container.Provide(repo.NewPromoRepository)
	container.Provide(service.NewCheckoutUsecase)

	if err := container.Invoke(controller.NewCheckoutHandler); err != nil {
		logrus.Fatal(err.Error())
	}

	if err := startApp(container); err != nil {
		logrus.Fatal(err.Error())
	}
}

func startApp(di *dig.Container) error {
	go func() {
		if err := di.Invoke(internal.Start); err != nil {
			log.Fatalf("start: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	if err := di.Invoke(internal.Shutdown); err != nil {
		return err
	}

	log.Println("Server exiting")

	return nil
}
