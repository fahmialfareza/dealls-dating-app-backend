package repository

import (
	"context"

	"github.com/fahmialfareza/deals-dating-app-backend/internal/domain"
	"github.com/fahmialfareza/deals-dating-app-backend/pkg/logger"
)

type IUserRepository interface {
	UpsertUser(ctx context.Context, data domain.User) (domain.User, error)
	GetUserProfile(ctx context.Context, notIn []uint) ([]domain.User, error)
	GetUserDetail(ctx context.Context, id *uint, email *string) (domain.User, error)
}

// UpsertUser implements IRepository.
func (r *Repository) UpsertUser(ctx context.Context, data domain.User) (domain.User, error) {
	segment := logger.StartSegment(ctx, "Repository.UpsertUser")
	defer segment.End()

	// Insert to db
	if err := r.postgres.UpsertUser(ctx, &data); err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"data": data,
		})
		return data, err
	}

	// set to redis
	if err := r.setUserDetailCache(ctx, data); err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"data": data,
		})
		return data, err
	}

	return data, nil
}

// GetUserProfile implements IRepository.
func (r *Repository) GetUserProfile(ctx context.Context, notIn []uint) ([]domain.User, error) {
	segment := logger.StartSegment(ctx, "Repository.GetUserProfile")
	defer segment.End()

	// get from postgres
	data, err := r.postgres.GetUserProfile(ctx, notIn)
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"not_in": notIn,
		})
		return data, err
	}

	return data, nil
}

// GetUserDetail implements IRepository.
func (r *Repository) GetUserDetail(ctx context.Context, id *uint, email *string) (domain.User, error) {
	segment := logger.StartSegment(ctx, "Repository.GetUserDetail")
	defer segment.End()

	// get from redis
	data, err := r.redis.GetUserDetail(ctx, id, email)
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"id":    id,
			"email": email,
		})
		return data, err
	}

	if data.ID != 0 {
		return data, nil
	}

	// get from postgres
	data, err = r.postgres.GetUserDetail(ctx, id, email)
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"id":    id,
			"email": email,
		})
		return data, err
	}

	// set to redis
	if err := r.setUserDetailCache(ctx, data); err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"data": data,
		})
		return data, err
	}

	return data, nil
}

func (u *Repository) setUserDetailCache(ctx context.Context, data domain.User) error {
	segment := logger.StartSegment(ctx, "Repository.setUserDetailCache")
	defer segment.End()

	if data.ID != 0 {
		// set to redis cache (By ID)
		if err := u.redis.SetUserDetail(ctx, &data.ID, nil, data); err != nil {
			logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
				"data": data,
			})
			return err
		}

		// set to redis cache (By Email)
		if err := u.redis.SetUserDetail(ctx, nil, &data.Email, data); err != nil {
			logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
				"data": data,
			})
			return err
		}
	}

	return nil
}
