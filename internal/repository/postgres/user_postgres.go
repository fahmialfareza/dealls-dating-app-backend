package postgres

import (
	"context"
	"strings"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/fahmialfareza/deals-dating-app-backend/internal/constant"
	"github.com/fahmialfareza/deals-dating-app-backend/internal/domain"
	"github.com/fahmialfareza/deals-dating-app-backend/pkg/converter"
	"github.com/fahmialfareza/deals-dating-app-backend/pkg/logger"
	"gorm.io/gorm"
)

type IUserPostgresRepository interface {
	UpsertUser(ctx context.Context, user *domain.User) error
	GetUserProfile(ctx context.Context, notIn []uint) ([]domain.User, error)
	GetUserDetail(ctx context.Context, id *uint, email *string) (domain.User, error)
}

func (u *PostgresRepository) UpsertUser(ctx context.Context, data *domain.User) (err error) {
	segment := logger.StartSegment(ctx, "PostgresRepository.UpsertUser")
	defer segment.End()

	if err := hystrix.DoC(ctx, constant.HystrixPostgres, func(ctx context.Context) error {
		// generate the password with bcrypt
		data.Password, err = converter.HashPassword(ctx, data.Password)
		if err != nil {
			logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
				"data": data,
			})
			return err
		}

		// create user
		if err := u.db.GORM(ctx).Save(data).Error; err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				err = constant.UserHasAlreadyExist
			}

			logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
				"data": data,
			})
			return err
		}

		return nil
	}, nil); err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"data": data,
		})
		return err
	}

	return nil
}

// GetUserProfile implements IPostgresRepository.
func (u *PostgresRepository) GetUserProfile(ctx context.Context, notIn []uint) (result []domain.User, err error) {
	segment := logger.StartSegment(ctx, "PostgresRepository.GetUserProfile")
	defer segment.End()

	if err := hystrix.DoC(ctx, constant.HystrixPostgres, func(ctx context.Context) error {
		if err := u.db.GORM(ctx).Model(&domain.User{}).Not(notIn).Preload("Profile", "deleted_at IS NULL").Find(&result).Error; err != nil {
			if err.Error() == gorm.ErrRecordNotFound.Error() {
				err = constant.UserAccountCanNotBeFound
			}

			logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
				"not_in": notIn,
			})
			return err
		}

		return nil
	}, nil); err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"not_in": notIn,
		})
		return result, err
	}

	return result, nil
}

// GetUserByIDOrEmail implements IPostgresRepository.
func (u *PostgresRepository) GetUserDetail(ctx context.Context, id *uint, email *string) (result domain.User, err error) {
	segment := logger.StartSegment(ctx, "PostgresRepository.GetUserDetail")
	defer segment.End()

	if err := hystrix.DoC(ctx, constant.HystrixPostgres, func(ctx context.Context) error {
		userQuery := u.db.GORM(ctx).Model(&domain.User{}).Preload("Profile", "deleted_at IS NULL")

		// Find by ID
		if id != nil {
			if err := userQuery.First(&result, *id).Error; err != nil {
				if err.Error() == gorm.ErrRecordNotFound.Error() {
					err = constant.UserAccountCanNotBeFound
				}

				logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
					"id":    id,
					"email": email,
				})
				return err
			}
			return nil
		} else if email != nil {
			// Find by Email
			if err := userQuery.Where("email = ?", *email).First(&result).Error; err != nil {
				if err.Error() == gorm.ErrRecordNotFound.Error() {
					err = constant.UserAccountCanNotBeFound
				}

				logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
					"id":    id,
					"email": email,
				})
				return err
			}
			return nil
		}

		return nil
	}, nil); err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"id":    id,
			"email": email,
		})
		return result, err
	}

	return result, nil
}
