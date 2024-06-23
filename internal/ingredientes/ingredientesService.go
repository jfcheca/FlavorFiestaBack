package ingredientes

import (
    "github.com/jfcheca/FlavorFiesta/internal/domain"
)

type Service interface {
    CrearIngredientes(ingredientes []domain.Ingredientes) error
    BuscarIngredientes(id int) (domain.Ingredientes, error)
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