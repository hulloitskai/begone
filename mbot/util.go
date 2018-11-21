package mbot

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/stevenxie/fbmsgr"
	ess "github.com/unixpickle/essentials"
)

// parseUserConvoID parses a conversation ID (found in a messenger.com URL) to
// an internal Facebook user ID (an 'fbid').
func (b *Bot) parseUserConvoID(id string) (fbid string, err error) {
	// Get corresponding messenger thread.
	url := fmt.Sprintf("%s/t/%s", fbmsgr.BaseURL, id)
	res, err := b.Session.Client.Get(url)
	if err != nil {
		return "", err
	}
	if res.StatusCode != 200 {
		return "", fmt.Errorf("recieved non-200 response (code: %d)",
			res.StatusCode)
	}
	defer res.Body.Close()

	// Read HTML from response body, and close body.
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", ess.AddCtx("reading response body", err)
	}
	if err = res.Body.Close(); err != nil {
		return "", ess.AddCtx("closing response body", err)
	}

	// Parse for fbid by iteratively narrowing down search scope.
	var (
		html         = string(data)
		definesIndex = strings.Index(html, definesToken)
	)
	if definesIndex == -1 {
		return "", errors.New("couldn't find \"defines index\" in HTML response")
	}
	html = html[definesIndex:]

	convoIDIndex := strings.Index(html, id)
	if convoIDIndex == -1 {
		return "", fmt.Errorf("couldn't find convoID '%s' in HTML response", id)
	}
	html = html[:convoIDIndex]

	fbidOpenIndex := strings.LastIndex(html, "\"id\":\"")
	if fbidOpenIndex == -1 {
		return "", fmt.Errorf("couldn't find fbid corresponding to convoID '%s' "+
			"in HTML response", id)
	}
	html = html[fbidOpenIndex+6:]

	fbidCloseIndex := strings.IndexByte(html, '"')
	if fbidCloseIndex == -1 {
		return "", errors.New("couldn't find closing quote of target fbid in " +
			"HTML response")
	}

	return html[:fbidCloseIndex], nil
}
