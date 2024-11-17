package app

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"tender_management/config"
	"tender_management/internal/controller"
	"tender_management/internal/controller/http"
	"tender_management/internal/usecase/cashing"
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

	rdb := cashing.NewRedisClient()

	controller1 := controller.NewController(db, logger1, rdb)

	path, err := os.Getwd()
	if err != nil {
		logger1.Error("Failed to get current working directory")
		return
	}

	casbinEnforcer, err := casbin.NewEnforcer(path+"/pkg/config/model.conf", path+"/pkg/config/policy.csv")
	if err != nil {
		panic(err)
	}

	engine := gin.Default()
	http.NewRouter(engine, logger1, casbinEnforcer, controller1)

	log.Fatal(engine.Run(cfg.RUN_PORT))
}
