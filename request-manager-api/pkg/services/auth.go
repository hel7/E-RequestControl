package services

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	Request_Manager "request_manager_api"
	"request_manager_api/pkg/repository"
	"sync"
	"time"
)

var (
	signingKey string
)

const (
	tokenTTL = 1 * time.Hour
)

func init() {
	// Load secrets from environment variables
	signingKey = os.Getenv("JWT_SIGNING_KEY")
	if signingKey == "" {
		log.Fatal("JWT_SIGNING_KEY environment variable is not set")
	}
}

type tokenClaims struct {
	jwt.StandardClaims
	UserID int `json:"UserID"`
	RoleID int `json:"RoleID"`
}

type AuthService struct {
	repo      repository.Authorization
	blacklist map[string]time.Time
	mu        sync.Mutex
}

func (s *AuthService) IsTokenValid(token string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	expiry, exists := s.blacklist[token]
	if exists && time.Now().Before(expiry) {
		return false
	}
	return true
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{
		repo:      repo,
		blacklist: make(map[string]time.Time),
	}
}

func (s *AuthService) CreateUser(user Request_Manager.User) (int, error) {
	if err := user.ValidatePassword(); err != nil {
		return 0, err
	}
	if err := user.ValidateEmail(); err != nil {
		return 0, err
	}
	hashedPassword, err := generatePasswordHash(user.Password)
	if err != nil {
		return 0, err
	}
	user.Password = hashedPassword
	return s.repo.CreateUser(user)
}

func (s *AuthService) CreateAdmin(user Request_Manager.User) (int, error) {
	if err := user.ValidatePassword(); err != nil {
		return 0, err
	}
	if err := user.ValidateEmail(); err != nil {
		return 0, err
	}
	hashedPassword, err := generatePasswordHash(user.Password)
	if err != nil {
		return 0, err
	}
	user.Password = hashedPassword
	return s.repo.CreateAdmin(user)
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return "", err
	}

	if err := comparePassword(user.Password, password); err != nil {
		return "", errors.New("invalid credentials")
	}

	claims := tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserID: user.UserID,
		RoleID: user.RoleID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) InvalidateToken(token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	claims := &tokenClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(signingKey), nil
	})

	if err != nil {
		return err
	}

	s.blacklist[token] = time.Unix(claims.ExpiresAt, 0)
	return nil
}

func (s *AuthService) CleanupBlacklist() {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		for token, expiry := range s.blacklist {
			if time.Now().After(expiry) {
				delete(s.blacklist, token)
			}
		}
		s.mu.Unlock()
	}
}

func (s *AuthService) ParseToken(accessToken string) (int, int, error) {
	if !s.IsTokenValid(accessToken) {
		return 0, 0, errors.New("token is invalidated")
	}
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserID, claims.RoleID, nil
}

func generatePasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
func comparePassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
