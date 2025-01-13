package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"
)

type logError struct {
	Message  string                 `json:"message"`
	Request  map[string]interface{} `json:"request"`
	FilePath string                 `json:"file_path"`
}

type logInfo struct {
	Message  string                 `json:"message"`
	Request  map[string]interface{} `json:"request"`
	FilePath string                 `json:"file_path"`
}

func log(ctx context.Context) *logrus.Entry {
	logCtx := ctx.Value("log")
	log, ok := logCtx.(*logrus.Entry)
	if !ok {
		logger := logrus.New()
		log = logrus.NewEntry(logger)
		log.Logger.SetFormatter(&logrus.JSONFormatter{})
		log.Logger.Out = os.Stdout
	}

	return log
}

func errorFormatter(logError logError) error {
	errorByte, _ := json.Marshal(logError)
	return fmt.Errorf("%s", string(errorByte))
}

func infoFormatter(infoLog logInfo) string {
	errorByte, _ := json.Marshal(infoLog)
	return string(errorByte)
}

func GetErrorFileLine() string {
	_, filePath, lineNumber, _ := runtime.Caller(1)
	absFilePath, _ := filepath.Abs(filePath)
	return fmt.Sprintf("%s:%v", absFilePath, lineNumber)
}

func PrintErrorLog(ctx context.Context, err error, filePath string, request map[string]interface{}) {
	log(ctx).Error(errorFormatter(logError{
		Message:  err.Error(),
		Request:  request,
		FilePath: filePath,
	}))
}

func PrintInfoLog(ctx context.Context, message string, filePath string, request map[string]interface{}) {
	log(ctx).Info(infoFormatter(logInfo{
		Message:  message,
		Request:  request,
		FilePath: filePath,
	}))
}
