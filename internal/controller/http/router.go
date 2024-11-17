package http

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "tender_management/docs"
	"tender_management/internal/controller"

	"log/slog"
)

// title Api For CRM
// version 1.0
// description Admin Panel
// @securityDefinitions.apiKey BearerAuth
// @in header
// @name Authorization
// @description Enter your bearer token here
func NewRouter(engine *gin.Engine, log *slog.Logger, casbin *casbin.Enforcer, ctr *controller.Controller) {

	engine.Use(CORSMiddleware())

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
