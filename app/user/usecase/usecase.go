package userusecase

import (
	"github.com/Kostich31/techpark_db/app/domain"
	"github.com/jackc/pgx"
)

type UseCase struct {
	Repository domain.UserRepository
}

func NewUseCase(repository domain.UserRepository) *UseCase {
	return &UseCase{Repository: repository}
}

func (uc *UseCase) CreateUser(user domain.User) ([]domain.User, error) {
	var resultArray []domain.User
	result, err := uc.Repository.AddUser(user)
	if err != nil {
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == domain.PgxUniqErrorCode {
			result, err1 := uc.Repository.GetUsersByNicknameOrEmail(user.Nickname, user.Email)
			if err1 != nil {
				return nil, err1
			}
			return result, err
		}
	}

	resultArray = append(resultArray, result)
	return resultArray, err
}

func (uc *UseCase) GetUserProfile(nickname string) (domain.User, *domain.CustomError) {
	user, err := uc.Repository.GetUser(nickname)
	if err == pgx.ErrNoRows {
		return domain.User{}, &domain.CustomError{Message: domain.NoUser}
	}
	return user, nil
}

func (uc *UseCase) UpdateUserProfile(user domain.User) (domain.User, *domain.CustomError) {
	userNew, err := uc.Repository.UpdateUser(user)
	if err != nil {
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == domain.PgxUniqErrorCode {
			return domain.User{}, &domain.CustomError{Message: domain.ConflictData}
		}
		if err == pgx.ErrNoRows {
			return domain.User{}, &domain.CustomError{Message: domain.NoUser}
		}

		return domain.User{}, &domain.CustomError{Message: err.Error()}
	}
	return userNew, nil
}
