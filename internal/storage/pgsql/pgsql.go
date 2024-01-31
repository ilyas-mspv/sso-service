package pgsql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx"
	_ "github.com/jackc/pgx"
	"sso-service/internal/domain/models"
	"sso-service/internal/storage"
)

type Storage struct {
	db *pgx.Conn
}

func New(databaseUrl string) (*Storage, error) {
	const op = "storage.pgsql.New"
	var conn *pgx.Conn
	connCfg, err := pgx.ParseConnectionString(databaseUrl)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	conn, err = pgx.Connect(connCfg)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: conn}, nil
}

// SaveUser saves user to db.
func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "storage.pgsql.SaveUser"
	var id int64

	tx, err := s.db.Begin()
	defer tx.Rollback()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	// todo send context
	err = tx.QueryRow("INSERT INTO users(email, password_hash) VALUES($1, $2)", email, passHash).Scan(&id)
	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// User returns user by email.
func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.pgsql.User"

	var user models.User
	row := s.db.QueryRow("SELECT id, email, password_hash FROM users WHERE email = $1", email)
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

// App returns app by id.
func (s *Storage) App(ctx context.Context, id int) (models.App, error) {
	const op = "storage.pgsql.App"
	row := s.db.QueryRow("SELECT id, name, secret FROM apps WHERE id = $1", id)
	var app models.App
	err := row.Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}

func (s *Storage) Close() error {
	const op = "storage.pgsql.Stop"
	err := s.db.Close()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
