package services

import "golang.org/x/crypto/bcrypt"

type BCryptService struct {
	cost int
}

func NewBCryptService(cost int) *BCryptService {
	return &BCryptService{cost}
}

func (s *BCryptService) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (s *BCryptService) Validate(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
