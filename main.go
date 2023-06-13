package main

import (
	// C "docker/controllers"
	D "docker/database"
	R "docker/routes"
	U "docker/utils"
	"time"

	"github.com/gofiber/fiber/v2/middleware/session"
)

func main() {
	// stablish connection to the database
	D.ConnectToDB(
		U.Env("DB_HOST"),
		U.Env("DB_USERNAME"),
		U.Env("DB_PASSWORD"),
		U.Env("DB_NAME"),
		U.Env("DB_PORT"),
	)
	// create woocommerce instance
	U.InitWoomeCommerce(U.Env("WC_CONSUMER_KEY"), U.Env("WC_CONSUMER_SECRET"), U.Env("WC_SHOP_NAME"))
	// C.Initilize() // initialize controllers value
	// ! session
	U.Store = session.New(session.Config{
		CookieHTTPOnly: true,
		Expiration:     time.Hour * 8760,
		KeyGenerator: func() string {
			// secretKey, err := C.GetEnvVar("SESSION_SECRET_KEY")
			secretKey := U.Env("SESSION_SECRET_KEY")
			var sessionID string
			sessionID, err := U.GenerateSessionID(secretKey)
			if err != nil {
				panic(err)
			}
			return sessionID
		},
	})
	R.RouterInit()

}
