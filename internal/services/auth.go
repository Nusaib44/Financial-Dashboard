package services

import (
	"github.com/agency-finance-reality/server/internal/repository"
)

type AuthService interface {
	EnsureUser(id string, email string) error
}

type authService struct {
	founderRepo repository.FounderRepository
}

func NewAuthService(founderRepo repository.FounderRepository) AuthService {
	return &authService{founderRepo: founderRepo}
}

func (s *authService) EnsureUser(id string, email string) error {
	return s.founderRepo.EnsureUser(id, email)
}
