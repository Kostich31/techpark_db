package threadusecase

import (
	"strconv"

	"github.com/Kostich31/techpark_db/app/domain"
	"github.com/Kostich31/techpark_db/app/tools"
	"github.com/jackc/pgx"
)

type UseCase struct {
	Repository      domain.ThreadRepository
	RepositoryUser  domain.UserRepository
	RepositoryForum domain.ForumRepository
}

func NewUseCase(repository domain.ThreadRepository, userRepository domain.UserRepository, forumRepository domain.ForumRepository) *UseCase {
	return &UseCase{Repository: repository, RepositoryUser: userRepository, RepositoryForum: forumRepository}
}

func (uc *UseCase) CreatePosts(slugOrId string, post []domain.Post) ([]domain.Post, *domain.CustomError) {
	var thread domain.Thread
	id, err := strconv.Atoi(slugOrId)
	if err != nil {
		thread, err = uc.Repository.GetThreadBySlug(slugOrId)
		if err != nil {
			if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == domain.PgxUniqErrorCode {
				return nil, &domain.CustomError{Message: domain.ConflictData}
			}

			if err == pgx.ErrNoRows {
				return nil, &domain.CustomError{Message: domain.NoUser}
			}

			return nil, &domain.CustomError{Message: err.Error()}
		}

	} else {
		thread, err = uc.Repository.GetThreadById(id)
		if err != nil {
			if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == domain.PgxUniqErrorCode {
				return nil, &domain.CustomError{Message: domain.ConflictData}
			}

			if err == pgx.ErrNoRows {
				return nil, &domain.CustomError{Message: domain.NoUser}
			}

			return nil, &domain.CustomError{Message: err.Error()}
		}
	}
	posts, err := uc.Repository.CreatePosts(int(thread.Id), thread.Forum, post)
	if err != nil {
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == domain.PgxUniqErrorCode {
			return nil, &domain.CustomError{Message: domain.ConflictData}
		}
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == domain.PgxBadParentErrorCode {
			return nil, &domain.CustomError{Message: domain.BadParentPost}
		}
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == domain.PgxNoFoundFieldErrorCode {
			return nil, &domain.CustomError{Message: domain.NoUser}
		}
		return nil, &domain.CustomError{Message: err.Error()}
	}

	return posts, nil
}

func (uc *UseCase) CreateVote(slugOrId string, vote domain.Vote) (domain.Thread, *domain.CustomError) {
	var thread domain.Thread
	err := uc.Repository.CreateVoteBySlugOrId(slugOrId, vote)
	if err != nil {
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == domain.PgxNoFoundFieldErrorCode {
			return domain.Thread{}, &domain.CustomError{Message: domain.NoUser}
		}
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == domain.PgxUniqErrorCode {
			err = uc.Repository.UpdateVoteBySlugOrId(slugOrId, vote)
			if err != nil {
				return domain.Thread{}, &domain.CustomError{Message: err.Error()}
			}
			thread, err = uc.Repository.GetThreadBySlugOrId(slugOrId)
			if err != nil {
				return domain.Thread{}, &domain.CustomError{Message: err.Error()}
			}

			return thread, nil
		}
		return domain.Thread{}, &domain.CustomError{Message: err.Error()}
	}

	thread, err = uc.Repository.GetThreadBySlugOrId(slugOrId)
	if err != nil {
		return domain.Thread{}, &domain.CustomError{Message: err.Error()}
	}

	return thread, nil
}

func (uc *UseCase) GetThreadDetails(slugOrId string) (domain.Thread, *domain.CustomError) {
	thread, err := uc.Repository.GetThreadBySlugOrId(slugOrId)
	if err != nil {
		return domain.Thread{}, &domain.CustomError{Message: err.Error()}
	}
	return thread, nil
}

func (uc *UseCase) GetPosts(slugOrId string, filter tools.FilterPosts) ([]*domain.Post, *domain.CustomError) {
	var result []*domain.Post
	var err error

	switch filter.Sort {
	case tools.SortParamFlatDefault:
		result, err = uc.Repository.GetPostsFlatSlugOrId(slugOrId, filter)
	case tools.SortParamParentTree:
		result, err = uc.Repository.GetPostsParentTreeSlugOrId(slugOrId, filter)
	case tools.SortParamTree:
		result, err = uc.Repository.GetPostsTreeSlugOrId(slugOrId, filter)
	}
	if err != nil {
		return nil, &domain.CustomError{Message: err.Error()}
	}

	if len(result) == 0 {
		_, err := uc.Repository.GetThreadBySlugOrId(slugOrId)
		if err != nil {
			return nil, &domain.CustomError{Message: domain.NoUser}
		}
		return []*domain.Post{}, nil
	}

	return result, nil
}

func (uc *UseCase) UpdateThread(slugOrId string, thread domain.Thread) (domain.Thread, *domain.CustomError) {
	thread, err := uc.Repository.UpdateThread(slugOrId, thread)
	if err != nil {
		if pgErr, ok := err.(pgx.PgError); ok && pgErr.Code == domain.PgxUniqErrorCode {
			return domain.Thread{}, &domain.CustomError{Message: domain.ConflictData}
		}
		if err == pgx.ErrNoRows {
			return domain.Thread{}, &domain.CustomError{Message: domain.NoUser}
		}

		return domain.Thread{}, &domain.CustomError{Message: err.Error()}
	}
	return thread, nil
}

func (uc *UseCase) GetPost(id string, filter tools.FilterOnePost) (domain.PostInfo, *domain.CustomError) {
	var result domain.PostInfo

	idNum, err := strconv.Atoi(id)
	if err != nil {
		return domain.PostInfo{}, &domain.CustomError{Message: err.Error()}
	}

	post, err := uc.Repository.GetPostById(idNum)
	if err != nil {
		return domain.PostInfo{}, &domain.CustomError{Message: err.Error()}
	}
	result.Post = post

	if filter.User {
		user, err := uc.RepositoryUser.GetUser(post.Author)
		if err != nil {
			return domain.PostInfo{}, &domain.CustomError{Message: err.Error()}
		}
		result.Author = &user
	}

	if filter.Thread {
		thread, err := uc.Repository.GetThreadById(int(post.Thread))
		if err != nil {
			return domain.PostInfo{}, &domain.CustomError{Message: err.Error()}
		}
		result.Thread = &thread
	}

	if filter.Forum {
		forum, err := uc.RepositoryForum.GetForumBySlug(post.Forum)
		if err != nil {
			return domain.PostInfo{}, &domain.CustomError{Message: err.Error()}
		}
		result.Forum = &forum
	}

	return result, nil
}

func (uc *UseCase) UpdatePost(id string, post domain.Post) (domain.Post, *domain.CustomError) {
	idNum, err := strconv.Atoi(id)
	if err != nil {
		return domain.Post{}, &domain.CustomError{Message: err.Error()}
	}
	if post.Message == "" {
		post, err = uc.Repository.GetPostById(idNum)
	} else {
		post, err = uc.Repository.UpdatePost(idNum, post)
	}
	if err != nil {
		return domain.Post{}, &domain.CustomError{Message: err.Error()}
	}

	return post, nil
}
