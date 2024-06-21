package tarjetas

import (
	"github.com/jfcheca/FlavorFiesta/internal/domain"
)

type Service interface {

	CargarTarjeta(tarjeta domain.Tarjetas) (domain.Tarjetas, error)
	BuscarTarjeta(id int) (domain.Tarjetas, error)
	DeleteTarjeta(id int) error
}

type service struct {
	r Repository
}

// NewService crea un nuevo servicio para estados
func NewService(r Repository) Service {
	return &service{r}
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> CREAR UN NUEVO ESTADO <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (s *service) CargarTarjeta(tarjeta domain.Tarjetas) (domain.Tarjetas, error) {
	tarjetaCreada, err := s.r.CargarTarjeta(tarjeta)
	if err != nil {
		return domain.Tarjetas{}, err
	}
	return tarjetaCreada, nil
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> OBTIENE PRODUCTO POR ID <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (s *service) BuscarTarjeta(id int) (domain.Tarjetas, error) {
	tarjeta, err := s.r.BuscarTarjeta(id)
	if err != nil {
		return domain.Tarjetas{}, err
	}
	return tarjeta, nil
}

func (s *service) DeleteTarjeta(id int) error {
	err := s.r.DeleteTarjeta(id)
	if err != nil {
		return err
	}
	return nil
}