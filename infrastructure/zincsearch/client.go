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

type ZincSearchClient struct {
	baseURL  string
	username string
	password string
}

// NewZincSearchClient crea una nueva instancia de ZincSearchClient
func NewZincSearchClient() *ZincSearchClient {
	baseURL := os.Getenv("ZINC_URL")
	username := os.Getenv("ZINC_USERNAME")
	password := os.Getenv("ZINC_PASSWORD")

	if baseURL == "" || username == "" || password == "" {
		fmt.Println("Error: Las variables de entorno ZINC_URL, ZINC_USERNAME y ZINC_PASSWORD deben estar definidas.")
		return nil
	}

	return &ZincSearchClient{
		baseURL:  baseURL,
		username: username,
		password: password,
	}
}

// SearchEmails realiza la búsqueda de emails en ZincSearch con paginación
func (zsc *ZincSearchClient) SearchEmailsWithPagination(searchQuery string) ([]model.Email, int, error) {
	req, err := http.NewRequest("POST", zsc.baseURL+"/emails_prueba/_search", bytes.NewBuffer([]byte(searchQuery)))
	if err != nil {
		return nil, 0, fmt.Errorf("error al crear la solicitud: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(zsc.username, zsc.password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("error al hacer la solicitud: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("error en la respuesta de ZincSearch: %v", resp.Status)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("error al leer el cuerpo de la respuesta: %v", err)
	}

	var results struct {
		Hits struct {
			Total struct {
				Value int `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source model.Email `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.Unmarshal(respBody, &results); err != nil {
		return nil, 0, fmt.Errorf("error al decodificar la respuesta: %v", err)
	}

	emails := make([]model.Email, len(results.Hits.Hits))
	for i, hit := range results.Hits.Hits {
		emails[i] = hit.Source
	}

	total := results.Hits.Total.Value
	return emails, total, nil
}
