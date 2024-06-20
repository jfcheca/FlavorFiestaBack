package usuarios

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"unicode"

	"github.com/jfcheca/FlavorFiesta/internal/domain"
	"gopkg.in/mail.v2"
	//"gopkg.in/mail.v2"
	//"gopkg.in/mail.v2"
)

type Service interface {
	ExisteEmail(email string) (bool, error)
	ExisteEmail2(email string) (domain.Usuarios, error)
	ExisteEmail3(email string) (domain.Usuarios, error)
	ExisteCelular(celular string) (bool, error)
	BuscarUsuario(id int) (domain.Usuarios, error)
	//  BuscarUsuarioPorEmailYPassword(email, password string) (domain.Usuarios, error)
	BuscarUsuarioPorEmailYPassword(email, password string) (bool, error)
	BuscarUsuarioPorEmailYPassword2(email, password string) (domain.Usuarios, error)
	BuscarUsuarioPorEmailYPassword3(email, password string) (bool, error, domain.Usuarios)
	BuscarTodosLosUsuarios() ([]domain.Usuarios, error)
	CrearUsuario(p domain.Usuarios) (domain.Usuarios, error)
	DeleteUsuario(id int) error
	Update(id int, p domain.Usuarios) (domain.Usuarios, error)
	UpdatePassword(id int, newPassword string) (domain.Usuarios, error)
	ActivarCuentaEstado2(id int, estadoCuenta string) (domain.Usuarios, error)

	ActivarCuenta(email string) error
}

type service struct {
	r Repository
}

// NewService crea un nuevo servicio
func NewService(r Repository) Service {
	return &service{r}
}

// ValidarUsuario valida los campos de un usuario
func ValidarUsuario(u domain.Usuarios) error {
	// Validar que el nombre no contenga símbolos ni números
	nameRegex := regexp.MustCompile(`^[a-zA-Z\s]+$`)
	if !nameRegex.MatchString(u.Nombre) {
		return errors.New("nombre no debe contener números o símbolos")
	}

	// Validar el formato del email
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(u.Email) {
		return errors.New("email no tiene un formato válido")
	}

	// Validar que el teléfono contenga solo números y un signo + opcional al inicio
	phoneRegex := regexp.MustCompile(`^\+?[0-9]+$`)
	if !phoneRegex.MatchString(u.Telefono) {
		return errors.New("teléfono debe contener solo números y puede comenzar con un signo +")
	}

	// Validar que la contraseña sea alfanumérica y tenga más de 6 caracteres
	if len(u.Password) < 6 {
		return errors.New("contraseña debe tener más de 6 caracteres")
	}
	if !esAlfanumerico(u.Password) {
		return errors.New("contraseña debe ser alfanumérica")
	}

	return nil
}

// Función auxiliar para verificar si una cadena es alfanumérica
func esAlfanumerico(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			return false
		}
	}
	return true
}

///////////////////>>>>>>>>>>>>>>>> CONFIRMACION DE EMAIL PARA CREAR USUARIO >>>>>>>>>>>>>>>>>>>>>>>>>

