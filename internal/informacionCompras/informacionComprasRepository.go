package informacioncompras

import (
	"errors"
	"fmt"
	"log"

	"github.com/jfcheca/FlavorFiesta/internal/domain"
	"github.com/jfcheca/FlavorFiesta/pkg/store"
)

type Repository interface {
	CrearInformacionCompra(ic domain.InformacionCompra) (domain.InformacionCompra, error)
	BuscarInformacionCompra(id int) (domain.InformacionCompra, error)
	UpdateInformacionCompra(id int, ic domain.InformacionCompra) (domain.InformacionCompra, error)
	DeleteInformacionCompra(id int) error
	ObtenerInformacionCompletaCompra(idOrden int) (domain.Orden, domain.InformacionCompra, domain.DatosEnvio, domain.Tarjetas, []domain.OrdenProducto, error)
}

type repositoryInformacionCompra struct {
	storage store.StoreInterfaceInformacionCompra
}

// NewRepositoryInformacionCompra crea un nuevo repositorio para InformacionCompra
func NewRepositoryInformacionCompras(storage store.StoreInterfaceInformacionCompra) Repository {
	return &repositoryInformacionCompra{storage: storage}
}

func (r *repositoryInformacionCompra) CrearInformacionCompra(ic domain.InformacionCompra) (domain.InformacionCompra, error) {
	createdIC, err := r.storage.CrearInformacionCompra(ic)
	if err != nil {
		log.Printf("Error al crear la InformacionCompra %v: %v\n", ic, err)
		return domain.InformacionCompra{}, fmt.Errorf("error creando InformacionCompra: %w", err)
	}
	return createdIC, nil
}

func (r *repositoryInformacionCompra) BuscarInformacionCompra(id int) (domain.InformacionCompra, error) {
	ic, err := r.storage.BuscarInformacionCompra(id)
	if err != nil {
		return domain.InformacionCompra{}, errors.New("InformacionCompra not found")
	}
	return ic, nil
}

func (r *repositoryInformacionCompra) UpdateInformacionCompra(id int, ic domain.InformacionCompra) (domain.InformacionCompra, error) {
	updatedIC, err := r.storage.UpdateInformacionCompra(id, ic)
	if err != nil {
		return domain.InformacionCompra{}, errors.New("error updating InformacionCompra")
	}
	return updatedIC, nil
}

func (r *repositoryInformacionCompra) DeleteInformacionCompra(id int) error {
	err := r.storage.DeleteInformacionCompra(id)
	if err != nil {
		return fmt.Errorf("error eliminando InformacionCompra: %w", err)
	}
	return nil
}

func (r *repositoryInformacionCompra) ObtenerInformacionCompletaCompra(idOrden int) (domain.Orden, domain.InformacionCompra, domain.DatosEnvio, domain.Tarjetas, []domain.OrdenProducto, error) {
	orden, ic, de, tarjeta, ordenesProductos, err := r.storage.ObtenerInformacionCompletaCompra(idOrden)
	if err != nil {
		log.Printf("Error al obtener la informacion completa de la compra con id_orden %d: %v\n", idOrden, err)
		return domain.Orden{}, domain.InformacionCompra{}, domain.DatosEnvio{}, domain.Tarjetas{}, nil, fmt.Errorf("error obteniendo la informacion completa de la compra: %w", err)
	}
	return orden, ic, de, tarjeta, ordenesProductos, nil
}
