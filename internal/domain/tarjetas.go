package domain


type Tarjetas struct {
    ID                    int    `json:"id"`
    Nombre                string `json:"nombre"`
    Apellido              string `json:"apellido"`
    Numero_Tarjeta        string `json:"numero_tarjeta"`
    Clave_Seguridad       string `json:"clave_seguridad"`
    Vencimiento           string `json:"vencimiento"`	
    Ultimos_Cuatro_Digitos string `json:"ultimos_cuatro_digitos"`
    ID_Usuario            int    `json:"id_usuario"`  // Agrega este campo
}