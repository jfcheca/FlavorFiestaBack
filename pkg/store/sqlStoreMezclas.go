package store

import (
	"database/sql"
	"errors"
	"log"

	"fmt"

	"github.com/jfcheca/FlavorFiesta/internal/domain"
)

type sqlStoreMezclas struct {
	db *sql.DB
}

func NewSqlStoreMezclas(db *sql.DB) StoreInterfaceMezclas {
	return &sqlStoreMezclas{
		db: db,
	}
}


// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> CREAR UN NUEVO PRODUCTO <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func (s *sqlStoreMezclas) CrearMezcla(mezcla domain.Mezclas) error {
    // Insertar la mezcla
    query := "INSERT INTO mezclas (nombre, descripcion) VALUES (?, ?);"
    stmt, err := s.db.Prepare(query)
    if err != nil {
        return fmt.Errorf("error preparing query: %w", err)
    }
    defer stmt.Close()

    res, err := stmt.Exec(mezcla.Nombre, mezcla.Descripcion)
    if err != nil {
        return fmt.Errorf("error executing query: %w", err)
    }

    // Obtener el ID generado para la nueva mezcla
    mezclaID, err := res.LastInsertId()
    if err != nil {
        return fmt.Errorf("error getting last insert ID: %w", err)
    }

    // Insertar ingredientes
    for _, ingrediente := range mezcla.Ingredientes {
        _, err := s.db.Exec("INSERT INTO ingredientes (mezcla_id, nombre) VALUES (?, ?);", mezclaID, ingrediente.Descripcion)
        if err != nil {
            return fmt.Errorf("error inserting ingredient: %w", err)
        }
    }

    // Insertar instrucciones
    for _, instruccion := range mezcla.Instrucciones {
        _, err := s.db.Exec("INSERT INTO instrucciones (mezcla_id, paso) VALUES (?, ?);", mezclaID, instruccion.Descripcion)
        if err != nil {
            return fmt.Errorf("error inserting instruction: %w", err)
        }
    }

 /*   // Insertar imagen (asumiendo que solo se guarda una imagen por mezcla)
    if len(mezcla.Imagenes) > 0 {
        _, err := s.db.Exec("INSERT INTO imagenes (mezcla_id, url) VALUES (?, ?);", mezclaID, mezcla.Imagenes[0].Url)
        if err != nil {
            return fmt.Errorf("error inserting image: %w", err)
        }
    }
*/
    return nil
}

func (s *sqlStoreMezclas) BuscarMezcla(id int) (domain.Mezclas, error) {
    var mezcla domain.Mezclas

    // Consulta para obtener la mezcla por su ID
    query := "SELECT id, nombre, descripcion FROM mezclas WHERE id = ?"
    err := s.db.QueryRow(query, id).Scan(&mezcla.ID, &mezcla.Nombre, &mezcla.Descripcion)
    if err != nil {
        if err == sql.ErrNoRows {
            return domain.Mezclas{}, errors.New("mezcla not found")
        }
        return domain.Mezclas{}, err
    }

    // Obtener los ingredientes de la mezcla
    ingredientesQuery := `
        SELECT 
            id, descripcion, id_mezclas
        FROM 
            ingredientes
        WHERE 
            id_mezclas = ?
    `
    ingRows, err := s.db.Query(ingredientesQuery, mezcla.ID)
    if err != nil {
        return domain.Mezclas{}, fmt.Errorf("error querying ingredientes: %w", err)
    }
    defer ingRows.Close()

    var ingredientes []domain.Ingredientes
    for ingRows.Next() {
        var ingrediente domain.Ingredientes
        if err := ingRows.Scan(&ingrediente.ID, &ingrediente.Descripcion, &ingrediente.Id_mezclas); err != nil {
            return domain.Mezclas{}, fmt.Errorf("error scanning ingrediente: %w", err)
        }
        ingredientes = append(ingredientes, ingrediente)
    }
    mezcla.Ingredientes = ingredientes

    // Obtener las instrucciones de la mezcla
    instruccionesQuery := `
        SELECT 
            id, descripcion, id_mezclas
        FROM 
            instrucciones
        WHERE 
            id_mezclas = ?
    `
    instRows, err := s.db.Query(instruccionesQuery, mezcla.ID)
    if err != nil {
        return domain.Mezclas{}, fmt.Errorf("error querying instrucciones: %w", err)
    }
    defer instRows.Close()

    var instrucciones []domain.Instrucciones
    for instRows.Next() {
        var instruccion domain.Instrucciones
        if err := instRows.Scan(&instruccion.ID, &instruccion.Descripcion, &instruccion.Id_mezclas); err != nil {
            return domain.Mezclas{}, fmt.Errorf("error scanning instruccion: %w", err)
        }
        instrucciones = append(instrucciones, instruccion)
    }
    mezcla.Instrucciones = instrucciones

    // Obtener las imágenes de la mezcla
    imagenesQuery := `
        SELECT 
            id, titulo, url, id_mezclas
        FROM 
            imgmezcla
        WHERE 
            id_mezclas = ?
    `
    imgRows, err := s.db.Query(imagenesQuery, mezcla.ID)
    if err != nil {
        return domain.Mezclas{}, fmt.Errorf("error querying imágenes: %w", err)
    }
    defer imgRows.Close()

    var imagenes []domain.Imagen
    for imgRows.Next() {
        var imagen domain.Imagen
        if err := imgRows.Scan(&imagen.ID, &imagen.Titulo, &imagen.Url, &imagen.Id_mezclas); err != nil {
            return domain.Mezclas{}, fmt.Errorf("error scanning imagen: %w", err)
        }
        imagenes = append(imagenes, imagen)
    }
    mezcla.Imagenes = imagenes

    return mezcla, nil
}

