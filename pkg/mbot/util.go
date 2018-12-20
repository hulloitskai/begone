package mbot

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	ess "github.com/unixpickle/essentials"
)

// fbEntity is either an FB user, or an FB group.
type fbEntity struct {
	ID      string
	IsGroup bool
}

// definesToken indicates the section of HTML that defines FB IDs in terms of
// convoIDs.
const definesToken = "require(\"ServerJSDefine\").handleDefines"

// extractFBEntity constructs an fbEntity using conversation URL (convoURL).
func (b *Bot) extractFBEntity(convoURL string) (*fbEntity, error) {
	var (
		id     = convoURL[strings.Index(convoURL, "/t/")+3:]
		entity fbEntity
	)

	// If convoID starts with a number, then it is corresponds to a group fbid.
	if fchar := id[0]; ('0' <= fchar) && (fchar <= '9') {
		entity.ID = id
		entity.IsGroup = !b.Cfg.AssumeUser
		return &entity, nil
	}

	// Get corresponding messenger thread.
	entity.IsGroup = false
	res, err := b.Session.Client.Get(convoURL)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("received non-200 response (code: %d)",
			res.StatusCode)
	}
	defer res.Body.Close()

	// Read HTML from response body, and close body.
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, ess.AddCtx("reading response body", err)
	}
	if err = res.Body.Close(); err != nil {
		return nil, ess.AddCtx("closing response body", err)
	}

	// Parse for fbid by iteratively narrowing down search scope.
	var (
		html         = string(data)
		definesIndex = strings.Index(html, definesToken)
	)
	if definesIndex == -1 {
		return nil, errors.New("couldn't find \"defines index\" in HTML response")
	}
	html = html[definesIndex:]

	convoIDIndex := strings.Index(html, id)
	if convoIDIndex == -1 {
		return nil, fmt.Errorf("couldn't find convoID '%s' in HTML response", id)
	}
	html = html[:convoIDIndex]

	fbidOpenIndex := strings.LastIndex(html, "\"id\":\"")
	if fbidOpenIndex == -1 {
		return nil, fmt.Errorf("couldn't find fbid corresponding to convoID '%s' "+
			"in HTML response", id)
	}
	html = html[fbidOpenIndex+6:]

	fbidCloseIndex := strings.IndexByte(html, '"')
	if fbidCloseIndex == -1 {
		return nil, errors.New("couldn't find closing quote of target fbid in " +
			"HTML response")
	}

	entity.ID = html[:fbidCloseIndex]
	return &entity, nil
}
