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


/*func ObtenerIDUsuarioDesdeContexto(c *gin.Context) (int, error) {
	authorizationHeader := c.GetHeader("Authorization")
	if authorizationHeader == "" {
		return 0, errors.New("authorization header not found")
	}

	tokenString := strings.Replace(authorizationHeader, "Bearer ", "", 1)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verifica la firma del token JWT aquí
		return []byte("your_secret_key"), nil
	})
	if err != nil {
		return 0, fmt.Errorf("error parsing JWT token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, errors.New("invalid JWT token")
	}

	idUsuario := int(claims["id"].(float64)) // Asumiendo que el ID de usuario está en los claims del token
	return idUsuario, nil
}*/

func (h *favoritosHandler) DeleteFavorito() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el ID del usuario y del producto de los parámetros de la URL
		idUsuarioParam := c.Param("id_usuario")
		idUsuario, err := strconv.Atoi(idUsuarioParam)
		if err != nil {
			web.Failure(c, 400, errors.New("invalid id_usuario"))
			return
		}

		idProductoParam := c.Param("id_producto")
		idProducto, err := strconv.Atoi(idProductoParam)
		if err != nil {
			web.Failure(c, 400, errors.New("invalid id_producto"))
			return
		}

		// Eliminar el favorito usando el servicio
		err = h.s.DeleteFavorito(idUsuario, idProducto)
		if err != nil {
			web.Failure(c, 404, err) // Not Found
			return
		}

		// Respuesta exitosa
		c.JSON(200, gin.H{"message": "El favorito se eliminó correctamente"})
	}
}