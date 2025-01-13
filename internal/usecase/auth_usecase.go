package usecase

import (
	"context"

	"github.com/fahmialfareza/deals-dating-app-backend/internal/constant"
	"github.com/fahmialfareza/deals-dating-app-backend/internal/domain"
	"github.com/fahmialfareza/deals-dating-app-backend/pkg/jwt"
	"github.com/fahmialfareza/deals-dating-app-backend/pkg/logger"
	"github.com/fahmialfareza/deals-dating-app-backend/pkg/object_storage/imagekit"
	"golang.org/x/crypto/bcrypt"
)

type IUserUsecase interface {
	Register(ctx context.Context, data domain.RegisterRequest) (domain.LoginResponse, error)
	Login(ctx context.Context, email, password string) (domain.LoginResponse, error)
	GetProfile(ctx context.Context, id uint) (domain.GetProfileResponse, error)
	CheckToken(ctx context.Context, token string) (uint, error)
}

// Register implements IUsecase.
func (u *Usecase) Register(ctx context.Context, data domain.RegisterRequest) (result domain.LoginResponse, err error) {
	segment := logger.StartSegment(ctx, "Usecase.Register")
	defer segment.End()

	userData := domain.User{
		Name:      data.Name,
		Email:     data.Email,
		Password:  data.Password,
		IsPremium: false,
		Profile: &domain.Profile{
			Bio: data.Bio,
		},
	}

	// upload image
	userData.Profile.Image, err = imagekit.UploadImage(ctx, data.ImageData, data.ImageFileName)
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), nil)
		return result, err
	}

	// create user
	user, err := u.repository.UpsertUser(ctx, userData)
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), nil)
		return result, err
	}

	// generate token
	token, err := jwt.SetToken(ctx, user.ID)
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), nil)
		return result, err
	}

	result = domain.LoginResponse{
		Token: token,
	}

	return result, nil
}

// Login implements IUsecase.
func (u *Usecase) Login(ctx context.Context, email, password string) (result domain.LoginResponse, err error) {
	segment := logger.StartSegment(ctx, "Usecase.Login")
	defer segment.End()

	// get user by email
	user, err := u.repository.GetUserDetail(ctx, nil, &email)
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"email":    email,
			"password": password,
		})
		return result, err
	}

	// Check the password
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		err = constant.PasswordIsWrong
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"email":    email,
			"password": password,
		})
		return result, err
	}

	// generate token
	token, err := jwt.SetToken(ctx, user.ID)
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"email":    email,
			"password": password,
		})
		return result, err
	}

	result = domain.LoginResponse{
		Token: token,
	}

	return result, nil
}

// GetProfile implements IUsecase.
func (u *Usecase) GetProfile(ctx context.Context, id uint) (result domain.GetProfileResponse, err error) {
	segment := logger.StartSegment(ctx, "Usecase.Login")
	defer segment.End()

	// get user by id
	user, err := u.repository.GetUserDetail(ctx, &id, nil)
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"id": id,
		})
		return result, err
	}

	result = domain.GetProfileResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
		Name:      user.Name,
		Email:     user.Email,
		IsPremium: user.IsPremium,
		Profile:   user.Profile,
	}

	return result, nil
}

func (a *Usecase) CheckToken(ctx context.Context, token string) (uint, error) {
	segment := logger.StartSegment(ctx, "Usecase.CheckToken")
	defer segment.End()

	// Check token validity and get user id
	userID, err := jwt.ExtractJWT(ctx, token)
	if err != nil {
		logger.PrintErrorLog(ctx, err, logger.GetErrorFileLine(), map[string]interface{}{
			"token": token,
		})
		return 0, err
	}

	return userID, nil
}
