package services

import "sistem-manajemen-armada/internal/repository"

type APIService struct {
	repo *repository.PostgreSQLRepository
}

func NewAPIService(repo *repository.PostgreSQLRepository) *APIService {
	return &APIService{repo: repo}
}
