package store

import (
	"github.com/jfcheca/FlavorFiesta/internal/domain"
)
//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> PRODUCTOS <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
type StoreInterfaceProducto interface {
	BuscarProducto(id int) (domain.Producto, error)
	BuscarTodosLosProductos() ([]domain.Producto, error)
	CrearProducto(producto domain.Producto) error
	UpdateProducto(id int, p domain.Producto) error
	Patch(id int, updatedFields map[string]interface{}) (domain.Producto, error)
	DeleteProducto(id int) error
	ExistsByID(id int) bool
	ObtenerNombreCategoria(id int) (string, error) // Añadir este método
}

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> IMAGENES <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
type StoreInterfaceImagenes interface {
	//	CrearImagen(imagen domain.Imagen) error
	CrearImagenes(imagenes []domain.Imagen) error
	CrearImagenesMezclas(imagenes []domain.Imagen) error
	BuscarImagen(id int) (domain.Imagen, error)
	UpdateImagen(id int, p domain.Imagen) (domain.Imagen, error)
	DeleteImagen(id int) error
	PatchImagen(id int, updatedFields map[string]interface{}) error
	ExistsByIDImagen(id int) bool
}

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> USUARIOS <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
type StoreInterfaceUsuarios interface {
	ExisteEmail(email string) (bool, error)
	ExisteEmail2(email string) (domain.Usuarios, error)
	ExisteCelular(celular string) (bool, error)
	CrearUsuario(usuario domain.Usuarios) error
	BuscarUsuario(id int) (domain.Usuarios, error)
	//BuscarUsuarioPorEmailYPassword(email, password string) (domain.Usuarios, error)
	BuscarUsuarioPorEmailYPassword(email, password string) (bool, error)
	BuscarUsuarioPorEmailYPassword2(email, password string) (domain.Usuarios, error)
	BuscarUsuarioPorEmailYPassword3(email, password string) (bool, error, domain.Usuarios)
	BuscarTodosLosUsuarios() ([]domain.Usuarios, error)
	DeleteUsuario(id int) error
	ExistsByIDUsuario(id int) (bool, error)
	ActivarCuenta(email string) error
	Update(usuario domain.Usuarios) error
	UpdatePassword(id int, newPassword string) (domain.Usuarios, error)
	ActivarCuentaEstado2(id int, estadoCuenta string) (domain.Usuarios, error)
	PatchUsuario(id int, updatedFields map[string]interface{}) error
}

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> CATEGORIAS <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
type StoreInterfaceCategorias interface {
	CrearCategoria(categoria domain.Categoria) error
	BuscarCategoria(id int) (domain.Categoria, error)
	BuscarTodosLasCategorias() ([]domain.Categoria, error)
	DeleteCategoria(id int) error
	ExistsByIDCategoria(id int) (bool, error)
	Update(categoria domain.Categoria) error
	PatchCategoria(id int, updatedFields map[string]interface{}) error
}

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ROLES <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
type StoreInterfaceRoles interface {
	CrearRol(rol domain.Rol) error
	BuscarTodosLosRoles() ([]domain.Rol, error)
	CambiarRol(usuarioID int, nuevoRol string) error
	// BuscarRol(id int) (domain.Rol, error)
}

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ESTADOS <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
type StoreInterfaceEstados interface {
	CrearEstados(estado domain.Estado) error
	BuscarTodosLosEstados() ([]domain.Estado, error)
	BuscarEstado(id int) (domain.Estado, error)
}

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ORDENES <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
type StoreInterfaceOrdenes interface {
	CrearOrden(orden domain.Orden) error
	BuscarOrden(id int) (domain.Orden, error)
	UpdateOrden(id int, orden domain.Orden) (domain.Orden, error)
	DeleteOrden(id int) error
	ExistsByIDOrden(id int) bool
	BuscarOrdenPorUsuarioYEstado(UserID, Estado string) (bool, error)
	BuscarOrdenPorUsuarioYEstado2(UserID, Estado string) (bool, error, domain.Orden)
	ObtenerOrdenesPorUsuarioYEstadoDiferenteA1(userID int) ([]domain.Orden, error)
}

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ORDEN PRODUCTO <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
type StoreInterfaceOrdenProducto interface {
	CrearOrdenProducto(op domain.OrdenProducto) (domain.OrdenProducto, error)
	BuscaOrdenProducto(id int) (domain.OrdenProducto, error)
	UpdateOrdenProducto(id int, p domain.OrdenProducto) error
	ExistsByID(id int) bool
	BuscarTodasLasOrdenesProducto() ([]domain.OrdenProducto, error)
	BuscarOrdenesProductoPorIDOrden(idOrden int) ([]domain.OrdenProducto, error)
	DeleteOrdenProducto(id int) error
}

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> FAVORITOS <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
type StoreInterfaceFavoritos interface {
	AgregarFavorito(favorito domain.Favoritos) error
	DeleteFavorito(id int) error
	BuscarFavorito(id int) (domain.Favoritos, error)
	BuscarFavoritosPorUsuario(idUsuario int) ([]domain.Favoritos, error) // Nueva función
}

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> DATOS TARJETAS <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
type StoreInterfaceDatosTarjetas interface {
	CargarTarjeta(tarjeta domain.Tarjetas) error
	BuscarTarjeta(id int) (domain.Tarjetas, error)
	DeleteTarjeta(id int) error
}

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> DATOS ENVIOS <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
type StoreInterfaceDatosEnvios interface {
    CrearDatosEnvio(datosEnvio domain.DatosEnvio) error
    BuscarTodosLosDatosEnvio() ([]domain.DatosEnvio, error)
    BuscarDatosEnvio(id int) (domain.DatosEnvio, error)
    EditarDatosEnvio(datosEnvio domain.DatosEnvio) error
    EliminarDatosEnvio(id int) error
}

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> INFORMACION COMPRA <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
type StoreInterfaceInformacionCompra interface {
    CrearInformacionCompra(ic domain.InformacionCompra) (domain.InformacionCompra, error)
    BuscarInformacionCompra(id int) (domain.InformacionCompra, error)
    UpdateInformacionCompra(id int, ic domain.InformacionCompra) (domain.InformacionCompra, error)
    DeleteInformacionCompra(id int) error
	ObtenerInformacionCompletaCompra(idOrden int) (domain.Orden, domain.InformacionCompra, domain.DatosEnvio, domain.Tarjetas, []domain.OrdenProducto, error)
}

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> INGREDIENTES <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
type StoreInterfaceIngredientes interface {
	CrearIngredientes(ingredientes []domain.Ingredientes) error
	BuscarIngredientes(id int) (domain.Ingredientes, error)
	DeleteIngredientes(id int) error
}

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> INSTRUCCIONES <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
type StoreInterfaceInstrucciones interface {
	CrearInstrucciones(instrucciones []domain.Instrucciones) error
	BuscarInstrucciones(id int) (domain.Instrucciones, error)
	DeleteInstrucciones(id int) error

}

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> MEZCLAS <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
type StoreInterfaceMezclas interface {
	CrearMezcla(mezcla domain.Mezclas) error
	BuscarMezcla(id int) (domain.Mezclas, error)
	DeleteMezclas(id int) error
	
	
}
