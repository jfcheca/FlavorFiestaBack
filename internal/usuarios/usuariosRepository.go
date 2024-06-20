package usuarios

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jfcheca/FlavorFiesta/internal/auth"
	"github.com/jfcheca/FlavorFiesta/internal/domain"
	"github.com/jfcheca/FlavorFiesta/pkg/store"
)

type Repository interface {
	ExisteEmail(email string) (bool, error)
	ExisteEmail2(email string) (domain.Usuarios, error)
	ExisteCelular(celular string) (bool, error)
	CrearUsuario(p domain.Usuarios) (domain.Usuarios, error)
	BuscarUsuario(id int) (domain.Usuarios, error)
	BuscarUsuarioPorEmailYPassword(email, password string) (bool, error)
	BuscarUsuarioPorEmailYPassword2(email, password string) (domain.Usuarios, error)
	BuscarUsuarioPorEmailYPassword3(email, password string) (bool, error, domain.Usuarios)
	BuscarTodosLosUsuarios() ([]domain.Usuarios, error)
	DeleteUsuario(id int) error
	Update(id int, p domain.Usuarios) (domain.Usuarios, error)
	UpdatePassword(id int, newPassword string) (domain.Usuarios, error)
	ActivarCuentaEstado2(id int, estadoCuenta string) (domain.Usuarios, error)

	ActivarCuenta(email string) error
}

type repository struct {
	storage  store.StoreInterfaceUsuarios
	authRepo auth.Repository
}

