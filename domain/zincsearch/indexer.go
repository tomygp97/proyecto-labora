package zincsearch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"project/domain/model"
)

// QueryEmails realiza una consulta a ZincSearch para obtener los correos electrónicos.
func QueryEmails(page, limit int, filters map[string]interface{}) ([]model.Email, error) {
	zincURL := os.Getenv("ZINC_URL")
	if zincURL == "" {
		return nil, fmt.Errorf("ZINC_URL no está configurado")
	}

	// Configurar índice
	indexName := "emails_prueba"
	url := fmt.Sprintf("%s/%s/_search", zincURL, indexName)

	// Construir la consulta para ZincSearch
	query := map[string]interface{}{
		"from": (page - 1) * limit,
		"size": limit,
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}

	jsonQuery, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("error serializando la consulta: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonQuery))
	if err != nil {
		return nil, fmt.Errorf("error creando solicitud HTTP: %w", err)
	}

	username := os.Getenv("ZINC_USERNAME")
	password := os.Getenv("ZINC_PASSWORD")
	if username == "" || password == "" {
		return nil, fmt.Errorf("ZINC_USERNAME o ZINC_PASSWORD no están configurados")
	}
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error enviando solicitud a ZincSearch: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error consultando ZincSearch, código de estado: %d, respuesta: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Hits struct {
			Hits []struct {
				Source model.Email `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return nil, fmt.Errorf("error deserializando la respuesta de ZincSearch: %w", err)
	}

	// Extraer los correos electrónicos de los resultados
	var emails []model.Email
	for _, hit := range result.Hits.Hits {
		emails = append(emails, hit.Source)
	}

	return emails, nil
}

func IndexToZinc(email model.Email) error {
	zincURL := os.Getenv("ZINC_URL")
	if zincURL == "" {
		return fmt.Errorf("ZINC_URL no está configurado")
	}

	// Configurar índice
	indexName := "emails_prueba"
	url := fmt.Sprintf("%s/%s/_doc", zincURL, indexName)

	emailData := map[string]interface{}{
		"message_id":   email.MessageID,
		"sender":       email.Sender,
		"receiver":     email.Receiver,
		"subject":      email.Subject,
		"mime_version": email.MimeVersion,
		"content_type": email.ContentType,
		"encoding":     email.Encoding,
		"folder":       email.Folder,
		"body":         email.Body,
		"date":         email.Date,
	}

	jsonEmail, err := json.Marshal(emailData)
	if err != nil {
		return fmt.Errorf("error serializando email: %w", err)
	}

	// Crear solicitud HTTP
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonEmail))
	if err != nil {
		return fmt.Errorf("error creando solicitud HTTP: %w", err)
	}

	username := os.Getenv("ZINC_USERNAME")
	password := os.Getenv("ZINC_PASSWORD")
	if username == "" || password == "" {
		return fmt.Errorf("ZINC_USERNAME o ZINC_PASSWORD no están configurados")
	}
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error enviando solicitud a ZincSearch: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error indexando correo en ZincSearch, código de estado: %d", resp.StatusCode)
	}

	// fmt.Println("Correo indexado correctamente en ZincSearch.")
	return nil
}
