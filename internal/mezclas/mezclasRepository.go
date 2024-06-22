package mezclas

import (

	"fmt"
	"log"

	"github.com/jfcheca/FlavorFiesta/internal/domain"
	"github.com/jfcheca/FlavorFiesta/pkg/store"
)

type Repository interface {

    CrearMezcla(p domain.Mezclas) (domain.Mezclas, error)
	
}

type repository struct {
	storage store.StoreInterfaceMezclas
}

// NewRepository crea un nuevo repositorio
func NewRepository(storage store.StoreInterfaceMezclas) Repository {
    return &repository{storage: storage}
}

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> CREAR PRODUCTO >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
func (r *repository) CrearMezcla(p domain.Mezclas) (domain.Mezclas, error) {
	err := r.storage.CrearMezcla(p)
	if err != nil {
		log.Printf("Error al crear el producto %v: %v\n", p, err)
		return domain.Mezclas{}, fmt.Errorf("error creando producto: %w", err)
	}
	return p, nil
}