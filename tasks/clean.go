package tasks

import (
	"time"

	"github.com/arken/arkstat/database"
	"github.com/arken/arkstat/mail"
)

func cleanUpOldNodes(db *database.DB, mailbox *mail.Mailbox) (err error) {
	// Grab a list of all of the old nodes in the DB not currently reporting.
	oldNodes, err := db.GetAllOlderThan(time.Hour * 24)
	if err != nil {
		return err
	}

	for _, node := range oldNodes {
		err = db.Remove(node.ID)
		if err != nil {
			return err
		}
		if node.Email != "" && mailbox != nil {
			err = mailbox.Send("mail/templates/alert.yml", node.Email)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
