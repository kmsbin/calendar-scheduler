package repositories

import (
	"calendar_scheduler/src/models"
	"database/sql"
	"errors"
	"time"
)

type ResetPasswordRepository struct {
	db *sql.DB
}

func NewResetPasswordRepository(db *sql.DB) ResetPasswordRepository {
	return ResetPasswordRepository{db}
}

func (r ResetPasswordRepository) SetResetPassword(resetPassword models.ResetPassword) error {
	_, err := r.db.Exec(
		"insert into reset_passwords values ($1, $2, $3, $4)",
		resetPassword.UserId,
		resetPassword.Email,
		resetPassword.Code,
		resetPassword.Expiry,
	)

	return err
}

func (r ResetPasswordRepository) GetResetPasswordByCode(code string) (*models.ResetPassword, error) {
	row := r.db.QueryRow("select * from reset_passwords where code = $1", code)
	resetPassword := models.ResetPassword{}
	var expiry string
	err := row.Scan(
		&resetPassword.UserId,
		&resetPassword.Email,
		&resetPassword.Code,
		&expiry,
	)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, ResetPasswordNotFound
		}
		return nil, err
	}
	resetPassword.Expiry, err = time.Parse(time.RFC3339, expiry)
	return &resetPassword, nil
}

func (r ResetPasswordRepository) DeleteResetPasswordData(code string) error {
	_, err := r.db.Exec(
		"delete from reset_passwords where code = $1",
		code,
	)

	return err
}

var ResetPasswordNotFound = errors.New("your reset passoword code is expired")
