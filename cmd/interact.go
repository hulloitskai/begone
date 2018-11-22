package cmd

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/howeyc/gopass"
	"github.com/stevenxie/begone/config"
)

// queryUsername requests and reads a username from stdin.
func queryUsername() (string, error) {
	var user string
	for user == "" {
		fmt.Print("Enter your FB Messenger email: ")
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
		if err != nil {
			return "", err
		}

		pass = string(passbytes)
		if pass == "" {
			fmt.Fprintln(os.Stderr, "Password must be non-empty!")
		}
	}
	return pass, nil
}

// queryMissing fills missing fields of cfg with values read from stdin.
func queryMissing(cfg *config.Config) error {
	if cfg.Username == "" {
		uname, err := queryUsername()
		if err != nil {
			return err
		}
		cfg.Username = uname
	}

	if cfg.Password == "" {
		pw, err := queryPassword()
		if err != nil {
			return err
		}
		cfg.Password = pw
	}

	return nil
}

// queryConvoID requests and reads a conversation ID from stdin.
func queryConvoID() (string, error) {
	var convoID string
	for convoID == "" {
		fmt.Println("Enter the target conversation URL " +
			"(https://messenger.com/t/...):")

		var rawurl string
		if _, err := fmt.Scanf("%s", &rawurl); err != nil {
			return "", err
		}
		convoID = parseConvoURL(rawurl)
	}
	return convoID, nil
}

// parseConvoURL will parse the url of an FB Messenger conversation (rawurl)
// into an conversation ID.
//
// If parseConvoURL fails, it will print an error message to os.Stderr, and
// return an empty string.
func parseConvoURL(rawurl string) (convoID string) {
	u, err := url.Parse(rawurl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Bad URL: %v\n", err)
		return ""
	}

	path := u.EscapedPath()
	slashIndex := strings.LastIndexByte(path, '/')
	if slashIndex == -1 {
		fmt.Fprintf(os.Stderr, "Bad URL: path does not contain any '/'\n")
		return ""
	}

	return path[slashIndex+1:]
}
