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
	"github.com/golang/mock/gomock"
)

func TestUsecase_validateSwipe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks_repository.NewMockIRepository(ctrl)
	swiperID := uint(1)
	now := time.Now()

	type fields struct {
		repository repository.IRepository
	}
	type args struct {
		ctx      context.Context
		swiperID uint
		date     time.Time
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantSwipes []domain.Swipe
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
				swiperID: swiperID,
				date:     now,
			},
			wantSwipes: []domain.Swipe{
				{
					ID:        1,
					CreatedAt: now,
					UpdatedAt: now,
					SwiperID:  1,
					SwipedID:  1,
					Type:      "like",
					Date:      now,
				},
			},
			wantErr: false,
			function: func(a args) {
				mockRepository.EXPECT().GetUserDetail(a.ctx, &a.swiperID, nil).Return(domain.User{
					ID:       1,
					Password: "a",
				}, nil)
				mockRepository.EXPECT().GetSwipe(a.ctx, swiperID, a.date).Return([]domain.Swipe{
					{
						ID:        1,
						CreatedAt: now,
						UpdatedAt: now,
						SwiperID:  1,
						SwipedID:  1,
						Type:      "like",
						Date:      now,
					},
				}, nil)
			},
		},
		{
			name: "error swipe",
			fields: fields{
				repository: mockRepository,
			},
			args: args{
				ctx:      context.Background(),
				swiperID: swiperID,
				date:     now,
			},
			wantSwipes: []domain.Swipe{},
			wantErr:    true,
			function: func(a args) {
				mockRepository.EXPECT().GetUserDetail(a.ctx, &a.swiperID, nil).Return(domain.User{
					ID:       1,
					Password: "a",
				}, nil)
				mockRepository.EXPECT().GetSwipe(a.ctx, swiperID, a.date).Return([]domain.Swipe{}, errors.New(""))
			},
		},
		{
			name: "error: get user detail",
			fields: fields{
				repository: mockRepository,
			},
			args: args{
				ctx:      context.Background(),
				swiperID: swiperID,
				date:     now,
			},
			wantSwipes: nil,
			wantErr:    true,
			function: func(a args) {
				mockRepository.EXPECT().GetUserDetail(a.ctx, &a.swiperID, nil).Return(domain.User{}, errors.New(""))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.function(tt.args)

			u := &Usecase{
				repository: tt.fields.repository,
			}
			gotSwipes, err := u.validateSwipe(tt.args.ctx, tt.args.swiperID, tt.args.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("Usecase.validateSwipe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSwipes, tt.wantSwipes) {
				t.Errorf("Usecase.validateSwipe() = %v, want %v", gotSwipes, tt.wantSwipes)
			}
		})
	}
}
