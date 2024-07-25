package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/yardnsm/gohever"
	"github.com/yardnsm/leeches/internal/bot"
	"github.com/yardnsm/leeches/internal/commands"
	"github.com/yardnsm/leeches/internal/config"
	"github.com/yardnsm/leeches/internal/db"
	"github.com/yardnsm/leeches/internal/model"

	tele "gopkg.in/telebot.v3"
)

const usage = `Usage:
    leeches [flags]

Options:
    --config              Path to config JSON file
    --credentials         Path to credentials parcel file
    --credit-card         Path to credit card parcel file
`

var (
	configPath            string
	credentialsParcelPath string
	creditCardParcelPath  string
)

func init() {
	flag.Usage = func() { fmt.Fprintf(os.Stderr, "%s\n", usage) }

	flag.StringVar(&configPath, "config", "", "path to config json file")
	flag.StringVar(&credentialsParcelPath, "credentials", "", "path to credentials parcel file")
	flag.StringVar(&creditCardParcelPath, "credit-card", "", "path to credit card parcel file")
	flag.Parse()

	if configPath == "" || credentialsParcelPath == "" || creditCardParcelPath == "" {
		log.Fatal("Please provide valid parameters")
		return
	}
}

func main() {
	log.Println("Starting bot")

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("Config loaded")

	credentials, err := config.LoadCredentialsConfig(credentialsParcelPath, []byte(cfg.CredentialsParcelPassword))
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("Credentials parcel decrypted and loaded")

	creditCard, err := config.LoadCreditCardConfig(creditCardParcelPath, []byte(cfg.CreditCardParcelPassword))
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("Credit Card parcel decrypted and loaded")

	botSettings := tele.Settings{
		Token:     cfg.TelegramToken,
		Poller:    &tele.LongPoller{Timeout: 10 * time.Second},
		ParseMode: tele.ModeMarkdown,
	}

	// Support for webhook poller. This is tailored for a running instance behind a revproxy / load
	// balancer, as the webhook will listen in plain HTTP and the LB should handle the TLS shit
	if cfg.Webhook.Port != "" {
		botSettings.Poller = &tele.Webhook{
			Listen: ":" + cfg.Webhook.Port,
			Endpoint: &tele.WebhookEndpoint{
				PublicURL: cfg.Webhook.PublicURL,
				Cert:      cfg.Webhook.Cert,
			},
		}
	}

	hvrClientConfig := gohever.Config{
		Credentials: gohever.BasicCredentials(credentials.Username, credentials.Password),
		CreditCard:  gohever.BasicCreditCard(creditCard.Number, creditCard.Month, creditCard.Year),
	}

	b, err := tele.NewBot(botSettings)
	if err != nil {
		log.Fatal(err)
		return
	}

	db := db.CreateDatabase(cfg.Database)
	usersRepository := model.NewUsersRepository(db)
	chargeRequestsRepository := model.NewChargeRequestsRepository(db)

	hvr := gohever.NewClient(hvrClientConfig)

	// Setup global middlewares
	b.Use(bot.AllowOnlyPrivateChatsMiddleware())
	b.Use(bot.SendErrorsToUsersChatMiddleware())

	router := bot.NewRouter().
		DefaultRestrictUser(bot.RestrictApproved).
		CreateContext(func(t tele.Context) bot.Context {

			// We don't care about errors here
			currentUser, _ := usersRepository.GetByTelegramID(t.Sender().ID)

			context := bot.NewContext(t)

			context.CurrentUser = currentUser
			context.Users = usersRepository
			context.ChargeRequests = chargeRequestsRepository
			context.Hever = hvr

			return context
		})

	// Default text command
	router.AddCommand(
		bot.NewCommand(tele.OnText).
			RestrictUser(bot.RestrictNone).
			Handle(func(c bot.Context, t tele.Context) error {
				if c.CurrentUser != nil {

					// Default behaviour: set user commands
					// For private chats with users, the chat id is the user id
					userCommands := router.GetUserCommands(*c.CurrentUser)
					bot.SetCommandsForChat(b, userCommands, c.CurrentUser.TelegramID)
				}

				// Secret text command to get the chat id and the user id
				if t.Message().Text == "leechmeup" {
					log.Printf("leechmeup was invoked for server: %v", t.Sender())
					t.Send(
						fmt.Sprintf("*User ID:* `%d`,\n*Chat ID:* `%d`", t.Sender().ID, t.Chat().ID),
						tele.ModeMarkdownV2,
					)
				}

				return nil
			}),
	)

	commands.Attach(router)
	router.Attach(b)

	// Delete the webhook if not used
	if cfg.Webhook.Port == "" {
		b.RemoveWebhook()
	}

	b.Start()
}
