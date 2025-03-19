package services

import (
    "errors"
    "go_app/models"
    "go_app/utils"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
)

type UserService struct {
    db *gorm.DB
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
    return s.db.Save(user).Error
}

func (s *UserService) DeleteUser(id uint) error {
    return s.db.Delete(&models.User{}, id).Error
}

func (s *UserService) ListUsers() ([]models.User, error) {
    var users []models.User
    err := s.db.Find(&users).Error
    return users, err
}

func (s *UserService) HashPassword(password string) (string, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(hashedPassword), err
}

func (s *UserService) ComparePasswords(hashedPassword, password string) error {
    return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (s *UserService) Login(email, password string) (*models.LoginResponse, error) {
    var user models.User
    if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
        return nil, errors.New("用户不存在")
    }

    // 验证密码
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return nil, errors.New("密码不正确")
    }

    // 生成 Token
    token, err := utils.GenerateToken(user.ID)
    if err != nil {
        return nil, err
    }

    return &models.LoginResponse{
        Token:    token,
        UserInfo: &user,  // 修改这里，传递指针
    }, nil
}