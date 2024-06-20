package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/jfcheca/FlavorFiesta/internal/domain"
	"github.com/jfcheca/FlavorFiesta/internal/imagenes"
	"github.com/jfcheca/FlavorFiesta/pkg/web"
)

type imagenesHandler struct {
	s imagenes.Service
}

// NewImagenHandler crea un nuevo controller de imagenes
func NewImagenHandler(s imagenes.Service) *imagenesHandler {
	return &imagenesHandler{
		s: s,
	}
}

var listaImagenes []domain.Imagen
var ultimaImagenID int = 1

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> CREA NUEVA IMAGEN <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (h *imagenesHandler) Post() gin.HandlerFunc {
    return func(c *gin.Context) {
        var imagenes []domain.Imagen
        err := c.ShouldBindJSON(&imagenes)
        if err != nil {
            web.Failure(c, 400, fmt.Errorf("invalid json: %w", err))
            return
        }

        // Intenta crear las imágenes utilizando el servicio
        err = h.s.CrearImagenes(imagenes)
        if err != nil {
            web.Failure(c, 500, fmt.Errorf("failed to create images: %w", err))
            return
        }

        // Si no hay error, devolver un mensaje de éxito
        web.Success(c, 200, "Imágenes creadas correctamente")
    }
}
// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> OBTIENE IMAGEN POR ID <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (h *imagenesHandler) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.Failure(c, 400, errors.New("Invalid id"))
			return
		}
		odonto, err := h.s.BuscarImagen(id)
		if err != nil {
			web.Failure(c, 404, errors.New("Imagen not found"))
			return
		}
		web.Success(c, 200, odonto)
	}
}

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ACTUALIZA IMAGEN <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (h *imagenesHandler) Put() gin.HandlerFunc {
    return func(c *gin.Context) {
        idString := c.Param("id")
        id, err := strconv.Atoi(idString)
        if err != nil {
            c.JSON(400, gin.H{"error": "invalid id"})
            return
        }
        var imagen domain.Imagen
        err = c.ShouldBindJSON(&imagen)
        if err != nil {
            c.JSON(400, gin.H{"error": "invalid json"})
            return
        }
        // Asignar el ID de la imagen a actualizar
        imagen.ID = id
        
        updatedImagen, err := h.s.UpdateImagen(id, imagen)
        if err != nil {
            c.JSON(400, gin.H{"error": err.Error()})
            return
        }
        // Devolver la imagen actualizada
        c.JSON(200, updatedImagen)
    }
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> PATCH ACTUALIZA ALGUNO/S DE SUS CAMPOS <<<<<<<<<<<zzzz<<<<<<<<<<<<<<<
func (h *imagenesHandler) Patch() gin.HandlerFunc {

    type Request struct {
        Titulo  string `json:"apellido"`
        Url    string `json:"nombre"`
        Matricula string `json:"matricula"`
    }

    return func(c *gin.Context) {
       
        var r Request
        if err := c.ShouldBindJSON(&r); err != nil {
            web.Failure(c, http.StatusBadRequest, errors.New("invalid JSON"))
            return
        }
        idParam := c.Param("id")
        id, err := strconv.Atoi(idParam)
        if err != nil {
            web.Failure(c, http.StatusBadRequest, errors.New("invalid ID"))
            return
        }

        // Verificar si la imagen existe antes de actualizar
        _, err = h.s.BuscarImagen(id)
        if err != nil {
            web.Failure(c, http.StatusNotFound, errors.New("imagen not found"))
            return
        }

        // Crear una estructura de actualización solo con los campos proporcionados
        update := domain.Imagen{}
        if r.Titulo != "" {
            update.Titulo = r.Titulo
        }
        if r.Url != "" {
            update.Url = r.Url
        }
        // Actualizar la imagen
        p, err := h.s.UpdateImagen(id, update)
        if err != nil {
            web.Failure(c, http.StatusConflict, err)
            return
        }

        web.Success(c, http.StatusOK, p)
    }
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ELIMINAR UNA IMAGEN <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (h *imagenesHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("TOKEN")
		if token == "123456" {
			// Permitir la eliminación de la imagen con el token correcto
			idParam := c.Param("id")
			id, err := strconv.Atoi(idParam)
			if err != nil {
				web.Failure(c, 400, errors.New("invalid id"))
				return
			}
			err = h.s.DeleteImagen(id)
			if err != nil {
				web.Failure(c, 404, err)
				return
			}
			// Se elimina la imagen correctamente, enviar mensaje de éxito
			c.JSON(200, gin.H{"message": "La imagen se elimino correctamente"})
		} else {
			// Token no válido
			web.Failure(c, 401, errors.New("invalid token"))
			return
		}
	}
}
