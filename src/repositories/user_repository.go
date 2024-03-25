package repositories

import (
	"calendar_scheduler/src/models"
	"database/sql"
	"errors"
)

type UserRepository struct {
	db *sql.DB
}

func (u *UserRepository) GetUserById(userId any) (*models.User, error) {
	row := u.db.QueryRow("select users.user_id, users.name, users.email from users where user_id = $1", userId)
	return prepareUser(row)
}

func prepareUser(rows *sql.Row) (*models.User, error) {
	var user models.User
	err := rows.Scan(&user.Id, &user.Name, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, UserNotFounded
		}
		return nil, err
	}
	return &user, nil
}

func (u *UserRepository) GetUserByEmail(email string) (*models.User, string, error) {
	var user models.User
	row := u.db.QueryRow("select user_id, email, name, password from users where email = $1", email)
	var userPassword string
	err := row.Scan(&user.Id, &user.Email, &user.Name, &userPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", UserNotFounded
		}
		return nil, "", err
	}
	return &user, userPassword, err
}

func (u *UserRepository) ResetPassword(userId int, password string) error {
	_, err := u.db.Exec(
		"update users set password = $2 where $1",
		userId,
		password,
	)

	return err
}
func (u *UserRepository) InsertUser(user *models.User, password []byte) error {
	_, err := u.db.Exec(
		"insert into users(name, email, password) values ($1, $2, $3)",
		user.Name,
		user.Email,
		password,
	)

	return err
}

func (u *UserRepository) DeleteUser(userId any) error {
	_, err := u.db.Exec("delete from users where user_id = $1", userId)
	return err
}

func NewUserRepository(db *sql.DB) UserRepository {
	return UserRepository{db}
}
