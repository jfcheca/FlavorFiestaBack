package domain

type Ingredientes struct {
    ID          int    `json:"id"`
    Descripcion string `json:"descripcion" binding:"required"`
    Id_mezclas int `json:"id_mezclas" binding:"required"`
}