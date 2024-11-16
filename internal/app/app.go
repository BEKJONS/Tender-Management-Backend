package app

import (
	"github.com/gin-gonic/gin"
	"log"
	"tender_management/config"
	"tender_management/internal/controller"
	"tender_management/internal/controller/http"
	"tender_management/internal/usecase/token"
	"tender_management/pkg/logger"
	"tender_management/pkg/postgres"
)

func Run(cfg config.Config) {

	logger1 := logger.NewLogger()

	db, err := postgres.Connection(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = token.ConfigToken(cfg)
	if err != nil {
		log.Fatal(err)
	}

	controller1 := controller.NewController(db, logger1)

	engine := gin.Default()
	http.NewRouter(engine, logger1, controller1)

	log.Fatal(engine.Run(cfg.RUN_PORT))
}
