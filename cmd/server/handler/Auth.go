package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jfcheca/FlavorFiesta/internal/auth"
	"github.com/jfcheca/FlavorFiesta/internal/usuarios"
	server "github.com/jfcheca/FlavorFiesta/internal/utils"
)

type AuthHandler struct {
	service auth.Service
	s       usuarios.Service
}

func NewAuthHandler(service auth.Service, s usuarios.Service) *AuthHandler {
	return &AuthHandler{
		service: service,
		s:       s,
	}
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> LOGIN >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
func (h *AuthHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Authenticate the user and get the token
		token, err := h.service.Authenticate(auth.Credentials{
			Email:    req.Email,
			Password: req.Password,
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		// Get user details for the response
		user, err := h.service.Login(req.Email, req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve user details"})
			return
		}

		// Return the user details and token in the response
		c.JSON(http.StatusOK, gin.H{
			"user":  user,
			"token": token,
		})
	}
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> FORGOT PASSWORD >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
func (h *AuthHandler) ForgotPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email string `json:"email"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userExists, err := h.s.ExisteEmail2(req.Email) // Assuming ExisteEmail2 returns a bool

		token, err := h.service.ForgotPassword(req.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not process forgot password request"})
			return
		}

		err = server.EnviarConfirmacionEmail(userExists, token)
		if err != nil {
			log.Printf("Error sending email: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not send recovery email"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token":   token,
			"user_id": userExists.ID,
		})
	}
}

func (h *AuthHandler) ActivarCuenta() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email string `json:"email"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := h.s.ExisteEmail2(req.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		}

		token, err := h.service.ActivarCuenta2(req.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not process account activation request"})
			return
		}

		err = server.EnviarConfirmacionCuentaEmail(user, token)
		if err != nil {
			log.Printf("Error sending email: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not send activation email"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": token,
			"email": req.Email,
		})
	}
}

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> CONFIRMACION DE MAIL >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
/*func enviarConfirmacionEmail(user domain.Usuarios, token string) error {
    smtpEmail := os.Getenv("SMTP_EMAIL")
    smtpPassword := os.Getenv("SMTP_PASSWORD")

    if smtpEmail == "" || smtpPassword == "" {
        log.Printf("SMTP_EMAIL o SMTP_PASSWORD no está definido")
        return fmt.Errorf("SMTP_EMAIL o SMTP_PASSWORD no está definido")
    }

    // Plantilla HTML como una cadena
	htmlTemplate := `
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<title>Recuperación de Contraseña en FlavorFiesta</title>
		<style>
			body {
				font-family: 'Poppins', sans-serif;
			}
		</style>
	</head>
	<body>
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; text-align: center;">
			<div style="background-color: #8FA206; border-radius: 50%; width: 60px; height: 60px; text-align: center; line-height: 60px; font-size: 30px; color: white; margin: 0 auto;">
				<span>&#128274;</span>
			</div>
			<h2 style="color: #CC2D4A; font-weight: bold;">RESTAURAR CONTRASEÑA</h2>
			<h2 style="color: #000000;">Hola, <strong style="color: #CC2D4A;">{{ .Nombre }}</strong></h2>
			<p>Nos has solicitado que se restablezca tu contraseña para tu cuenta. Para continuar, haz clic en el botón de abajo:</p>
			<p style="text-align: center;">
				<a href="{{ .ResetURL }}" style="display: inline-block; background-color: #8FA206; color: white; text-decoration: none; padding: 10px 20px; border-radius: 5px;">Establecer nueva contraseña</a>
			</p>
			<p>Recuerda que este enlace es válido durante 24 horas.</p>
			<p>Gracias,<br>El equipo de FlavorFiesta</p>
		</div>
	</body>
	</html>
`
    // Cargar y parsear la plantilla HTML
    tpl, err := template.New("emailTemplate").Parse(htmlTemplate)
    if err != nil {
        log.Printf("Error parsing template: %v", err)
        return err
    }

    // Construir la URL de restablecimiento
	resetURL := fmt.Sprintf("http://localhost:5173/ResetPassword/%d/%s", user.ID, token)

    // Crear un buffer para renderizar la plantilla
    var bodyContent bytes.Buffer
    data := struct {
        Nombre   string
        ResetURL string
    }{
        Nombre:   user.Nombre, // Ajusta esto según el campo correcto en tu estructura de usuario
        ResetURL: resetURL,
    }
    err = tpl.Execute(&bodyContent, data)
    if err != nil {
        log.Printf("Error executing template: %v", err)
        return err
    }

    // Configurar el mensaje de correo
    m := mail.NewMessage()
    m.SetHeader("From", smtpEmail)
    m.SetHeader("To", user.Email)
    m.SetHeader("Subject", "Recuperación de Contraseña en FlavorFiesta")
    m.SetBody("text/html", bodyContent.String())

    // Configurar el Dialer para enviar el correo
    d := mail.NewDialer("smtp.gmail.com", 587, smtpEmail, smtpPassword)

    // Enviar el correo
    if err := d.DialAndSend(m); err != nil {
        return err
    }
    return nil
}*/

///////////////////////////////////////////////////////////////////////////////////////////

/*func enviarConfirmacionCuentaEmail(user domain.Usuarios, token string) error {
    smtpEmail := os.Getenv("SMTP_EMAIL")
    smtpPassword := os.Getenv("SMTP_PASSWORD")

    if smtpEmail == "" || smtpPassword == "" {
        log.Printf("SMTP_EMAIL o SMTP_PASSWORD no está definido")
        return fmt.Errorf("SMTP_EMAIL o SMTP_PASSWORD no está definido")
    }

    // Plantilla HTML como una cadena
    htmlTemplate := `
    <!DOCTYPE html>
    <html>
    <head>
        <meta charset="UTF-8">
        <title>Activación de Cuenta en FlavorFiesta</title>
        <style>
            body {
                font-family: 'Poppins', sans-serif;
            }
        </style>
    </head>
    <body>
        <div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; text-align: center;">
            <div style="background-color: #8FA206; border-radius: 50%; width: 60px; height: 60px; text-align: center; line-height: 60px; font-size: 30px; color: white; margin: 0 auto;">
                <span>&#128274;</span>
            </div>
            <h2 style="color: #CC2D4A; font-weight: bold;">ACTIVACIÓN DE CUENTA</h2>
            <h2 style="color: #000000;">Hola, <strong style="color: #CC2D4A;">{{ .Nombre }}</strong></h2>
            <p>Gracias por registrarte en FlavorFiesta. Para activar tu cuenta, por favor haz clic en el botón de abajo:</p>
            <p style="text-align: center;">
                <a href="{{ .ActivationURL }}" style="display: inline-block; background-color: #8FA206; color: white; text-decoration: none; padding: 10px 20px; border-radius: 5px;">Activar Cuenta</a>
            </p>
            <p>Recuerda que este enlace es válido durante 24 horas.</p>
            <p>Gracias,<br>El equipo de FlavorFiesta</p>
        </div>
    </body>
    </html>
    `

    // Cargar y parsear la plantilla HTML
    tpl, err := template.New("emailTemplate").Parse(htmlTemplate)
    if err != nil {
        log.Printf("Error parsing template: %v", err)
        return err
    }

    // Construir la URL de activación
    activationURL := fmt.Sprintf("http://localhost:5173/ActivateAccount/%d/%s", user.ID, token)

    // Crear un buffer para renderizar la plantilla
    var bodyContent bytes.Buffer
    data := struct {
        Nombre       string
        ActivationURL string
    }{
        Nombre:       user.Nombre,
        ActivationURL: activationURL,
    }
    err = tpl.Execute(&bodyContent, data)
    if err != nil {
        log.Printf("Error executing template: %v", err)
        return err
    }

    // Configurar el mensaje de correo
    m := mail.NewMessage()
    m.SetHeader("From", smtpEmail)
    m.SetHeader("To", user.Email)
    m.SetHeader("Subject", "Activación de Cuenta en FlavorFiesta")
    m.SetBody("text/html", bodyContent.String())

    // Configurar el Dialer para enviar el correo
    d := mail.NewDialer("smtp.gmail.com", 587, smtpEmail, smtpPassword)

    // Enviar el correo
    if err := d.DialAndSend(m); err != nil {
        return err
    }
    return nil
}*/

//////////////////////
