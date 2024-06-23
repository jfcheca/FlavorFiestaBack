package store

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

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
	fmt.Println("SQL Store: CrearIngredientes llamado con:", ingredientes)
	query := "INSERT INTO ingredientes (descripcion, id_mezclas) VALUES (?, ?);"
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("error al preparar la consulta: %w", err)
	}
	defer stmt.Close()

	for _, ingrediente := range ingredientes {
		fmt.Println("Ejecutando consulta para el ingrediente:", ingrediente)
		_, err := stmt.Exec(ingrediente.Descripcion, ingrediente.Id_mezclas)
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

func (s *sqlStoreIngredientes) DeleteIngredientes(id int) error {
	query := "DELETE FROM ingredientes WHERE id = ?;"
	stmt, err := s.db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(id)
	if err != nil {
		log.Fatal(err)
	}
	_, err = res.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}