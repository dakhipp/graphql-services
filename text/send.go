package text

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// configuration struct created from environment variables
type envConfig struct {
	FromPhone  string `envconfig:"FROM_PHONE"`
	TwilioSID  string `envconfig:"TWILIO_SSID"`
	TwilioAuth string `envconfig:"TWILIO_AUTH"`
}

// send is a base function used to send out different email templates
func send(to, body string) {
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

	// API URL endpoint
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + cfg.TwilioSID + "/Messages.json"

	// build message
	msgData := url.Values{}
	msgData.Set("To", to)
	msgData.Set("From", cfg.FromPhone)
	msgData.Set("Body", body)
	msgDataReader := *strings.NewReader(msgData.Encode())

	// create HTTP request client
	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.SetBasicAuth(cfg.TwilioSID, cfg.TwilioAuth)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// make HTTP POST request
	resp, _ := client.Do(req)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err == nil {
			fmt.Println("Successfully sent Twilio SMS")
		} else {
			fmt.Printf("Error sending Twilio SMS %v\n", resp.Status)
		}
	} else {
		fmt.Printf("Error sending Twilio SMS %v\n", resp.Status)
	}
}
