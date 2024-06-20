package store

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jfcheca/FlavorFiesta/internal/domain"
	"log"
)

type sqlStoreOrdenes struct {
	db *sql.DB
}

// StoreInterfaceOrdenes defines the methods for interacting with the `ordenes` table.


// NewSqlStoreOrden creates a new sqlStore for ordenes.
func NewSqlStoreOrden(db *sql.DB) StoreInterfaceOrdenes {
	return &sqlStore{
		db: db,
	}
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> CREAR UNA NUEVA ORDEN <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (s *sqlStore) CrearOrden(orden domain.Orden) error {
    query := "INSERT INTO ordenes (id_usuario, id_estado, fechaOrden, total) VALUES (?, ?, ?, ?);"
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("error preparing query: %w", err)
	}
	defer stmt.Close()

    res, err := stmt.Exec(orden.ID_Usuario, orden.ID_Estado, orden.FechaOrden, orden.Total)
	if err != nil {
		return fmt.Errorf("error executing query: %w", err)
	}

	_, err = res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error fetching rows affected: %w", err)
	}

	return nil
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>  BUSCAR ORDEN POR ID <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (s *sqlStore) BuscarOrden(id int) (domain.Orden, error) {
	var orden domain.Orden
	query := "SELECT id, id_usuario, id_estado, fechaOrden, total FROM ordenes WHERE id = ?"

	err := s.db.QueryRow(query, id).Scan(&orden.ID, &orden.ID_Usuario, &orden.ID_Estado, &orden.FechaOrden, &orden.Total)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Orden{}, errors.New("orden not found")
		}
		return domain.Orden{}, err
	}

	return orden, nil
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>  BUSCAR ORDEN POR USUARIO Y ESTADO <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (s *sqlStore) BuscarOrdenPorUsuarioYEstado(userID, estadoID string) (bool, error) {
    var orden domain.Orden
    query := "SELECT id FROM ordenes WHERE id_usuario = ? AND id_estado = ?"

    err := s.db.QueryRow(query, userID, estadoID).Scan(&orden.ID)
    if err != nil {
        if err == sql.ErrNoRows {
            return false, nil
        }
        return false, err
    }

    return true, nil
}
// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>  BUSCAR ORDEN POR USUARIO Y ESTADO CON TODOS LOS DATOS <<<<<<<<<<<<<
func (s *sqlStore) BuscarOrdenPorUsuarioYEstado2(userID, estadoID string) (bool, error, domain.Orden) {
    var orden domain.Orden
    query := "SELECT * FROM ordenes WHERE id_usuario = ? AND id_estado = ?"

    err := s.db.QueryRow(query, userID, estadoID).Scan(&orden.ID, &orden.ID_Usuario, &orden.ID_Estado, &orden.FechaOrden, &orden.Total)
    if err != nil {
        if err == sql.ErrNoRows {
            return false, err, domain.Orden{}
        }
        return false, err, domain.Orden{}
    }
    return true, nil, orden
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ACTUALIZA UNA ORDEN <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (s *sqlStore) UpdateOrden(id int, orden domain.Orden) (domain.Orden, error) {
    query := "UPDATE ordenes SET id_usuario = ?, id_estado = ?, fechaOrden = ?, total = ? WHERE id = ?"

    // Ejecutar la actualización en la base de datos
    result, err := s.db.Exec(query, orden.ID_Usuario, orden.ID_Estado, orden.FechaOrden, orden.Total, id)
    if err != nil {
        return domain.Orden{}, err
    }

    // Verificar cuántas filas fueron afectadas por la actualización
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return domain.Orden{}, err
    }

    // Si no se actualizó ninguna fila, devolver un error
    if rowsAffected == 0 {
        return domain.Orden{}, fmt.Errorf("Orden con ID %d no encontrada", id)
    }

    // Consultar nuevamente la orden actualizada desde la base de datos
    updatedOrder, err := s.BuscarOrden(id)
    if err != nil {
        return domain.Orden{}, fmt.Errorf("Error al buscar la orden actualizada: %v", err)
    }

    // Devolver la orden actualizada
    return updatedOrder, nil
}

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> PATCH ORDEN >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
func (s *sqlStore) PatchOrden(id int, updatedFields map[string]interface{}) error {
    if len(updatedFields) == 0 {
        return errors.New("no fields provided for patching")
    }

    query := "UPDATE ordenes SET"
    values := make([]interface{}, 0)
    index := 0
    for field, value := range updatedFields {
        if index > 0 {
            query += ","
        }
        query += fmt.Sprintf(" %s = ?", field)
        values = append(values, value)
        index++
    }
    query += " WHERE id = ?"
    values = append(values, id)

    stmt, err := s.db.Prepare(query)
    if err != nil {
        return err
    }
    defer stmt.Close()

    _, err = stmt.Exec(values...)
    if err != nil {
        return err
    }

    return nil
}
// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>ELIMINAR UNA ORDEN <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (s *sqlStore) DeleteOrden(id int) error {
	query := "DELETE FROM ordenes WHERE id = ?;"
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

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> VERIFICA SI EXISTE ORDEN CON ESE ID <<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (s *sqlStore) ExistsByIDOrden(id int) bool {
	query := "SELECT COUNT(*) FROM ordenes WHERE id = ?"
	var count int
	err := s.db.QueryRow(query, id).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}