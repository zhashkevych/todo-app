package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"
	"todo-app"
	"todo-app/pkg/repository"

	"github.com/dgrijalva/jwt-go"
)

const tokenTTL = 30 * time.Hour

const JWT_SECRET = "rkjk#4#%35FSFJlja#4353KSFjH"
const SOLT = "hjqrhjqw124617ajfhajs"

type AuthService struct {
	repo repository.Authorization
}

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user todo.User) (int, error) {
	user.Password = generatePasswordHash(user.Password) // Перезаписываем пароль на хэшированный
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(username, password string) (string, error) { // Генерация токена по имени и паролю
	user, err := s.repo.GetUser(username, generatePasswordHash(password)) // поиск пользователя в БД
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{ // генерация токена
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(), // время действия токена
			IssuedAt:  time.Now().Unix(),               //время создания
		},
		user.Id,
	})

	return token.SignedString([]byte(JWT_SECRET))
}

func (s *AuthService) ParseToken(accesstoken string) (int, error) { //Парс токена (получаем из токена id)
	token, err := jwt.ParseWithClaims(accesstoken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(JWT_SECRET), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}

func generatePasswordHash(password string) string { // Хэширование пароля
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(SOLT)))
}
