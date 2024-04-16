package storage

import (
	"context"
	"errors"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/kirilltitov/go-shortener/internal/logger"
)

type PgSQL struct {
	C *pgx.Conn
}

func NewPgSQLStorage(ctx context.Context, DSN string) (*PgSQL, error) {
	conn, err := pgx.Connect(ctx, DSN)
	if err != nil {
		return nil, err
	}

	logger.Log.Infof("Connected to PgSQL with DSN %s", DSN)

	return &PgSQL{C: conn}, nil
}

type DBRow struct {
	ID        int       `db:"id"`
	ShortURL  string    `db:"short_url"`
	URL       string    `db:"url"`
	CreatedAt time.Time `db:"created_at"`
}

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

	return row.URL, nil
}

func (p PgSQL) Set(ctx context.Context, URL string) (string, error) {
	tx, err := p.C.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	shortURL, err := txInsert(ctx, tx, URL)
	if err != nil {
		return "", err
	}
	if err := tx.Commit(ctx); err != nil {
		return "", err
	}

	return shortURL, nil
}

func (p PgSQL) MultiSet(ctx context.Context, items Items) (Items, error) {
	tx, err := p.C.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var result Items
	for _, item := range items {
		shortURL, err := txInsert(ctx, tx, item.URL)
		if err != nil {
			return nil, err
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

// я не смог разобраться с Goose за приемлемое время :(
func (p PgSQL) MigrateUp(ctx context.Context) error {
	_, err := p.C.Exec(ctx, `
	create table if not exists url
	(
		id         serial primary key,
		short_url  varchar                             not null,
		url        varchar                             not null,
		created_at timestamp default CURRENT_TIMESTAMP not null
	)
	`)

	return err
}

func txInsert(ctx context.Context, tx pgx.Tx, URL string) (string, error) {
	var cur int
	if err := pgxscan.Get(ctx, tx, &cur, `
		insert into public.url (short_url, url) values ($1, $2) returning id
		`, "", URL); err != nil {
		return "", err
	}

	shortURL := intToShortURL(cur)

	if _, err := tx.Exec(ctx, `update public.url set short_url = $1 where id = $2`, shortURL, cur); err != nil {
		return "", err
	}

	return shortURL, nil
}