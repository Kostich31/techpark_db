package serviceusecase

import (
	"github.com/Kostich31/techpark_db/app/domain"
)

type UseCase struct {
	Repository domain.ServiceRepository
}

func NewUseCase(repository domain.ServiceRepository) *UseCase {
	return &UseCase{Repository: repository}
}

func (uc *UseCase) GetStatus() (domain.Status, error) {
	return uc.Repository.GetStatus()
}

func (uc *UseCase) Clear() error {
	return uc.Repository.Clear()
}
