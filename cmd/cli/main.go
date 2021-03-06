package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/dish.io/cmd/cli/seed"
	"github.com/dish.io/internal/database/postgres"
	"github.com/dish.io/internal/services/user"
	_ "github.com/lib/pq"
	"os"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// DB initialisation
	db, err := sql.Open("postgres", "user=user dbname=dish-db password=longpassword port=7001 sslmode=disable")
	if err != nil {
		fmt.Printf("Could not open db: %s", err.Error())
		os.Exit(0)
	}
	if err := db.PingContext(ctx); err != nil {
		fmt.Printf("Could not ping db: %s", err.Error())
		os.Exit(0)
	}
	store := &postgres.Store{DB: db}
	// Service initialisation
	us := &user.Service{Store: store}
	// --------------- Commands ----------------
	// seeding database command
	seedCmd := flag.NewFlagSet("seed", flag.ExitOnError)
	// seed command arguments
	seedFile := seedCmd.String("filename", "seed/user.json", "Provide a json file containing seeding data to seed database")
	tableArgument := seedCmd.String("table", "users", "Specify the tables name to seed")

	// checking length of arguments
	if len(os.Args) < 2 {
		fmt.Println("use --help for more information")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "seed":
		seed.HandleSeed(ctx, seedCmd, tableArgument, seedFile, us)
	}
}
