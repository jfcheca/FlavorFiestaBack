package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jfcheca/FlavorFiesta/cmd/server/handler"
	//"github.com/jfcheca/FlavorFiesta/cmd/server/middleware"
    "github.com/jfcheca/FlavorFiesta/internal/favoritos"
    "github.com/jfcheca/FlavorFiesta/internal/auth"
	"github.com/jfcheca/FlavorFiesta/internal/categorias"
	"github.com/jfcheca/FlavorFiesta/internal/estados"
	"github.com/jfcheca/FlavorFiesta/internal/imagenes"
	"github.com/jfcheca/FlavorFiesta/internal/ordenProducto"
	"github.com/jfcheca/FlavorFiesta/internal/ordenes"
	"github.com/jfcheca/FlavorFiesta/internal/productos"
	"github.com/jfcheca/FlavorFiesta/internal/roles"
	"github.com/jfcheca/FlavorFiesta/internal/usuarios"
	"github.com/jfcheca/FlavorFiesta/internal/tarjetas"	
	"github.com/jfcheca/FlavorFiesta/internal/informacioncompras"
	"github.com/jfcheca/FlavorFiesta/internal/datosenvio"
	"github.com/jfcheca/FlavorFiesta/internal/ingredientes"
	"github.com/jfcheca/FlavorFiesta/internal/instrucciones"
	"github.com/jfcheca/FlavorFiesta/internal/mezclas"
 //   "github.com/jfcheca/FlavorFiesta/internal/favoritos"
	"github.com/jfcheca/FlavorFiesta/pkg/store"
	"github.com/joho/godotenv"
	// "github.com/jfcheca/FlavorFiesta/internal/auth"
	//	"gopkg.in/mail.v2"

	"io/ioutil"
	"strings"
)


func main() {

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> CARGAMOS LAS VARIABLES DE ENTORNO DEL ARCHIVO .ENV >>>>>>>>>>>>>>>>>>>>
    err := godotenv.Load(".env")
    if err != nil {
        log.Fatal("Error al cargar el archivo .env:", err)
    }

    log.Printf("SMTP_EMAIL: %s", os.Getenv("SMTP_EMAIL"))
    log.Printf("SMTP_PASSWORD: %s", os.Getenv("SMTP_PASSWORD"))
    
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbName := os.Getenv("DB_NAME")

    // Abrir una conexión temporal a MySQL para ejecutar comandos administrativos
    dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/"
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal("Error al conectar con MySQL:", err)
    }
    defer db.Close()

    // Eliminar la base de datos si ya existe
    _, err = db.Exec("DROP DATABASE IF EXISTS " + dbName)
    if err != nil {
        log.Fatal("Error al eliminar la base de datos '" + dbName + "':", err)
    }

    // Crear la base de datos
    _, err = db.Exec("CREATE DATABASE " + dbName)
    if err != nil {
        log.Fatal("Error al crear la base de datos '" + dbName + "':", err)
    }

    // Conectar a la base de datos
    dsn = dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName
    bd, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal("Error al conectar con la base de datos '" + dbName + "':", err)
    }
    defer bd.Close()

    // Cargar contenido del archivo schema.sql
    sqlFile, err := ioutil.ReadFile("schema.sql")
    if err != nil {
        log.Fatal("Error al leer el archivo schema.sql:", err)
    }

    // Dividir el contenido en sentencias SQL individuales
    sqlStatements := strings.Split(string(sqlFile), ";")

    // Ejecutar cada sentencia SQL en el archivo schema.sql
    for _, statement := range sqlStatements {
        cleanedStatement := strings.TrimSpace(statement)
        if cleanedStatement == "" {
            continue
        }

        _, err := bd.Exec(cleanedStatement)
        if err != nil {
            log.Fatal("Error al ejecutar la sentencia SQL:", err)
        }
    }

    // Configurar el enrutador Gin
    r := gin.Default()
    r.Static("/Probando", "./public")

    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:5173"}, // URL del frontend
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        AllowCredentials: true,
    }))

    // Definir rutas
    r.GET("/api/ping", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "pong"})
    })

