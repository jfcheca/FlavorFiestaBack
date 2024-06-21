package tarjetas

import (
	"errors"
	"fmt"
	"log"

	"github.com/jfcheca/FlavorFiesta/internal/domain"
	"github.com/jfcheca/FlavorFiesta/pkg/store"
)

type Repository interface {
	CargarTarjeta(tarjeta domain.Tarjetas) (domain.Tarjetas, error)
	BuscarTarjeta(id int) (domain.Tarjetas, error)
	DeleteTarjeta(id int) error
	
}

type repository struct {
	storage store.StoreInterfaceDatosTarjetas
}

// NewRepository crea un nuevo repositorio para estados
func NewRepository(storage store.StoreInterfaceDatosTarjetas) Repository {
	return &repository{storage}
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> CREAR ESTADO >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
func (r *repository) CargarTarjeta(tarjeta domain.Tarjetas) (domain.Tarjetas, error) {
	err := r.storage.CargarTarjeta(tarjeta)
	if err != nil {
		log.Printf("Error al crear el estado %v: %v\n", tarjeta, err)
		return domain.Tarjetas{}, fmt.Errorf("error creando estado: %w", err)
	}
	return tarjeta, nil
}

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> BUSCAR ESTADOS >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
func (r *repository) BuscarTarjeta(id int) (domain.Tarjetas, error) {
	tarjeta, err := r.storage.BuscarTarjeta(id)
	if err != nil {
		return domain.Tarjetas{}, errors.New("producto not found")
	}
	return tarjeta, nil
}

func (r *repository) DeleteTarjeta(id int) error {
	err := r.storage.DeleteTarjeta(id)
	if err != nil {
		return err
	}
	return nil
}