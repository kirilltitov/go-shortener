package storage

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kirilltitov/go-shortener/internal/logger"
)

// ErrDuplicate является ошибкой о дубликате ссылки.
var ErrDuplicate = errors.New("duplicate URL found")

// PgSQL является хранилищем сокращенных ссылок в PostgreSQL.
type PgSQL struct {
	C *pgxpool.Pool
}

// NewPgSQLStorage создает, конфигурирует, соединяется с БД и возвращает объект хранилища PostgreSQL.
// Возвращает ошибку, если не удалось подключиться к БД.
func NewPgSQLStorage(ctx context.Context, DSN string) (*PgSQL, error) {
	conf, err := pgxpool.ParseConfig(DSN)
	if err != nil {
		return nil, err
	}
	conf.MaxConns = 10

	pool, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		return nil, err
	}

	logger.Log.Infof("Connected to PgSQL with DSN %s", DSN)

	return &PgSQL{C: pool}, nil
}

// DBRow являет собою запись из таблицы url.
type DBRow struct {
	ID         int       `db:"id"`
	UserID     uuid.UUID `db:"user_id"`
	ShortURL   string    `db:"short_url"`
	URL        string    `db:"url"`
	IsDeleted  bool      `db:"is_deleted"`
	Duplicates int       `db:"duplicates"`
	CreatedAt  time.Time `db:"created_at"`
}

// Get загружает из хранилища информацию о сокращенной ссылке.
func (p PgSQL) Get(ctx context.Context, shortURL string) (string, error) {
	var row DBRow
	if err := pgxscan.Get(ctx, p.C, &row, `select * from public.url where short_url = $1`, shortURL); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrNotFound
		} else {
			return "", err
		}
	}

	logger.Log.Infof("Queried row %v", row)

	if row.IsDeleted {
		return "", ErrDeleted
	}

	return row.URL, nil
}

// Set записывает в хранилище информацию о сокращенной ссылке.
func (p PgSQL) Set(ctx context.Context, userID uuid.UUID, URL string) (string, error) {
	tx, err := p.C.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	shortURL, err := txInsert(ctx, tx, userID, URL)
	if err != nil {
		return shortURL, err
	}
	if err := tx.Commit(ctx); err != nil {
		return "", err
	}

	return shortURL, nil
}

// MultiSet записывает в хранилище информацию о нескольких сокращенных ссылках.
func (p PgSQL) MultiSet(ctx context.Context, userID uuid.UUID, items Items) (Items, error) {
	tx, err := p.C.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var result Items
	for _, item := range items {
		shortURL, err2 := txInsert(ctx, tx, userID, item.URL)
		if err2 != nil {
			return nil, err2
		}
		result = append(result, Item{
			UUID: item.UUID,
			URL:  shortURL,
		})
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetByUser загружает из хранилища все сокращенные ссылки пользователя.
func (p PgSQL) GetByUser(ctx context.Context, userID uuid.UUID) (Items, error) {
	var result Items

	var rows []*DBRow
	err := pgxscan.Select(ctx, p.C, &rows, `select * from url where user_id = $1`, userID)
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		result = append(result, Item{
			UUID:     strconv.Itoa(row.ID),
			URL:      row.URL,
			ShortURL: row.ShortURL,
		})
	}

	return result, nil
}

// DeleteByUser удаляет из хранилища все сокращенные ссылки пользователя.
func (p PgSQL) DeleteByUser(ctx context.Context, userID uuid.UUID, shortURL string) error {
	tx, err := p.C.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	res, err := tx.Exec(ctx, `update url set is_deleted = true where short_url = $1 and user_id = $2`, shortURL, userID)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	logger.Log.Infof("Deleting URL '%s' for user %s: %d rows affected", shortURL, userID, res.RowsAffected())

	return nil
}

// GetStats возвращает статистику хранилища.
func (p PgSQL) GetStats(ctx context.Context) (*Stats, error) {
	var stats Stats

	if err := pgxscan.Get(ctx, p.C, &stats, `select count(distinct user_id) users, count(id) urls from public.url`); err != nil {
		return nil, err
	}

	return &stats, nil
}

// MigrateUp выполняет миграции в хранилище.
func (p PgSQL) MigrateUp(ctx context.Context) error {
	if _, err := p.C.Exec(ctx, `
		create table if not exists url
		(
			id         serial primary key,
			user_id    uuid not null,
			short_url  varchar not null,
			url        varchar not null constraint url_pk unique,
			is_deleted bool not null default false,
			duplicates int not null default 0,
			created_at timestamp default CURRENT_TIMESTAMP not null
		)
	`); err != nil {
		return err
	}

	return nil
}

func txInsert(ctx context.Context, tx pgx.Tx, userID uuid.UUID, URL string) (string, error) {
	var inserted struct {
		Cur        int    `db:"id"`
		ShortURL   string `db:"short_url"`
		Duplicates int    `db:"duplicates"`
	}
	if err := pgxscan.Get(ctx, tx, &inserted, `
		insert into public.url (user_id, short_url, url) values ($1, $2, $3)
		on conflict on constraint url_pk do update set duplicates = url.duplicates + 1 returning id, short_url, duplicates
		`, userID, "", URL); err != nil {
		logger.Log.Infof("Could not insert new row: %+v", err)
		return "", err
	}

	var shortURL string
	var err error
	if inserted.ShortURL != "" {
		shortURL = inserted.ShortURL
		err = ErrDuplicate
		logger.Log.Infof(
			"Found duplicate for URL '%s', returning pre-existing short URL '%s' (this is %dth duplicate)",
			URL, shortURL, inserted.Duplicates)
	} else {
		shortURL = intToShortURL(inserted.Cur)
		if _, err2 := tx.Exec(ctx, `update public.url set short_url = $1 where id = $2`, shortURL, inserted.Cur); err2 != nil {
			return "", err2
		}
	}

	return shortURL, err
}
