package handler

import (
	"errors"
	"fmt"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/jfcheca/FlavorFiesta/internal/domain"
	"github.com/jfcheca/FlavorFiesta/internal/ingredientes"
	"github.com/jfcheca/FlavorFiesta/pkg/web"
)

type ingredientesHandler struct {
	s ingredientes.Service
}

// NewImagenHandler crea un nuevo controller de imagenes
func NewIngredientesHandler(s ingredientes.Service) *ingredientesHandler {
	return &ingredientesHandler{
		s: s,
	}
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> CREA NUEVOS INGREDIENTES <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (h *ingredientesHandler) Post() gin.HandlerFunc {
    return func(c *gin.Context) {
        var ingredientes []domain.Ingredientes
        err := c.ShouldBindJSON(&ingredientes)
        if err != nil {
            web.Failure(c, 400, errors.New("invalid json: "+err.Error()))
            fmt.Println("Error al hacer bind del JSON:", err)
            return
        }

        err = h.s.CrearIngredientes(ingredientes)
        if err != nil {
            web.Failure(c, 500, fmt.Errorf("failed to create ingredients: %w", err))
            return
        }

        web.Success(c, 200, "Ingredientes creados correctamente")
    }
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> OBTIENE IMAGEN POR ID <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (h *ingredientesHandler) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.Failure(c, 400, errors.New("Invalid id"))
			return
		}
		ingredientes, err := h.s.BuscarIngredientes(id)
		if err != nil {
			web.Failure(c, 404, errors.New("Imagen not found"))
			return
		}
		web.Success(c, 200, ingredientes)
	}
}
func (h *ingredientesHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.Failure(c, 400, errors.New("invalid id"))
			return
		}
		err = h.s.DeleteIngredientes(id)
		if err != nil {
			web.Failure(c, 404, err)
			return
		}
		// Se elimina el producto correctamente, enviar mensaje de éxito
		c.JSON(200, gin.H{"message": "El ingrediente se elimino correctamente"})
	}
}