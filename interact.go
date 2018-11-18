package main

import (
	"fmt"

	"github.com/howeyc/gopass"
)

// queryUser gets an FB Messenger username from os.Stdin.
func queryUser() (string, error) {
	var user string

	fmt.Print("Enter your FB Messenger username / email: ")
	if _, err := fmt.Scanf("%s", &user); err != nil {
		return "", fmt.Errorf("begone: failed to read username: %v", err)
	}

	return user, nil
}

func queryPass() (string, error) {
	fmt.Print("Enter your FB Messenger password: ")
	passbytes, err := gopass.GetPasswdMasked()
	if err != nil {
		if err == gopass.ErrInterrupted {
			return "", nil
		}

		return "", fmt.Errorf("begone: failed to read password: %v", err)
	}

	return string(passbytes), nil
}
