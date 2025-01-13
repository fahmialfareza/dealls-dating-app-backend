package imagekit

import (
	"context"

	"github.com/fahmialfareza/deals-dating-app-backend/pkg/logger"
	"github.com/imagekit-developer/imagekit-go/api/uploader"
)

func UploadImage(ctx context.Context, file string, fileName string) (string, error) {
	segment := logger.StartSegment(ctx, "imagekit.UploadImage")
	defer segment.End()

	useUniqueFileName := true
	resp, err := ik.Uploader.Upload(ctx, file, uploader.UploadParam{
		UseUniqueFileName: &useUniqueFileName,
		FileName:          fileName,
	})
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"file_name": fileName,
		})
		return "", err
	}

	return resp.Data.Url, nil
}
