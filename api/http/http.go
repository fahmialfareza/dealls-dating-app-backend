package http

import (
	"fmt"

	"github.com/fahmialfareza/deals-dating-app-backend/configs/environment"
	"github.com/fahmialfareza/deals-dating-app-backend/configs/logger/newrelic"
	"github.com/fahmialfareza/deals-dating-app-backend/internal/constant"
	"github.com/fahmialfareza/deals-dating-app-backend/internal/delivery/http"
	"github.com/fahmialfareza/deals-dating-app-backend/internal/usecase"
	"github.com/gin-gonic/gin"
)

func StartHTTPServer(newRelic *newrelic.NewRelic, usecase usecase.IUsecase) {
	handler := http.NewHTTPHandler(newRelic.App, newRelic.Log, usecase)

	if environment.Environment == constant.ProductionEnvironment {
		gin.SetMode(gin.ReleaseMode)
	}

	// Set up Gin router
	router := gin.Default()

	// Router
	v1 := router.Group("/api/v1")
	{
		authV1 := v1.Group("/auth")
		{
			authV1.POST("/register", handler.Register)
			authV1.POST("/login", handler.Login)
			authV1.GET("/profile", handler.GetProfile)
		}

		swipeV1 := v1.Group("/swipe")
		{
			swipeV1.GET("/", handler.GetSwipe)
			swipeV1.POST("/", handler.Swipe)
		}

		purchaseV1 := v1.Group("/purchase")
		{
			purchaseV1.POST("/", handler.Purchase)
		}
	}

	router.Run(fmt.Sprintf(":%s", environment.Port))
}
