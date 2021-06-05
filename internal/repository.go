package internal

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	pgx "github.com/jackc/pgx/v4"
)

// Repository represents
type PostRepository struct {
	logger *log.Logger
	dbconn *pgx.Conn
}

// NewPostRepository
func NewPostRepository(dbconn *pgx.Conn) *PostRepository {
	return &PostRepository{
		log.New(os.Stderr, repositoryLogPrefix, loggerFlags),
		dbconn,
	}
}

// Save
func (r *PostRepository) Save(ctx context.Context, posts <-chan Post, done chan<- struct{}) {
	r.logger.Println("Start saving posts")

	const (
		batchMaxSize = 10
		argsCount    = 3
		lockTimeout  = 100 * time.Millisecond
	)

	var (
		valPlaceholders = make([]string, 0, batchMaxSize)
		valToInsert     = make([]interface{}, 0, argsCount*batchMaxSize)
		valCounter      = 1
	)

	for post := range posts {
		linkPH, titlePH, descriptionPH := valCounter, valCounter+1, valCounter+2
		valCounter += argsCount

		valPlaceholders = append(
			valPlaceholders,
			fmt.Sprintf("($%d, $%d, $%d)", linkPH, titlePH, descriptionPH),
		)

		valToInsert = append(valToInsert, post.Link, post.Title, post.Description)

		if len(valPlaceholders) == batchMaxSize {
			err := r.batchInsert(ctx, valPlaceholders, valToInsert)
			if err != nil {
				r.logger.Printf("could not insert batch: %s", err)
			}

			valPlaceholders = valPlaceholders[:0]
			valToInsert = valToInsert[:0]
			valCounter = 1
		}
	}

	r.batchInsert(ctx, valPlaceholders, valToInsert)

	done <- struct{}{}
}

// batchInsert
func (r *PostRepository) batchInsert(ctx context.Context, placeholders []string, values []interface{}) error {
	query := fmt.Sprintf(
		"INSERT INTO post (link, title, description) VALUES %s ON CONFLICT DO NOTHING",
		strings.Join(placeholders, ","),
	)

	ct, err := r.dbconn.Exec(ctx, query, values...)
	if err != nil {
		r.logger.Printf("could not insert row: %s", err)
	}

	r.logger.Printf("%d new lines inserted", ct.RowsAffected())

	return nil
}

// Search
func (r *PostRepository) Search(ctx context.Context, param string, limit int) ([]Post, error) {
	query := "SELECT title, link, description FROM post WHERE titletsv @@ to_tsquery('russian', $1) LIMIT $2"

	rows, err := r.dbconn.Query(ctx, query, param, limit)
	if err != nil {
		return nil, fmt.Errorf("could not query post: %w", err)
	}
	defer rows.Close()

	posts := make([]Post, 0, limit)

	for rows.Next() {
		p := Post{}

		err = rows.Scan(&p.Title, &p.Link, &p.Description)
		if err != nil {
			return nil, fmt.Errorf("could not get result values: %w", err)
		}

		posts = append(posts, p)
	}

	return posts, nil
}
