package favoritos

import (
	"fmt"
	"log"

	"github.com/jfcheca/FlavorFiesta/internal/domain"
	"github.com/jfcheca/FlavorFiesta/pkg/store"
)

type Repository interface {
    
    AgregarFavorito(favorito domain.Favoritos) (domain.Favoritos, error)
    DeleteFavorito(idUsuario, idProducto int) error
    BuscarFavorito(id int) (domain.Favoritos, error)
    BuscarFavoritosPorUsuario(idUsuario int) ([]domain.Favoritos, error)
}

type repository struct {
    storage store.StoreInterfaceFavoritos
}

func NewRepository(storage store.StoreInterfaceFavoritos) Repository {
    return &repository{storage}
}



func (r *repository) AgregarFavorito(p domain.Favoritos) (domain.Favoritos, error) {
    err := r.storage.AgregarFavorito(p)
    if err != nil {
        log.Printf("Error al agregar el favorito %v: %v\n", p, err)
        return domain.Favoritos{}, fmt.Errorf("error creando favorito: %w", err)
    }
    return p, nil
}

func (r *repository) BuscarFavorito(id int) (domain.Favoritos, error) {
    favorito, err := r.storage.BuscarFavorito(id)
    if err != nil {
        return domain.Favoritos{}, err
    }
    return favorito, nil
}

func (r *repository) BuscarFavoritosPorUsuario(idUsuario int) ([]domain.Favoritos, error) {
    favoritos, err := r.storage.BuscarFavoritosPorUsuario(idUsuario)
    if err != nil {
        return nil, err
    }
    return favoritos, nil
}


func (r *repository) DeleteFavorito(idUsuario, idProducto int) error {
	err := r.storage.DeleteFavorito(idUsuario, idProducto)
	if err != nil {
		return err
	}
	return nil
}

