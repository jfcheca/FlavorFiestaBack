package handler

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/jfcheca/FlavorFiesta/internal/domain"
	"github.com/jfcheca/FlavorFiesta/internal/datosenvio"
	"github.com/jfcheca/FlavorFiesta/pkg/web"
)

type datosEnvioHandler struct {
	s datosenvio.Service
}

// NewDatosEnvioHandler crea un nuevo handler para DatosEnvio
func NewDatosEnvioHandler(s datosenvio.Service) *datosEnvioHandler {
	return &datosEnvioHandler{
		s: s,
	}
}
var datoEnvio domain.DatosEnvio
var ultimoDatoEnvioID int = 1
// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> CREA NUEVO DATOSENVIO <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (h *datosEnvioHandler) Post() gin.HandlerFunc {
	return func(c *gin.Context) {
		var datosEnvio domain.DatosEnvio
		datosEnvio.ID = ultimoDatoEnvioID
		ultimoDatoEnvioID++
		err := c.ShouldBindJSON(&datosEnvio)
		if err != nil {
			web.Failure(c, 400, errors.New("invalid json: "+err.Error()))
			return
		}

		datosEnvio, err = h.s.CrearDatosEnvio(datosEnvio)
		if err != nil {
			web.Failure(c, 500, errors.New("failed to create datos de envio: "+err.Error()))
			return
		}

		web.Success(c, 201, datosEnvio)
	}
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> BUSCAR DATOSENVIO POR ID <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (h *datosEnvioHandler) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.Failure(c, 400, errors.New("invalid id"))
			return
		}

		log.Printf("Buscando DatosEnvio con ID: %d", id)
		datosEnvio, err := h.s.BuscarDatosEnvio(id)
		if err != nil {
			web.Failure(c, 404, errors.New("datos de envio not found"))
			return
		}

		web.Success(c, 200, datosEnvio)
	}
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> BUSCAR TODOS LOS DATOSENVIO <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (h *datosEnvioHandler) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		datosEnvios, err := h.s.BuscarTodosLosDatosEnvio()
		if err != nil {
			web.Failure(c, 500, fmt.Errorf("error buscando todos los datos de envio: %w", err))
			return
		}
		web.Success(c, 200, datosEnvios)
	}
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> EDITAR DATOSENVIO <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (h *datosEnvioHandler) Put() gin.HandlerFunc {
	return func(c *gin.Context) {
		idString := c.Param("id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			web.Failure(c, 400, errors.New("invalid id"))
			return
		}

		var datosEnvio domain.DatosEnvio
		err = c.ShouldBindJSON(&datosEnvio)
		if err != nil {
			web.Failure(c, 400, errors.New("invalid json: "+err.Error()))
			return
		}

		datosEnvio.ID = id
		datosEnvio, err = h.s.EditarDatosEnvio(datosEnvio)
		if err != nil {
			web.Failure(c, 500, fmt.Errorf("error editando datos de envio: %w", err))
			return
		}

		web.Success(c, 200, datosEnvio)
	}
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ELIMINAR DATOSENVIO <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (h *datosEnvioHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("TOKEN")
		if token == "123456" {
			idParam := c.Param("id")
			id, err := strconv.Atoi(idParam)
			if err != nil {
				web.Failure(c, 400, errors.New("invalid id"))
				return
			}

			err = h.s.EliminarDatosEnvio(id)
			if err != nil {
				web.Failure(c, 404, err)
				return
			}

			web.Success(c, 200, gin.H{"message": "datos de envio eliminado correctamente"})
		} else {
			web.Failure(c, 401, errors.New("invalid token"))
		}
	}
}