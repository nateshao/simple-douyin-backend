package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/YOJIA-yukino/simple-douyin-backend/internal/cache"
	"github.com/YOJIA-yukino/simple-douyin-backend/internal/config"
	"github.com/YOJIA-yukino/simple-douyin-backend/internal/database"
	"github.com/YOJIA-yukino/simple-douyin-backend/internal/model"
	"github.com/YOJIA-yukino/simple-douyin-backend/internal/utils"
	"github.com/YOJIA-yukino/simple-douyin-backend/internal/utils/logger"
	"gorm.io/gorm"
)

// UserService 用户服务接口
type UserService interface {
	Register(ctx context.Context, username, password string) (*model.User, error)
	Login(ctx context.Context, username, password string) (*model.User, error)
	GetUserByID(ctx context.Context, userID int64) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, userID int64) error
}

// userServiceImpl 用户服务实现
type userServiceImpl struct {
	db         *gorm.DB
	cache      cache.Cache
	hasher     utils.PasswordHasher
	jwtManager *utils.JWTManager
	config     *config.Config
	userCache  *UserCache
}

// UserCache 用户缓存
type UserCache struct {
	cache map[string]*model.User
	mu    sync.RWMutex
	ttl   time.Duration
}

// NewUserCache 创建用户缓存
func NewUserCache(ttl time.Duration) *UserCache {
	return &UserCache{
		cache: make(map[string]*model.User),
		ttl:   ttl,
	}
}

// NewUserService 创建用户服务
func NewUserService(cfg *config.Config) UserService {
	hasher := utils.NewBCryptHasher(cfg.Security.BCryptCost)
	jwtManager := utils.NewJWTManager(cfg.JWT.Secret, time.Duration(cfg.JWT.ExpireTime)*time.Hour)

	return &userServiceImpl{
		db:         database.GetDB(),
		cache:      cache.GetRedisClient(),
		hasher:     hasher,
		jwtManager: jwtManager,
		config:     cfg,
		userCache:  NewUserCache(30 * time.Minute),
	}
}

// Register 用户注册
func (s *userServiceImpl) Register(ctx context.Context, username, password string) (*model.User, error) {
	// 输入验证
	if err := utils.ValidateUsername(username); err != nil {
		return nil, fmt.Errorf("invalid username: %w", err)
	}
	if err := utils.ValidatePassword(password); err != nil {
		return nil, fmt.Errorf("invalid password: %w", err)
	}

	// 检查用户是否已存在
	existingUser, err := s.GetUserByUsername(ctx, username)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user already exists")
	}

	// 哈希密码
	hashedPassword, err := s.hasher.Hash(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 创建用户
	user := &model.User{
		UserID:        generateUserID(),
		UserName:      username,
		PassWord:      hashedPassword,
		FollowCount:   0,
		FollowerCount: 0,
	}

	// 保存到数据库
	if err := s.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// 缓存用户信息
	s.cacheUser(ctx, user)

	logger.GlobalLogger.Printf("User registered successfully: %s", username)
	return user, nil
}

// Login 用户登录
func (s *userServiceImpl) Login(ctx context.Context, username, password string) (*model.User, error) {
	// 输入验证
	if err := utils.ValidateUsername(username); err != nil {
		return nil, fmt.Errorf("invalid username: %w", err)
	}

	// 获取用户信息
	user, err := s.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// 验证密码
	if !s.hasher.Verify(password, user.PassWord) {
		return nil, fmt.Errorf("invalid password")
	}

	logger.GlobalLogger.Printf("User logged in successfully: %s", username)
	return user, nil
}

// GetUserByID 根据ID获取用户
func (s *userServiceImpl) GetUserByID(ctx context.Context, userID int64) (*model.User, error) {
	// 先从缓存获取
	cacheKey := fmt.Sprintf("user:%d", userID)
	if user := s.getUserFromCache(cacheKey); user != nil {
		return user, nil
	}

	// 从数据库获取
	var user model.User
	if err := s.db.WithContext(ctx).Where("user_id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// 缓存用户信息
	s.cacheUser(ctx, &user)

	return &user, nil
}

// GetUserByUsername 根据用户名获取用户
func (s *userServiceImpl) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	// 先从缓存获取
	cacheKey := fmt.Sprintf("user:username:%s", username)
	if user := s.getUserFromCache(cacheKey); user != nil {
		return user, nil
	}

	// 从数据库获取
	var user model.User
	if err := s.db.WithContext(ctx).Where("user_name = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// 缓存用户信息
	s.cacheUser(ctx, &user)

	return &user, nil
}

// UpdateUser 更新用户信息
func (s *userServiceImpl) UpdateUser(ctx context.Context, user *model.User) error {
	if err := s.db.WithContext(ctx).Save(user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// 清除缓存
	s.clearUserCache(user.UserID, user.UserName)

	logger.GlobalLogger.Printf("User updated successfully: %d", user.UserID)
	return nil
}

// DeleteUser 删除用户
func (s *userServiceImpl) DeleteUser(ctx context.Context, userID int64) error {
	// 先获取用户信息
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	// 软删除
	if err := s.db.WithContext(ctx).Delete(user).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	// 清除缓存
	s.clearUserCache(userID, user.UserName)

	logger.GlobalLogger.Printf("User deleted successfully: %d", userID)
	return nil
}

// cacheUser 缓存用户信息
func (s *userServiceImpl) cacheUser(ctx context.Context, user *model.User) {
	// 缓存到Redis
	userData := map[string]interface{}{
		"user_id":        user.UserID,
		"username":       user.UserName,
		"follow_count":   user.FollowCount,
		"follower_count": user.FollowerCount,
	}

	cacheKey := fmt.Sprintf("user:%d", user.UserID)
	if err := s.cache.HMSet(ctx, cacheKey, userData); err != nil {
		logger.GlobalLogger.Printf("Failed to cache user: %v", err)
		return
	}

	// 设置过期时间
	s.cache.Expire(ctx, cacheKey, 30*time.Minute)

	// 缓存用户名映射
	usernameKey := fmt.Sprintf("user:username:%s", user.UserName)
	s.cache.SetWithExpiration(ctx, usernameKey, user.UserID, 30*time.Minute)
}

// getUserFromCache 从缓存获取用户
func (s *userServiceImpl) getUserFromCache(cacheKey string) *model.User {
	s.userCache.mu.RLock()
	defer s.userCache.mu.RUnlock()

	if user, exists := s.userCache.cache[cacheKey]; exists {
		return user
	}
	return nil
}

// clearUserCache 清除用户缓存
func (s *userServiceImpl) clearUserCache(userID int64, username string) {
	ctx := context.Background()

	// 清除Redis缓存
	cacheKey := fmt.Sprintf("user:%d", userID)
	usernameKey := fmt.Sprintf("user:username:%s", username)

	s.cache.Del(ctx, cacheKey, usernameKey)

	// 清除内存缓存
	s.userCache.mu.Lock()
	defer s.userCache.mu.Unlock()

	delete(s.userCache.cache, cacheKey)
	delete(s.userCache.cache, usernameKey)
}

// generateUserID 生成用户ID
func generateUserID() int64 {
	return time.Now().UnixNano()
}
