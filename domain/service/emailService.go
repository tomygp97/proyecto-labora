package service

import (
	"project/domain/model"
	zincSearchClient "project/infrastructure/zincsearch"

	"database/sql"
	"fmt"
)

// EmailService es el servicio que maneja las operaciones sobre los correos electrónicos.
type EmailService struct {
	db *sql.DB
}

// NewEmailService crea una nueva instancia de EmailService.
func NewEmailService(db *sql.DB) *EmailService {
	return &EmailService{db: db}
}

func (es *EmailService) GetEmailsWithPagination(offset int, limit int) ([]model.Email, int, error) {
	var emails []model.Email

	query := `SELECT id, message_id, sender, receiver, subject, mime_version, content_type, encoding, folder, body, date 
              FROM emails LIMIT ? OFFSET ?`

	// Ejecutar la consulta con los parámetros limit y offset
	rows, err := es.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// Leer los resultados de la consulta y mapearlos a la estructura Email
	for rows.Next() {
		var email model.Email
		err := rows.Scan(&email.ID, &email.MessageID, &email.Sender, &email.Receiver, &email.Subject,
			&email.MimeVersion, &email.ContentType, &email.Encoding, &email.Folder, &email.Body, &email.Date)
		if err != nil {
			return nil, 0, err
		}
		emails = append(emails, email)
	}

	// Obtener el total de correos
	var total int
	err = es.db.QueryRow("SELECT COUNT(*) FROM emails").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return emails, total, nil
}

// GetEmailByID obtiene un correo electrónico por ID.
func (es *EmailService) GetEmailByID(id int) (*model.Email, error) {
	fmt.Printf("Iniciando consulta para obtener correo con ID: %d\n", id)
	var email model.Email

	err := es.db.QueryRow(`
		SELECT id, message_id, sender, receiver, subject, mime_version, content_type, encoding, folder, body, date 
		FROM emails 
		WHERE id = ?`, id).
		Scan(&email.ID, &email.MessageID, &email.Sender, &email.Receiver, &email.Subject, &email.MimeVersion,
			&email.ContentType, &email.Encoding, &email.Folder, &email.Body, &email.Date)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("Correo no encontrado para ID: %d\n", id)
			return nil, nil
		}
		fmt.Printf("Error al ejecutar consulta para ID %d: %v\n", id, err)
		return nil, fmt.Errorf("error al obtener el correo: %w", err)
	}

	fmt.Printf("Correo encontrado para ID: %d\n", id)
	return &email, nil
}

// función que hace la búsqueda en ZincSearch
func (es *EmailService) SearchEmailsWithPagination(query string, offset int, limit int) ([]model.Email, int, error) {
	if query == "" {
		return nil, 0, fmt.Errorf("el término de búsqueda no puede estar vacío")
	}

	searchQuery := fmt.Sprintf(`{
        "query": {
            "match": {
                "subject": "%s"
            }
        },
        "from": %d,
        "size": %d
    }`, query, offset, limit)

	emails, total, err := zincSearchClient.NewZincSearchClient().SearchEmailsWithPagination(searchQuery)
	if err != nil {
		return nil, 0, fmt.Errorf("error al buscar en ZincSearch: %v", err)
	}

	return emails, total, nil
}
