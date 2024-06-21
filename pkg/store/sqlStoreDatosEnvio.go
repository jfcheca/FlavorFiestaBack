package store

import (
	"database/sql"
	"fmt"
	"errors"
	"github.com/jfcheca/FlavorFiesta/internal/domain"
)

type sqlStoreDatosEnvios struct {
	db *sql.DB
}

func NewSqlStoreDatosEnvios(db *sql.DB) StoreInterfaceDatosEnvios {
	return &sqlStoreDatosEnvios{
		db: db,
	}
}

// CrearDatosEnvio - Crea un nuevo registro en la tabla DatosEnvio
func (s *sqlStoreDatosEnvios) CrearDatosEnvio(datosEnvio domain.DatosEnvio) error {
    query := `INSERT INTO DatosEnvio (id_usuario, nombre, apellido, direccion, apartamento, ciudad, codigo_postal, estado)
              VALUES (?, ?, ?, ?, ?, ?, ?, ?);`
    stmt, err := s.db.Prepare(query)
    if err != nil {
        return fmt.Errorf("error preparing query: %w", err)
    }
    defer stmt.Close()

    _, err = stmt.Exec(datosEnvio.IDUsuario, datosEnvio.Nombre, datosEnvio.Apellido, datosEnvio.Direccion, datosEnvio.Apartamento, datosEnvio.Ciudad, datosEnvio.CodigoPostal, datosEnvio.Estado)
    if err != nil {
        return fmt.Errorf("error executing query: %w", err)
    }

    return nil
}

// BuscarTodosLosDatosEnvio - Retorna todos los registros de la tabla DatosEnvio
func (s *sqlStoreDatosEnvios) BuscarTodosLosDatosEnvio() ([]domain.DatosEnvio, error) {
    var datosEnvios []domain.DatosEnvio
    query := "SELECT id, id_usuario, nombre, apellido, direccion, apartamento, ciudad, codigo_postal, estado FROM DatosEnvio"
    rows, err := s.db.Query(query)
    if err != nil {
        return nil, fmt.Errorf("error querying datos_envio: %w", err)
    }
    defer rows.Close()

    for rows.Next() {
        var datosEnvio domain.DatosEnvio
        if err := rows.Scan(&datosEnvio.ID, &datosEnvio.IDUsuario, &datosEnvio.Nombre, &datosEnvio.Apellido, &datosEnvio.Direccion, &datosEnvio.Apartamento, &datosEnvio.Ciudad, &datosEnvio.CodigoPostal, &datosEnvio.Estado); err != nil {
            return nil, fmt.Errorf("error scanning datos_envio: %w", err)
        }
        datosEnvios = append(datosEnvios, datosEnvio)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating rows: %w", err)
    }

    return datosEnvios, nil
}

// BuscarDatosEnvio - Retorna un registro de DatosEnvio por ID
func (s *sqlStoreDatosEnvios) BuscarDatosEnvio(id int) (domain.DatosEnvio, error) {
    var datosEnvio domain.DatosEnvio
    query := "SELECT id, id_usuario, nombre, apellido, direccion, apartamento, ciudad, codigo_postal, estado FROM DatosEnvio WHERE id = ?"

    err := s.db.QueryRow(query, id).Scan(&datosEnvio.ID, &datosEnvio.IDUsuario, &datosEnvio.Nombre, &datosEnvio.Apellido, &datosEnvio.Direccion, &datosEnvio.Apartamento, &datosEnvio.Ciudad, &datosEnvio.CodigoPostal, &datosEnvio.Estado)
    if err != nil {
        if err == sql.ErrNoRows {
            return domain.DatosEnvio{}, errors.New("datos_envio not found")
        }
        return domain.DatosEnvio{}, err
    }

    return datosEnvio, nil
}

// EditarDatosEnvio - Edita un registro existente en la tabla DatosEnvio
func (s *sqlStoreDatosEnvios) EditarDatosEnvio(datosEnvio domain.DatosEnvio) error {
    query := `UPDATE DatosEnvio SET id_usuario = ?, nombre = ?, apellido = ?, direccion = ?, apartamento = ?, ciudad = ?, codigo_postal = ?, estado = ? WHERE id = ?`
    stmt, err := s.db.Prepare(query)
    if err != nil {
        return fmt.Errorf("error preparing query: %w", err)
    }
    defer stmt.Close()

    _, err = stmt.Exec(datosEnvio.IDUsuario, datosEnvio.Nombre, datosEnvio.Apellido, datosEnvio.Direccion, datosEnvio.Apartamento, datosEnvio.Ciudad, datosEnvio.CodigoPostal, datosEnvio.Estado, datosEnvio.ID)
    if err != nil {
        return fmt.Errorf("error executing query: %w", err)
    }

    return nil
}

// EliminarDatosEnvio - Elimina un registro de la tabla DatosEnvio por ID
func (s *sqlStoreDatosEnvios) EliminarDatosEnvio(id int) error {
    query := "DELETE FROM DatosEnvio WHERE id = ?"
    stmt, err := s.db.Prepare(query)
    if err != nil {
        return fmt.Errorf("error preparing query: %w", err)
    }
    defer stmt.Close()

    _, err = stmt.Exec(id)
    if err != nil {
        return fmt.Errorf("error executing query: %w", err)
    }

    return nil
}