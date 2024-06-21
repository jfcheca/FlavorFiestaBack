package informacioncompra

import (

	"github.com/jfcheca/FlavorFiesta/internal/domain"
)

type service struct {
	r Repository
}
type Service interface {
	CrearInformacionCompra(ic domain.InformacionCompra) (domain.InformacionCompra, error)
	BuscarInformacionCompra(id int) (domain.InformacionCompra, error)
	UpdateInformacionCompra(id int, ic domain.InformacionCompra) (domain.InformacionCompra, error)
	DeleteInformacionCompra(id int) error
}

// NewService crea un nuevo servicio para InformacionCompra
func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) CrearInformacionCompra(ic domain.InformacionCompra) (domain.InformacionCompra, error) {
	createdIC, err := s.r.CrearInformacionCompra(ic)
	if err != nil {
		return domain.InformacionCompra{}, err
	}
	return createdIC, nil
}

func (s *service) BuscarInformacionCompra(id int) (domain.InformacionCompra, error) {
	ic, err := s.r.BuscarInformacionCompra(id)
	if err != nil {
		return domain.InformacionCompra{}, err
	}
	return ic, nil
}

func (s *service) UpdateInformacionCompra(id int, ic domain.InformacionCompra) (domain.InformacionCompra, error) {
	updatedIC, err := s.r.UpdateInformacionCompra(id, ic)
	if err != nil {
		return domain.InformacionCompra{}, err
	}
	return updatedIC, nil
}

func (s *service) DeleteInformacionCompra(id int) error {
	err := s.r.DeleteInformacionCompra(id)
	if err != nil {
		return err
	}
	return nil
}
