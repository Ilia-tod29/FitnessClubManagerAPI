package main

import (
	"FitnessClubManagerAPI/api"
	db "FitnessClubManagerAPI/db/sqlc"
	"FitnessClubManagerAPI/util"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/stripe/stripe-go/v76"
	"log"
)

//const (
//	dbDriver      = "postgres"
//	dbSource      = "postgresql://root:root@localhost:5432/fitness_club_manager?sslmode=disable"
//	serverAddress = "0.0.0.0:8080"
//)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config")
	}
	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db", err)
	}

	stripe.Key = config.StripesSecretKey

	store := db.NewStore(connPool)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot initiate the server", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