////////////////////////////////////////// >>>>>>>>>>>>>> TODO LO REFERIDO A LA AUTENTICACION >>>>>>>>>>>>>>>>>>>>>>>>>>>>
 /*   r.POST("/api/login", func(c *gin.Context) {
        var credentials auth.Credentials
        if err := c.ShouldBindJSON(&credentials); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        token, err := auth.Authenticate(credentials)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"token": token})
    })
*/


 //Middleware de autenticación
  // r.Use(middleware.AuthMiddleware())

	// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> PRODUCTOS <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
	storageProducto := store.NewSqlStoreProductos(bd)
	repoProducto := productos.NewRepository(storageProducto)
	serviceProducto := productos.NewService(repoProducto)
	productoHandler := handler.NewProductHandler(serviceProducto)

	// Rutas para el manejo de productos
	productosGroup := r.Group("/productos")
	{
		productosGroup.GET("/:id", productoHandler.BuscarProducto())
		productosGroup.GET("/", productoHandler.GetAll())
		productosGroup.POST("/crear", productoHandler.Post())
		productosGroup.DELETE("/:id", productoHandler.Delete())
		productosGroup.PATCH("/:id", productoHandler.Patch())
		productosGroup.PUT("/:id", productoHandler.Put())
	}

	// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> IMÁGENES <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
	storageImagen := store.NewSqlStoreImagen(bd)
	repoImagen := imagenes.NewRepository(storageImagen)
	serviceImagen := imagenes.NewService(repoImagen)
	imagenHandler := handler.NewImagenHandler(serviceImagen)

	// Rutas para el manejo de imágenes
	imagenesGroup := r.Group("/imagenes")
	{
		imagenesGroup.GET("/:id", imagenHandler.GetByID())
		imagenesGroup.POST("/crear", imagenHandler.Post())
		imagenesGroup.POST("/crearmezcla", imagenHandler.PostMezcla())
		imagenesGroup.DELETE("/:id", imagenHandler.Delete())
		imagenesGroup.PATCH("/:id", imagenHandler.Patch())
		imagenesGroup.PUT("/:id", imagenHandler.Put())
	}

	/// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> USUARIOS <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
	// Authent
	storageAuth := store.NewSqlStoreUsuarios(bd) // Reuse user store for authentication
	repoAuth := auth.NewRepository(storageAuth)
	serviceAuth := auth.NewService(repoAuth)
	storageUsuario := store.NewSqlStoreUsuarios(bd)
	repoUsuario := usuarios.NewRepository(storageUsuario, auth.NewRepository(storageUsuario))
	serviceUsuario := usuarios.NewService(repoUsuario)
	usuariosHandler := handler.NewUsuarioHandler(serviceUsuario, serviceAuth)
	authHandler := handler.NewAuthHandler(serviceAuth, serviceUsuario)

	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/login", authHandler.Login())
		authRoutes.POST("/forgotPassword", authHandler.ForgotPassword())
		authRoutes.POST("/tokenActivarCuenta", authHandler.ActivarCuenta())
	}

	// Rutas para el manejo de usuarios
	usuariosGroup := r.Group("/usuarios")
	{
		usuariosGroup.GET("/:id", usuariosHandler.GetByID())
		usuariosGroup.GET("/email&pass", usuariosHandler.GetByEmailAndPassword())
		usuariosGroup.GET("/email&passdatos", usuariosHandler.GetByEmailAndPasswordConDatos())
		usuariosGroup.GET("/", usuariosHandler.GetAll())
		usuariosGroup.POST("/crear", usuariosHandler.Post())
		usuariosGroup.DELETE("/:id", usuariosHandler.DeleteUsuario())
		usuariosGroup.PUT("/:id", usuariosHandler.Put())
		usuariosGroup.PATCH("/:id", usuariosHandler.Patch())
		usuariosGroup.PUT("/forgotPassword/:id", usuariosHandler.UpdatePassword())
		usuariosGroup.PUT("/ActivarCuentaEstado/:id", usuariosHandler.ActivarCuentaEstado2())

	}

	// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> CATEGORÍAS <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
	storageCategoria := store.NewSqlStoreCategorias(bd)
	repoCategoria := categorias.NewRepository(storageCategoria)
	serviceCategoria := categorias.NewService(repoCategoria)
	categoriasHandler := handler.NewCategoriaHandler(serviceCategoria)

	// Rutas para el manejo de categorías
	categoriasGroup := r.Group("/categorias")
	{
		categoriasGroup.GET("/:id", categoriasHandler.GetByID())
		categoriasGroup.GET("/", categoriasHandler.GetAll())
		categoriasGroup.POST("/crear", categoriasHandler.Post())
		categoriasGroup.DELETE("/:id", categoriasHandler.DeleteCategoria())
		categoriasGroup.PATCH("/:id", categoriasHandler.Patch())
		categoriasGroup.PUT("/:id", categoriasHandler.Put())
	}

	// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ÓRDENES <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
	storageOrden := store.NewSqlStoreOrden(bd)
	repoOrden := ordenes.NewRepository(storageOrden)
	serviceOrden := ordenes.NewService(repoOrden)
	ordenHandler := handler.NewOrdenHandler(serviceOrden)

	// Rutas para el manejo de órdenes
	ordenesGroup := r.Group("/ordenes")
	{
		ordenesGroup.GET("/:id", ordenHandler.GetOrdenByID())
		ordenesGroup.GET("/user/:userID/estado-diferente-a-1", ordenHandler.ObtenerOrdenesPorUsuarioYEstadoDiferenteA1())
		ordenesGroup.GET("/usuario&estado", ordenHandler.GetOrdenByUserIDyOrden())
		ordenesGroup.GET("/usuario&estadoConDatos", ordenHandler.GetOrdenByUsuarioYEstadoConDatos())
		ordenesGroup.POST("/crear", ordenHandler.Post())
		ordenesGroup.PUT("/:id", ordenHandler.Put())
		ordenesGroup.DELETE("/:id", ordenHandler.Delete())
	}

	// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ORDEN PRODUCTO <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
	storageOrdenProducto := store.NewSqlStoreOrdenProducto(bd)
	repoOrdenProducto := ordenProductos.NewRepository(storageOrdenProducto)
	serviceOrdenProducto := ordenProductos.NewService(repoOrdenProducto)
	ordenProductoHandler := handler.NewOrdenProductoHandler(serviceOrdenProducto)

	// Rutas para el manejo de órdenes de productos
	ordenProductosGroup := r.Group("/ordenProductos")
	{
		ordenProductosGroup.GET("/:id", ordenProductoHandler.GetByID())
		ordenProductosGroup.GET("/orden/:idOrden", ordenProductoHandler.BuscarPorIDOrden())
		ordenProductosGroup.GET("/", ordenProductoHandler.GetAll())
		ordenProductosGroup.POST("/crear", ordenProductoHandler.Post())
		ordenProductosGroup.PUT("/:id", ordenProductoHandler.Put())
		ordenProductosGroup.DELETE("/:id", ordenProductoHandler.Delete())
	}

	// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ESTADOS <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
	storageEstado := store.NewSqlStoreEstados(bd)
	repoEstado := estados.NewRepository(storageEstado)
	serviceEstado := estados.NewService(repoEstado)
	estadoHandler := handler.NewEstadoHandler(serviceEstado)

	// Rutas para el manejo de estados
	estadosGroup := r.Group("/estados")
	{
		estadosGroup.GET("/:id", estadoHandler.BuscarEstado())
		estadosGroup.GET("/", estadoHandler.GetAll())
		estadosGroup.POST("/crear", estadoHandler.Post())
	}

	//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>FAVORITOS>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

	// Crear el almacenamiento SQL con la base de datos 'FlavorFiesta'
	storageFavoritos := store.NewSqlStoreFavoritos(db)
	repoFavoritos := favoritos.NewRepository(storageFavoritos)
	serviceFavoritos := favoritos.NewServiceFavoritos(repoFavoritos)
	favoritoHandler := handler.NewFavoritosHandler(serviceFavoritos)

	// Rutas para el manejo de órdenes
	favoritos := r.Group("/favoritos")
	{

		favoritos.POST("/agregar", favoritoHandler.Post())
		favoritos.GET("/:id", favoritoHandler.GetByID())
		favoritos.GET("/usuario/:id", favoritoHandler.GetFavoritosPorUsuario())
		favoritos.DELETE("/user/:id_usuario/producto/:id_producto", favoritoHandler.DeleteFavorito())
		
	}

	// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> ROLES <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
	storageRol := store.NewSqlStoreRoles(bd)
	repoRol := roles.NewRepository(storageRol)
	serviceRol := roles.NewService(repoRol)
	rolHandler := handler.NewRolHandler(serviceRol)

	// Rutas para el manejo de roles
	roles := r.Group("/roles")
	{
		roles.GET("/", rolHandler.GetAll())
		roles.POST("/crear", rolHandler.Post())
	}

	// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> TARJETAS <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
	storageTarjetas := store.NewSqlStoreTarjetas(bd)
	repoTarjetas := tarjetas.NewRepository(storageTarjetas)
	serviceTarjetas := tarjetas.NewService(repoTarjetas)
	tarjetasHandler := handler.NewTarjetaHandler(serviceTarjetas)

	// Rutas para el manejo de categorías
	tarjetasGroup := r.Group("/tarjetas")
	{
		tarjetasGroup.POST("/crear", tarjetasHandler.Post())
		tarjetasGroup.GET("/:id", tarjetasHandler.GetByID())
		tarjetasGroup.DELETE("/:id", tarjetasHandler.DeleteTarjeta())

	
	}

		// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> DATOS ENVÍO <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
		storageDatosEnvio := store.NewSqlStoreDatosEnvios(bd)
		repoDatosEnvio := datosenvio.NewRepository(storageDatosEnvio)
		serviceDatosEnvio := datosenvio.NewService(repoDatosEnvio)
		datosEnvioHandler := handler.NewDatosEnvioHandler(serviceDatosEnvio)
	
		// Rutas para el manejo de datos de envío
		datosEnvioGroup := r.Group("/datosEnvio")
		{
			datosEnvioGroup.GET("/:id", datosEnvioHandler.GetByID())
			datosEnvioGroup.GET("/", datosEnvioHandler.GetAll())
			datosEnvioGroup.POST("/crear", datosEnvioHandler.Post())
			datosEnvioGroup.DELETE("/:id", datosEnvioHandler.Delete())
			datosEnvioGroup.PUT("/:id", datosEnvioHandler.Put())
		}

			// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> INFORMACION COMPRA <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
	storageInformacionCompra := store.NewSqlStoreInformacionCompra(bd)
	repoInformacionCompra := informacioncompras.NewRepositoryInformacionCompras(storageInformacionCompra)
	serviceInformacionCompra := informacioncompras.NewService(repoInformacionCompra)
	estadoInformacionCompra := handler.NewInformacionCompraHandler(serviceInformacionCompra)

	// Rutas para el manejo de estados
	informacionCompra := r.Group("/informacionCompra")
	{
		informacionCompra.GET("/:id", estadoInformacionCompra.GetByID())
		informacionCompra.GET("/informacionCompleta/:id_orden", estadoInformacionCompra.ObtenerInformacionCompletaCompraByIDOrden()) 
		informacionCompra.POST("/crear", estadoInformacionCompra.Post())
		informacionCompra.PUT("/:id", estadoInformacionCompra.Put())
		informacionCompra.DELETE("/:id", estadoInformacionCompra.Delete())


	}

	// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> INGREDIENTES <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
	storageIngredientes := store.NewSqlStoreIngredientes(bd)
	repoIngredientes := ingredientes.NewRepository(storageIngredientes)
	serviceIngredientes := ingredientes.NewService(repoIngredientes)
	ingredientesHandler := handler.NewIngredientesHandler(serviceIngredientes)

	// Rutas para el manejo de imágenes
	ingredientesGroup := r.Group("/ingredientes")
	{
		ingredientesGroup.POST("/crear", ingredientesHandler.Post())
		ingredientesGroup.GET("/:id", ingredientesHandler.GetByID())
		ingredientesGroup.DELETE("/:id", ingredientesHandler.Delete())
		
//		imagenesGroup.DELETE("/:id", imagenHandler.Delete())
//		imagenesGroup.PATCH("/:id", imagenHandler.Patch())
//		imagenesGroup.PUT("/:id", imagenHandler.Put())
	}

	// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> INSTRUCCIONES <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
	storageInstrucciones := store.NewSqlStoreInstrucciones(bd)
	repoInstrucciones := instrucciones.NewRepository(storageInstrucciones)
	serviceInstrucciones := instrucciones.NewService(repoInstrucciones)
	instruccionesHandler := handler.NewInstruccionesHandler(serviceInstrucciones)

	// Rutas para el manejo de imágenes
	instruccionesGroup := r.Group("/instrucciones")
	{
		instruccionesGroup.POST("/crear", instruccionesHandler.Post())
		instruccionesGroup.GET("/:id", instruccionesHandler.GetByID())
		instruccionesGroup.DELETE("/:id", instruccionesHandler.Delete())
		
//		imagenesGroup.DELETE("/:id", imagenHandler.Delete())
//		imagenesGroup.PATCH("/:id", imagenHandler.Patch())
//		imagenesGroup.PUT("/:id", imagenHandler.Put())
	}

	// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> MEZCLAS <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
	storageMezclas := store.NewSqlStoreMezclas(bd)
	repoMezclas := mezclas.NewRepository(storageMezclas)
	serviceMezclas := mezclas.NewService(repoMezclas)
	mezclasHandler := handler.NewMezclasHandler(serviceMezclas)
	

	// Rutas para el manejo de imágenes
	mezclasGroup := r.Group("/mezclas")
	{
		mezclasGroup.POST("/crear", mezclasHandler.Post())
		mezclasGroup.GET("/:id", mezclasHandler.GetByID())
		mezclasGroup.DELETE("/:id", mezclasHandler.Delete())
		mezclasGroup.GET("/", mezclasHandler.GetAll())
		
//		imagenesGroup.DELETE("/:id", imagenHandler.Delete())
//		imagenesGroup.PATCH("/:id", imagenHandler.Patch())
//		imagenesGroup.PUT("/:id", imagenHandler.Put())
	}

	/*   // Endpoints protegidos con middleware de rol ADMIN
	adminRoutes := r.Group("/admin")
	adminRoutes.Use(middleware.AdminRoleMiddleware())
	{
	    adminRoutes.PUT("/roles/cambiar", rolHandler.CambiarRol())
	}
	*/

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>ACA ES LA CONEXION PARA EL >
    // Ejecutar el servidor en el puerto 8080
//    r.Run(":8080")

 // >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>ACA ES LA CONEXION PARA LA>
 if err := http.ListenAndServe(":8080", r); err != nil {
    panic(err)
}
}
