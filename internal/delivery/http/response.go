package http

import (
	"net/http"

	"github.com/fahmialfareza/deals-dating-app-backend/internal/domain"
	"github.com/gin-gonic/gin"
	newrelic "github.com/newrelic/go-agent/v3/newrelic"
)

type httpResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func (h *HTTPHandler) sendResponse(c *gin.Context, txn *newrelic.Transaction, data interface{}, err error) {
	if err != nil {
		if txn != nil {
			txn.NoticeError(err)
		}

		err, ok := domain.ExtractCustomError(err)
		if ok {
			c.JSON(err.Code, httpResponse{
				Data:    data,
				Message: err.Message,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, httpResponse{
			Data:    data,
			Message: http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	c.JSON(http.StatusOK, httpResponse{
		Data:    data,
		Message: "Success!",
	})
}
