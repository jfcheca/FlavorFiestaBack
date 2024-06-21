package domain

type DatosEnvio struct {
	ID           int    `json:"id"`
	IDUsuario    int    `json:"id_usuario"`
	Nombre       string `json:"nombre"`
	Apellido     string `json:"apellido"`
	Direccion    string `json:"direccion"`
	Apartamento  string `json:"apartamento"`
	Ciudad       string `json:"ciudad"`
	CodigoPostal string `json:"codigo_postal"`
	Estado       string `json:"estado"`
}
