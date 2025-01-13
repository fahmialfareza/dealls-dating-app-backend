package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/fahmialfareza/deals-dating-app-backend/configs/cache/redis"
	mocks_redis "github.com/fahmialfareza/deals-dating-app-backend/configs/cache/redis/mocks"
	"github.com/fahmialfareza/deals-dating-app-backend/configs/environment"
	"github.com/fahmialfareza/deals-dating-app-backend/internal/domain"
	"github.com/golang/mock/gomock"
	goRedis "github.com/redis/go-redis/v9"
)

func TestRedisRepository_DeleteUserDetail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := mocks_redis.NewMockIRedis(ctrl)

	id := uint(1)
	email := "a@a.com"

	type fields struct {
		redis redis.IRedis
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
		wantErr  bool
		function func(args)
	}{
		{
			name: "success: id",
			fields: fields{
				redis: mockRedis,
			},
			args: args{
				ctx:   context.Background(),
				id:    &id,
				email: &email,
			},
			wantErr: false,
			function: func(a args) {
				key := fmt.Sprintf(userDetail, *a.id)
				mockRedis.EXPECT().Del(a.ctx, key).Return(nil)
			},
		},
		{
			name: "success: email",
			fields: fields{
				redis: mockRedis,
			},
			args: args{
				ctx:   context.Background(),
				id:    nil,
				email: &email,
			},
			wantErr: false,
			function: func(a args) {
				key := fmt.Sprintf(userDetail, *a.email)
				mockRedis.EXPECT().Del(a.ctx, key).Return(nil)
			},
		},
		{
			name: "error",
			fields: fields{
				redis: mockRedis,
			},
			args: args{
				ctx:   context.Background(),
				id:    nil,
				email: &email,
			},
			wantErr: true,
			function: func(a args) {
				key := fmt.Sprintf(userDetail, *a.email)
				mockRedis.EXPECT().Del(a.ctx, key).Return(errors.New(""))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.function(tt.args)

			r := &RedisRepository{
				redis: tt.fields.redis,
			}
			if err := r.DeleteUserDetail(tt.args.ctx, tt.args.id, tt.args.email); (err != nil) != tt.wantErr {
				t.Errorf("RedisRepository.DeleteUserDetail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRedisRepository_GetUserDetail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := mocks_redis.NewMockIRedis(ctrl)
	id := uint(1)
	email := "a@a.com"

	type fields struct {
		redis redis.IRedis
	}
	type args struct {
		ctx   context.Context
		id    *uint
		email *string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult domain.User
		wantErr    bool
		function   func(args)
	}{
		{
			name: "error",
			fields: fields{
				redis: mockRedis,
			},
			args: args{
				ctx:   context.Background(),
				id:    &id,
				email: &email,
			},
			wantResult: domain.User{},
			wantErr:    true,
			function: func(a args) {
				key := fmt.Sprintf(userDetail, *a.id)
				mockRedis.EXPECT().Get(a.ctx, key).Return("", errors.New(""))
			},
		},
		{
			name: "error redis nil",
			fields: fields{
				redis: mockRedis,
			},
			args: args{
				ctx:   context.Background(),
				id:    nil,
				email: &email,
			},
			wantResult: domain.User{},
			wantErr:    false,
			function: func(a args) {
				key := fmt.Sprintf(userDetail, *a.email)
				mockRedis.EXPECT().Get(a.ctx, key).Return("", goRedis.Nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.function(tt.args)

			r := &RedisRepository{
				redis: tt.fields.redis,
			}
			gotResult, err := r.GetUserDetail(tt.args.ctx, tt.args.id, tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisRepository.GetUserDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("RedisRepository.GetUserDetail() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestRedisRepository_SetUserDetail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := mocks_redis.NewMockIRedis(ctrl)
	id := uint(1)
	email := "a@a.com"
	now := time.Now()
	data := domain.User{
		ID:        1,
		CreatedAt: now,
		UpdatedAt: now,
		Name:      "a",
		Email:     "a@a.com",
		Password:  "a",
		IsPremium: false,
	}

	type fields struct {
		redis redis.IRedis
	}
	type args struct {
		ctx   context.Context
		id    *uint
		email *string
		data  domain.User
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
				redis: mockRedis,
			},
			args: args{
				ctx:   context.Background(),
				id:    &id,
				email: &email,
				data:  data,
			},
			wantErr: false,
			function: func(a args) {
				key := fmt.Sprintf(userDetail, *a.id)
				dataBytes, _ := json.Marshal(data)
				mockRedis.EXPECT().Set(a.ctx, key, string(dataBytes), environment.RedisExpireTime).Return(nil)
			},
		},
		{
			name: "error",
			fields: fields{
				redis: mockRedis,
			},
			args: args{
				ctx:   context.Background(),
				id:    nil,
				email: &email,
				data:  data,
			},
			wantErr: true,
			function: func(a args) {
				key := fmt.Sprintf(userDetail, *a.email)
				dataBytes, _ := json.Marshal(data)
				mockRedis.EXPECT().Set(a.ctx, key, string(dataBytes), environment.RedisExpireTime).Return(errors.New(""))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.function(tt.args)

			r := &RedisRepository{
				redis: tt.fields.redis,
			}
			if err := r.SetUserDetail(tt.args.ctx, tt.args.id, tt.args.email, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("RedisRepository.SetUserDetail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
