package http

import (
	"net/http"

	"github.com/fahmialfareza/deals-dating-app-backend/internal/domain"
	"github.com/fahmialfareza/deals-dating-app-backend/internal/dto"
	"github.com/fahmialfareza/deals-dating-app-backend/pkg/logger"
	"github.com/gin-gonic/gin"
)

type ISwipeHandler interface {
	GetSwipe(c *gin.Context)
	Swipe(c *gin.Context)
}

// GetSwipe implements IHTTPHandler.
func (h *HTTPHandler) GetSwipe(c *gin.Context) {
	txn, ctx := logger.StartTransaction(h.newRelicApp, h.log, c, "HTTPHandler.GetSwipe")
	defer txn.End()

	segment := logger.StartSegment(ctx, "HTTPHandler.GetSwipe")
	defer segment.End()

	authHeader := c.Request.Header.Get("Authorization")
	user, err := h.authMiddleware(ctx, authHeader)
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), nil)
		h.sendResponse(c, txn, nil, err)
		return
	}

	data, err := h.usecase.GetSwipe(ctx, user.ID)
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), nil)
		h.sendResponse(c, txn, nil, err)
		return
	}

	h.sendResponse(c, txn, data, nil)
}

// Swipe implements IHTTPHandler.
func (h *HTTPHandler) Swipe(c *gin.Context) {
	txn, ctx := logger.StartTransaction(h.newRelicApp, h.log, c, "HTTPHandler.Login")
	defer txn.End()

	segment := logger.StartSegment(ctx, "HTTPHandler.Login")
	defer segment.End()

	authHeader := c.Request.Header.Get("Authorization")
	user, err := h.authMiddleware(ctx, authHeader)
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), nil)
		h.sendResponse(c, txn, nil, err)
		return
	}

	var swipe domain.SwipeRequest
	if err := c.ShouldBindJSON(&swipe); err != nil {
		err := domain.WrapError(http.StatusBadRequest, err.Error())
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), nil)
		h.sendResponse(c, txn, nil, err)
		return
	}

	// validate the request
	if err := dto.ValidateSwipe(ctx, swipe); err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), nil)
		h.sendResponse(c, txn, nil, err)
		return
	}

	if err := h.usecase.Swipe(c, user.ID, swipe.SwipedID, swipe.SwipeType); err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), nil)
		h.sendResponse(c, txn, nil, err)
		return
	}

	h.sendResponse(c, txn, nil, nil)
}
