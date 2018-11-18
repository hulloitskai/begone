package main

import (
	"fmt"
	"os"

	"github.com/stevenxie/begone/messenger"
)

// loadConfig loads a messenger.Config by reading it from a file based on
// the program name.
//
// An empty Config will be returned if no such file exists.
func loadConfig() (cfg *messenger.Config, err error) {
	cfg = new(messenger.Config)
	if err := messenger.ReadConfig(Namespace, cfg); err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}

	return cfg, nil
}

// queryMissing fills missing config values from cfg with values read from
// os.Stdin.
func queryMissing(cfg *messenger.Config) error {
	for cfg.User == "" {
		user, err := queryUser()
		if err != nil {
			return err
		}

		cfg.User = user
	}

	for cfg.Pass == "" {
		pass, err := queryPass()
		if err != nil {
			return err
		}

		cfg.Pass = pass
	}

	for cfg.ConvoID == "" {
		fmt.Println("A convoID is found in the Messenger URL as " +
			"https://messenger.com/t/{{convoID}}.")
		fmt.Print("Enter the target convoID: ")
		if _, err := fmt.Scanf("%s", &cfg.ConvoID); err != nil {
			return fmt.Errorf("begone: failed to read convoID: %s", err)
		}
	}

	return nil
}
