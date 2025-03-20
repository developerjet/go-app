package services

import (
	"errors"
	"fmt"
	"time"
	"go_app/models"
	"go_app/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

// 保留这个已实现的版本
func (s *UserService) VerifyPassword(hashedPassword, password string) error {
    return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) CreateUser(user *models.User) error {
	return s.db.Create(user).Error
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := s.db.First(&user, id).Error
	return &user, err
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := s.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (s *UserService) UpdateUser(user *models.User) error {
	result := s.db.Save(user)
	return result.Error
}

func (s *UserService) DeleteUser(id uint) error {
	return s.db.Delete(&models.User{}, id).Error
}

func (s *UserService) ListUsers() ([]models.User, error) {
	var users []models.User
	err := s.db.Find(&users).Error
	return users, err
}

// HashPassword 密码加密
func (s *UserService) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

// Login 登录
func (s *UserService) Login(email, password string) (*models.User, string, error) {
    var user models.User
    if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
        return nil, "", fmt.Errorf("用户不存在")
    }

    if err := s.VerifyPassword(user.Password, password); err != nil {
        return nil, "", fmt.Errorf("密码错误")
    }

    // 生成新token
    token, err := utils.GenerateToken(user.ID)
    if err != nil {
        return nil, "", err
    }

    // 删除该用户的旧token
    if err := s.db.Where("user_id = ?", user.ID).Delete(&models.UserToken{}).Error; err != nil {
        return nil, "", err
    }

    // 创建新token记录
    userToken := &models.UserToken{
        UserID:    user.ID,
        Token:     token,
        ExpiredAt: time.Now().Add(24 * time.Hour),
    }

    if err := s.db.Create(userToken).Error; err != nil {
        return nil, "", err
    }

    user.Token = token
    return &user, token, nil
}

// GetUserByIDSafe 安全地获取用户信息，不返回敏感字段
func (s *UserService) GetUserByIDSafe(id uint) (*models.User, error) {
    var user models.User
    err := s.db.Select("id, username, email, created_at, updated_at").First(&user, id).Error
    if err != nil {
        return nil, errors.New("用户不存在")
    }
    return &user, nil
}

// ListUsersSafe 安全地获取用户列表，不返回敏感字段
func (s *UserService) ListUsersSafe() ([]models.User, error) {
    var users []models.User
    err := s.db.Select("id, username, email, created_at, updated_at").Find(&users).Error
    if err != nil {
        return nil, errors.New("获取用户列表失败")
    }
    return users, err
}

// UpdateUserSafe 安全地更新用户信息，只允许更新指定字段
func (s *UserService) UpdateUserSafe(user *models.User) error {
    result := s.db.Model(user).Select("username", "email").Updates(user)
    if result.Error != nil {
        return errors.New("更新用户信息失败")
    }
    return nil
}

// ListUsers 获取用户列表
// ListUsersDetail 获取用户列表(不含敏感数据)
func (s *UserService) ListUsersDetail() ([]models.User, error) {
    var users []models.User
    err := s.db.Select("id, username, email, created_at, updated_at").Find(&users).Error
    if err != nil {
        return nil, errors.New("获取用户列表失败")
    }
    return users, err
}


// 保留新的 UpdateUser 方法，这个实现更安全，只允许更新特定字段
// UpdateUserInfo 更新用户基本信息，只允许更新用户名和邮箱
func (s *UserService) UpdateUserInfo(user *models.User) error {
    result := s.db.Model(user).Select("username", "email").Updates(user)
    if result.Error != nil {
        return errors.New("更新用户信息失败")
    }
    return nil
}

// ChangePassword 修改用户密码
func (s *UserService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return errors.New("用户不存在")
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("旧密码不正确")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("密码加密失败")
	}

	// 更新密码
	if err := s.db.Model(&user).Update("password", string(hashedPassword)).Error; err != nil {
		return errors.New("密码更新失败")
	}

	return nil
}

// Logout 用户退出
func (s *UserService) Logout(userID uint) error {
    // 这里可以添加一些退出时的清理工作
    // 比如：清除用户的token记录、更新最后登出时间等
    
    // 如果使用 Redis 存储 token，可以在这里删除
    // 目前简单返回成功，因为 JWT 是无状态的
    return nil
}
