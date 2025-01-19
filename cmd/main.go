package main

import (
	"project/api/routes"
	"project/domain/model"
	"project/domain/service"
	"project/domain/zincsearch"
	db "project/infrastructure/mysql"
	zincSearchClient "project/infrastructure/zincsearch"

	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func withSemaphore(sem chan struct{}, f func()) {
	sem <- struct{}{}
	defer func() { <-sem }()
	f()
}

func main() {
	// Cargar variables de entorno
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error cargando archivo .env: %v", err)
	}

	// Inicializar conexión a la base de datos
	dbConn, err := db.InitMySQL()
	if err != nil {
		fmt.Printf("Error inicializando base de datos: %v\n", err)
		return
	}
	defer dbConn.Close()

	baseDir := os.Getenv("BASE_DIR")
	if baseDir == "" {
		log.Fatal("BASE_DIR no está definida en las variables de entorno.")
	}

	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		fmt.Printf("El directorio %s no existe\n", baseDir)
		return
	}

	// Llamar al cliente de ZincSearch
	client := zincSearchClient.NewZincSearchClient()
	if client == nil {
		log.Fatal("Error al inicializar el cliente de ZincSearch")
	}

	const maxGoroutines = 100
	sem := make(chan struct{}, maxGoroutines)

	var wg sync.WaitGroup
	results := make(chan model.Email, 100)

	// Captura el tiempo antes de comenzar el procesamiento
	startTime := time.Now()

	// Goroutine para recolectar, guardar en MySQL y ZincSearch
	emailsSaved := 0
	emailsIndexed := 0
	go func() {
		for email := range results {
			// Guardar en MySQL
			if err := saveEmailToDB(dbConn, email); err != nil {
				fmt.Printf("Error guardando email en MySQL: %v\n", err)
			} else {
				emailsSaved++
			}

			// Indexar en ZincSearch
			if err := zincsearch.IndexToZinc(email); err != nil {
				fmt.Printf("Error indexando email en ZincSearch: %v\n", err)
			} else {
				emailsIndexed++
			}
		}
	}()

	// Recorrer los archivos del directorio
	fmt.Println("Procesando los datos...")
	err = filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			wg.Add(1)

			go withSemaphore(sem, func() {
				defer wg.Done()
				err := service.ProcessFile(path, results, &wg)
				if err != nil {
					fmt.Printf("Error procesando archivo %s: %v\n", path, err)
				}
			})
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error recorriendo el directorio: %v\n", err)
		return
	}

	wg.Wait()

	close(results)

	// Calcular tiempo total de ejecución del procesamiento
	elapsedTime := time.Since(startTime)
	fmt.Printf("Todos los correos fueron procesados en: %v\n", elapsedTime)
	fmt.Printf("Se guardaron %d correos en la base de datos correctamente.\n", emailsSaved)
	fmt.Printf("Se indexaron %d correos en ZincSearch correctamente.\n", emailsIndexed)

	// Iniciar el servidor de Gin
	r := gin.Default()
	routes.SetupEmailRoutes(r)

	err = r.Run(":8080")
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}

func saveEmailToDB(dbConn *sql.DB, email model.Email) error {
	query := `
        INSERT INTO emails (message_id, sender, receiver, subject, mime_version, content_type, encoding, folder, body, date)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `

	_, err := dbConn.Exec(query, email.MessageID, email.Sender, email.Receiver, email.Subject, email.MimeVersion,
		email.ContentType, email.Encoding, email.Folder, email.Body, email.Date)

	if err != nil {
		fmt.Printf("Error al guardar email: %v\n", err)
		return err
	}

	return nil
}
