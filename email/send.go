package email

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/matcornic/hermes"
	gomail "gopkg.in/gomail.v2"
)

// configuration struct created from environment variables
type envConfig struct {
	FromEmail string `envconfig:"FROM_EMAIL"`
	SmtpHost  string `envconfig:"SMTP_HOST"`
	SmtpPort  int    `envconfig:"SMTP_PORT"`
	SmtpUser  string `envconfig:"SMTP_USER"`
	SmptPass  string `envconfig:"SMTP_PASS"`
}

// send is a base function used to send out different email templates
func send(to, subject, text, html string) {
	// load in secret env vars from file in development
	if os.Getenv("ENV") == "dev" {
		_ = godotenv.Load(".env.dev")
	}

	// attempt to cast env variables into envConfig struct
	var cfg envConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(cfg.FromEmail)

	// build out the email
	m := gomail.NewMessage()
	m.SetHeader("From", cfg.FromEmail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", text)
	m.AddAlternative("text/html", html)

	// create a new SMTP connection
	d := gomail.NewPlainDialer(cfg.SmtpHost, cfg.SmtpPort, cfg.SmtpUser, cfg.SmptPass)

	// send the created email
	if err := d.DialAndSend(m); err != nil {
		fmt.Printf("Failed to send email %v", err)
	}
}

// getHeader sets config values that are global to each email sent
func getHeader() hermes.Hermes {
	return hermes.Hermes{
		// Theme: new(Salted),
		Product: hermes.Product{
			Name: "Company Name",
			// Link: "https://dakotahipp.com/",
			// Logo: "http://www.duchess-france.org/wp-content/uploads/2016/01/gopher.png",
			Copyright: "Copyright © 2018",
		},
	}
}
