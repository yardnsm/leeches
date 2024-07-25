package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/yardnsm/leeches/internal/config"
	"github.com/yardnsm/leeches/internal/db"
	"github.com/yardnsm/leeches/internal/model"
)

const usage = `Usage:
    adduser [flags]

Options:
    -c, --config          Path to config JSON file
    -d, --display-name    Display name for the user.
    -t, --telegram-id     The Telegram user ID.
    -a, --admin	          Whether the user is an admin.
`

var (
	configPath  string
	displayName string
	telegramId  int64
	admin       bool
)

func init() {
	flag.Usage = func() { fmt.Fprintf(os.Stderr, "%s\n", usage) }

	flag.StringVar(&configPath, "c", "", "path to config json file")
	flag.StringVar(&configPath, "config", "", "path to config json file")

	flag.StringVar(&displayName, "d", "", "display name for the user")
	flag.StringVar(&displayName, "display-name", "", "display name for the user")

	flag.Int64Var(&telegramId, "t", 0, "the telegram user id")
	flag.Int64Var(&telegramId, "telegram-id", 0, "the telegram user id")

	flag.BoolVar(&admin, "a", false, "whether the user is an admin")
	flag.BoolVar(&admin, "admin", false, "whether the user is an admin")

	flag.Parse()

	if configPath == "" || displayName == "" || telegramId == 0 {
		log.Fatal("Please provide valid parameters")
		return
	}
}

func main() {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	db := db.CreateDatabase(cfg.Database)
	usersRepository := model.NewUsersRepository(db)

	existing, _ := usersRepository.GetByTelegramID(telegramId)
	if existing != nil {
		log.Fatalf("User with TelegramID %d already exists", telegramId)
		return
	}

	user := &model.User{
		TelegramID:  telegramId,
		DisplayName: displayName,
		IsApproved:  true,
		IsAdmin:     admin,
	}

	err = usersRepository.Create(user)
	if err != nil {
		log.Fatalf("Unable to create a new user: %w", err)
		return
	}

	log.Printf("New user created! %+v", *user)
}
