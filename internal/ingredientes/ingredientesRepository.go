package ingredientes

import (
    "errors"

    "github.com/jfcheca/FlavorFiesta/internal/domain"
    "github.com/jfcheca/FlavorFiesta/pkg/store"
)

type Repository interface {
    CrearIngredientes(ingredientes []domain.Ingredientes) error
    BuscarIngredientes(id int) (domain.Ingredientes, error)
}

type repository struct {
    storage store.StoreInterfaceIngredientes
}

func NewRepository(storage store.StoreInterfaceIngredientes) Repository {
    return &repository{storage}
}

func (r *repository) CrearIngredientes(ingredientes []domain.Ingredientes) error {
    err := r.storage.CrearIngredientes(ingredientes)
    if err != nil {
        return errors.New("Error creando Ingredientes, producto inexistente")
    }
    return nil
}

func (r *repository) BuscarIngredientes(id int) (domain.Ingredientes, error) {
    ingrediente, err := r.storage.BuscarIngredientes(id)
    if err != nil {
        return domain.Ingredientes{}, errors.New("ingrediente no encontrado")
    }
    return ingrediente, nil
}