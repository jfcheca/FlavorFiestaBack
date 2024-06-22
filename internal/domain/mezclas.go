package domain

type Mezclas struct {
	ID           int    `json:"id"`
	Nombre 	     string   `json:"nombre" binding:"required"`
	Descripcion  string `json:"descripcion" binding:"required"`
	Ingredientes           []Ingredientes        `json:"ingredientes"`
	Instrucciones          []Instrucciones        `json:"instrucciones"`
	Imagenes           []Imagen        `json:"imagenes"`

}