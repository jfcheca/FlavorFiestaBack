package ingredientes

import (
    "github.com/jfcheca/FlavorFiesta/internal/domain"
)

type Service interface {
    CrearIngredientes(ingredientes []domain.Ingredientes) error
    BuscarIngredientes(id int) (domain.Ingredientes, error)
    DeleteIngredientes(id int) error
}

type service struct {
    r Repository
}

func NewService(r Repository) Service {
    return &service{r}
}

func (s *service) CrearIngredientes(ingredientes []domain.Ingredientes) error {
    return s.r.CrearIngredientes(ingredientes)
}

func (s *service) BuscarIngredientes(id int) (domain.Ingredientes, error) {
    p, err := s.r.BuscarIngredientes(id)
    if err != nil {
        return domain.Ingredientes{}, err
    }
    return p, nil
}

func (s *service) DeleteIngredientes(id int) error {
    err := s.r.DeleteIngredientes(id)
    if err != nil {
        return err
    }
    return nil
}