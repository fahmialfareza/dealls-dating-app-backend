package http

import (
	"database/sql"

	"github.com/fahmialfareza/deals-dating-app-backend/pkg/database/postgres"
	"github.com/fahmialfareza/deals-dating-app-backend/pkg/logger"
	"github.com/gin-gonic/gin"
)

type IPurchaseHandler interface {
	Purchase(c *gin.Context)
}

// Purchase implements IHTTPHandler.
func (h *HTTPHandler) Purchase(c *gin.Context) {
	txn, ctx := logger.StartTransaction(h.newRelicApp, h.log, c, "HTTPHandler.Purchase")
	defer txn.End()

	segment := logger.StartSegment(ctx, "HTTPHandler.Purchase")
	defer segment.End()

	authHeader := c.Request.Header.Get("Authorization")
	user, err := h.authMiddleware(ctx, authHeader)
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), nil)
		h.sendResponse(c, txn, nil, err)
		return
	}

	ctx, err = postgres.StartTransaction(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	})
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), nil)
		h.sendResponse(c, txn, nil, err)
		return
	}
	defer func() {
		if err != nil {
			logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), nil)
			err := postgres.Rollback(ctx)
			if err != nil {
				logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), nil)
			}
		}
	}()

	if err := h.usecase.Purchase(ctx, user.ID); err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), nil)
		h.sendResponse(c, txn, nil, err)
		return
	}

	if err = postgres.Commit(ctx); err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), nil)
		h.sendResponse(c, txn, nil, err)
		return
	}

	h.sendResponse(c, txn, nil, nil)
}
