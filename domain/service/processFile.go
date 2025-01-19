package service

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"project/domain/model"
	"strings"
	"sync"
)

// Función para procesar un archivo y extraer los datos
func ProcessFile(filePath string, results chan<- model.Email, wg *sync.WaitGroup) error {
	// Abrir el archivo
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error abriendo el archivo %s: %v\n", filePath, err)
		return err
	}
	defer file.Close()

	// Crear un scanner para leer línea por línea
	var email model.Email
	scanner := bufio.NewScanner(file)

	// Recorrer las líneas del archivo
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Message-ID:") {
			email.MessageID = strings.TrimSpace(strings.TrimPrefix(line, "Message-ID:"))
		} else if strings.HasPrefix(line, "Date:") {
			email.Date = sql.NullString{
				String: strings.TrimSpace(strings.TrimPrefix(line, "Date:")),
				Valid:  true,
			}
		} else if strings.HasPrefix(line, "From:") {
			email.Sender = strings.TrimSpace(strings.TrimPrefix(line, "From:"))
		} else if strings.HasPrefix(line, "To:") {
			email.Receiver = strings.TrimSpace(strings.TrimPrefix(line, "To:"))
		} else if strings.HasPrefix(line, "Subject:") {
			email.Subject = strings.TrimSpace(strings.TrimPrefix(line, "Subject:"))
		} else if strings.HasPrefix(line, "Mime-Version:") {
			email.MimeVersion = strings.TrimSpace(strings.TrimPrefix(line, "Mime-Version:"))
		} else if strings.HasPrefix(line, "Content-Type:") {
			email.ContentType = strings.TrimSpace(strings.TrimPrefix(line, "Content-Type:"))
		} else if strings.HasPrefix(line, "Content-Transfer-Encoding:") {
			email.Encoding = strings.TrimSpace(strings.TrimPrefix(line, "Content-Transfer-Encoding:"))
		} else if strings.HasPrefix(line, "X-Folder:") {
			email.Folder = strings.TrimSpace(strings.TrimPrefix(line, "X-Folder:"))
		} else if len(line) == 0 {
			// Cuando encontramos una línea vacía, significa que el cuerpo del correo comienza aquí
			email.Body = scanner.Text()
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("error leyendo archivo %s: %v", filePath, err)
		return err
	}

	results <- email
	return nil
}
