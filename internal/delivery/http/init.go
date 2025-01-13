package http

import (
	"github.com/fahmialfareza/deals-dating-app-backend/internal/usecase"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
)

type IHTTPHandler interface {
	IAuthHandler
	ISwipeHandler
	IPurchaseHandler
}

type HTTPHandler struct {
	newRelicApp *newrelic.Application
	log         *logrus.Logger
	usecase     usecase.IUsecase
}

func NewHTTPHandler(newRelicApp *newrelic.Application, log *logrus.Logger, usecase usecase.IUsecase) IHTTPHandler {
	return &HTTPHandler{
		newRelicApp: newRelicApp,
		log:         log,
		usecase:     usecase,
	}
}
