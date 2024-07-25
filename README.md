# Leeches 🪱

There is this consumer club in Israel called "Hever". It's pretty awesome - you get a chargeable
credit card which have discounts in TONS of businesses, but... The website used to charge it? Sucks.
Terrible. Slow.

I built this Telegram bot in order to simplify my experience with it. It allows me to get the
current card's status and to charge it, in the simplest way possible. I can also give access to my
family and friends, and they'll be able to use the bot the same as me. However, when want to charge
the card, the bot will notify me and ask me to confirm the action.

Sweet, isn't it?

> [!WARNING]
> This project was meant to be used for educational purposes only. I am not affiliated with Hever in
> any way.

## How?

The bot is written on Golang. I used the amazing [telebot](https://github.com/tucnak/telebot)
framework to interact with Telegram, and built a simple [wrapper](./internal/bot) around it so I can pass
around my extended context when handling incoming commands. I've also built an API client for Hever,
called [gohever](http://github.com/yardnsm/gohever).

## Security

Two problems came in mind when I initially designed this bot:

- I don't want my credentials to be plainly stored on the disk, but I'm fine with them stored in
  memory.
- I don't want my credit card details plainly stored on the disk, and I rather minimize the time of
  them stored in memory.

The credentials and the credit card details saved in an encrypted JSON file on disk. It uses AES-GCM
with a key derived using a scrypt key derivation function, which is damn slow and memory intensive.
The credentials and credit card details will be loaded and decrypted when the bot starts.

Regarding the credit card details, I originally wanted to load them to memory only when needed,
while asking the user each time for the password. I realized this is hard to implement due to the
way I implemented the "state machine" of the bot, so...

## Setting up

### Locally

1. Clone this repo on your machine:

    ```
    $ git clone https://github.com/yardnsm/leeches
    ```

1. Create a parcel containing the credentials for the Hever website:

    ```
    $ echo '{"username": "username", "password": "password"}' | go run ./cmd/parcel -e > ./configs/credentials.parcel
    ```

1. Create a parcel containing the credit card info:

    ```
    $ echo '{"number": "", "year": "", "month": ""}' | go run ./cmd/parcel -e > ./configs/credit_card.parcel
    ```

1. Create a configuration file and fill the relevant info:

    ```
    $ cp configs/config{.sample,}.json
    ```

1. Start the bot

    ```
    $ go run ./cmd/leeches \
        --config ./configs/config.json \
        --credentials ./configs/credentials.parcel \
        --credit-card ./configs/credit_card.parcel
    ```

1. Next, you'll need to create a user for yourself. Start a Telegram chat with your bot, and type
   the following secret command: `leechmeup`.

   It'll print out your user's id. Next, create a new user and restart the bot:

    ```
    $ go run ./cmd/adduser \
        --config config.json \
        --admin \
        --display-name "Israel Israeli" \
        --telegram-id 123454321
    ```

---

## License

MIT © [Yarden Sod-Moriah](https://ysm.sh/)
