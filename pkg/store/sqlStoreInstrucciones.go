package store

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

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
    fmt.Println("SQL Store: CrearInstrucciones llamado con:", instrucciones)
    query := "INSERT INTO instrucciones (descripcion, id_mezclas) VALUES (?, ?);"
    stmt, err := s.db.Prepare(query)
    if err != nil {
        return fmt.Errorf("error al preparar la consulta: %w", err)
    }
    defer stmt.Close()

    for _, instruccion := range instrucciones {
        fmt.Println("Ejecutando consulta para la instrucción:", instruccion)
        _, err := stmt.Exec(instruccion.Descripcion, instruccion.Id_mezclas)
        if err != nil {
            return fmt.Errorf("error al ejecutar la consulta para la instrucción %v: %w", instruccion, err)
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

func (s *sqlStoreInstrucciones) DeleteInstrucciones(id int) error {
	query := "DELETE FROM instrucciones WHERE id = ?;"
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