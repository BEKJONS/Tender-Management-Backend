package http

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "tender_management/docs"
	"tender_management/internal/controller"
	rate_limiting "tender_management/internal/usecase/redis/rate-limiting"

	"log/slog"
)

// @title CRM API
// @version 1.0
// @description Admin Panel for managing the CRM
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your bearer token here
// @BasePath /api/v1

func NewRouter(engine *gin.Engine, log *slog.Logger, casbin *casbin.Enforcer, ctr *controller.Controller, limiting *rate_limiting.RateLimiter) {

	engine.Use(CORSMiddleware())
	engine.Use(RateLimitingMiddleware(limiting))

	engine.GET("/swagger/*eny", ginSwagger.WrapHandler(swaggerFiles.Handler))

	user := engine.Group("/auth")
	tend := engine.Group("/tenders")
	bid := engine.Group("/")
	usr := engine.Group("/users")

	newUserRoutes(user, ctr.Auth, log)
	newTenderRoutes(tend, ctr.Tend, casbin, log)
	newBidRoutes(bid, ctr.Bid, casbin, log)
	newUserController(usr, ctr.Tend, casbin, ctr.Bid, log)
}
