package domain

type Favoritos struct {
	ID               int       `json:"id"`
	Id_producto      int       `json:"id_producto"`
	Id_usuario       int       `json:"id_usuario"`
	Producto         Producto  `json:"producto"`
	Imagenes         []Imagen  `json:"imagenes"`
}