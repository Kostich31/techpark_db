package userrepository

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

func (repository *Repository) AddUser(user domain.User) (domain.User, error) {
	_, err := repository.db.Exec(`INSERT INTO Users (nickname, fullname, about, email) VALUES ($1, $2, $3, $4)`,
		user.Nickname, user.FullName, user.About, user.Email)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (repository *Repository) GetUser(nickname string) (domain.User, error) {
	var result domain.User
	row := repository.db.QueryRow(`SELECT nickname, fullname, about, email 
		FROM Users WHERE nickname=$1`, nickname)

	err := row.Scan(&result.Nickname, &result.FullName, &result.About, &result.Email)
	if err != nil {
		return domain.User{}, err
	}
	return result, nil
}

func (repository *Repository) UpdateUser(user domain.User) (domain.User, error) {
	query := repository.db.QueryRow(`UPDATE users SET 
		fullname=COALESCE(NULLIF($1, ''), fullname), 
		about=COALESCE(NULLIF($2, ''), about),
		email=COALESCE(NULLIF($3, ''), email) 
		where nickname=$4 returning nickname, fullname, about, email`, user.FullName, user.About, user.Email, user.Nickname)

	err := query.Scan(
		&user.Nickname,
		&user.FullName,
		&user.About,
		&user.Email)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (repository *Repository) GetUsersByNicknameOrEmail(nickname string, email string) ([]domain.User, error) {
	rows, err := repository.db.Query(`SELECT nickname, fullname, about, email 
		FROM users WHERE nickname=$1 or email=$2`, nickname, email)
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
