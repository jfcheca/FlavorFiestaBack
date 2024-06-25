package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jfcheca/FlavorFiesta/internal/auth"
	"github.com/jfcheca/FlavorFiesta/internal/domain"
	"github.com/jfcheca/FlavorFiesta/internal/usuarios"
	"github.com/jfcheca/FlavorFiesta/pkg/web"
)

type usuariosHandler struct {
	s       usuarios.Service
	authSvc auth.Service
}

// NewImagenHandler crea un nuevo controller de imagenes
func NewUsuarioHandler(s usuarios.Service, authSvc auth.Service) *usuariosHandler {
	return &usuariosHandler{
		s:       s,
		authSvc: authSvc,
	}
}

var listaUsuarios []domain.Usuarios
var ultimoUsuarioID int = 1

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> CREA NUEVA USUARIO <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (h *usuariosHandler) Post() gin.HandlerFunc {
	return func(c *gin.Context) {
		var usuario domain.Usuarios
		err := c.ShouldBindJSON(&usuario)
		if err != nil {
			web.Failure(c, 400, errors.New("invalid json: "+err.Error()))
			fmt.Println("Error al hacer bind del JSON:", err)
			return
		}

		// Asignar el ID adecuadamente (parece que lo gestionas con una variable global)
		usuario.ID = ultimoUsuarioID
		ultimoUsuarioID++

		// Llamar al servicio para crear el usuario
		createdUsuario, err := h.s.CrearUsuario(usuario)
		if err != nil {
			log.Printf("Error creating user: %v", err)
			web.Failure(c, 500, errors.New("Falló la creación de usuario, revise los datos ingresados"))
			return
		}

		// Devolver el usuario creado con su ID asignado
		c.JSON(http.StatusOK, createdUsuario)
	}
}

func (h *usuariosHandler) ActivarCuentaEstado(c *gin.Context) {
	idStr := c.Param("id")
	_, err := strconv.Atoi(idStr)
	if err != nil {
		web.Failure(c, http.StatusBadRequest, errors.New("invalid user ID"))
		return
	}

	// Obtener el email validado del contexto
	email := c.GetString("email")

	var estadoCuenta struct {
		EstadoCuenta string `json:"estado_cuenta"`
	}
	if err := c.ShouldBindJSON(&estadoCuenta); err != nil {
		web.Failure(c, http.StatusBadRequest, errors.New("invalid JSON"))
		return
	}

	// Activar la cuenta usando el email
	err = h.s.ActivarCuenta(email)
	if err != nil {
		web.Failure(c, http.StatusInternalServerError, errors.New("could not activate account"))
		return
	}

	web.Success(c, http.StatusOK, "account activated")
}

