package store

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jfcheca/FlavorFiesta/internal/domain"
)

type sqlStoreIngredientes struct {
	db *sql.DB
}

func NewSqlStoreIngredientes(db *sql.DB) StoreInterfaceIngredientes {
	return &sqlStoreIngredientes{
		db: db,
	}
}

func (s *sqlStoreIngredientes) CrearIngredientes(ingredientes []domain.Ingredientes) error {
	query := "INSERT INTO ingredientes (descripcion) VALUES (?);"
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("error al preparar la consulta: %w", err)
	}
	defer stmt.Close()

	for _, ingrediente := range ingredientes {
		_, err := stmt.Exec(ingrediente.Descripcion)
		if err != nil {
			return fmt.Errorf("error al ejecutar la consulta para el ingrediente %v: %w", ingrediente, err)
		}
	}

	return nil
}

func (s *sqlStoreIngredientes) BuscarIngredientes(id int) (domain.Ingredientes, error) {
	var ingrediente domain.Ingredientes
	query := "SELECT id, descripcion FROM ingredientes WHERE id = ?"

	err := s.db.QueryRow(query, id).Scan(&ingrediente.ID, &ingrediente.Descripcion)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Ingredientes{}, errors.New("ingrediente no encontrado")
		}
		return domain.Ingredientes{}, err
	}

	return ingrediente, nil
}