func enviarConfirmacionEmail(user domain.Usuarios) error {
	m := mail.NewMessage()
	m.SetHeader("From", os.Getenv("SMTP_EMAIL"))
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "¡Confirma tu registro en FlavorFiesta!")

	body := fmt.Sprintf(`
        <html>
        <head>
            <style>
                .container {
                    font-family: Arial, sans-serif;
                    color: #333;
                    max-width: 600px;
                    margin: 0 auto;
                    padding: 20px;
                    border: 1px solid #ddd;
                    border-radius: 10px;
                }
                .header {
                    background-color:#8FA206;
                    color: white;
                    padding: 10px 0;
                    text-align: center;
                    border-radius: 10px 10px 0 0;
                }
                .content {
                    padding: 20px;
                }
                .button {
                    display: inline-block;
                    padding: 10px 20px;
                    font-size: 16px;
                    background-color:#8FA206;
                    text-decoration: none;
                    border-radius: 5px;
                }
                .footer {
                    text-align: center;
                    margin-top: 20px;
                    color: #777;
                    font-size: 12px;
                }
            </style>
        </head>
        <body>
            <div class="container">
                <div class="header">
                    <h1>Bienvenido a FlavorFiesta, %s!</h1>
                </div>
                <div class="content">
                    <p>Hola %s,</p>
                    <p>Gracias por registrarte en FlavorFiesta. Estamos encantados de tenerte con nosotros. Por favor, confirma tu correo electrónico haciendo clic en el siguiente enlace:</p>
                    <p style="text-align: center;">
                        <a  class="button" href='http://localhost:5173/confirm-email/:token' >Confirmar Email</a>
                    </p>
                    <p>Si no te has registrado en FlavorFiesta, por favor ignora este correo.</p>
                </div>
                <div class="footer">
                    <p>&copy; 2024 FlavorFiesta. Todos los derechos reservados.</p>
                </div>
            </div>
        </body>
        </html>
    `, user.Nombre, user.Nombre)

	m.SetBody("text/html", body)

	d := mail.NewDialer("smtp.gmail.com", 587, os.Getenv("SMTP_EMAIL"), os.Getenv("SMTP_PASSWORD"))

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

