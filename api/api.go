package api

import (
	"github.com/fahmialfareza/deals-dating-app-backend/api/http"
	"github.com/fahmialfareza/deals-dating-app-backend/internal/repository"
	"github.com/fahmialfareza/deals-dating-app-backend/internal/usecase"
)

func StartAPI() {
	config := loadConfig(service)

	repository := repository.NewRepository(config.postgres, config.redis)
	usecase := usecase.NewUsecase(repository)

	if runServices["http"] {
		http.StartHTTPServer(config.newRelic, usecase)
	}
}
