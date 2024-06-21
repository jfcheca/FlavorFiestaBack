package domain

type InformacionCompra struct {
    ID           int `json:"id"`
    IDOrden      int `json:"id_orden"`
    IDDatosEnvio int `json:"id_datosenvio"`
    IDTarjeta    int `json:"id_tarjeta"`
}