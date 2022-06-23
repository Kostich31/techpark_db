package threadrepository

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

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

func (repository *Repository) CreatePosts(threadId int, threadForum string, posts []domain.Post) ([]domain.Post, error) {
	query := `INSERT INTO post(parent, author, message, thread, forum) VALUES `
	var values []interface{}
	if len(posts) == 0 {
		query += fmt.Sprintf(`(0, null, null, %d, '%s')`, threadId, threadForum)
	}
	for i, post := range posts {
		value := fmt.Sprintf("($%d, $%d, $%d, $%d, $%d),",
			i*5+1, i*5+2, i*5+3, i*5+4, i*5+5)
		query += value
		values = append(values, post.Parent, post.Author, post.Message, threadId, threadForum)
	}
	query = strings.TrimSuffix(query, ",")
	query += ` RETURNING id, parent, author, message, isEdited, forum, thread, created;`

	rows, err := repository.db.Query(query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []domain.Post{}
	if len(posts) != 0 {
		for rows.Next() {
			var post domain.Post
			err := rows.Scan(&post.Id, &post.Parent, &post.Author, &post.Message,
				&post.IsEdited, &post.Forum, &post.Thread, &post.Created)
			if err != nil {
				return nil, err
			}
			result = append(result, post)
		}

		if rows.Err() != nil {
			return nil, rows.Err()
		}
	}
	return result, nil
}

func (repository *Repository) GetThreadBySlug(slug string) (domain.Thread, error) {
	var result domain.Thread
	row := repository.db.QueryRow(`SELECT id, title, author, forum, message, votes, slug, created 
		FROM thread WHERE slug=$1`, slug)

	err := row.Scan(&result.Id, &result.Title, &result.Author, &result.Forum, &result.Message, &result.Votes,
		&result.Slug, &result.Created)
	if err != nil {
		return domain.Thread{}, err
	}

	return result, nil
}

func (repository *Repository) GetThreadById(id int) (domain.Thread, error) {
	var result domain.Thread
	row := repository.db.QueryRow(`SELECT id, title, author, forum, message, votes, slug, created 
		FROM thread WHERE id=$1`, id)
	var nullSlug sql.NullString
	err := row.Scan(&result.Id, &result.Title, &result.Author, &result.Forum, &result.Message, &result.Votes,
		&nullSlug, &result.Created)
	if err != nil {
		return domain.Thread{}, err
	}
	result.Slug = nullSlug.String
	return result, nil
}

func (repository *Repository) CreateVoteBySlugOrId(slugOrId string, vote domain.Vote) error {
	id, err := strconv.Atoi(slugOrId)
	if err != nil {
		_, err = repository.db.Exec(`
				INSERT INTO 
				vote(nickname, voice, thread) 
				VALUES ($1, $2, (select id from thread where slug = $3));`,
			vote.NickName,
			vote.Voice,
			slugOrId)
	} else {
		_, err = repository.db.Exec(`
				INSERT INTO 
				vote(nickname, voice, thread) 
				VALUES ($1, $2, $3);`,
			vote.NickName,
			vote.Voice,
			id)
	}

	if err != nil {
		return err
	}
	return nil
}

func (repository *Repository) UpdateVoteBySlugOrId(slugOrId string, vote domain.Vote) error {
	id, err := strconv.Atoi(slugOrId)
	if err != nil {
		_, err = repository.db.Exec(`
			UPDATE vote 
			SET voice=$1
			WHERE nickname=$2 and thread=(select id from thread where slug = $3)`,
			vote.Voice,
			vote.NickName,
			slugOrId)
	} else {
		_, err = repository.db.Exec(`
			UPDATE vote 
			SET voice=$1
			WHERE nickname=$2 and thread=$3`,
			vote.Voice,
			vote.NickName,
			id)
	}

	if err != nil {
		return err
	}

	return nil
}

func (repository *Repository) GetThreadBySlugOrId(slugOrId string) (domain.Thread, error) {
	var result domain.Thread
	var row *pgx.Row
	id, err := strconv.Atoi(slugOrId)
	if err != nil {
		row = repository.db.QueryRow(`SELECT id, title, author, forum, message, votes, slug, created 
			FROM thread WHERE slug=$1`, slugOrId)
	} else {
		row = repository.db.QueryRow(`SELECT id, title, author, forum, message, votes, slug, created 
			FROM thread WHERE id=$1`, id)
	}
	var nullSlug sql.NullString
	err = row.Scan(&result.Id, &result.Title, &result.Author, &result.Forum, &result.Message, &result.Votes,
		&nullSlug, &result.Created)
	if err != nil {
		return domain.Thread{}, err
	}
	result.Slug = nullSlug.String
	return result, nil
}

func (repository *Repository) GetPostsFlatSlugOrId(slugOrId string, filter tools.FilterPosts) ([]*domain.Post, error) {
	var rows *pgx.Rows
	var err error
	id, err := strconv.Atoi(slugOrId)
	if err != nil {
		var tmpId sql.NullInt64
		row := repository.db.QueryRow(`select id from thread where slug = $1`, slugOrId)
		err = row.Scan(&tmpId)
		if filter.Since == tools.SinceParamDefault {
			rows, err = repository.db.Query(`
											select id, parent, author, message, isEdited, forum,  
											thread, created from post where 
											thread = $1 
											order by id `+filter.Desc+` limit $2`, tmpId, filter.Limit)
		} else {
			if filter.Desc == tools.SortParamTrue {
				rows, err = repository.db.Query(`
											select id, parent, author, message, isEdited, forum,  
											thread, created from post where 
											thread = $1 and id < $2 
											order by id desc limit $3`, tmpId, filter.Since, filter.Limit)
			} else {
				rows, err = repository.db.Query(`
											select id, parent, author, message, isEdited, forum,  
											thread, created from post where 
											thread = $1 and id > $2 
											order by id asc limit $3`, tmpId, filter.Since, filter.Limit)
			}
		}
	} else {
		if filter.Since == tools.SinceParamDefault {
			rows, err = repository.db.Query(`
											select id, parent, author, message, isEdited, forum,  
											thread, created from post where 
											thread = $1 
											order by id `+filter.Desc+` limit $2`, id, filter.Limit)
		} else {
			if filter.Desc == tools.SortParamTrue {
				rows, err = repository.db.Query(`
											select id, parent, author, message, isEdited, forum,  
											thread, created from post where 
											thread = $1 and id < $2 
											order by id desc limit $3`, id, filter.Since, filter.Limit)
			} else {
				rows, err = repository.db.Query(`
											select id, parent, author, message, isEdited, forum,  
											thread, created from post where 
											thread = $1 and id > $2 
											order by id asc limit $3`, id, filter.Since, filter.Limit)
			}
		}
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*domain.Post
	for rows.Next() {
		post := &domain.Post{}

		err = rows.Scan(
			&post.Id,
			&post.Parent,
			&post.Author,
			&post.Message,
			&post.IsEdited,
			&post.Forum,
			&post.Thread,
			&post.Created)
		if err != nil {
			return nil, err
		}

		result = append(result, post)
	}

	return result, err
}

func (repository *Repository) GetPostsTreeSlugOrId(slugOrId string, filter tools.FilterPosts) ([]*domain.Post, error) {
	var rows *pgx.Rows
	var err error
	id, err := strconv.Atoi(slugOrId)
	if err != nil {
		var tmpId sql.NullInt64
		row := repository.db.QueryRow("select id from thread where slug = $1", slugOrId)
		err = row.Scan(&tmpId)
		if filter.Since == tools.SinceParamDefault {
			rows, err = repository.db.Query(`
											select id, parent, author, message, isEdited, forum,  
											thread, created from post where 
											thread = $1 
											order by paths `+filter.Desc+`, id `+filter.Desc+` limit $2`, tmpId, filter.Limit)
		} else {
			if filter.Desc == tools.SortParamTrue {
				rows, err = repository.db.Query(`
											select id, parent, author, message, isEdited, forum,  
											thread, created from post where 
											thread = $1 and paths < (select paths from post where id=$2) 
											order by paths desc, id desc limit $3`, tmpId, filter.Since, filter.Limit)

			} else {
				rows, err = repository.db.Query(`
											select id, parent, author, message, isEdited, forum,  
											thread, created from post where 
											thread = $1 and paths > (select paths from post where id=$2) 
											order by paths asc, id asc limit $3`, tmpId, filter.Since, filter.Limit)
			}
		}
	} else {
		if filter.Since == tools.SinceParamDefault {
			rows, err = repository.db.Query(`
											select id, parent, author, message, isEdited, forum,  
											thread, created from post where 
											thread = $1 
											order by paths `+filter.Desc+`, id `+filter.Desc+` limit $2`, id, filter.Limit)
		} else {
			if filter.Desc == tools.SortParamTrue {
				rows, err = repository.db.Query(`
											select id, parent, author, message, isEdited, forum,  
											thread, created from post where 
											thread = $1 and paths < (select paths from post where id=$2) 
											order by paths desc, id desc limit $3`, id, filter.Since, filter.Limit)
			} else {
				rows, err = repository.db.Query(`
											select id, parent, author, message, isEdited, forum,  
											thread, created from post where 
											thread = $1 and paths > (select paths from post where id=$2) 
											order by paths asc, id asc limit $3`, id, filter.Since, filter.Limit)
			}
		}
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*domain.Post
	for rows.Next() {
		post := &domain.Post{}

		err = rows.Scan(
			&post.Id,
			&post.Parent,
			&post.Author,
			&post.Message,
			&post.IsEdited,
			&post.Forum,
			&post.Thread,
			&post.Created)
		if err != nil {
			return nil, err
		}

		result = append(result, post)
	}

	return result, err
}

func (repository *Repository) GetPostsParentTreeSlugOrId(slugOrId string, filter tools.FilterPosts) ([]*domain.Post, error) {
	var rows *pgx.Rows
	var err error
	id, err := strconv.Atoi(slugOrId)
	if err != nil {
		var tmpId sql.NullInt64
		row := repository.db.QueryRow("select id from thread where slug = $1", slugOrId)
		err = row.Scan(&tmpId)
		if filter.Since == tools.SinceParamDefault {
			if filter.Desc == tools.SortParamTrue {
				rows, err = repository.db.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created FROM post
					WHERE paths[1] IN (SELECT id FROM post WHERE thread = $1 
					AND parent = 0 ORDER BY id DESC LIMIT $2)
					ORDER BY paths[1] DESC, paths ASC, id ASC;`,
					tmpId,
					filter.Limit)
			} else {
				rows, err = repository.db.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created FROM post
					WHERE paths[1] IN (SELECT id FROM post WHERE thread = $1 
					AND parent = 0 ORDER BY id ASC LIMIT $2)
					ORDER BY paths ASC, id ASC;`,
					tmpId,
					filter.Limit)
			}
		} else {
			if filter.Desc == tools.SortParamTrue {
				rows, err = repository.db.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created FROM post
					WHERE paths[1] IN (SELECT id FROM post WHERE thread = $1 
					AND parent = 0 AND paths[1] <
					(SELECT paths[1] FROM post WHERE id = $2) ORDER BY id DESC LIMIT $3)
					ORDER BY paths[1] DESC, paths ASC, id ASC;`,
					tmpId,
					filter.Since,
					filter.Limit)
			} else {
				rows, err = repository.db.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created FROM post
					WHERE paths[1] IN (SELECT id FROM post WHERE thread = $1
					AND parent = 0 AND paths[1] >
					(SELECT paths[1] FROM post WHERE id = $2) ORDER BY id ASC LIMIT $3) 
					ORDER BY paths ASC, id ASC;`,
					tmpId,
					filter.Since,
					filter.Limit)
			}
		}
	} else {
		if filter.Since == tools.SinceParamDefault {
			if filter.Desc == tools.SortParamTrue {
				rows, err = repository.db.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created FROM post
					WHERE paths[1] IN (SELECT id FROM post WHERE thread = $1 AND parent = 0 ORDER BY id DESC LIMIT $2)
					ORDER BY paths[1] DESC, paths ASC, id ASC;`,
					id,
					filter.Limit)
			} else {
				rows, err = repository.db.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created FROM post
					WHERE paths[1] IN (SELECT id FROM post WHERE thread = $1 AND parent = 0 ORDER BY id ASC LIMIT $2)
					ORDER BY paths ASC, id ASC;`,
					id,
					filter.Limit)
			}
		} else {
			if filter.Desc == tools.SortParamTrue {
				rows, err = repository.db.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created FROM post
					WHERE paths[1] IN (SELECT id FROM post WHERE thread = $1 AND parent = 0 AND paths[1] <
					(SELECT paths[1] FROM post WHERE id = $2) ORDER BY id DESC LIMIT $3)
					ORDER BY paths[1] DESC, paths ASC, id ASC;`,
					id,
					filter.Since,
					filter.Limit)
			} else {
				rows, err = repository.db.Query(`
					SELECT id, parent, author, message, isEdited, forum, thread, created FROM post
					WHERE paths[1] IN (SELECT id FROM post WHERE thread = $1 AND parent = 0 AND paths[1] >
					(SELECT paths[1] FROM post WHERE id = $2) ORDER BY id ASC LIMIT $3) 
					ORDER BY paths ASC, id ASC;`,
					id,
					filter.Since,
					filter.Limit)
			}
		}
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*domain.Post
	for rows.Next() {
		post := &domain.Post{}

		err = rows.Scan(
			&post.Id,
			&post.Parent,
			&post.Author,
			&post.Message,
			&post.IsEdited,
			&post.Forum,
			&post.Thread,
			&post.Created)
		if err != nil {
			return nil, err
		}

		result = append(result, post)
	}

	return result, err
}

func (repository *Repository) UpdateThread(slugOrId string, thread domain.Thread) (domain.Thread, error) {
	var row *pgx.Row
	var err error
	id, err := strconv.Atoi(slugOrId)
	if err != nil {
		row = repository.db.QueryRow(`UPDATE thread SET 
			title=COALESCE(NULLIF($1, ''), title), 
			author=COALESCE(NULLIF($2, ''), author), 
			forum=COALESCE(NULLIF($3, ''), forum), 
			message=COALESCE(NULLIF($4, ''), message) 
			where slug=$5 returning id, title, author, forum, message, votes, slug, created`,
			thread.Title, thread.Author, thread.Forum, thread.Message, slugOrId)
	} else {
		row = repository.db.QueryRow(`UPDATE thread SET 
			title=COALESCE(NULLIF($1, ''), title), 
			author=COALESCE(NULLIF($2, ''), author), 
			forum=COALESCE(NULLIF($3, ''), forum),
			message=COALESCE(NULLIF($4, ''), message) 
			where id=$5 returning id, title, author, forum, message, votes, slug, created`,
			thread.Title, thread.Author, thread.Forum, thread.Message, id)
	}

	err = row.Scan(
		&thread.Id,
		&thread.Title,
		&thread.Author,
		&thread.Forum,
		&thread.Message,
		&thread.Votes,
		&thread.Slug,
		&thread.Created)
	if err != nil {
		return domain.Thread{}, err
	}

	return thread, nil
}

func (repository *Repository) GetPostById(id int) (domain.Post, error) {
	var result domain.Post
	row := repository.db.QueryRow(`SELECT id, parent, author, message, isEdited,
		forum, thread, created 
		FROM post WHERE id=$1`, id)

	err := row.Scan(&result.Id, &result.Parent, &result.Author, &result.Message, &result.IsEdited,
		&result.Forum, &result.Thread, &result.Created)
	if err != nil {
		return domain.Post{}, err
	}
	return result, nil
}

func (repository *Repository) UpdatePost(id int, post domain.Post) (domain.Post, error) {
	query := repository.db.QueryRow(`UPDATE post SET
		message=$1,
		isedited= case when message = $1 then isedited else true end 
		where id=$2 
		returning id, parent, author, message, isedited, forum, thread, created`,
		post.Message, id)

	err := query.Scan(
		&post.Id,
		&post.Parent,
		&post.Author,
		&post.Message,
		&post.IsEdited,
		&post.Forum,
		&post.Thread,
		&post.Created,
	)
	if err != nil {
		return domain.Post{}, err
	}

	return post, nil
}
