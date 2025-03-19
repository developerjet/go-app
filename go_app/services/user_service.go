package services

import (
    "time"
    "go_app/models"
    "go_app/middleware"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
    "github.com/golang-jwt/jwt/v4"
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

// 在 UserService 结构体中添加以下方法
func (s *UserService) GenerateToken(userID uint) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(time.Hour * 24).Unix(),
    })
    return token.SignedString([]byte(middleware.JWTSecret))
}