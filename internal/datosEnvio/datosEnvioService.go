package datosenvio

import (
	"fmt"
	"log"

	"github.com/jfcheca/FlavorFiesta/internal/domain"
)

type Service interface {
    CrearDatosEnvio(de domain.DatosEnvio) (domain.DatosEnvio, error)
    BuscarDatosEnvio(id int) (domain.DatosEnvio, error)
    BuscarTodosLosDatosEnvio() ([]domain.DatosEnvio, error)
    EditarDatosEnvio(de domain.DatosEnvio) (domain.DatosEnvio, error)
    EliminarDatosEnvio(id int) error
}

type service struct {
    r Repository
}

// NewService crea un nuevo servicio
func NewService(r Repository) Service {
    return &service{r}
}


// CrearDatosEnvio - Crea un nuevo registro en la tabla DatosEnvio
func (s *service) CrearDatosEnvio(de domain.DatosEnvio) (domain.DatosEnvio, error) {
    datosEnvioCreado, err := s.r.CrearDatosEnvio(de)
    if err != nil {
        return domain.DatosEnvio{}, err
    }
    return datosEnvioCreado, nil
}

// BuscarDatosEnvio - Retorna un registro de DatosEnvio por ID
func (s *service) BuscarDatosEnvio(id int) (domain.DatosEnvio, error) {
    log.Printf("Service: Buscando DatosEnvio con ID: %d", id)
    datosEnvio, err := s.r.BuscarDatosEnvio(id)
    if err != nil {
        log.Printf("Service: Error al buscar DatosEnvio: %v", err)
        return domain.DatosEnvio{}, err
    }
    return datosEnvio, nil
}

// BuscarTodosLosDatosEnvio - Retorna todos los registros de la tabla DatosEnvio
func (s *service) BuscarTodosLosDatosEnvio() ([]domain.DatosEnvio, error) {
    datosEnvios, err := s.r.BuscarTodosLosDatosEnvio()
    if err != nil {
        return nil, fmt.Errorf("error buscando todos los datos de envio: %w", err)
    }
    return datosEnvios, nil
}

// EditarDatosEnvio - Edita un registro existente en la tabla DatosEnvio
func (s *service) EditarDatosEnvio(de domain.DatosEnvio) (domain.DatosEnvio, error) {
    datosEnvioEditado, err := s.r.EditarDatosEnvio(de)
    if err != nil {
        return domain.DatosEnvio{}, err
    }
    return datosEnvioEditado, nil
}

// EliminarDatosEnvio - Elimina un registro de la tabla DatosEnvio por ID
func (s *service) EliminarDatosEnvio(id int) error {
    err := s.r.EliminarDatosEnvio(id)
    if err != nil {
        return err
    }
    return nil
}