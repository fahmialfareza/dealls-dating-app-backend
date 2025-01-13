package http

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	"github.com/fahmialfareza/deals-dating-app-backend/internal/domain"
	"github.com/fahmialfareza/deals-dating-app-backend/internal/dto"
	"github.com/fahmialfareza/deals-dating-app-backend/pkg/database/postgres"
	"github.com/fahmialfareza/deals-dating-app-backend/pkg/logger"
	"github.com/gin-gonic/gin"
)

type IAuthHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	GetProfile(c *gin.Context)
}

func (h *HTTPHandler) Register(c *gin.Context) {
	txn, ctx := logger.StartTransaction(h.newRelicApp, h.log, c, "HTTPHandler.Register")
	defer txn.End()

	segment := logger.StartSegment(ctx, "HTTPHandler.Register")
	defer segment.End()

	var user domain.RegisterRequest

	// Get the file from the form-data
	file, err := c.FormFile("image")
	if err != nil {
		err := domain.WrapError(http.StatusBadRequest, "Failed to read image file")
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), nil)
		h.sendResponse(c, txn, nil, err)
		return
	}

	// Open the uploaded file
	fileContent, err := file.Open()
	if err != nil {
		err := domain.WrapError(http.StatusBadRequest, "Failed to open image file")
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), nil)
		h.sendResponse(c, txn, nil, err)
		return
	}
	defer fileContent.Close()

	// Read the file content into memory
	fileBytes, err := io.ReadAll(fileContent)
	if err != nil {
		err := domain.WrapError(http.StatusInternalServerError, "Failed to read file content")
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), nil)
		h.sendResponse(c, txn, nil, err)
		return
	}

	// Encode the file content to Base64
	base64String := base64.StdEncoding.EncodeToString(fileBytes)

	// Prepare the Base64 data with a prefix (e.g., data:image/jpeg;base64,)
	fileType := file.Header.Get("Content-Type")
	if fileType == "" {
		fileType = "application/octet-stream" // Default fallback
	}
	user.ImageData = fmt.Sprintf("data:%s;base64,%s", fileType, base64String)
	user.ImageFileName = file.Filename
	user.Name = c.Request.FormValue("name")
	user.Email = c.Request.FormValue("email")
	user.Password = c.Request.FormValue("password")
	user.Bio = c.Request.FormValue("bio")

	// validate the request
	if err := dto.ValidateAuthRegister(ctx, user); err != nil {
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

	var data domain.LoginResponse
	data, err = h.usecase.Register(ctx, user)
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), nil)
		h.sendResponse(c, txn, nil, err)
		return
	}

	if err = postgres.Commit(ctx); err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), nil)
		h.sendResponse(c, txn, nil, err)
		return
	}

	h.sendResponse(c, txn, data, nil)
}

// Login implements IHTTPHandler.
func (h *HTTPHandler) Login(c *gin.Context) {
	txn, ctx := logger.StartTransaction(h.newRelicApp, h.log, c, "HTTPHandler.Login")
	defer txn.End()

	segment := logger.StartSegment(ctx, "HTTPHandler.Login")
	defer segment.End()

	var user domain.LoginRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		err := domain.WrapError(http.StatusBadRequest, err.Error())
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), nil)
		h.sendResponse(c, txn, nil, err)
		return
	}

	// validate the request
	if err := dto.ValidateAuthLogin(ctx, user); err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), nil)
		h.sendResponse(c, txn, nil, err)
		return
	}

	data, err := h.usecase.Login(c, user.Email, user.Password)
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), nil)
		h.sendResponse(c, txn, nil, err)
		return
	}

	h.sendResponse(c, txn, data, nil)
}

// GetProfile implements IHTTPHandler.
func (h *HTTPHandler) GetProfile(c *gin.Context) {
	txn, ctx := logger.StartTransaction(h.newRelicApp, h.log, c, "HTTPHandler.GetProfile")
	defer txn.End()

	segment := logger.StartSegment(ctx, "HTTPHandler.GetProfile")
	defer segment.End()

	authHeader := c.Request.Header.Get("Authorization")
	user, err := h.authMiddleware(ctx, authHeader)
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), nil)
		h.sendResponse(c, txn, nil, err)
		return
	}

	h.sendResponse(c, txn, user, nil)
}
