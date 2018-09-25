package text

import (
	"encoding/json"
	"fmt"
)

// ConfirmPhoneArgs is a struct containing all values needed to send the Confirm Phone text message
type ConfirmPhoneArgs struct {
	ToPhone          string
	VerificationCode string
}

// SendConfirmPhone is a function that is ran when a kafka message is received with the key "confirm-phone"
func SendConfirmPhone(b []byte) {
	// attempt to unmarshal ConfirmAccountArgs
	var args ConfirmPhoneArgs
	err := json.Unmarshal(b, &args)
	if err != nil {
		fmt.Println("There was an error unmarshalling ConfirmPhoneArgs")
	}

	to := args.ToPhone
	body := "Your verification code is " + args.VerificationCode

	send(to, body)
}
