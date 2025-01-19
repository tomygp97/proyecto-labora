package controllers

import (
	"fmt"
	"net/http"
	"project/domain/model"
	"project/domain/service"
	"project/infrastructure/mysql"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EmailResponse struct {
	Emails     []model.Email `json:"emails"`
	HasNext    bool          `json:"has_next"`
	HasPrev    bool          `json:"has_prev"`
	Page       int           `json:"page"`
	Total      int           `json:"total"`
	TotalPages int           `json:"total_pages"`
}

type EmailController struct {
	emailService *service.EmailService
}

func NewEmailController() *EmailController {
	mysqlDB, err := mysql.InitMySQL()
	if err != nil {
		panic("Error al conectar a la base de datos")
	}
	return &EmailController{
		emailService: service.NewEmailService(mysqlDB),
	}
}

// GetEmails maneja la ruta GET /emails y devuelve los correos electrónicos con paginación.
func (ec *EmailController) GetEmails(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Página inválida"})
		return
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Límite inválido"})
		return
	}

	const maxLimit = 100
	if limitInt > maxLimit {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":           "El límite máximo de registros es 100",
			"max_limit":       maxLimit,
			"requested_limit": limitInt,
		})
		return
	}

	offset := (pageInt - 1) * limitInt

	// Obtener los correos electrónicos con paginación
	emails, total, err := ec.emailService.GetEmailsWithPagination(offset, limitInt)
	if err != nil {
		fmt.Printf("Error en GetEmailsHandler: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los correos electrónicos"})
		return
	}

	totalPages := (total + limitInt - 1) / limitInt
	if pageInt > totalPages {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Página fuera de rango", "max_pages": totalPages})
		return
	}

	response := EmailResponse{
		Emails:     emails,
		Total:      total,
		Page:       pageInt,
		TotalPages: totalPages,
		HasPrev:    pageInt > 1,
		HasNext:    pageInt < totalPages,
	}

	c.JSON(http.StatusOK, response)
}

// GetEmailByID maneja la ruta GET /emails/:id y devuelve un correo electrónico específico por ID.
func (ec *EmailController) GetEmailByID(c *gin.Context) {
	id := c.Param("id")
	fmt.Printf("Solicitando correo con ID: %s desde el controlador\n", id)

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	email, err := ec.emailService.GetEmailByID(idInt)
	if err != nil {
		fmt.Printf("Error en GetEmailByIDHandler: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener el correo electrónico"})
		return
	}

	if email == nil {
		fmt.Printf("Correo no encontrado con ID: %d\n", idInt)
		c.JSON(http.StatusNotFound, gin.H{"error": "Correo no encontrado"})
		return
	}

	fmt.Printf("Correo encontrado con ID: %d\n", idInt)
	c.JSON(http.StatusOK, email)
}

// Realiza la búsqueda de correos electrónicos en zincsearch
func (ec *EmailController) SearchEmails(c *gin.Context) {
	query := c.DefaultQuery("query", "")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El término de búsqueda no puede estar vacío"})
		return
	}

	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Página inválida"})
		return
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Límite inválido"})
		return
	}

	const maxLimit = 100
	if limitInt > maxLimit {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":           "El límite máximo de registros es 100",
			"max_limit":       maxLimit,
			"requested_limit": limitInt,
		})
		return
	}

	offset := (pageInt - 1) * limitInt

	// Realizar la búsqueda en el servicio con paginación
	emails, total, err := ec.emailService.SearchEmailsWithPagination(query, offset, limitInt)
	if err != nil {
		fmt.Printf("Error en SearchEmailsHandler: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los correos electrónicos"})
		return
	}

	totalPages := (total + limitInt - 1) / limitInt

	response := EmailResponse{
		Emails:     emails,
		Total:      total,
		Page:       pageInt,
		TotalPages: totalPages,
		HasPrev:    pageInt > 1,
		HasNext:    pageInt < totalPages,
	}

	c.JSON(http.StatusOK, response)
}
