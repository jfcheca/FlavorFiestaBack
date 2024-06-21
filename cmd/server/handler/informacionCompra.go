package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jfcheca/FlavorFiesta/internal/domain"
	"github.com/jfcheca/FlavorFiesta/internal/informacioncompras"
	"github.com/jfcheca/FlavorFiesta/pkg/web"
)

type informacionCompraHandler struct {
	s informacioncompra.Service
}

// NewInformacionCompraHandler crea un nuevo controlador de InformacionCompra
func NewInformacionCompraHandler(s informacioncompra.Service) *informacionCompraHandler {
	return &informacionCompraHandler{
		s: s,
	}
}

var infoCompra domain.InformacionCompra
var ultimoInfoCompraID int = 1

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> CREA UNA NUEVA INFORMACION COMPRA <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (h *informacionCompraHandler) Post() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ic domain.InformacionCompra
		infoCompra.ID = ultimoInfoCompraID
		ultimoInfoCompraID++
		err := c.ShouldBindJSON(&ic)
		if err != nil {
			web.Failure(c, http.StatusBadRequest, errors.New("invalid json: "+err.Error()))
			fmt.Println("Error al hacer bind del JSON:", err)
			return
		}

		// Crear la información de compra utilizando el servicio
		createdIC, err := h.s.CrearInformacionCompra(ic)
		if err != nil {
			web.Failure(c, http.StatusInternalServerError, errors.New("failed to create informacion compra"))
			return
		}

		// Devolver la información de compra creada con su ID asignado
		c.JSON(http.StatusOK, createdIC)
	}
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> OBTIENE INFORMACION COMPRA POR ID <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (h *informacionCompraHandler) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.Failure(c, http.StatusBadRequest, errors.New("Invalid id"))
			return
		}
		ic, err := h.s.BuscarInformacionCompra(id)
		if err != nil {
			web.Failure(c, http.StatusNotFound, errors.New("InformacionCompra not found"))
			return
		}
		web.Success(c, http.StatusOK, ic)
	}
}


//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ACTUALIZA INFORMACION COMPRA >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
func (h *informacionCompraHandler) Put() gin.HandlerFunc {
	return func(c *gin.Context) {
		idString := c.Param("id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		var ic domain.InformacionCompra
		err = c.ShouldBindJSON(&ic)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}
		updatedIC, err := h.s.UpdateInformacionCompra(id, ic)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Devolver solo la información de compra actualizada
		c.JSON(http.StatusOK, updatedIC)
	}
}



// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ELIMINAR INFORMACION DE COMPRA >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
func (h *informacionCompraHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.Failure(c, http.StatusBadRequest, errors.New("invalid id"))
			return
		}
		err = h.s.DeleteInformacionCompra(id)
		if err != nil {
			web.Failure(c, http.StatusNotFound, err)
			return
		}
		// Informacion de compra eliminada correctamente, enviar mensaje de éxito
		c.JSON(http.StatusOK, gin.H{"message": "La informacion de compra se elimino correctamente"})
	}
}
