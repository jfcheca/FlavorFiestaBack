package store

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jfcheca/FlavorFiesta/internal/domain"
)

type sqlStoreInstrucciones struct {
	db *sql.DB
}

func NewSqlStoreInstrucciones(db *sql.DB) StoreInterfaceInstrucciones {
	return &sqlStoreInstrucciones{
		db: db,
	}
}

func (s *sqlStoreInstrucciones) CrearInstrucciones(instrucciones []domain.Instrucciones) error {
	query := "INSERT INTO instrucciones (descripcion) VALUES (?);"
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("error al preparar la consulta: %w", err)
	}
	defer stmt.Close()

	for _, instrucciones := range instrucciones {
		_, err := stmt.Exec(instrucciones.Descripcion)
		if err != nil {
			return fmt.Errorf("error al ejecutar la consulta para el ingrediente %v: %w", instrucciones, err)
		}
	}

	return nil
}

func (s *sqlStoreInstrucciones) BuscarInstrucciones(id int) (domain.Instrucciones, error) {
	var instrucciones domain.Instrucciones
	query := "SELECT id, descripcion FROM instrucciones WHERE id = ?"

	err := s.db.QueryRow(query, id).Scan(&instrucciones.ID, &instrucciones.Descripcion)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Instrucciones{}, errors.New("ingrediente no encontrado")
		}
		return domain.Instrucciones{}, err
	}

	return instrucciones, nil
}