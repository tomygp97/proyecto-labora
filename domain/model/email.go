package model

import (
	"database/sql"
)

type Email struct {
	ID          int
	MessageID   string         // Message-ID
	Sender      string         // From
	Receiver    string         // To
	Subject     string         // Subject
	MimeVersion string         // Mime-Version
	ContentType string         // Content-Type
	Encoding    string         // Content-Transfer-Encoding
	Folder      string         // X-Folder
	Body        string         // Contenido del email
	Date        sql.NullString // Date
}