func (h *usuariosHandler) ActivarCuenta() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			web.Failure(c, http.StatusUnauthorized, errors.New("missing authorization header"))
			return
		}

		const prefix = "Bearer "
		if !strings.HasPrefix(token, prefix) {
			web.Failure(c, http.StatusUnauthorized, errors.New("invalid authorization format"))
			return
		}

		authToken := token[len(prefix):]

		email, err := h.authSvc.ValidateToken(authToken)
		if err != nil {
			web.Failure(c, http.StatusUnauthorized, errors.New("invalid token"))
			return
		}

		var req struct {
			EstadoCuenta string `json:"estado_cuenta"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			web.Failure(c, http.StatusBadRequest, errors.New("invalid JSON"))
			return
		}

		err = h.s.ActivarCuenta(email)
		if err != nil {
			web.Failure(c, http.StatusInternalServerError, errors.New("could not activate account"))
			return
		}

		web.Success(c, http.StatusOK, "account activated")
	}
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> OBTIENE USUARIO POR ID <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (h *usuariosHandler) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.Failure(c, 400, errors.New("Invalid id"))
			return
		}
		usuario, err := h.s.BuscarUsuario(id)
		if err != nil {
			web.Failure(c, 404, errors.New("No se encuentra"))
			return
		}
		web.Success(c, 200, usuario)
	}
}

// OBTIENE USUARIO POR ID Y PW Y DEVUELVE UN BOOLEANO Y UN MENSAJE
func (h *usuariosHandler) GetByEmailAndPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.Query("email")
		password := c.Query("password")

		if email == "" || password == "" {
			web.Failure(c, 400, errors.New("Email and password are required"))
			return
		}

		exists, err := h.s.BuscarUsuarioPorEmailYPassword(email, password)
		if err != nil {
			web.Failure(c, 404, errors.New("User not found"))
			return
		}

		if exists {
			c.JSON(200, gin.H{
				"success": true,
				"message": "Usuario encontrado",
			})
		} else {
			c.JSON(200, gin.H{
				"success": false,
				"message": "Usuario no encontrado",
			})
		}
	}
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> OBTIENE USUARIO POR EMAIL Y CLAVE <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (h *usuariosHandler) GetByEmailAndPasswordConDatos() gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.Query("email")
		password := c.Query("password")

		if email == "" || password == "" {
			web.Failure(c, 400, errors.New("Email and password are required"))
			return
		}

		exists, err, usuario := h.s.BuscarUsuarioPorEmailYPassword3(email, password)
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

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> OBTIENE TODOS LOS USUARIOS <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (h *usuariosHandler) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		usuarios, err := h.s.BuscarTodosLosUsuarios()
		if err != nil {
			web.Failure(c, 500, fmt.Errorf("error buscando todos los usuarios: %w", err))
			return
		}
		web.Success(c, 200, usuarios)
	}
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ELIMINAR UN USUARIO >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
func (h *usuariosHandler) DeleteUsuario() gin.HandlerFunc {
	return func(c *gin.Context) {

			// Permitir la eliminación del usuario con el token correcto
			idParam := c.Param("id")
			id, err := strconv.Atoi(idParam)
			if err != nil {
				web.Failure(c, 400, errors.New("invalid id"))
				return
			}

			// Verificar si el usuario existe antes de intentar eliminarlo
			_, err = h.s.BuscarUsuario(id)
			if err != nil {
				web.Failure(c, 404, errors.New("El usuario no existe"))
				return
			}

			// Intentar eliminar el usuario
			err = h.s.DeleteUsuario(id)
			if err != nil {
				web.Failure(c, 500, err)
				return
			}

			// Usuario eliminado correctamente, enviar mensaje de éxito
			c.JSON(200, gin.H{"message": "El usuario se eliminó correctamente"})
		} 
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ACTUALIZA UN USUARIO >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
func (h *usuariosHandler) Put() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			web.Failure(c, 400, errors.New("invalid id"))
			return
		}

		_, err = h.s.BuscarUsuario(id)
		if err != nil {
			web.Failure(c, 404, errors.New("user not found"))
			return
		}

		var usuario domain.Usuarios
		err = c.ShouldBindJSON(&usuario)
		if err != nil {
			web.Failure(c, 400, errors.New("invalid json"))
			return
		}

		updatedUsuario, err := h.s.Update(id, usuario)
		if err != nil {
			web.Failure(c, 500, errors.New("could not update user"))
			return
		}

		web.Success(c, 200, updatedUsuario)
	}
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> PATCH ACTUALIZA UN USUARIO O ALGUNO DE SUS CAMPOS <<<<<<<<<<<<<<<<<<<<<<<<<<
func (h *usuariosHandler) Patch() gin.HandlerFunc {

	type Request struct {
		Nombre   string `json:"nombre"`
		Email    string `json:"email"`
		Telefono string `json:"telefono"`
		Password string `json:"password"`
	}

	return func(c *gin.Context) {

		/*token := c.GetHeader("TOKEN")
		if token == "" || token != os.Getenv("TOKEN") {
			web.Failure(c, http.StatusUnauthorized, errors.New("invalid token"))
			return
		}*/

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

		// Verificar si el odontólogo existe antes de actualizar
		_, err = h.s.BuscarUsuario(id)
		if err != nil {
			web.Failure(c, http.StatusNotFound, errors.New("odontologo not found"))
			return
		}

		// Crear una estructura de actualización solo con los campos proporcionados
		update := domain.Usuarios{}
		if r.Nombre != "" {
			update.Nombre = r.Nombre
		}
		if r.Email != "" {
			update.Email = r.Email
		}
		if r.Telefono != "" {
			update.Telefono = r.Telefono
		}
		if r.Password != "" {
			update.Password = r.Password
		}

		// Actualizar el usuario
		p, err := h.s.Update(id, update)
		if err != nil {
			web.Failure(c, http.StatusConflict, err)
			return
		}

		web.Success(c, http.StatusOK, p)
	}
}

func (h *usuariosHandler) UpdatePassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			web.Failure(c, http.StatusUnauthorized, errors.New("missing authorization header"))
			return
		}

		const prefix = "Bearer "
		if !strings.HasPrefix(token, prefix) {
			web.Failure(c, http.StatusUnauthorized, errors.New("invalid authorization format"))
			return
		}

		authToken := token[len(prefix):]

		email, err := h.authSvc.ValidateToken(authToken)
		if err != nil {
			web.Failure(c, http.StatusUnauthorized, errors.New("invalid token"))
			return
		}

		var newPassword struct {
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&newPassword); err != nil {
			web.Failure(c, http.StatusBadRequest, errors.New("invalid JSON"))
			return
		}

		usuario, err := h.s.ExisteEmail2(email)
		if err != nil {
			web.Failure(c, http.StatusNotFound, errors.New("user not found"))
			return
		}

		updatedUser, err := h.s.UpdatePassword(usuario.ID, newPassword.Password)
		if err != nil {
			web.Failure(c, http.StatusInternalServerError, errors.New("could not update password"))
			return
		}

		web.Success(c, http.StatusOK, updatedUser.Password)
	}
}

func (h *usuariosHandler) ActivarCuentaEstado2() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			web.Failure(c, http.StatusUnauthorized, errors.New("missing authorization header"))
			return
		}

		const prefix = "Bearer "
		if !strings.HasPrefix(token, prefix) {
			web.Failure(c, http.StatusUnauthorized, errors.New("invalid authorization format"))
			return
		}

		authToken := token[len(prefix):]

		email, err := h.authSvc.ValidateToken(authToken)
		if err != nil {
			web.Failure(c, http.StatusUnauthorized, errors.New("invalid token"))
			return
		}

		var estadoCuenta struct {
			Estado_Cuenta string `json:"estado_cuenta"`
		}
		if err := c.ShouldBindJSON(&estadoCuenta); err != nil {
			web.Failure(c, http.StatusBadRequest, errors.New("invalid JSON"))
			return
		}

		usuario, err := h.s.ExisteEmail2(email)
		if err != nil {
			web.Failure(c, http.StatusNotFound, errors.New("user not found"))
			return
		}

		updatedUser, err := h.s.ActivarCuentaEstado2(usuario.ID, estadoCuenta.Estado_Cuenta)
		if err != nil {
			web.Failure(c, http.StatusInternalServerError, errors.New("could not update password"))
			return
		}

		web.Success(c, http.StatusOK, updatedUser.Password)
	}
}
