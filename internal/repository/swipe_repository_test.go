package repository

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/fahmialfareza/deals-dating-app-backend/internal/domain"
	postgresRepo "github.com/fahmialfareza/deals-dating-app-backend/internal/repository/postgres"
	mocks_postgres_repository "github.com/fahmialfareza/deals-dating-app-backend/internal/repository/postgres/mocks"
	redisRepo "github.com/fahmialfareza/deals-dating-app-backend/internal/repository/redis"
	mocks_redis_respository "github.com/fahmialfareza/deals-dating-app-backend/internal/repository/redis/mocks"
	"github.com/golang/mock/gomock"
)

func TestRepository_CreateSwipe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostgres := mocks_postgres_repository.NewMockIPostgresRepository(ctrl)
	mockRedis := mocks_redis_respository.NewMockIRedisRepository(ctrl)

	now := time.Now()

	type fields struct {
		postgres postgresRepo.IPostgresRepository
		redis    redisRepo.IRedisRepository
	}
	type args struct {
		ctx  context.Context
		data *domain.Swipe
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantErr  bool
		function func(args)
	}{
		{
			name: "success",
			fields: fields{
				postgres: mockPostgres,
				redis:    mockRedis,
			},
			args: args{
				ctx: context.Background(),
				data: &domain.Swipe{
					SwiperID: 1,
					SwipedID: 2,
					Type:     "like",
					Date:     now,
				},
			},
			wantErr: false,
			function: func(a args) {
				mockPostgres.EXPECT().InsertSwipe(a.ctx, a.data).Return(nil)
			},
		},
		{
			name: "error",
			fields: fields{
				postgres: mockPostgres,
				redis:    mockRedis,
			},
			args: args{
				ctx: context.Background(),
				data: &domain.Swipe{
					SwiperID: 1,
					SwipedID: 2,
					Type:     "like",
					Date:     now,
				},
			},
			wantErr: true,
			function: func(a args) {
				mockPostgres.EXPECT().InsertSwipe(a.ctx, a.data).Return(errors.New(""))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.function(tt.args)

			r := &Repository{
				postgres: tt.fields.postgres,
				redis:    tt.fields.redis,
			}
			if err := r.CreateSwipe(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Repository.CreateSwipe() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_GetSwipe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostgres := mocks_postgres_repository.NewMockIPostgresRepository(ctrl)
	mockRedis := mocks_redis_respository.NewMockIRedisRepository(ctrl)

	now := time.Now()

	type fields struct {
		postgres postgresRepo.IPostgresRepository
		redis    redisRepo.IRedisRepository
	}
	type args struct {
		ctx      context.Context
		swiperID uint
		date     time.Time
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		want     []domain.Swipe
		wantErr  bool
		function func(args)
	}{
		{
			name: "success",
			fields: fields{
				postgres: mockPostgres,
				redis:    mockRedis,
			},
			args: args{
				ctx:      context.Background(),
				swiperID: 1,
				date:     now,
			},
			want: []domain.Swipe{
				{
					ID:        1,
					SwiperID:  1,
					SwipedID:  1,
					Type:      "like",
					CreatedAt: now,
					UpdatedAt: now,
					Date:      now,
				},
			},
			wantErr: false,
			function: func(a args) {
				mockPostgres.EXPECT().GetSwipe(a.ctx, a.swiperID, a.date).Return([]domain.Swipe{
					{
						ID:        1,
						SwiperID:  1,
						SwipedID:  1,
						Type:      "like",
						CreatedAt: now,
						UpdatedAt: now,
						Date:      now,
					},
				}, nil)
			},
		},
		{
			name: "error",
			fields: fields{
				postgres: mockPostgres,
				redis:    mockRedis,
			},
			args: args{
				ctx:      context.Background(),
				swiperID: 1,
				date:     now,
			},
			want:    []domain.Swipe{},
			wantErr: true,
			function: func(a args) {
				mockPostgres.EXPECT().GetSwipe(a.ctx, a.swiperID, a.date).Return([]domain.Swipe{}, errors.New(""))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.function(tt.args)

			r := &Repository{
				postgres: tt.fields.postgres,
				redis:    tt.fields.redis,
			}
			got, err := r.GetSwipe(tt.args.ctx, tt.args.swiperID, tt.args.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetSwipe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.GetSwipe() = %v, want %v", got, tt.want)
			}
		})
	}
}
