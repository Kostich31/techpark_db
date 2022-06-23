package forumusecase

import (
	"github.com/Kostich31/techpark_db/app/domain"
	"github.com/Kostich31/techpark_db/app/tools"
	"github.com/jackc/pgx"
)

type UseCase struct {
	RepositoryForum  domain.ForumRepository
	RepositoryThread domain.ThreadRepository
}

func NewUseCase(repositoryForum domain.ForumRepository, repository domain.ThreadRepository) *UseCase {
	return &UseCase{RepositoryForum: repositoryForum, RepositoryThread: repository}
}

func (uc *UseCase) CreateForum(forumGet domain.Forum) (domain.Forum, *domain.CustomError) {
	forum, err := uc.RepositoryForum.AddForum(forumGet)
	if err != nil {
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == domain.PgxNoFoundFieldErrorCode {
			return domain.Forum{}, &domain.CustomError{Message: domain.NoUser}
		}
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == domain.PgxUniqErrorCode {
			forum, err = uc.RepositoryForum.GetForumBySlug(forumGet.Slug)
			if err != nil {
				return domain.Forum{}, &domain.CustomError{Message: err.Error()}
			}
			return forum, &domain.CustomError{Message: domain.ConflictData}
		}
		return domain.Forum{}, &domain.CustomError{Message: err.Error()}
	}

	return forum, nil
}

func (uc *UseCase) GetDetailsForum(slug string) (domain.Forum, *domain.CustomError) {
	forum, err := uc.RepositoryForum.GetDetailsForum(slug)
	if err == pgx.ErrNoRows {
		return domain.Forum{}, &domain.CustomError{Message: domain.NoSlug}
	}
	return forum, nil
}

func (uc *UseCase) CreateThread(threadGet domain.Thread) (domain.Thread, *domain.CustomError) {
	var randomSlug bool
	if threadGet.Slug == "" {
		randomSlug = true
	}
	thread, err := uc.RepositoryForum.AddThread(threadGet)
	if err != nil {
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == domain.PgxNoFoundFieldErrorCode {
			return domain.Thread{}, &domain.CustomError{Message: domain.NoUser}
		}
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == domain.PgxUniqErrorCode {
			thread, err = uc.RepositoryThread.GetThreadBySlug(threadGet.Slug)
			if err != nil {
				return domain.Thread{}, &domain.CustomError{Message: err.Error()}
			}
			return thread, &domain.CustomError{Message: domain.ConflictData}
		}
		return domain.Thread{}, &domain.CustomError{Message: err.Error()}
	}

	if randomSlug == true {
		thread.Slug = ""
	}
	return thread, nil
}

func (uc *UseCase) GetUsersForum(slug string, filter tools.FilterUser) ([]domain.User, *domain.CustomError) {
	users, err := uc.RepositoryForum.GetUsersForum(slug, filter)
	if users == nil {
		_, err = uc.RepositoryForum.GetForumBySlug(slug)
		if err != nil {
			return nil, &domain.CustomError{Message: domain.NoSlug}
		}
		return []domain.User{}, nil
	}
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, &domain.CustomError{Message: domain.NoSlug}
		}
		return nil, &domain.CustomError{Message: err.Error()}
	}

	return users, nil
}

func (uc *UseCase) GetForumThreads(slug string, filter tools.FilterThread) ([]domain.Thread, *domain.CustomError) {
	threads, err := uc.RepositoryForum.GetForumThreads(slug, filter)
	if threads == nil {
		_, err := uc.RepositoryForum.GetForumBySlug(slug)
		if err != nil {
			return nil, &domain.CustomError{Message: err.Error()}
		}
		return []domain.Thread{}, nil
	}
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, &domain.CustomError{Message: domain.NoSlug}
		}
		return nil, &domain.CustomError{Message: err.Error()}
	}

	return threads, nil
}
