package repositories

import (
	"calendar_scheduler/src/database"
	"calendar_scheduler/src/models"
	"database/sql"
	"errors"
)

type UserRepository struct {
	db *sql.DB
}

func (u *UserRepository) GetUserById(userId any) (*models.User, error) {
	rows := u.db.QueryRow("select users.user_id, users.name, users.email from users where user_id = $1", userId)
	return prepareUser(rows)
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

func (u *UserRepository) InsertUser(user *models.User, password []byte) error {
	_, err := u.db.Exec(
		"insert into users(name, email, password) values ($1, $2, $3)",
		user.Name,
		user.Email,
		password,
	)
	return err
}
func NewUserRepository() UserRepository {
	db, err := database.OpenConnection()
	if err != nil {
		panic(err)
	}
	return UserRepository{db}
}
