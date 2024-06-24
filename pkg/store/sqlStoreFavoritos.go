package store

import (
	"database/sql"
	"errors"
	"fmt"

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
    query := "INSERT INTO flavorfiesta.fav (id_usuario, id_producto) VALUES (?, ?);"
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

func (s *sqlStoreFavoritos) BuscarFavorito(id int) (domain.Favoritos, error) {
    var favorito domain.Favoritos
    query := "SELECT id, id_usuario, id_producto FROM flavorfiesta.fav WHERE id = ?"

    err := s.db.QueryRow(query, id).Scan(&favorito.ID, &favorito.Id_usuario, &favorito.Id_producto)
    if err != nil {
        return domain.Favoritos{}, err
    }

    return favorito, nil
}

func (s *sqlStoreFavoritos) BuscarFavoritosPorUsuario(idUsuario int) ([]domain.Favoritos, error) {
    var favoritos []domain.Favoritos
    query := "SELECT id, id_usuario, id_producto FROM flavorfiesta.fav WHERE id_usuario = ?"

    rows, err := s.db.Query(query, idUsuario)
    if err != nil {
        return nil, fmt.Errorf("error executing query: %w", err)
    }
    defer rows.Close()

    for rows.Next() {
        var favorito domain.Favoritos
        err := rows.Scan(&favorito.ID, &favorito.Id_usuario, &favorito.Id_producto)
        if err != nil {
            return nil, fmt.Errorf("error scanning row: %w", err)
        }
        favoritos = append(favoritos, favorito)
    }

    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("error with rows: %w", err)
    }

    return favoritos, nil
}

func (s *sqlStoreFavoritos) DeleteFavorito(idUsuario, idProducto int) error {
	query := "DELETE FROM flavorfiesta.fav WHERE id_usuario = ? AND id_producto = ?;"
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(idUsuario, idProducto)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("favorito no encontrado")
	}

	return nil
}


// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>  BUSCAR FAVORITO POR ID <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
/*func (s *sqlStoreFavoritos) BuscarFavorito(id int) (domain.Favoritos, error) {
    var favorito domain.Favoritos

    // Consulta para obtener el ID del producto asociado al favorito
    query := "SELECT id_producto FROM fav WHERE id = ?"

    // Ejecutar la consulta y escanear el resultado
    err := s.db.QueryRow(query, id).Scan(&favorito.Id_producto)
    if err != nil {
        return domain.Favoritos{}, err
    }

    // Consulta para obtener las imágenes asociadas al producto
    queryImagenes := "SELECT i.id AS imagen_id, i.titulo AS imagen_titulo, i.url AS imagen_url FROM imagenes i WHERE i.id_producto = ?"


    // Ejecutar la consulta de imágenes y procesar los resultados
    rows, err := s.db.Query(queryImagenes, favorito.Id_producto)
    if err != nil {
        return domain.Favoritos{}, err
    }
    defer rows.Close()

    // Iterar sobre las filas de imágenes y agregarlas al slice de imágenes del producto en favorito
    for rows.Next() {
        var imagen domain.Imagen
        err := rows.Scan(
            &imagen.ID,
            &imagen.Titulo,
            &imagen.Url,
        )
        if err != nil {
            return domain.Favoritos{}, err
        }
        favorito.Producto.Imagenes = append(favorito.Producto.Imagenes, imagen)
    }

    // Manejar cualquier error que ocurra durante el procesamiento de filas
    err = rows.Err()
    if err != nil {
        return domain.Favoritos{}, err
    }

    return favorito, nil
}*/



/*func (s *sqlStoreFavoritos) BuscarFavorito(id int) (domain.Favoritos, error) {
    var favorito domain.Favoritos
    query := `
        SELECT 
            f.id, 
            f.id_usuario, 
            f.id_producto,
            p.id AS producto_id,
            p.nombre AS producto_nombre,
            p.descripcion AS producto_descripcion,
            p.precio AS producto_precio,
            p.stock AS producto_stock,
            p.ranking AS producto_ranking,
            p.id_categoria AS producto_id_categoria,
            c.nombre AS categoria
        FROM 
            fav f
            JOIN productos p ON f.id_producto = p.id
        JOIN 
            categorias c ON p.id_categoria = c.id
        WHERE 
            f.id = ?
    `

    // Ejecutar la consulta
    err := s.db.QueryRow(query, id).Scan(
        &favorito.ID,
        &favorito.Id_usuario,
        &favorito.Id_producto,
        &favorito.Producto.ID,
        &favorito.Producto.Nombre,
        &favorito.Producto.Descripcion,
        &favorito.Producto.Precio,
        &favorito.Producto.Stock,
        &favorito.Producto.Ranking,
        &favorito.Producto.Id_categoria,
        &favorito.Producto.Categoria,
    )
    if err != nil {
        return domain.Favoritos{}, err
    }

    // Consulta para obtener las imágenes del producto asociado
    queryImagenes := `
        SELECT 
            i.id AS imagen_id,
            i.titulo AS imagen_titulo,
            i.url AS imagen_url
        FROM 
            imagenes i
        WHERE 
            i.id_producto = ?
    `

    // Ejecutar la consulta de imágenes
    rows, err := s.db.Query(queryImagenes, favorito.Producto.ID)
    if err != nil {
        return domain.Favoritos{}, err
    }
    defer rows.Close()

    // Iterar sobre las filas de imágenes y agregarlas al favorito
    for rows.Next() {
        var imagen domain.Imagen
        err := rows.Scan(
            &imagen.ID,
            &imagen.Titulo,
            &imagen.Url,
        )
        if err != nil {
            return domain.Favoritos{}, err
        }

        // Agregar la imagen al slice de imágenes del producto en favorito
        favorito.Producto.Imagenes = append(favorito.Producto.Imagenes, imagen)
    }

    // Manejar cualquier error que ocurra durante el procesamiento de filas
    err = rows.Err()
    if err != nil {
        return domain.Favoritos{}, err
    }

    return favorito, nil
}
*/

