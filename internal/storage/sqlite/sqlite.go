package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Woland-prj/microtasks_sso/internal/domain/cerrors"
	"github.com/Woland-prj/microtasks_sso/internal/domain/entities"
	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.aqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %w", op,
			cerrors.NewCriticalInternalError("sql.Open", err),
		)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(ctx context.Context, user *entities.User) (int64, error) {
	const op = "storage.sqlite.SaveUser"

	stmt, err := s.db.PrepareContext(ctx, "INSERT INTO users (email, pass_hash) VALUES (?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, cerrors.NewCriticalInternalError("s.db.PrepareContext", err))
	}

	res, err := stmt.ExecContext(ctx, user.Email, user.PassHash)
	if err != nil {
		var sqliteErr sqlite3.Error
		// if errors.As(err, &sqliteErr) {
		// 	fmt.Println(fmt.Sprintf("its sqlite error %v", sqliteErr))
		// } else {
		// 	fmt.Printf("Not sqlier error %s", err.Error())
		// }
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, cerrors.NewAlreadyExistsError(fmt.Sprintf("user %d", user.UID)))
		}
		return 0, fmt.Errorf("%s: %w", op, cerrors.NewCriticalInternalError("stmt.ExecContext", err))
	}

	uid, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, cerrors.NewCriticalInternalError("res.LastInsertId", err))
	}

	return uid, nil
}

func (s *Storage) GetUser(ctx context.Context, email string) (*entities.User, error) {
	const op = "storage.sqlite.GetUser"

	stmt, err := s.db.PrepareContext(ctx, "SELECT id, email, pass_hash FROM users WHERE email = ?")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, cerrors.NewCriticalInternalError("s.db.PrepareContext", err))
	}

	row := stmt.QueryRowContext(ctx, email)

	var user entities.User
	err = row.Scan(&user.UID, &user.Email, &user.PassHash)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, cerrors.NewNotFoundError(fmt.Sprintf("user %s", email)))
		}
		return nil, fmt.Errorf("%s: %w", op, cerrors.NewCriticalInternalError("stmt.ExecContext", err))
	}

	return &user, nil
}

func (s *Storage) GetApp(ctx context.Context, id int64) (*entities.App, error) {
	const op = "storage.sqlite.GetApp"

	stmt, err := s.db.PrepareContext(ctx, "SELECT id, name, auth_secret, refresh_secret FROM apps WHERE id = ?")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, cerrors.NewCriticalInternalError("s.db.PrepareContext", err))
	}

	row := stmt.QueryRowContext(ctx, id)

	var app entities.App
	err = row.Scan(&app.ID, &app.Name, &app.AuthSecret, &app.RefreshSecret)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, cerrors.NewNotFoundError(fmt.Sprintf("app %d", id)))
		}
		return nil, fmt.Errorf("%s: %w", op, cerrors.NewCriticalInternalError("stmt.ExecContext", err))
	}

	return &app, nil
}
