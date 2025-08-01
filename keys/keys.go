package keys

import (
	"github.com/zalando/go-keyring"
)

var service = "monktypes-stream-commands"

func AddKey(label string, secret string) error {
	// set secret
	err := keyring.Set(service, label, secret)
	if err != nil {
		return err
	}
	return nil
}

func GetKey(label string) (string, error) {
	secret, err := keyring.Get(service, label)
	if err != nil {
		return "", err
	}

	return secret, nil
}
