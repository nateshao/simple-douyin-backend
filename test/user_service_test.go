package test

import (
	"context"
	"testing"

	"github.com/YOJIA-yukino/simple-douyin-backend/internal/config"
	"github.com/YOJIA-yukino/simple-douyin-backend/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserService_Register(t *testing.T) {
	// 加载测试配置
	cfg, err := config.LoadConfig("../../configs/config.yaml")
	require.NoError(t, err)

	// 创建用户服务
	userService := service.NewUserService(cfg)

	ctx := context.Background()

	tests := []struct {
		name     string
		username string
		password string
		wantErr  bool
	}{
		{
			name:     "valid registration",
			username: "testuser1",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "username too short",
			username: "ab",
			password: "password123",
			wantErr:  true,
		},
		{
			name:     "password too short",
			username: "testuser2",
			password: "123",
			wantErr:  true,
		},
		{
			name:     "empty username",
			username: "",
			password: "password123",
			wantErr:  true,
		},
		{
			name:     "empty password",
			username: "testuser3",
			password: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := userService.Register(ctx, tt.username, tt.password)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.username, user.UserName)
				assert.NotEmpty(t, user.UserID)
				assert.NotEmpty(t, user.PassWord)
			}
		})
	}
}

func TestUserService_Login(t *testing.T) {
	cfg, err := config.LoadConfig("../../configs/config.yaml")
	require.NoError(t, err)

	userService := service.NewUserService(cfg)
	ctx := context.Background()

	// 先注册一个用户
	username := "logintest"
	password := "password123"

	registeredUser, err := userService.Register(ctx, username, password)
	require.NoError(t, err)
	require.NotNil(t, registeredUser)

	tests := []struct {
		name     string
		username string
		password string
		wantErr  bool
	}{
		{
			name:     "valid login",
			username: username,
			password: password,
			wantErr:  false,
		},
		{
			name:     "wrong password",
			username: username,
			password: "wrongpassword",
			wantErr:  true,
		},
		{
			name:     "user not exists",
			username: "nonexistent",
			password: password,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := userService.Login(ctx, tt.username, tt.password)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.username, user.UserName)
				assert.Equal(t, registeredUser.UserID, user.UserID)
			}
		})
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	cfg, err := config.LoadConfig("../../configs/config.yaml")
	require.NoError(t, err)

	userService := service.NewUserService(cfg)
	ctx := context.Background()

	// 先注册一个用户
	username := "getbyidtest"
	password := "password123"

	registeredUser, err := userService.Register(ctx, username, password)
	require.NoError(t, err)
	require.NotNil(t, registeredUser)

	tests := []struct {
		name    string
		userID  int64
		wantErr bool
	}{
		{
			name:    "valid user id",
			userID:  registeredUser.UserID,
			wantErr: false,
		},
		{
			name:    "invalid user id",
			userID:  999999999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := userService.GetUserByID(ctx, tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.userID, user.UserID)
				assert.Equal(t, username, user.UserName)
			}
		})
	}
}

func TestUserService_GetUserByUsername(t *testing.T) {
	cfg, err := config.LoadConfig("../../configs/config.yaml")
	require.NoError(t, err)

	userService := service.NewUserService(cfg)
	ctx := context.Background()

	// 先注册一个用户
	username := "getbyusernametest"
	password := "password123"

	registeredUser, err := userService.Register(ctx, username, password)
	require.NoError(t, err)
	require.NotNil(t, registeredUser)

	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{
			name:     "valid username",
			username: username,
			wantErr:  false,
		},
		{
			name:     "invalid username",
			username: "nonexistent",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := userService.GetUserByUsername(ctx, tt.username)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.username, user.UserName)
				assert.Equal(t, registeredUser.UserID, user.UserID)
			}
		})
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	cfg, err := config.LoadConfig("../../configs/config.yaml")
	require.NoError(t, err)

	userService := service.NewUserService(cfg)
	ctx := context.Background()

	// 先注册一个用户
	username := "updatetest"
	password := "password123"

	user, err := userService.Register(ctx, username, password)
	require.NoError(t, err)
	require.NotNil(t, user)

	// 更新用户信息
	user.FollowCount = 10
	user.FollowerCount = 20

	err = userService.UpdateUser(ctx, user)
	assert.NoError(t, err)

	// 验证更新是否成功
	updatedUser, err := userService.GetUserByID(ctx, user.UserID)
	assert.NoError(t, err)
	assert.Equal(t, int64(10), updatedUser.FollowCount)
	assert.Equal(t, int64(20), updatedUser.FollowerCount)
}

func TestUserService_DeleteUser(t *testing.T) {
	cfg, err := config.LoadConfig("../../configs/config.yaml")
	require.NoError(t, err)

	userService := service.NewUserService(cfg)
	ctx := context.Background()

	// 先注册一个用户
	username := "deletetest"
	password := "password123"

	user, err := userService.Register(ctx, username, password)
	require.NoError(t, err)
	require.NotNil(t, user)

	// 删除用户
	err = userService.DeleteUser(ctx, user.UserID)
	assert.NoError(t, err)

	// 验证用户是否已被删除
	deletedUser, err := userService.GetUserByID(ctx, user.UserID)
	assert.Error(t, err)
	assert.Nil(t, deletedUser)
}
