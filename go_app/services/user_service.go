package services

import (
	"errors"
	"fmt"
	"go_app/models"
	"go_app/utils"
	"mime/multipart"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

// SaveAvatar 保存用户头像到图床并更新数据库
func (s *UserService) SaveAvatar(userID uint, file *multipart.FileHeader) (string, error) {
    // 打开文件
    src, err := file.Open()
    if err != nil {
        return "", fmt.Errorf("打开文件失败: %v", err)
    }
    defer src.Close()

    // 调用图床 API 上传图片
    imageURL, err := utils.UploadToImageHost(file)
    if err != nil {
        return "", fmt.Errorf("上传图片失败: %v", err)
    }

    // 更新用户头像URL
    var user models.User
    if err := s.db.First(&user, userID).Error; err != nil {
        return "", fmt.Errorf("查找用户失败: %v", err)
    }

    user.AvatarURL = imageURL
    if err := s.db.Save(&user).Error; err != nil {
        return "", fmt.Errorf("更新用户头像失败: %v", err)
    }

    return imageURL, nil
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
    // 确保查询包含所有必要字段
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
// GenerateToken 生成新的 token
func (s *UserService) GenerateToken(user *models.User) (string, error) {
	// 增加 token 版本号
	user.TokenVersion++
	if err := s.db.Save(user).Error; err != nil {
		return "", err
	}

	// 生成包含版本号的 token
	token, err := utils.GenerateToken(user.ID, user.TokenVersion)
	if err != nil {
		return "", err
	}

	return token, nil
}

// Login 用户登录
func (s *UserService) Login(email, password string) (*models.User, string, error) {
	var user models.User
	// 先查询所有必要字段，包括密码
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, "", errors.New("用户不存在")
	}

	// 验证密码
	if err := s.VerifyPassword(user.Password, password); err != nil {
		return nil, "", errors.New("密码错误")
	}

	// 生成新的 token
	token, err := s.GenerateToken(&user)
	if err != nil {
		return nil, "", fmt.Errorf("生成token失败: %v", err)
	}

	// 将 token 设置到用户对象中，并清除敏感信息
	user.Token = token
	user.Password = ""

	// 重新查询需要返回的字段
	var safeUser models.User
	if err := s.db.Select("id, username, email, avatar_url, created_at, updated_at").
		First(&safeUser, user.ID).Error; err != nil {
		return nil, "", fmt.Errorf("获取用户信息失败: %v", err)
	}
	safeUser.Token = token

	return &safeUser, token, nil
}

// GetUserByIDSafe 安全地获取用户信息，不返回敏感字段
func (s *UserService) GetUserByIDSafe(id uint) (*models.User, error) {
    var user models.User
    err := s.db.Select("id, username, email, avatar_url, created_at, updated_at").First(&user, id).Error
    if err != nil {
        return nil, errors.New("用户不存在")
    }
    return &user, nil
}

// ListUsersSafe 安全地获取用户列表，不返回敏感字段
func (s *UserService) ListUsersSafe() ([]models.User, error) {
    var users []models.User
    err := s.db.Select("id, username, email, avatar_url, created_at, updated_at").Find(&users).Error
    if err != nil {
        return nil, errors.New("获取用户列表失败")
    }
    return users, err
}

// ListUsersWithPage 分页获取用户列表
func (s *UserService) ListUsersWithPage(page, pageSize int) ([]*models.User, int64, error) {
    var users []*models.User
    var total int64

    // 计算偏移量
    offset := (page - 1) * pageSize

    // 获取总记录数
    if err := s.db.Model(&models.User{}).Count(&total).Error; err != nil {
        return nil, 0, err
    }

    // 获取分页数据，添加 avatar_url 字段
    if err := s.db.Select("id, username, email, avatar_url, created_at, updated_at").
        Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
        return nil, 0, err
    }

    return users, total, nil
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
    // 确保查询包含所有字段，特别是密码字段
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
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return errors.New("用户不存在")
	}

	// 增加 token 版本号使当前 token 失效
	user.TokenVersion++
	if err := s.db.Save(&user).Error; err != nil {
		return errors.New("登出失败")
	}

	return nil
}

// 添加新方法
func (s *UserService) IsEmailExists(email string) bool {
	var count int64
	s.db.Model(&models.User{}).Where("email = ?", email).Count(&count)
	return count > 0
}
