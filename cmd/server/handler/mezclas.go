package handler

import (
	"errors"
	"fmt"
	"strconv"

	//	"os"

	"github.com/gin-gonic/gin"
	"github.com/jfcheca/FlavorFiesta/internal/domain"
	"github.com/jfcheca/FlavorFiesta/internal/mezclas"
	"github.com/jfcheca/FlavorFiesta/pkg/web"
	//	"strings"
)

type mezclasHandler struct {
	s mezclas.Service
}

// NewProductHandler crea un nuevo controller de productos
func NewMezclasHandler(s mezclas.Service) *mezclasHandler {
	return &mezclasHandler{
		s: s,
	}
}

var listaMezclas []domain.Mezclas
var ultimaMezclaID int = 1

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> CREA UNA NUEVA BEBIDA <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

func (h *mezclasHandler) Post() gin.HandlerFunc {
    return func(c *gin.Context) {
        var mezcla domain.Mezclas
        err := c.ShouldBindJSON(&mezcla)
        if err != nil {
            web.Failure(c, 400, errors.New("invalid json: " + err.Error()))
            fmt.Println("Error al hacer bind del JSON:", err)
            return
        }

        // Crear la mezcla utilizando el servicio
        mezclaCreada, err := h.s.CrearMezcla(mezcla)
        if err != nil {
            web.Failure(c, 500, errors.New("failed to create mezcla: " + err.Error()))
            return
        }

        // Devolver la mezcla creada con su ID asignado
        c.JSON(200, mezclaCreada)
    }
}

func (h *mezclasHandler) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.Failure(c, 400, errors.New("Invalid id"))
			return
		}
		mezcla, err := h.s.BuscarMezcla(id)
		if err != nil {
			web.Failure(c, 404, errors.New("No se encuentra"))
			return
		}
		web.Success(c, 200, mezcla)
	}
}