package handler

import (
	"errors"
	"fmt"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/jfcheca/FlavorFiesta/internal/domain"
	"github.com/jfcheca/FlavorFiesta/internal/ordenes"
	"github.com/jfcheca/FlavorFiesta/pkg/web"
)

type ordenHandler struct {
	s ordenes.Service
}

func NewOrdenHandler(s ordenes.Service) *ordenHandler {
	return &ordenHandler{
		s: s,
	}
}

var orden domain.Orden
var ultimaOrdenID int = 1
//

func (h *ordenHandler) Post() gin.HandlerFunc {
	return func(c *gin.Context) {
        var orden domain.Orden
        orden.ID = ultimaOrdenID
        ultimaOrdenID++
        err := c.ShouldBindJSON(&orden)
        if err != nil {
            web.Failure(c, 400, errors.New("invalid json: " + err.Error()))
            fmt.Println("Error al hacer bind del JSON:", err)
            return
        }

		// Crear la orden utilizando el servicio
		createdOrd, err := h.s.CrearOrden(orden)
		if err != nil {
			web.Failure(c, 500, errors.New("failed to create order"))
			return
		}
		// Devolver la orden creada con su ID asignado a la base de datos
		c.JSON(200, createdOrd)
	}
}

func (h *ordenHandler) GetOrdenByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.Failure(c, 400, errors.New("Invalid id"))
			return
		}
		ord, err := h.s.BuscarOrden(id)
		if err != nil {
			web.Failure(c, 404, errors.New("No se encuentra la orden"))
			return
		}
		web.Success(c, 200, ord)
	}
}

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> OBTIENE ORDEN POR ID Y PW Y DEVUELVE UN BOOLEANO Y UN MENSAJE <<<<<<<<<<<<<<<<<<<<<<
func (h *ordenHandler) GetOrdenByUserIDyOrden() gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.Query("UserID")
        estadoID := c.Query("EstadoID") 

        if userID == "" || estadoID == "" {
            web.Failure(c, 400, errors.New("UserID and EstadoID are required"))
            return
        }

        exists, err := h.s.BuscarOrdenPorUsuarioYEstado(userID, estadoID)
        if err != nil {
            web.Failure(c, 404, errors.New("Order not found"))
            return
        }

        if exists {
            c.JSON(200, gin.H{
                "success": true,
                "message": "Orden encontrada",
            })
        } else {
            c.JSON(200, gin.H{
                "success": false,
                "message": "Orden no encontrada",
            })
        }
    }
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> OBTIENE ORDEN POR USER Y ESTADO CON DATOS <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (h *ordenHandler) GetOrdenByUsuarioYEstadoConDatos() gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.Query("UserID")
        estadoID := c.Query("EstadoID")

        if userID == "" || estadoID == "" {
            web.Failure(c, 400, errors.New("UserID and EstadoID are required"))
            return
        }

        exists, err, usuario := h.s.BuscarOrdenPorUsuarioYEstado2(userID, estadoID)
        if err != nil {
            if err.Error() == "usuario not found" {
                web.Failure(c, 404, errors.New("User not found"))
            } else {
                web.Failure(c, 500, errors.New("Error retrieving user details"))
            }
            return
        }

        if exists {
            c.JSON(200, gin.H{
                "success": true,
                "message": "Usuario encontrado",
                "usuario": usuario,
            })
        } else {
            c.JSON(200, gin.H{
                "success": false,
                "message": "Usuario no encontrado",
            })
        }
    }
}
//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ACTUALIZA LA ORDEN <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (h *ordenHandler) Put() gin.HandlerFunc {
    return func(c *gin.Context) {
        idString := c.Param("id")
        id, err := strconv.Atoi(idString)
        if err != nil {
            c.JSON(400, gin.H{"error": "invalid id"})
            return
        }

        var updatedOrder domain.Orden
        err = c.ShouldBindJSON(&updatedOrder)
        if err != nil {
            c.JSON(400, gin.H{"error": "invalid json"})
            return
        }

        // Llamar al servicio para actualizar la orden
        updatedOrder, err = h.s.UpdateOrden(id, updatedOrder)  // Aquí usamos `=` en lugar de `:=`
        if err != nil {
            c.JSON(400, gin.H{"error": err.Error()})
            return
        }

        // Devolver la orden actualizada en la respuesta
        c.JSON(200, updatedOrder)
    }
}
//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ELIMINA LA ORDEN <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (h *ordenHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {

			idParam := c.Param("id")
			id, err := strconv.Atoi(idParam)
			if err != nil {
				web.Failure(c, 400, errors.New("invalid id"))
				return
			}
			err = h.s.DeleteOrden(id)
			if err != nil {
				web.Failure(c, 404, err)
				return
			}
			c.JSON(200, gin.H{"message": "La orden se eliminó correctamente"})
		} 
}

func (h *ordenHandler) ObtenerOrdenesPorUsuarioYEstadoDiferenteA1() gin.HandlerFunc {
    return func(c *gin.Context) {
        userIDStr := c.Param("userID")
        userID, err := strconv.Atoi(userIDStr)
        if err != nil {
            web.Failure(c, 400, errors.New("invalid userID"))
            return
        }

        ordenes, err := h.s.ObtenerOrdenesPorUsuarioYEstadoDiferenteA1(userID)
        if err != nil {
            web.Failure(c, 500, errors.New("failed to retrieve orders"))
            return
        }

        web.Success(c, 200, ordenes)
    }
}