package main

import (
	"context"
	"database/sql"
	"github.com/11011san/gator/internal/database"
	_ "github.com/lib/pq"
)
import (
	"github.com/11011san/gator/internal/config"
	"log"
	"os"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	conf, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	db, err := sql.Open("postgres", conf.DBURL)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
	dbQueries := database.New(db)
	programState := state{cfg: &conf, db: dbQueries}
	cmds := getCommands()

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
		return
	}
	var comArgs []string
	if len(os.Args) > 2 {
		comArgs = os.Args[2:]
	} else {
		comArgs = make([]string, 0)
	}
	com := command{Name: os.Args[1], Args: comArgs}
	err = cmds.run(&programState, com)
	if err != nil {
		log.Fatal(err)
	}

}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}
