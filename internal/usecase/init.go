package usecase

import (
	"github.com/fahmialfareza/deals-dating-app-backend/internal/repository"
)

type IUsecase interface {
	IUserUsecase
	ISwipeUsecase
	IPurchaseUsecase
}

type Usecase struct {
	repository repository.IRepository
}

func NewUsecase(repository repository.IRepository) IUsecase {
	return &Usecase{
		repository: repository,
	}
}
