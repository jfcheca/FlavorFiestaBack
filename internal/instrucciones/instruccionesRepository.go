package instrucciones

import (
	"errors"
	"fmt"

	"github.com/jfcheca/FlavorFiesta/internal/domain"
	"github.com/jfcheca/FlavorFiesta/pkg/store"
)

type Repository interface {
    CrearInstrucciones(instrucciones []domain.Instrucciones) error
    BuscarInstrucciones(id int) (domain.Instrucciones, error)
}

type repository struct {
    storage store.StoreInterfaceInstrucciones
}

func NewRepository(storage store.StoreInterfaceInstrucciones) Repository {
    return &repository{storage}
}

func (r *repository) CrearInstrucciones(instrucciones []domain.Instrucciones) error {
    fmt.Println("Repositorio: CrearInstrucciones llamado con:", instrucciones)
    err := r.storage.CrearInstrucciones(instrucciones)
    if err != nil {
        fmt.Println("Error en el storage:", err)
        return errors.New("Error creando Instrucciones, producto inexistente")
    }
    return nil
}

func (r *repository) BuscarInstrucciones(id int) (domain.Instrucciones, error) {
    instrucciones, err := r.storage.BuscarInstrucciones(id)
    if err != nil {
        return domain.Instrucciones{}, errors.New("instrucciones no encontradas")
    }
    return instrucciones, nil
}