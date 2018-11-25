package interact

import (
	"github.com/howeyc/gopass"
	"github.com/stevenxie/begone/internal/config"
)

// QueryUsername requests and reads a username on os.Stdin.
func (p *Prompter) QueryUsername() (string, error) {
	var user string
	for user == "" {
		p.Println("Enter your FB Messenger email:")
		if _, err := p.scanf("%s", &user); err != nil {
			return "", err
		}

		if user == "" {
			p.Errln("Email cannot be empty!")
			p.Println()
		}
	}
	return user, nil
}

// QueryPassword requests and reads a password from os.Stdin.
func (p *Prompter) QueryPassword() (string, error) {
	var pass string
	for pass == "" {
		p.Println("Enter your FB Messenger password:")
		passbytes, err := gopass.GetPasswdMasked()
		if err != nil {
			return "", err
		}

		pass = string(passbytes)
		if pass == "" {
			p.Errln("Password cannot be empty!")
			p.Println()
		}
	}
	return pass, nil
}

// QueryMissing fills missing fields of cfg with values read from os.Stdin.
func (p *Prompter) QueryMissing(cfg *config.Config, skipPass bool) error {
	if cfg.Username == "" {
		user, err := p.QueryUsername()
		if err != nil {
			return err
		}
		cfg.Username = user
	}

	if !skipPass && (cfg.Password == "") {
		p.Println()
		pass, err := p.QueryPassword()
		if err != nil {
			return err
		}
		cfg.Password = pass
	}

	return nil
}
