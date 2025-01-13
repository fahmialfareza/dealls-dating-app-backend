package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/fahmialfareza/deals-dating-app-backend/internal/domain"
	"github.com/fahmialfareza/deals-dating-app-backend/internal/repository"
	mocks_repository "github.com/fahmialfareza/deals-dating-app-backend/internal/repository/mocks"
	"github.com/fahmialfareza/deals-dating-app-backend/pkg/jwt"
	"github.com/golang/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestUsecase_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks_repository.NewMockIRepository(ctrl)
	passwordByte, _ := bcrypt.GenerateFromPassword([]byte("a"), 10)
	token, _ := jwt.SetToken(context.Background(), 1)

	type fields struct {
		repository repository.IRepository
	}
	type args struct {
		ctx      context.Context
		email    string
		password string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult domain.LoginResponse
		wantErr    bool
		function   func(args)
	}{
		{
			name: "success",
			fields: fields{
				repository: mockRepository,
			},
			args: args{
				ctx:      context.Background(),
				email:    "a@a.com",
				password: "a",
			},
			wantResult: domain.LoginResponse{
				Token: token,
			},
			wantErr: false,
			function: func(a args) {
				mockRepository.EXPECT().GetUserDetail(a.ctx, nil, &a.email).Return(domain.User{ID: 1, Password: string(passwordByte)}, nil)
			},
		},
		{
			name: "error",
			fields: fields{
				repository: mockRepository,
			},
			args: args{
				ctx:      context.Background(),
				email:    "a@a.com",
				password: "a",
			},
			wantResult: domain.LoginResponse{},
			wantErr:    true,
			function: func(a args) {
				mockRepository.EXPECT().GetUserDetail(a.ctx, nil, &a.email).Return(domain.User{ID: 1, Password: string(passwordByte)}, errors.New(""))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.function(tt.args)

			u := &Usecase{
				repository: tt.fields.repository,
			}
			gotResult, err := u.Login(tt.args.ctx, tt.args.email, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Usecase.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("Usecase.Login() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestUsecase_GetProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks_repository.NewMockIRepository(ctrl)
	now := time.Now()

	type fields struct {
		repository repository.IRepository
	}
	type args struct {
		ctx context.Context
		id  uint
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		want     domain.User
		wantErr  bool
		function func(args)
	}{
		{
			name: "success",
			fields: fields{
				repository: mockRepository,
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want: domain.User{
				ID:        1,
				CreatedAt: now,
				UpdatedAt: now,
				Name:      "a",
				Email:     "a@a.com",
				Password:  "a",
				IsPremium: false,
			},
			wantErr: false,
			function: func(a args) {
				mockRepository.EXPECT().GetUserDetail(a.ctx, &a.id, nil).Return(domain.User{
					ID:        1,
					CreatedAt: now,
					UpdatedAt: now,
					Name:      "a",
					Email:     "a@a.com",
					Password:  "a",
					IsPremium: false,
				}, nil)
			},
		},
		{
			name: "error",
			fields: fields{
				repository: mockRepository,
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want:    domain.User{},
			wantErr: true,
			function: func(a args) {
				mockRepository.EXPECT().GetUserDetail(a.ctx, &a.id, nil).Return(domain.User{}, errors.New(""))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.function(tt.args)

			u := &Usecase{
				repository: tt.fields.repository,
			}
			got, err := u.GetProfile(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Usecase.GetProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Usecase.GetProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}
