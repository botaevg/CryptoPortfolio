package db

import (
	"context"
	"github.com/jackc/pgx/v4"
	"log"
)

func NewClient(dsn string) (*pgx.Conn, error) {
	con, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Print("DB connect error")

	}

	q := `CREATE TABLE users(
	id serial primary key,
    username VARCHAR(100),
	password VARCHAR(100)
	);`
	_, err = con.Exec(context.Background(), q)
	if err != nil {
		log.Print("ТАБЛИЦА НЕ СОЗДАНА")
		log.Print(err)
	}

	q = `CREATE UNIQUE INDEX users_unique ON users USING btree(username);`
	_, err = con.Exec(context.Background(), q)
	if err != nil {
		log.Print("Unique не создана")
		log.Print(err)
	}

	q = `CREATE TABLE portfolio(
	id serial primary key,
    userid int,
	name VARCHAR(100)
	);`
	_, err = con.Exec(context.Background(), q)
	if err != nil {
		log.Print("ТАБЛИЦА НЕ СОЗДАНА")
		log.Print(err)
	}

	q = `CREATE TABLE coinaction(
	id serial primary key,
    portfolioid int,
	name VARCHAR(100),
	amount double precision,
	pricepurchase double precision
	);`
	_, err = con.Exec(context.Background(), q)
	if err != nil {
		log.Print("ТАБЛИЦА НЕ СОЗДАНА")
		log.Print(err)
	}

	return con, nil
}
