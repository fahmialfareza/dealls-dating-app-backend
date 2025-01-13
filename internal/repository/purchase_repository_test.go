package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/fahmialfareza/deals-dating-app-backend/internal/domain"
	postgresRepo "github.com/fahmialfareza/deals-dating-app-backend/internal/repository/postgres"
	mocks_postgres_repository "github.com/fahmialfareza/deals-dating-app-backend/internal/repository/postgres/mocks"
	redisRepo "github.com/fahmialfareza/deals-dating-app-backend/internal/repository/redis"
	mocks_redis_respository "github.com/fahmialfareza/deals-dating-app-backend/internal/repository/redis/mocks"
	"github.com/golang/mock/gomock"
)

func TestRepository_CreatePurchase(t *testing.T) {
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
		data *domain.Purchase
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantErr  bool
		function func(args args)
	}{
		{
			name: "success",
			fields: fields{
				postgres: mockPostgres,
				redis:    mockRedis,
			},
			args: args{
				ctx: context.Background(),
				data: &domain.Purchase{
					UserID:       1,
					PackageType:  "premium",
					PurchaseDate: now,
				},
			},
			wantErr: false,
			function: func(args args) {
				mockPostgres.EXPECT().InsertPurchase(args.ctx, args.data).Return(nil)
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
				data: &domain.Purchase{
					UserID:       1,
					PackageType:  "premium",
					PurchaseDate: now,
				},
			},
			wantErr: true,
			function: func(args args) {
				mockPostgres.EXPECT().InsertPurchase(args.ctx, args.data).Return(errors.New(""))
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
			if err := r.CreatePurchase(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Repository.CreatePurchase() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
