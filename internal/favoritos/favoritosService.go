package favoritos

import (
    "github.com/jfcheca/FlavorFiesta/internal/domain"
)

type Service interface {
    AgregarFavorito(favorito domain.Favoritos) (domain.Favoritos, error)
    DeleteFavorito(idUsuario, idProducto int) error
    BuscarFavorito(id int) (domain.Favoritos, error)
    BuscarFavoritosPorUsuario(idUsuario int) ([]domain.Favoritos, error)
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

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> OBTIENE FAVORITO POR ID <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (s *service) BuscarFavorito(id int) (domain.Favoritos, error) {
    favorito, err := s.r.BuscarFavorito(id)
    if err != nil {
        return domain.Favoritos{}, err
    }
    return favorito, nil
}

func (s *service) BuscarFavoritosPorUsuario(idUsuario int) ([]domain.Favoritos, error) {
    favoritos, err := s.r.BuscarFavoritosPorUsuario(idUsuario)
    if err != nil {
        return nil, err
    }
    return favoritos, nil
}

func (s *service) DeleteFavorito(idUsuario, idProducto int) error {
	err := s.r.DeleteFavorito(idUsuario, idProducto)
	if err != nil {
		return err
	}
	return nil
}

