package servicerepository

import (
	"github.com/Kostich31/techpark_db/app/domain"
	"github.com/jackc/pgx"
)

type Repository struct {
	db *pgx.ConnPool
}

func NewRepository(db *pgx.ConnPool) *Repository {
	return &Repository{db: db}
}

func (repository *Repository) GetStatus() (domain.Status, error) {
	var result domain.Status
	row := repository.db.QueryRow(
		`SELECT * FROM
		(SELECT COUNT(*) FROM users) as u,
 		(SELECT COUNT(*) FROM forum) as f,
		(SELECT COUNT(*) FROM thread) as t,
		(SELECT COUNT(*) FROM post) as p;`)

	err := row.Scan(
		&result.User,
		&result.Forum,
		&result.Thread,
		&result.Post)
	if err != nil {
		return domain.Status{}, err
	}

	return result, nil
}
func (repository *Repository) Clear() error {
	_, err := repository.db.Exec(`TRUNCATE users, forum, thread, post, vote, users_forum;`)
	if err != nil {
		return err
	}

	return nil
}
