package store

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jfcheca/FlavorFiesta/internal/domain"
)

type sqlStoreInformacionCompra struct {
	db *sql.DB
}

func NewSqlStoreInformacionCompra(db *sql.DB) StoreInterfaceInformacionCompra {
	return &sqlStoreInformacionCompra{
		db: db,
	}
}


func (s *sqlStoreInformacionCompra) CrearInformacionCompra(ic domain.InformacionCompra) (domain.InformacionCompra, error) {
	query := "INSERT INTO InformacionCompra (id_orden, id_datosenvio, id_tarjetas) VALUES (?, ?, ?);"
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return domain.InformacionCompra{}, fmt.Errorf("error preparing query: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(ic.IDOrden, ic.IDDatosEnvio, ic.IDTarjeta)
	if err != nil {
		return domain.InformacionCompra{}, fmt.Errorf("error executing query: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return domain.InformacionCompra{}, fmt.Errorf("error getting last insert ID: %w", err)
	}

	ic.ID = int(id)
	return ic, nil
}

func (s *sqlStoreInformacionCompra) BuscarInformacionCompra(id int) (domain.InformacionCompra, error) {
	var ic domain.InformacionCompra
	query := "SELECT id, id_orden, id_datosenvio, id_tarjetas FROM InformacionCompra WHERE id = ?"

	err := s.db.QueryRow(query, id).Scan(&ic.ID, &ic.IDOrden, &ic.IDDatosEnvio, &ic.IDTarjeta)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.InformacionCompra{}, errors.New("InformacionCompra not found")
		}
		return domain.InformacionCompra{}, err
	}

	return ic, nil
}

func (s *sqlStoreInformacionCompra) UpdateInformacionCompra(id int, ic domain.InformacionCompra) (domain.InformacionCompra, error) {
	query := "UPDATE InformacionCompra SET id_orden = ?, id_datosenvio = ?, id_tarjetas = ? WHERE id = ?;"
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return domain.InformacionCompra{}, fmt.Errorf("error preparing query: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(ic.IDOrden, ic.IDDatosEnvio, ic.IDTarjeta, id)
	if err != nil {
		return domain.InformacionCompra{}, fmt.Errorf("error executing query: %w", err)
	}

	ic.ID = id
	return ic, nil
}


func (s *sqlStoreInformacionCompra) DeleteInformacionCompra(id int) error {
	query := "DELETE FROM InformacionCompra WHERE id = ?;"
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("error preparing delete query: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("error executing delete query: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("InformacionCompra with ID %d not found", id)
	}

	return nil
}


