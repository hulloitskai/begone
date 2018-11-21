package cmd

import (
	"fmt"
	"os"

	"github.com/howeyc/gopass"
)

// queryUsername requests and reads a username from stdin.
func queryUsername() (string, error) {
	var user string
	for user == "" {
		fmt.Print("Enter your FB Messenger username / email: ")
		if _, err := fmt.Scanf("%s", &user); err != nil {
			return "", err
		}

		if user == "" {
			fmt.Fprintln(os.Stderr, "Username must be non-empty!")
		}
	}
	return user, nil
}

// queryPassword requests and reads a password from stdin.
func queryPassword() (string, error) {
	var pass string
	for pass == "" {
		fmt.Print("Enter your FB Messenger password: ")
		passbytes, err := gopass.GetPasswdMasked()
		pass = string(passbytes)

		if err != nil {
			if err == gopass.ErrInterrupted {
				return pass, nil
			}
			return "", err
		}

		if pass == "" {
			fmt.Fprintln(os.Stderr, "Password must be non-empty!")
		}
	}
	return pass, nil
}
