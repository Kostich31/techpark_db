package forumrepository

import (
	"database/sql"
	"errors"

	"github.com/Kostich31/techpark_db/app/domain"
	"github.com/Kostich31/techpark_db/app/tools"
	"github.com/jackc/pgx"
)

type Repository struct {
	db *pgx.ConnPool
}

func NewRepository(db *pgx.ConnPool) *Repository {
	return &Repository{db: db}
}

func (repository *Repository) AddForum(forum domain.Forum) (domain.Forum, error) {
	row := repository.db.QueryRow(`INSERT INTO forum (title, "user", slug) 
		VALUES ($1,COALESCE((SELECT nickname FROM users WHERE nickname = $2), $2), $3)
		RETURNING title, "user", slug`,
		forum.Title, forum.User, forum.Slug)

	err := row.Scan(
		&forum.Title,
		&forum.User,
		&forum.Slug)
	if err != nil {
		return domain.Forum{}, err
	}

	return forum, nil
}

func (repository *Repository) GetDetailsForum(slug string) (domain.Forum, error) {
	var result domain.Forum
	row := repository.db.QueryRow(`SELECT title, "user", slug, posts, threads 
		FROM Forum WHERE slug=$1`, slug)

	err := row.Scan(&result.Title, &result.User, &result.Slug, &result.Posts, &result.Threads)
	if err != nil {
		return domain.Forum{}, err
	}
	return result, nil
}

func (repository *Repository) AddThread(thread domain.Thread) (domain.Thread, error) {
	row := repository.db.QueryRow(`INSERT INTO thread (title, author, forum, message, slug, created)
		VALUES ($1, $2, COALESCE((SELECT slug from forum where slug = $3), $3), $4, coalesce(nullif($5,'')), $6) 
		returning id, title, author, forum, message, slug, created`,
		thread.Title, thread.Author, thread.Forum, thread.Message, thread.Slug, thread.Created)

	var nullSlug sql.NullString
	err := row.Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &nullSlug, &thread.Created)
	if err != nil {
		return domain.Thread{}, err
	}
	thread.Slug = nullSlug.String
	return thread, nil
}

func (repository *Repository) GetUsersForum(slug string, filter tools.FilterUser) ([]domain.User, error) {
	var rows *pgx.Rows
	var err error
	if filter.Since == tools.SinceParamDefault {
		rows, err = repository.db.Query(`SELECT u.nickname, fullname, about, email 
			FROM users_forum as u inner join users on u.nickname = users.nickname where u.slug = $1 
			order by u.nickname COLLATE "C" `+filter.Desc+` limit $2`,
			slug,
			filter.Limit,
		)
	} else {
		if filter.Desc == tools.SortParamTrue {
			rows, err = repository.db.Query(`SELECT u.nickname, fullname, about, email
				FROM users_forum as u inner join users on u.nickname = users.nickname 
				where u.slug = $1 and u.nickname < ($2 collate "C") 
				order by u.nickname COLLATE "C" desc limit $3`,
				slug,
				filter.Since,
				filter.Limit,
			)
		} else {
			rows, err = repository.db.Query(`SELECT u.nickname, fullname, about, email 
				FROM users_forum as u inner join users on u.nickname = users.nickname 
				where u.slug = $1 and u.nickname > ($2 collate "C") 
				order by u.nickname COLLATE "C" asc limit $3`,
				slug,
				filter.Since,
				filter.Limit,
			)
		}
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		err = rows.Scan(&user.Nickname, &user.FullName, &user.About, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return users, nil
}

func (repository *Repository) GetForumThreads(slug string, filter tools.FilterThread) ([]domain.Thread, error) {
	var rows *pgx.Rows
	var err error

	if filter.Sort != tools.SortParamDefault && filter.Sort != tools.SortParamTrue {
		return nil, errors.New("sql attack")
	}

	if filter.Since == "" {
		rows, err = repository.db.Query(`select id, title, author, forum, message, votes, slug, created `+
			`from thread where forum = $1 order by created `+filter.Sort+` limit $2`, slug, filter.Limit)
	} else {
		if filter.Sort == tools.SortParamTrue {
			rows, err = repository.db.Query(`select id, title, author, forum, message, votes, slug, created 
				from thread where forum = $1 and created <= $3 order by created `+filter.Sort+` limit $2`, slug, filter.Limit, filter.Since)
		} else {
			rows, err = repository.db.Query(`select id, title, author, forum, message, votes, slug, created 
				from thread where forum = $1 and created >= $3 order by created `+filter.Sort+` limit $2`, slug, filter.Limit, filter.Since)
		}
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nullSlug sql.NullString
	var threads []domain.Thread
	for rows.Next() {
		var thread domain.Thread
		err = rows.Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message,
			&thread.Votes, &nullSlug, &thread.Created)
		if err != nil {
			return nil, err
		}
		thread.Slug = nullSlug.String
		threads = append(threads, thread)
	}

	return threads, nil
}

func (repository *Repository) GetForumBySlug(slug string) (domain.Forum, error) {
	var result domain.Forum
	row := repository.db.QueryRow(`SELECT slug, title, "user", posts, threads
		FROM forum WHERE slug=$1`, slug)

	err := row.Scan(&result.Slug, &result.Title, &result.User, &result.Posts, &result.Threads)
	if err != nil {
		return domain.Forum{}, err
	}

	return result, nil
}
