package store

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/jfcheca/FlavorFiesta/internal/domain"
)

type sqlStoreTarjetas struct {
	db *sql.DB
}

func NewSqlStoreTarjetas(db *sql.DB) StoreInterfaceDatosTarjetas {
	return &sqlStoreTarjetas{
		db: db,
	}
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> CREAR UN NUEVO ESTADO <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (s *sqlStoreTarjetas) CargarTarjeta(tarjeta domain.Tarjetas) error {
    query := "INSERT INTO datostarjetas (nombre, apellido, numero_tarjeta, clave_seguridad, vencimiento, ultimos_cuatro_digitos, id_usuario) VALUES (?, ?, ?, ?, ?, ?, ?);"
    stmt, err := s.db.Prepare(query)
    if err != nil {
        return fmt.Errorf("error preparing query: %w", err)
    }
    defer stmt.Close()

    _, err = stmt.Exec(tarjeta.Nombre, tarjeta.Apellido, tarjeta.Numero_Tarjeta, tarjeta.Clave_Seguridad, tarjeta.Vencimiento, tarjeta.Ultimos_Cuatro_Digitos, tarjeta.ID_Usuario)
    if err != nil {
        return fmt.Errorf("error executing query: %w", err)
    }

    return nil
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>  BUSCAR ESTADO POR ID <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (s *sqlStoreTarjetas) BuscarTarjeta(id int) (domain.Tarjetas, error) {
    var tarjeta domain.Tarjetas
    query := "SELECT id, nombre, apellido, numero_tarjeta, clave_seguridad, vencimiento, ultimos_cuatro_digitos, id_usuario FROM datostarjetas WHERE id = ?"

    err := s.db.QueryRow(query, id).Scan(&tarjeta.ID, &tarjeta.Nombre, &tarjeta.Apellido, &tarjeta.Numero_Tarjeta, &tarjeta.Clave_Seguridad, &tarjeta.Vencimiento, &tarjeta.Ultimos_Cuatro_Digitos, &tarjeta.ID_Usuario)
    if err != nil {
        if err == sql.ErrNoRows {
            return domain.Tarjetas{}, errors.New("estado not found")
        }
        return domain.Tarjetas{}, err
    }

    return tarjeta, nil
}

func (s *sqlStoreTarjetas) DeleteTarjeta(id int) error {
	query := "DELETE FROM datostarjetas WHERE id = ?;"
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