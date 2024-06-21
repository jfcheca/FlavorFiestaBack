package favoritos

import (
    "github.com/jfcheca/FlavorFiesta/internal/domain"
)

type Service interface {
    AgregarFavorito(favorito domain.Favoritos) (domain.Favoritos, error)
    DeleteFavorito(id int) error
}

type service struct {
    r Repository
}

func NewServiceFavoritos(r Repository) Service {
    return &service{r}
}

func (s *service) AgregarFavorito(f domain.Favoritos) (domain.Favoritos, error) {
    favorito, err := s.r.AgregarFavorito(f)
    if err != nil {
        return domain.Favoritos{}, err
    }
    return favorito, nil
}

func (s *service) DeleteFavorito(id int) error {
	err := s.r.DeleteFavorito(id)
	if err != nil {
		return err
	}
	return nil
}