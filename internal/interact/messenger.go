package interact

import (
	"net/url"
	"strings"
)

// QueryConvoURL requests and reads a conversation URL from os.Stdin.
func (p *Prompter) QueryConvoURL() (string, error) {
	var rawurl string
	for rawurl == "" {
		p.Println("Enter the target conversation URL " +
			"(https://www.messenger.com/t/...):")

		if _, err := p.scanf("%s", &rawurl); err != nil {
			return "", err
		}

		if rawurl == "" {
			p.Errln("Conversation URL cannot be empty.")
			continue
		}

		// Ensure that rawurl is valid.
		u, err := url.Parse(rawurl)
		if err != nil {
			p.Errf("Invalid URL: %s\n", err)
			rawurl = ""
			continue
		}

		if (u.Hostname() != "www.messenger.com") ||
			!strings.ContainsRune(u.EscapedPath(), 't') {
			p.Errf("URL is not of the form 'https://www.messenger.com/t/...'.\n")
			rawurl = ""
		}
	}
	return rawurl, nil
}
