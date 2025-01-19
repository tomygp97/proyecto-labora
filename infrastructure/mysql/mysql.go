package mysql

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func InitMySQL() (*sql.DB, error) {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		return nil, fmt.Errorf("la variable de entorno MYSQL_DSN no está configurada")
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error al abrir la conexión a MySQL: %w", err)
	}

	// Configuración del pool de conexiones
	db.SetMaxOpenConns(10)                 // Máximo de conexiones abiertas
	db.SetMaxIdleConns(5)                  // Máximo de conexiones inactivas
	db.SetConnMaxLifetime(time.Minute * 5) // Tiempo máximo de vida de una conexión

	// Verificar conexión
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error al verificar la conexión a MySQL: %w", err)
	}

	fmt.Println("Conexión exitosa a MySQL")
	return db, nil
}
