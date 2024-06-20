package auth

import (
	"errors"
	"sync"

	"github.com/jfcheca/FlavorFiesta/internal/domain"
)

// Service define los métodos que debe implementar el servicio de autenticación
type Service interface {
	Login(email, password string) (domain.Usuarios, error)
	Authenticate(credentials Credentials) (string, error)
	ForgotPassword(email string) (string, error)
	ActivarCuenta2(email string) (string, error)
	ValidateToken(token string) (string, error)
	ActivarCuenta(email string) (string, error)
}

type service struct {
	repo   Repository
	tokens map[string]string
	mu     sync.Mutex
}

// NewService crea un nuevo servicio de autenticación
func NewService(repo Repository) Service {
	return &service{
		repo:   repo,
		tokens: make(map[string]string),
	}
}

// El método Login en el servicio (service) delega la lógica de inicio de sesión al repositorio (repo)
func (s *service) Login(email, password string) (domain.Usuarios, error) {
	return s.repo.Login(email, password)
}

//Delega la autenticación de las credenciales del usuario al repositorio (repo)

func (s *service) Authenticate(credentials Credentials) (string, error) {
	return s.repo.Authenticate(credentials)
}

// Maneja la lógica de generación de un token de recuperación de contraseña y la verificación de que el email exista en la base de datos
func (s *service) ForgotPassword(email string) (string, error) {
	token, err := s.repo.ForgotPassword(email)
	if err != nil {
		return "", err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.tokens[token] = email

	return token, nil
}

func (s *service) ActivarCuenta2(email string) (string, error) {
	token, err := s.repo.ActivarCuenta2(email)
	if err != nil {
		return "", err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.tokens[token] = email

	return token, nil
}

//Maneja la validación de un token de recuperación de contraseña
//Verifica si el token proporcionado existe en el mapa s.tokens.
//Si el token no existe, devuelve un error indicando que el token es inválido

func (s *service) ValidateToken(token string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	email, exists := s.tokens[token]
	if !exists {
		return "", errors.New("invalid token")
	}

	return email, nil
}

// Maneja la lógica de generación de un token de recuperación de contraseña y la verificación de que el email exista en la base de datos
func (s *service) ActivarCuenta(email string) (string, error) {
	token, err := s.repo.ActivarCuenta(email)
	if err != nil {
		return "", err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.tokens[token] = email

	return token, nil
}