// NewRepository crea un nuevo repositorio
func NewRepository(storage store.StoreInterfaceUsuarios, authRepo auth.Repository) Repository {
	return &repository{
		storage:  storage,
		authRepo: authRepo, // Asignar el repositorio de autenticación aquí
	}
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>VALIDACIONES >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
func (r *repository) CrearUsuario(u domain.Usuarios) (domain.Usuarios, error) {
	// Verificar si el email ya existe
	exists, err := r.storage.ExisteEmail(u.Email)
	if err != nil {
		return domain.Usuarios{}, err
	}
	if exists {
		return domain.Usuarios{}, errors.New("email already exists")
	}

	// Verificar si el número de teléfono ya existe
	exists, err = r.storage.ExisteCelular(u.Telefono)
	if err != nil {
		return domain.Usuarios{}, err
	}
	if exists {
		return domain.Usuarios{}, errors.New("phone number already exists")
	}

	// Establecer el estado de cuenta como "pendiente"
	u.Estado_Cuenta = "pendiente"

	// Generar y asignar el token al usuario
	token, err := r.authRepo.GenerateToken(u.Email) // Llamar a la función de generación de token desde el repositorio de auth
	if err != nil {
		return domain.Usuarios{}, err
	}
	u.Token = token

	// Si el email y el número de teléfono no existen, continuar con la creación del usuario
	err = r.storage.CrearUsuario(u)
	if err != nil {
		return domain.Usuarios{}, err
	}

	// Devolver el usuario creado
	return u, nil
}

func (r *repository) ActivarCuenta(email string) error {
	usuario, err := r.storage.ExisteEmail2(email)
	if err != nil {
		return fmt.Errorf("error al verificar el email: %v", err)
	}
	if usuario.ID == 0 {
		return errors.New("usuario no encontrado")
	}

	return r.storage.ActivarCuenta(email)
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> BUSCAR USUARIO >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
func (r *repository) BuscarUsuario(id int) (domain.Usuarios, error) {
	usuario, err := r.storage.BuscarUsuario(id)
	if err != nil {
		return domain.Usuarios{}, errors.New("usuario not found")
	}
	return usuario, nil
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> BUSCAR USUARIO POR MAIL Y CLAVE >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
func (r *repository) BuscarUsuarioPorEmailYPassword(email, password string) (bool, error) {
	exists, err := r.storage.BuscarUsuarioPorEmailYPassword(email, password)
	if err != nil {
		return false, errors.New("user not found")
	}
	return exists, nil
}

// este trae todos los datos completos
func (r *repository) BuscarUsuarioPorEmailYPassword2(email, password string) (domain.Usuarios, error) {
	usuario, err := r.storage.BuscarUsuarioPorEmailYPassword2(email, password)
	if err != nil {
		return domain.Usuarios{}, errors.New("usuario not found")
	}
	return usuario, nil
}

// ////////////////////////////////////////
func (r *repository) BuscarUsuarioPorEmailYPassword3(email, password string) (bool, error, domain.Usuarios) {
	exists, err, usuario := r.storage.BuscarUsuarioPorEmailYPassword3(email, password)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, errors.New("usuario not found"), domain.Usuarios{}
		}
		return false, err, domain.Usuarios{}
	}
	return exists, nil, usuario
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> BUSCAR TODOS LOS USUARIOS >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
func (r *repository) BuscarTodosLosUsuarios() ([]domain.Usuarios, error) {
	usuarios, err := r.storage.BuscarTodosLosUsuarios()
	if err != nil {
		return nil, fmt.Errorf("error buscando todos los usuarios: %w", err)
	}
	return usuarios, nil
}

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> DELETE USUARIO >>>>>>>>>>>>>>>>>>>>>>>>>>>

// DeleteUsuario elimina un usuario del repositorio
func (r *repository) DeleteUsuario(id int) error {
	err := r.storage.DeleteUsuario(id)
	if err != nil {
		return err
	}
	return nil
}

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> VERIFICACIONES A LA DB >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

// Implementación de los métodos ExisteEmail y ExisteCelular
func (r *repository) ExisteEmail(email string) (bool, error) {
	return r.storage.ExisteEmail(email)
}

func (r *repository) ExisteEmail2(email string) (domain.Usuarios, error) {
	usuario, err := r.storage.ExisteEmail2(email)
	if err != nil {
		return domain.Usuarios{}, errors.New("usuario not found")
	}
	return usuario, nil
}

func (r *repository) ExisteCelular(celular string) (bool, error) {
	return r.storage.ExisteCelular(celular)
}

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ACTUALIZAR USUARIO >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

// UpdateUsuario actualiza un usuario en el almacenamiento.
func (r *repository) Update(id int, p domain.Usuarios) (domain.Usuarios, error) {

	err := r.storage.Update(p)
	if err != nil {
		return domain.Usuarios{}, errors.New("error updating product")
	}
	return p, nil
}

func (r *repository) UpdatePassword(id int, newPassword string) (domain.Usuarios, error) {
	// Primero, buscar el usuario por su ID
	p, err := r.BuscarUsuario(id)
	if err != nil {
		return domain.Usuarios{}, err
	}

	// Actualizar la contraseña si se proporcionó una nueva contraseña válida
	if newPassword != "" {
		p.Password = newPassword
	}

	// Llamar al método en storage para actualizar la contraseña
	updatedUsuario, err := r.storage.UpdatePassword(id, p.Password)
	if err != nil {
		return domain.Usuarios{}, err
	}

	return updatedUsuario, nil
}

func (r *repository) ActivarCuentaEstado2(id int, estadoCuenta string) (domain.Usuarios, error) {
	// Primero, buscar el usuario por su ID
	p, err := r.BuscarUsuario(id)
	if err != nil {
		return domain.Usuarios{}, err
	}

	// Actualizar la contraseña si se proporcionó una nueva contraseña válida
	if estadoCuenta != "" {
		p.Estado_Cuenta = estadoCuenta
	}

	// Llamar al método en storage para actualizar la contraseña
	updatedUsuario, err := r.storage.ActivarCuentaEstado2(id, p.Estado_Cuenta)
	if err != nil {
		return domain.Usuarios{}, err
	}

	return updatedUsuario, nil
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> PATCH USUARIO >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
func (r *repository) Patch(id int, updatedFields map[string]interface{}) (domain.Usuarios, error) {
	// Obtener el odontólogo por su ID
	usuario, err := r.BuscarUsuario(id)
	if err != nil {
		return domain.Usuarios{}, err
	}

	// Actualizar los campos proporcionados en updatedFields
	for field, value := range updatedFields {
		switch field {
		case "Nombre":
			if nombre, ok := value.(string); ok {
				usuario.Nombre = nombre
			}
		case "Email":
			if email, ok := value.(string); ok {
				usuario.Email = email
			}
		case "Telefono":
			if telefono, ok := value.(string); ok {
				usuario.Telefono = telefono
			}
		case "Password":
			if password, ok := value.(string); ok {
				usuario.Password = password
			}
		// Puedes añadir más campos aquí según sea necesario
		default:
			return domain.Usuarios{}, errors.New("campo desconocido: " + field)
		}
	}

	// Actualizar el odontólogo en el almacenamiento
	updatedUsuario, err := r.Update(id, usuario)
	if err != nil {
		return domain.Usuarios{}, err
	}

	return updatedUsuario, nil
}
