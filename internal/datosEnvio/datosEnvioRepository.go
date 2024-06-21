package datosenvio

import (
    "errors"
    "fmt"
	"github.com/jfcheca/FlavorFiesta/internal/domain"
	"github.com/jfcheca/FlavorFiesta/pkg/store"
)

type Repository interface {
    CrearDatosEnvio(de domain.DatosEnvio) (domain.DatosEnvio, error)
    BuscarDatosEnvio(id int) (domain.DatosEnvio, error)
    BuscarTodosLosDatosEnvio() ([]domain.DatosEnvio, error)
    EditarDatosEnvio(de domain.DatosEnvio) (domain.DatosEnvio, error)
    EliminarDatosEnvio(id int) error
}

type repository struct {
    storage store.StoreInterfaceDatosEnvios
}

// NewRepository crea un nuevo repositorio
func NewRepository(storage store.StoreInterfaceDatosEnvios) Repository {
    return &repository{storage}
}


// CrearDatosEnvio - Crea un nuevo registro en la tabla DatosEnvio
func (r *repository) CrearDatosEnvio(de domain.DatosEnvio) (domain.DatosEnvio, error) {
    err := r.storage.CrearDatosEnvio(de)
    if err != nil {
        return domain.DatosEnvio{}, errors.New("error creando datos de envio")
    }
    return de, nil
}

// BuscarDatosEnvio - Retorna un registro de DatosEnvio por ID
func (r *repository) BuscarDatosEnvio(id int) (domain.DatosEnvio, error) {
    datosEnvio, err := r.storage.BuscarDatosEnvio(id)
    if err != nil {
        return domain.DatosEnvio{}, errors.New("datos de envio no encontrados")
    }
    return datosEnvio, nil
}

// BuscarTodosLosDatosEnvio - Retorna todos los registros de la tabla DatosEnvio
func (r *repository) BuscarTodosLosDatosEnvio() ([]domain.DatosEnvio, error) {
    datosEnvios, err := r.storage.BuscarTodosLosDatosEnvio()
    if err != nil {
        return nil, errors.New("error al buscar todos los datos de envio")
    }
    return datosEnvios, nil
}

// EditarDatosEnvio - Edita un registro existente en la tabla DatosEnvio
func (r *repository) EditarDatosEnvio(de domain.DatosEnvio) (domain.DatosEnvio, error) {
    // Verificar si el registro existe por su ID
    _, err := r.storage.BuscarDatosEnvio(de.ID)
    if err != nil {
        return domain.DatosEnvio{}, fmt.Errorf("datos de envio con ID %d no encontrados", de.ID)
    }

    // Actualizar el registro en el almacenamiento
    err = r.storage.EditarDatosEnvio(de)
    if err != nil {
        return domain.DatosEnvio{}, err
    }

    return de, nil
}

// EliminarDatosEnvio - Elimina un registro de la tabla DatosEnvio por ID
func (r *repository) EliminarDatosEnvio(id int) error {
    err := r.storage.EliminarDatosEnvio(id)
    if err != nil {
        return err
    }
    return nil
}
