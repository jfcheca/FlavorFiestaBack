package ingredientes

import (
	"errors"
	"fmt"

	"github.com/jfcheca/FlavorFiesta/internal/domain"
	"github.com/jfcheca/FlavorFiesta/pkg/store"
)

type Repository interface {
    CrearIngredientes(ingredientes []domain.Ingredientes) error
    BuscarIngredientes(id int) (domain.Ingredientes, error)
    DeleteIngredientes(id int) error
}

type repository struct {
    storage store.StoreInterfaceIngredientes
}

func NewRepository(storage store.StoreInterfaceIngredientes) Repository {
    return &repository{storage}
}

func (r *repository) CrearIngredientes(ingredientes []domain.Ingredientes) error {
    fmt.Println("Repositorio: CrearIngredientes llamado con:", ingredientes)
    err := r.storage.CrearIngredientes(ingredientes)
    if err != nil {
        fmt.Println("Error en el storage:", err)
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

func (r *repository) DeleteIngredientes(id int) error {
    err := r.storage.DeleteIngredientes(id)
    if err != nil {
        return err
    }
    return nil
}