func (s *sqlStoreMezclas) DeleteMezclas(id int) error {
	query := "DELETE FROM mezclas WHERE id = ?;"
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

// Función para buscar todas las mezclas con sus ingredientes, instrucciones e imágenes
func (s *sqlStoreMezclas) BuscarTodasLasMezclas() ([]domain.Mezclas, error) {
    var mezclas []domain.Mezclas

    // Consulta para obtener todas las mezclas
    query := `
        SELECT 
            m.id, m.nombre, m.descripcion
        FROM 
            mezclas m
    `

    rows, err := s.db.Query(query)
    if err != nil {
        return nil, fmt.Errorf("error querying mezclas: %w", err)
    }
    defer rows.Close()

    for rows.Next() {
        var mezcla domain.Mezclas

        if err := rows.Scan(&mezcla.ID, &mezcla.Nombre, &mezcla.Descripcion); err != nil {
            return nil, fmt.Errorf("error scanning mezcla: %w", err)
        }

        // Obtener los ingredientes de la mezcla
        ingredientesQuery := `
            SELECT 
                i.id, i.descripcion, i.id_mezclas
            FROM 
                ingredientes i
            WHERE 
                i.id_mezclas = ?
        `
        ingredientesRows, err := s.db.Query(ingredientesQuery, mezcla.ID)
        if err != nil {
            return nil, fmt.Errorf("error querying ingredientes: %w", err)
        }
        defer ingredientesRows.Close()

        var ingredientes []domain.Ingredientes
        for ingredientesRows.Next() {
            var ingrediente domain.Ingredientes
            if err := ingredientesRows.Scan(&ingrediente.ID, &ingrediente.Descripcion, &ingrediente.Id_mezclas); err != nil {
                return nil, fmt.Errorf("error scanning ingrediente: %w", err)
            }
            ingredientes = append(ingredientes, ingrediente)
        }
        mezcla.Ingredientes = ingredientes

        // Obtener las instrucciones de la mezcla
        instruccionesQuery := `
            SELECT 
                ins.id, ins.descripcion, ins.id_mezclas
            FROM 
                instrucciones ins
            WHERE 
                ins.id_mezclas = ?
        `
        instruccionesRows, err := s.db.Query(instruccionesQuery, mezcla.ID)
        if err != nil {
            return nil, fmt.Errorf("error querying instrucciones: %w", err)
        }
        defer instruccionesRows.Close()

        var instrucciones []domain.Instrucciones
        for instruccionesRows.Next() {
            var instruccion domain.Instrucciones
            if err := instruccionesRows.Scan(&instruccion.ID, &instruccion.Descripcion, &instruccion.Id_mezclas); err != nil {
                return nil, fmt.Errorf("error scanning instruccion: %w", err)
            }
            instrucciones = append(instrucciones, instruccion)
        }
        mezcla.Instrucciones = instrucciones

        // Obtener las imágenes de la mezcla
        imagenesQuery := `
            SELECT 
                img.id, img.titulo, img.url, img.id_mezclas
            FROM 
                imgmezcla img
            WHERE 
                img.id_mezclas = ?
        `
        imagenesRows, err := s.db.Query(imagenesQuery, mezcla.ID)
        if err != nil {
            return nil, fmt.Errorf("error querying imágenes: %w", err)
        }
        defer imagenesRows.Close()

        var imagenes []domain.Imagen
        for imagenesRows.Next() {
            var imagen domain.Imagen
            if err := imagenesRows.Scan(&imagen.ID, &imagen.Titulo, &imagen.Url, &imagen.Id_mezclas); err != nil {
                return nil, fmt.Errorf("error scanning imagen: %w", err)
            }
            imagenes = append(imagenes, imagen)
        }
        mezcla.Imagenes = imagenes

        mezclas = append(mezclas, mezcla)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating rows: %w", err)
    }

    return mezclas, nil
}