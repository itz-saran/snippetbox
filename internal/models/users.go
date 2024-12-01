package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

// Insert new user into table
func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := "INSERT INTO users (name, email, hashed_password, created) VALUES ($1, $2, $3, $4) RETURNING id"
	var result int
	err = m.DB.QueryRow(stmt, name, email, hashedPassword, pq.FormatTimestamp(time.Now())).Scan(&result)
	if err != nil {
		var pgError *pq.Error
		if errors.As(err, &pgError) && string(pgError.Code) == "23505" {
			return ErrDuplicateEmail
		}
		return err
	}
	return nil
}

// Verify whether the user exists with provided email address and password. This will return
// relevant user ID if they exist
func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	stmt := "SELECT id, hashed_password from users WHERE email=$1"
	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool
	stmt := "SELECT EXISTS (SELECT * from users WHERE id=$1)"
	err := m.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}
