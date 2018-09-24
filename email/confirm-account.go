package email

import (
	"encoding/json"
	"fmt"

	"github.com/matcornic/hermes"
)

// ConfirmAccountArgs is a struct containing all values needed to the Confirm Account email template
type ConfirmAccountArgs struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

// SendConfirmAccountEmail is a function that is ran when a kafka message is received with the key "confirm-account"
func SendConfirmAccountEmail(b []byte) {
	// attempt to unmarshal ConfirmAccountArgs
	var args ConfirmAccountArgs
	err := json.Unmarshal(b, &args)
	if err != nil {
		fmt.Println("There was an error unmarshalling ConfirmAccountArgs")
	}

	// set up basic email config
	h := getHeader()
	e := getConfirmAccountEmail(args)
	to := args.Email
	subject := "Confirm your account"

	// attempt to create html body from the tempalate
	html, err1 := h.GenerateHTML(e)
	if err1 != nil {
		fmt.Println("Error HTML email template")
	}

	// attempt to create text body from the tempalate
	text, err2 := h.GeneratePlainText(e)
	if err2 != nil {
		fmt.Println("Error text email template")
	}

	// if there are no errors creating the tempalate, attempt to send the email
	if err1 == nil && err2 == nil {
		send(to, subject, text, html)
	}
}

// getConfirmAccountEmail takes ConfirmAccountArgs and returns a hermes email template
func getConfirmAccountEmail(args ConfirmAccountArgs) hermes.Email {
	return hermes.Email{
		Body: hermes.Body{
			Name: args.FirstName + " " + args.LastName,
			Intros: []string{
				"Welcome! We're very excited to have you on board.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "To get started with Hermes, please click here:",
					Button: hermes.Button{
						Color:     "#22BC66", // Optional action button color
						Text:      "Confirm your account",
						TextColor: "#FFF",
						Link:      "https://hermes-example.com/confirm?token=d9729feb74992cc3482b350163a1a010",
					},
				},
			},
			Outros: []string{
				"Need help, or have questions? Just reply to this email, we'd love to help.",
			},
		},
	}
}
