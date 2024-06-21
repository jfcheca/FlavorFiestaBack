package store

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jfcheca/FlavorFiesta/internal/domain"
)

type sqlStoreFavoritos struct {
    db *sql.DB
}

func NewSqlStoreFavoritos(db *sql.DB) StoreInterfaceFavoritos {
    return &sqlStoreFavoritos{
        db: db,
    }
}

func (s *sqlStoreFavoritos) AgregarFavorito(favorito domain.Favoritos) error {
    query := "INSERT INTO flavorfiesta.favoritos (id_usuario, id_producto) VALUES (?, ?);"
    stmt, err := s.db.Prepare(query)
    if err != nil {
        return fmt.Errorf("error preparing query: %w", err)
    }
    defer stmt.Close()

    res, err := stmt.Exec(favorito.Id_usuario, favorito.Id_producto)
    if err != nil {
        return fmt.Errorf("error executing query: %w", err)
    }

    _, err = res.RowsAffected()
    if err != nil {
        return fmt.Errorf("error fetching rows affected: %w", err)
    }

    return nil
}

func (s *sqlStoreFavoritos) DeleteFavorito(id int) error {
	query := "DELETE FROM favoritos WHERE id = ?;"
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