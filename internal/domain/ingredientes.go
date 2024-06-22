package domain

type Ingredientes struct {
    ID          int    `json:"id"`
    Descripcion string `json:"descripcion" binding:"required"`
}