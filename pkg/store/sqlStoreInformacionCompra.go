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
func (s *sqlStoreInformacionCompra) ObtenerInformacionCompletaCompra(idOrden int) (domain.Orden, domain.InformacionCompra, domain.DatosEnvio, domain.Tarjetas, []domain.OrdenProducto, error) {
    var orden domain.Orden
    var ic domain.InformacionCompra
    var de domain.DatosEnvio
    var tarjeta domain.Tarjetas
    var ordenesProductos []domain.OrdenProducto

    query := `
SELECT
    o.id AS orden_id,
    o.id_usuario,
    o.id_estado,
    o.fechaOrden,
    o.total AS total_orden,
    de.id AS datosenvio_id,
    de.nombre AS nombre_envio,
    de.apellido AS apellido_envio,
    de.direccion AS direccion_envio,
    de.apartamento AS apartamento_envio,
    de.ciudad AS ciudad_envio,
    de.codigo_postal AS codigo_postal_envio,
    de.estado AS estado_envio,
    t.id AS tarjeta_id,
    t.nombre AS nombre_tarjeta,
    t.apellido AS apellido_tarjeta,
    t.numero_tarjeta,
    t.clave_seguridad,
    t.vencimiento,
    t.ultimos_cuatro_digitos,
    op.id AS ordenproducto_id,
    op.id_producto AS producto_id,
    op.cantidad AS cantidad_producto,
    op.total AS total_producto,
    p.id AS producto_id,
    p.nombre AS producto_nombre,
    p.descripcion AS producto_descripcion,
    p.precio AS producto_precio,
    p.stock AS producto_stock,
    p.ranking AS producto_ranking,
    p.id_categoria AS producto_id_categoria,
    c.nombre AS producto_categoria
FROM
    InformacionCompra ic
    INNER JOIN ordenes o ON ic.id_orden = o.id
    INNER JOIN DatosEnvio de ON ic.id_datosenvio = de.id
    INNER JOIN tarjetas t ON ic.id_tarjetas = t.id
    INNER JOIN ordenproducto op ON ic.id_orden = op.id_orden
    INNER JOIN productos p ON op.id_producto = p.id
    INNER JOIN categorias c ON p.id_categoria = c.id
WHERE
    ic.id_orden = ?
`

    // Utilizamos Query para ejecutar la consulta y obtener los resultados
    rows, err := s.db.Query(query, idOrden)
    if err != nil {
        return domain.Orden{}, domain.InformacionCompra{}, domain.DatosEnvio{}, domain.Tarjetas{}, nil, fmt.Errorf("error executing query: %w", err)
    }
    defer rows.Close()

    // Iteramos sobre las filas para escanear los resultados
    for rows.Next() {
        var op domain.OrdenProducto
        var producto domain.Producto
        err := rows.Scan(
            &orden.ID,
            &orden.ID_Usuario,
            &orden.ID_Estado,
            &orden.FechaOrden,
            &orden.Total,
            &de.ID,
            &de.Nombre,
            &de.Apellido,
            &de.Direccion,
            &de.Apartamento,
            &de.Ciudad,
            &de.CodigoPostal,
            &de.Estado,
            &tarjeta.ID,
            &tarjeta.Nombre,
            &tarjeta.Apellido,
            &tarjeta.Numero_Tarjeta,
            &tarjeta.Clave_Seguridad,
            &tarjeta.Vencimiento,
            &tarjeta.Ultimos_Cuatro_Digitos,
            &op.ID,
            &op.ID_Producto,
            &op.Cantidad,
            &op.Total,
            &producto.ID,
            &producto.Nombre,
            &producto.Descripcion,
            &producto.Precio,
            &producto.Stock,
            &producto.Ranking,
            &producto.Id_categoria,
            &producto.Categoria,
        )
        if err != nil {
            return domain.Orden{}, domain.InformacionCompra{}, domain.DatosEnvio{}, domain.Tarjetas{}, nil, fmt.Errorf("error scanning row: %w", err)
        }
        op.Producto = producto
        ordenesProductos = append(ordenesProductos, op)
    }

    // Verificamos si hubo algún error al iterar sobre las filas
    if err := rows.Err(); err != nil {
        return domain.Orden{}, domain.InformacionCompra{}, domain.DatosEnvio{}, domain.Tarjetas{}, nil, fmt.Errorf("error iterating over rows: %w", err)
    }

    // Si no se encontraron resultados
    if len(ordenesProductos) == 0 {
        return domain.Orden{}, domain.InformacionCompra{}, domain.DatosEnvio{}, domain.Tarjetas{}, nil, errors.New("InformacionCompra not found")
    }

    // Asignamos el id de la orden a la estructura de InformacionCompra
    ic.IDOrden = idOrden

    return orden, ic, de, tarjeta, ordenesProductos, nil
}

