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

func TestRepository_UpsertUser(t *testing.T) {
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
		data domain.User
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
				postgres: mockPostgres,
				redis:    mockRedis,
			},
			args: args{
				ctx: context.Background(),
				data: domain.User{
					ID:        1,
					CreatedAt: now,
					UpdatedAt: now,
					Name:      "a",
					Email:     "a@a.com",
					Password:  "a",
					IsPremium: false,
				},
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
				mockPostgres.EXPECT().UpsertUser(a.ctx, &a.data).Return(nil)
				mockRedis.EXPECT().SetUserDetail(a.ctx, &a.data.ID, nil, a.data).Return(nil)
				mockRedis.EXPECT().SetUserDetail(a.ctx, nil, &a.data.Email, a.data).Return(nil)
			},
		},
		{
			name: "upsert user error",
			fields: fields{
				postgres: mockPostgres,
				redis:    mockRedis,
			},
			args: args{
				ctx: context.Background(),
				data: domain.User{
					ID:        1,
					CreatedAt: now,
					UpdatedAt: now,
					Name:      "a",
					Email:     "a@a.com",
					Password:  "a",
					IsPremium: false,
				},
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
			wantErr: true,
			function: func(a args) {
				mockPostgres.EXPECT().UpsertUser(a.ctx, &a.data).Return(errors.New(""))
			},
		},
		{
			name: "upsert user redis by email error",
			fields: fields{
				postgres: mockPostgres,
				redis:    mockRedis,
			},
			args: args{
				ctx: context.Background(),
				data: domain.User{
					ID:        1,
					CreatedAt: now,
					UpdatedAt: now,
					Name:      "a",
					Email:     "a@a.com",
					Password:  "a",
					IsPremium: false,
				},
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
			wantErr: true,
			function: func(a args) {
				mockPostgres.EXPECT().UpsertUser(a.ctx, &a.data).Return(nil)
				mockRedis.EXPECT().SetUserDetail(a.ctx, &a.data.ID, nil, a.data).Return(nil)
				mockRedis.EXPECT().SetUserDetail(a.ctx, nil, &a.data.Email, a.data).Return(errors.New(""))
			},
		},
		{
			name: "upsert user redis by id error",
			fields: fields{
				postgres: mockPostgres,
				redis:    mockRedis,
			},
			args: args{
				ctx: context.Background(),
				data: domain.User{
					ID:        1,
					CreatedAt: now,
					UpdatedAt: now,
					Name:      "a",
					Email:     "a@a.com",
					Password:  "a",
					IsPremium: false,
				},
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
			wantErr: true,
			function: func(a args) {
				mockPostgres.EXPECT().UpsertUser(a.ctx, &a.data).Return(nil)
				mockRedis.EXPECT().SetUserDetail(a.ctx, &a.data.ID, nil, a.data).Return(errors.New(""))
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
			got, err := r.UpsertUser(tt.args.ctx, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.UpsertUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.UpsertUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_GetUserProfile(t *testing.T) {
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
		ctx   context.Context
		notIn []uint
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		want     []domain.User
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
				ctx:   context.Background(),
				notIn: []uint{1},
			},
			want: []domain.User{
				{
					ID:        2,
					CreatedAt: now,
					UpdatedAt: now,
					Name:      "a",
					Email:     "a@a.com",
					Password:  "a",
					IsPremium: false,
				},
			},
			wantErr: false,
			function: func(a args) {
				mockPostgres.EXPECT().GetUserProfile(a.ctx, a.notIn).Return([]domain.User{
					{
						ID:        2,
						CreatedAt: now,
						UpdatedAt: now,
						Name:      "a",
						Email:     "a@a.com",
						Password:  "a",
						IsPremium: false,
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
				ctx:   context.Background(),
				notIn: []uint{1},
			},
			want:    []domain.User{},
			wantErr: true,
			function: func(a args) {
				mockPostgres.EXPECT().GetUserProfile(a.ctx, a.notIn).Return([]domain.User{}, errors.New(""))
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
			got, err := r.GetUserProfile(tt.args.ctx, tt.args.notIn)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetUserProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.GetUserProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_GetUserDetail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostgres := mocks_postgres_repository.NewMockIPostgresRepository(ctrl)
	mockRedis := mocks_redis_respository.NewMockIRedisRepository(ctrl)

	now := time.Now()
	id := uint(1)
	email := "a@a.com"

	type fields struct {
		postgres postgresRepo.IPostgresRepository
		redis    redisRepo.IRedisRepository
	}
	type args struct {
		ctx   context.Context
		id    *uint
		email *string
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
			name: "success: redis",
			fields: fields{
				postgres: mockPostgres,
				redis:    mockRedis,
			},
			args: args{
				ctx:   context.Background(),
				id:    &id,
				email: &email,
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
				mockRedis.EXPECT().GetUserDetail(a.ctx, a.id, a.email).Return(domain.User{
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
			name: "success: database",
			fields: fields{
				postgres: mockPostgres,
				redis:    mockRedis,
			},
			args: args{
				ctx:   context.Background(),
				id:    &id,
				email: &email,
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
				mockRedis.EXPECT().GetUserDetail(a.ctx, a.id, a.email).Return(domain.User{}, nil)
				data := domain.User{
					ID:        1,
					CreatedAt: now,
					UpdatedAt: now,
					Name:      "a",
					Email:     "a@a.com",
					Password:  "a",
					IsPremium: false,
				}
				mockPostgres.EXPECT().GetUserDetail(a.ctx, a.id, a.email).Return(data, nil)
				mockRedis.EXPECT().SetUserDetail(a.ctx, &data.ID, nil, data).Return(nil)
				mockRedis.EXPECT().SetUserDetail(a.ctx, nil, &data.Email, data).Return(nil)
			},
		},
		{
			name: "error cache",
			fields: fields{
				postgres: mockPostgres,
				redis:    mockRedis,
			},
			args: args{
				ctx:   context.Background(),
				id:    &id,
				email: &email,
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
			wantErr: true,
			function: func(a args) {
				mockRedis.EXPECT().GetUserDetail(a.ctx, a.id, a.email).Return(domain.User{}, nil)
				data := domain.User{
					ID:        1,
					CreatedAt: now,
					UpdatedAt: now,
					Name:      "a",
					Email:     "a@a.com",
					Password:  "a",
					IsPremium: false,
				}
				mockPostgres.EXPECT().GetUserDetail(a.ctx, a.id, a.email).Return(data, nil)
				mockRedis.EXPECT().SetUserDetail(a.ctx, &data.ID, nil, data).Return(nil)
				mockRedis.EXPECT().SetUserDetail(a.ctx, nil, &data.Email, data).Return(errors.New(""))
			},
		},
		{
			name: "error database",
			fields: fields{
				postgres: mockPostgres,
				redis:    mockRedis,
			},
			args: args{
				ctx:   context.Background(),
				id:    &id,
				email: &email,
			},
			want:    domain.User{},
			wantErr: true,
			function: func(a args) {
				mockRedis.EXPECT().GetUserDetail(a.ctx, a.id, a.email).Return(domain.User{}, nil)
				mockPostgres.EXPECT().GetUserDetail(a.ctx, a.id, a.email).Return(domain.User{}, errors.New(""))
			},
		},
		{
			name: "error cache",
			fields: fields{
				postgres: mockPostgres,
				redis:    mockRedis,
			},
			args: args{
				ctx:   context.Background(),
				id:    &id,
				email: &email,
			},
			want:    domain.User{},
			wantErr: true,
			function: func(a args) {
				mockRedis.EXPECT().GetUserDetail(a.ctx, a.id, a.email).Return(domain.User{}, errors.New(""))
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
			got, err := r.GetUserDetail(tt.args.ctx, tt.args.id, tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetUserDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.GetUserDetail() = %v, want %v", got, tt.want)
			}
		})
	}
}