// CrearUsuario crea un nuevo usuario utilizando el repositorio y devuelve el usuario creado
func (s *service) CrearUsuario(u domain.Usuarios) (domain.Usuarios, error) {
	// Validar el usuario
	if err := ValidarUsuario(u); err != nil {
		return domain.Usuarios{}, err
	}

	// Verificar si el email ya existe
	exists, err := s.r.ExisteEmail(u.Email)
	if err != nil {
		return domain.Usuarios{}, err
	}
	if exists {
		return domain.Usuarios{}, errors.New("email already exists")
	}

	// Verificar si el número de teléfono ya existe
	exists, err = s.r.ExisteCelular(u.Telefono)
	if err != nil {
		return domain.Usuarios{}, err
	}
	if exists {
		return domain.Usuarios{}, errors.New("phone number already exists")
	}

	// Si el email y el número de teléfono no existen, continuar con la creación del usuario
	usuarioCreado, err := s.r.CrearUsuario(u)
	if err != nil {
		return domain.Usuarios{}, err
	}

	// Enviar correo de confirmación
	/*if err := enviarConfirmacionEmail(usuarioCreado); err != nil {
		log.Printf("Error sending confirmation email: %v", err)
		// Puedes devolver un error específico o envolver el error original para mejor trazabilidad
		return domain.Usuarios{}, fmt.Errorf("error sending confirmation email: %w", err)
	} */

	return usuarioCreado, nil
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> BUSCAR USUARIO POR ID >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
func (s *service) BuscarUsuario(id int) (domain.Usuarios, error) {
	p, err := s.r.BuscarUsuario(id)
	if err != nil {
		return domain.Usuarios{}, err
	}
	return p, nil
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> BUSCAR USUARIO POR EMAIL Y CLAVE >>>>>>>>>>>>>>>>>>>>>>>>>>
func (s *service) BuscarUsuarioPorEmailYPassword(email, password string) (bool, error) {
	exists, err := s.r.BuscarUsuarioPorEmailYPassword(email, password)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// ESTE TRAE TODOS LOS DATOS COMPLETOS
func (s *service) BuscarUsuarioPorEmailYPassword2(email, password string) (domain.Usuarios, error) {
	p, err := s.r.BuscarUsuarioPorEmailYPassword2(email, password)
	if err != nil {
		return domain.Usuarios{}, err
	}
	return p, nil
}

// //////////////////////////////////////////////
func (s *service) BuscarUsuarioPorEmailYPassword3(email, password string) (bool, error, domain.Usuarios) {
	exists, err, usuario := s.r.BuscarUsuarioPorEmailYPassword3(email, password)
	if err != nil {
		return false, err, domain.Usuarios{}
	}
	return exists, nil, usuario
}

//////////////////
/*return false, errors.New("usuario not found"), domain.Usuarios{}
}
return false, err, domain.Usuarios{}
}
return exists, nil, usuario

*/

func (s *service) BuscarTodosLosUsuarios() ([]domain.Usuarios, error) {
	usuarios, err := s.r.BuscarTodosLosUsuarios()
	if err != nil {
		return nil, fmt.Errorf("error buscando todos los usuarios: %w", err)
	}
	return usuarios, nil
}

func (s *service) DeleteUsuario(id int) error {
	err := s.r.DeleteUsuario(id)
	if err != nil {
		return err
	}
	return nil
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ACTUALIZA  UN  PRODUCTO <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (s *service) Update(id int, u domain.Usuarios) (domain.Usuarios, error) {
	p, err := s.r.BuscarUsuario(id)
	if err != nil {
		return domain.Usuarios{}, err
	}

	if u.Nombre != "" {
		p.Nombre = u.Nombre
	}
	if u.Email != "" {
		p.Email = u.Email
	}
	if u.Telefono != "" {
		p.Telefono = u.Telefono
	}
	if u.Password != "" {
		p.Password = u.Password
	}

	p, err = s.r.Update(id, p)
	if err != nil {
		return domain.Usuarios{}, err
	}

	return p, nil

}

//OJO ACA Q LO ESTOY LLAMANDO COMO EL CONTROLER Y NO COMO LA INTERFAZ

func (s *service) Patch(id int, updatedFields map[string]interface{}) (domain.Usuarios, error) {
	// Obtener el odontólogo por su ID
	usuario, err := s.r.BuscarUsuario(id)
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
			return domain.Usuarios{}, fmt.Errorf("campo desconocido: %s", field)
		}
	}

	// Actualizar el odontólogo en el repositorio
	updatedUsuario, err := s.r.Update(id, usuario)
	if err != nil {
		return domain.Usuarios{}, err
	}

	return updatedUsuario, nil
}

func (s *service) UpdatePassword(id int, newPassword string) (domain.Usuarios, error) {
	p, err := s.r.BuscarUsuario(id)
	if err != nil {
		return domain.Usuarios{}, err
	}

	if newPassword != "" {
		updatedUser := p
		updatedUser.Password = newPassword

		pUpdated, err := s.r.UpdatePassword(id, updatedUser.Password)
		if err != nil {
			return domain.Usuarios{}, err
		}

		return pUpdated, nil
	}

	return p, nil
}

func (s *service) ActivarCuentaEstado2(id int, estadoCuenta string) (domain.Usuarios, error) {
	p, err := s.r.BuscarUsuario(id)
	if err != nil {
		return domain.Usuarios{}, err
	}

	if estadoCuenta != "" {
		updatedUser := p
		updatedUser.Estado_Cuenta = estadoCuenta

		pUpdated, err := s.r.ActivarCuentaEstado2(id, updatedUser.Estado_Cuenta)
		if err != nil {
			return domain.Usuarios{}, err
		}

		return pUpdated, nil
	}

	return p, nil
}

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>  VALIDACIONES DE MAIL Y CELULAR >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

// Métodos para ExisteEmail y ExisteCelular
func (s *service) ExisteEmail(email string) (bool, error) {
	return s.r.ExisteEmail(email)
}

func (s *service) ExisteEmail2(email string) (domain.Usuarios, error) {
	return s.r.ExisteEmail2(email)
}

func (s *service) ExisteEmail3(email string) (domain.Usuarios, error) {
	return s.r.ExisteEmail2(email)
}

func (s *service) ExisteCelular(celular string) (bool, error) {
	return s.r.ExisteCelular(celular)
}

func (s *service) ActivarCuenta(email string) error {
	// Verificar si el email existe en la base de datos
	usuario, err := s.r.ExisteEmail2(email)
	if err != nil {
		return fmt.Errorf("error al verificar el email: %v", err)
	}
	if usuario.ID == 0 {
		return errors.New("usuario no encontrado")
	}

	return s.r.ActivarCuenta(email)
}
