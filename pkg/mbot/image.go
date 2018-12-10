package mbot

import (
	"os"
	"path/filepath"

	"github.com/stevenxie/fbmsgr"
	ess "github.com/unixpickle/essentials"
)

type imageCycler struct {
	Image   *fbmsgr.UploadResult
	FBID    string
	IsGroup bool
}

// Cycle implements Cycler's Cycle method.
func (ic *imageCycler) Cycle(session *fbmsgr.Session) error {
	var err error
	if ic.IsGroup {
		_, err = session.SendGroupAttachment(ic.FBID, ic.Image)
	} else {
		_, err = session.SendAttachment(ic.FBID, ic.Image)
	}
	return err
}

// Finished implements Cycler's Continue method.
func (ic *imageCycler) Finished() bool {
	return false
}

// CycleImage repeatedly sends an image file located at fpath to a conversation
// identified by convoURL.
func (b *Bot) CycleImage(convoURL, fpath string) error {
	file, err := os.Open(fpath)
	if err != nil {
		return ess.AddCtx("mbot: opening image file", err)
	}

	// Upload file to Facebook.
	res, err := b.Upload(filepath.Base(file.Name()), file)
	if err != nil {
		return err
	}

	// Close file.
	if err = file.Close(); err != nil {
		return ess.AddCtx("mbot: closing image file", err)
	}

	entity, err := b.extractFBEntity(convoURL)
	if err != nil {
		return ess.AddCtx("mbot: looking up fbid", err)
	}

	return b.CycleUsing(&imageCycler{
		Image:   res,
		FBID:    entity.ID,
		IsGroup: entity.IsGroup,
	})
}
