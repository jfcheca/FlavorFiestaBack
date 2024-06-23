package mezclas

import (
	"errors"
	"fmt"

	"github.com/jfcheca/FlavorFiesta/internal/domain"
	"github.com/jfcheca/FlavorFiesta/pkg/store"
)

type Repository interface {

    CrearMezcla(p domain.Mezclas) (domain.Mezclas, error)
    BuscarMezcla(id int) (domain.Mezclas, error)
	
}

type repository struct {
	storage store.StoreInterfaceMezclas
}

// NewRepository crea un nuevo repositorio
func NewRepository(storage store.StoreInterfaceMezclas) Repository {
    return &repository{storage: storage}
}

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> CREAR PRODUCTO >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
func (r *repository) CrearMezcla(mezcla domain.Mezclas) (domain.Mezclas, error) {
    err := r.storage.CrearMezcla(mezcla)
    if err != nil {
        return domain.Mezclas{}, fmt.Errorf("error creando mezcla: %w", err)
    }
    return mezcla, nil
}

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> BUSCAR CATEGORIA >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
func (r *repository) BuscarMezcla(id int) (domain.Mezclas, error) {
	producto, err := r.storage.BuscarMezcla(id)
	if err != nil {
		return domain.Mezclas{}, errors.New("categoria not found")
	}
	return producto, nil

}