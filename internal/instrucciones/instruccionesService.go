package instrucciones

import (
	"fmt"

	"github.com/jfcheca/FlavorFiesta/internal/domain"
)

type Service interface {
    CrearInstrucciones(instrucciones []domain.Instrucciones) error
    BuscarInstrucciones(id int) (domain.Instrucciones, error)
    DeleteInstrucciones(id int) error
}

type service struct {
    r Repository
}

func NewService(r Repository) Service {
    return &service{r}
}

func (s *service) CrearInstrucciones(instrucciones []domain.Instrucciones) error {
    err := s.r.CrearInstrucciones(instrucciones)
    if err != nil {
        return fmt.Errorf("error creando instrucciones: %w", err)
    }
    return nil
}

func (s *service) BuscarInstrucciones(id int) (domain.Instrucciones, error) {
    p, err := s.r.BuscarInstrucciones(id)
    if err != nil {	
        return domain.Instrucciones{}, err
    }
    return p, nil
}

func (s *service) DeleteInstrucciones(id int) error {
    err := s.r.DeleteInstrucciones(id)
    if err != nil {
        return err
    }
    return nil
}