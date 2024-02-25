package main

import (
	"context"
	"database/sql"
	"flag"
	"github.com/Georgy27/go-grpc/week_2/postgres/internal/config"
	"github.com/Georgy27/go-grpc/week_2/postgres/internal/config/env"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v4/pgxpool"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

func main() {
	flag.Parse()
	ctx := context.Background()

	// Считываем переменные окружения
	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	pgConfig, err := env.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	// Создаем пул соединений с базой данных
	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Делаем запрос на вставку записи в таблицу note
	builderInsert := sq.Insert("note").
		PlaceholderFormat(sq.Dollar).
		Columns("title", "body").
		Values(gofakeit.City(), gofakeit.Address().Street).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	var noteID int
	err = pool.QueryRow(ctx, query, args...).Scan(&noteID)
	if err != nil {
		log.Fatalf("failed to insert note: %v", err)
	}

	log.Printf("inserted note with id: %d", noteID)

	// Делаем запрос на выборку записей из таблицы note
	builderSelect := sq.Select("id", "title", "body", "created_at", "updated_at").
		From("note").
		PlaceholderFormat(sq.Dollar).
		OrderBy("id ASC").
		Limit(10)

	query, args, err = builderSelect.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to select notes: %v", err)
	}

	var id int
	var title, body string
	var createdAt time.Time
	var updatedAt sql.NullTime

	for rows.Next() {
		err = rows.Scan(&id, &title, &body, &createdAt, &updatedAt)
		if err != nil {
			log.Fatalf("failed to scan note: %v", err)
		}

		log.Printf("id: %d, title: %s, body: %s, created_at: %v, updated_at: %v\n", id, title, body, createdAt, updatedAt)
	}

	// Делаем запрос на обновление записи в таблице note
	builderUpdate := sq.Update("note").
		PlaceholderFormat(sq.Dollar).
		Set("title", gofakeit.City()).
		Set("body", gofakeit.Address().Street).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": noteID})

	query, args, err = builderUpdate.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	res, err := pool.Exec(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to update note: %v", err)
	}

	log.Printf("updated %d rows", res.RowsAffected())

	// Делаем запрос на получение измененной записи из таблицы note
	builderSelectOne := sq.Select("id", "title", "body", "created_at", "updated_at").
		From("note").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": noteID}).
		Limit(1)

	query, args, err = builderSelectOne.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	err = pool.QueryRow(ctx, query, args...).Scan(&id, &title, &body, &createdAt, &updatedAt)
	if err != nil {
		log.Fatalf("failed to select notes: %v", err)
	}

	log.Printf("id: %d, title: %s, body: %s, created_at: %v, updated_at: %v\n", id, title, body, createdAt, updatedAt)
}
