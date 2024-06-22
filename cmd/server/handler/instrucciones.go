package handler

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jfcheca/FlavorFiesta/internal/domain"
	"github.com/jfcheca/FlavorFiesta/internal/instrucciones"
	"github.com/jfcheca/FlavorFiesta/pkg/web"
)

type instruccionesHandler struct {
	s instrucciones.Service
}

// NewImagenHandler crea un nuevo controller de imagenes
func NewInstruccionesHandler(s instrucciones.Service) *instruccionesHandler {
	return &instruccionesHandler{
		s: s,
	}
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> CREA NUEVA INSTRUCCION <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (h *instruccionesHandler) Post() gin.HandlerFunc {
	return func(c *gin.Context) {
		var instrucciones []domain.Instrucciones
		err := c.ShouldBindJSON(&instrucciones)
		if err != nil {
			web.Failure(c, 400, errors.New("invalid json: "+err.Error()))
			fmt.Println("Error al hacer bind del JSON:", err)
			return
		}

		err = h.s.CrearInstrucciones(instrucciones)
		if err != nil {
			web.Failure(c, 500, fmt.Errorf("failed to create ingredients: %w", err))
			return
		}

		web.Success(c, 200, "Ingredientes creados correctamente")
	}
}
// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> OBTIENE IMAGEN POR ID <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (h *instruccionesHandler) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.Failure(c, 400, errors.New("Invalid id"))
			return
		}
		instrucciones, err := h.s.BuscarInstrucciones(id)
		if err != nil {
			web.Failure(c, 404, errors.New("Imagen not found"))
			return
		}
		web.Success(c, 200, instrucciones)
	}
}