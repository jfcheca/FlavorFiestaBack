package handler

import (
	"errors"
	"fmt"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/jfcheca/FlavorFiesta/internal/domain"
	"github.com/jfcheca/FlavorFiesta/pkg/web"
	"github.com/jfcheca/FlavorFiesta/internal/tarjetas"

)

type tarjetasHandler struct {
	s tarjetas.Service
}

// NewEstadoHandler crea un nuevo controlador de estados
func NewTarjetaHandler(s tarjetas.Service) *tarjetasHandler {
	return &tarjetasHandler{
		s: s,
	}
}

var tarjeta domain.Tarjetas
var ultimaTarjetaID int = 1
// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> CREAR NUEVO ESTADO <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (h *tarjetasHandler) Post() gin.HandlerFunc {
    return func(c *gin.Context) {
        var tarjeta domain.Tarjetas
        tarjeta.ID = ultimaTarjetaID
        ultimaTarjetaID++
        err := c.ShouldBindJSON(&tarjeta)
        if err != nil {
            web.Failure(c, 400, errors.New("invalid json: " + err.Error()))
            fmt.Println("Error al hacer bind del JSON:", err)
            return
        }

        // Validar que id_usuario esté presente
        if tarjeta.ID_Usuario == 0 {
            web.Failure(c, 400, errors.New("id_usuario es requerido"))
            return
        }

        // Crear la tarjeta utilizando el servicio
        createdTarjeta, err := h.s.CargarTarjeta(tarjeta)
        if err != nil {
            fmt.Println("Error al crear la tarjeta:", err)
            web.Failure(c, 500, errors.New("fallo la creacion del estado, revise que los datos ingresados sean correctos"))
            return
        }
        // Devolver la tarjeta creada con su ID asignado por la base de datos
        c.JSON(200, createdTarjeta)
    }
}
// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> OBTIENE ESTADO POR ID <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

func (h *tarjetasHandler) GetByID() gin.HandlerFunc {
    return func(c *gin.Context) {
        idParam := c.Param("id")
        id, err := strconv.Atoi(idParam)
        if err != nil {
            web.Failure(c, 400, errors.New("Invalid id"))
            return
        }
        fmt.Printf("Fetching tarjeta with ID: %d\n", id) // Agrega este registro

        tarjeta, err := h.s.BuscarTarjeta(id)
        if err != nil {
            fmt.Printf("Error fetching tarjeta: %v\n", err) // Agrega este registro
            web.Failure(c, 404, errors.New("No se encuentra"))
            return
        }
        web.Success(c, 200, tarjeta)
    }
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ELIMINAR UNA TARJETA >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
func (h *tarjetasHandler) DeleteTarjeta() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.Failure(c, 400, errors.New("invalid id"))
			return
		}
		err = h.s.DeleteTarjeta(id)
		if err != nil {
			web.Failure(c, 404, err)
			return
		}
		// Se elimina la categoria correctamente, enviar mensaje de éxito
		c.JSON(200, gin.H{"message": "La categoria se elimino correctamente"})
	}
}
