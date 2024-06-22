package handler

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jfcheca/FlavorFiesta/internal/domain"
	"github.com/jfcheca/FlavorFiesta/internal/favoritos"
	"github.com/jfcheca/FlavorFiesta/pkg/web"
)

type favoritosHandler struct {
    s favoritos.Service
}

// NewFavoritosHandler crea un nuevo controlador de favoritos
func NewFavoritosHandler(s favoritos.Service) *favoritosHandler {
    return &favoritosHandler{
        s: s,
    }
}

var favorito domain.Favoritos
var ultimofavoritoID int = 1

// Post maneja la solicitud para agregar un nuevo favorito
func (h *favoritosHandler) Post() gin.HandlerFunc {
    return func(c *gin.Context) {
        var fav domain.Favoritos
        fav.ID = ultimofavoritoID
        ultimofavoritoID++
        err := c.ShouldBindJSON(&fav)
        if err != nil {
            web.Failure(c, 400, errors.New("invalid json: " + err.Error()))
            fmt.Println("Error al hacer bind del JSON:", err)
            return
        }

        // Crear el favorito utilizando el servicio
        createdFavorito, err := h.s.AgregarFavorito(fav)
        if err != nil {
            web.Failure(c, 500, errors.New("fallo la creacion del favorito, revise que los datos ingresados sean correctos"))
            fmt.Println("Error al agregar el favorito", fav, ":", err)
            return
        }

        // Devolver el favorito creado con su ID asignado
        c.JSON(200, createdFavorito)
    }
}

func (h *favoritosHandler) GetByID() gin.HandlerFunc {
    return func(c *gin.Context) {
        idParam := c.Param("id")
        id, err := strconv.Atoi(idParam)
        if err != nil {
            web.Failure(c, 400, errors.New("Invalid id"))
            return
        }
        fmt.Printf("Fetching favorito with ID: %d\n", id)

        favorito, err := h.s.BuscarFavorito(id)
        if err != nil {
            fmt.Printf("Error fetching favorito: %v\n", err)
            web.Failure(c, 404, errors.New("Favorito no encontrado"))
            return
        }
        web.Success(c, 200, favorito)
    }
}

func (h *favoritosHandler) GetFavoritosPorUsuario() gin.HandlerFunc {
    return func(c *gin.Context) {
        idParam := c.Param("id")  // Cambio a "id"
        if idParam == "" {
            fmt.Println("id parameter is missing")
            web.Failure(c, 400, errors.New("Invalid id"))
            return
        }

        fmt.Printf("Received id param: %s\n", idParam)
        idUsuario, err := strconv.Atoi(idParam)
        if err != nil {
            fmt.Printf("Error converting id: %v\n", err)
            web.Failure(c, 400, errors.New("Invalid id"))
            return
        }
        fmt.Printf("Fetching favoritos for user with ID: %d\n", idUsuario)

        favoritos, err := h.s.BuscarFavoritosPorUsuario(idUsuario)
        if err != nil {
            fmt.Printf("Error fetching favoritos: %v\n", err)
            web.Failure(c, 404, errors.New("No se encontraron favoritos para el usuario"))
            return
        }
        web.Success(c, 200, favoritos)
    }
}

func (h *favoritosHandler) DeleteFavorito() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.Failure(c, 400, errors.New("invalid id"))
			return
		}
		err = h.s.DeleteFavorito(id)
		if err != nil {
			web.Failure(c, 404, err)
			return
		}
		// Se elimina la categoria correctamente, enviar mensaje de éxito
		c.JSON(200, gin.H{"message": "El favorito se elimino correctamente"})
	}
}
