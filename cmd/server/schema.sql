CREATE TABLE IF NOT EXISTS categorias (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nombre VARCHAR(255) NOT NULL,
    descripcion VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS productos (
    id INT AUTO_INCREMENT PRIMARY KEY,
    id_categoria INT,
    nombre VARCHAR(255) NOT NULL,
    descripcion VARCHAR(255) NOT NULL,
--    categoria VARCHAR(255) NOT NULL,
    precio DECIMAL(50,2) NOT NULL,
    stock INT NOT NULL,
    ranking DECIMAL(20,2) NOT NULL,
    FOREIGN KEY (id_categoria) REFERENCES categorias(id)
);

CREATE TABLE IF NOT EXISTS imagenes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    id_producto INT NOT NULL,  -- Agregar este campo para la relación
    titulo VARCHAR(255) NOT NULL,
    url VARCHAR(255) NOT NULL,
    FOREIGN KEY (id_producto) REFERENCES productos(id)
);

CREATE TABLE IF NOT EXISTS roles (
  id INT AUTO_INCREMENT PRIMARY KEY,
  nombre VARCHAR(20) NOT NULL
);


CREATE TABLE IF NOT EXISTS usuarios (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nombre VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    telefono VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    id_rol INT,
    estado_cuenta VARCHAR(50), 
    FOREIGN KEY (id_rol) REFERENCES roles(id)
);

CREATE TABLE IF NOT EXISTS fav (
    id INT AUTO_INCREMENT PRIMARY KEY,
    id_usuario INT NOT NULL,
    id_producto INT NOT NULL,
    FOREIGN KEY (id_usuario) REFERENCES usuarios(id),
    FOREIGN KEY (id_producto) REFERENCES productos(id)
);


CREATE TABLE IF NOT EXISTS estados (
  id INT AUTO_INCREMENT PRIMARY KEY,
  nombre VARCHAR(20) NOT NULL
);

CREATE TABLE IF NOT EXISTS ordenes (
  id INT AUTO_INCREMENT PRIMARY KEY,
  id_usuario INT NOT NULL,
  id_estado INT NOT NULL,
  fechaOrden VARCHAR(255)  NOT NULL,
  total DECIMAL(10,2) NOT NULL,
  FOREIGN KEY (id_usuario) REFERENCES usuarios(id),
  FOREIGN KEY (id_estado) REFERENCES estados(id)
);

CREATE TABLE IF NOT EXISTS OrdenProducto (
  id INT AUTO_INCREMENT PRIMARY KEY,
  id_orden INT NOT NULL,
  id_producto INT NOT NULL,
  cantidad INT NOT NULL,
  total DECIMAL(50,2) NOT NULL,
  FOREIGN KEY (id_orden) REFERENCES ordenes(id),
  FOREIGN KEY (id_producto) REFERENCES productos(id)
);

CREATE TABLE IF NOT EXISTS tarjetas (
  id INT AUTO_INCREMENT PRIMARY KEY,
  id_usuario INT NOT NULL,
  nombre VARCHAR(255)  NOT NULL,
  apellido VARCHAR(255)  NOT NULL,
  numero_tarjeta VARCHAR(255)  NOT NULL,
  clave_seguridad VARCHAR(255)  NOT NULL,
  vencimiento VARCHAR(255)  NOT NULL,
  ultimos_cuatro_digitos VARCHAR(255)  NOT NULL,
  FOREIGN KEY (id_usuario) REFERENCES usuarios(id)
);

CREATE TABLE IF NOT EXISTS DatosEnvio (
  id INT AUTO_INCREMENT PRIMARY KEY,
  id_usuario INT NOT NULL,
  nombre VARCHAR(255) NOT NULL,
  apellido VARCHAR(255) NOT NULL,
  direccion VARCHAR(255)  NOT NULL,
  apartamento VARCHAR(255) NOT NULL,
  ciudad VARCHAR(255) NOT NULL,
  codigo_postal VARCHAR(255) NOT NULL,
  estado VARCHAR(255) NOT NULL,
  FOREIGN KEY (id_usuario) REFERENCES usuarios(id)
);
CREATE TABLE IF NOT EXISTS InformacionCompra (
  id INT AUTO_INCREMENT PRIMARY KEY,
  id_orden INT NOT NULL,
  id_datosenvio INT NOT NULL,
  id_tarjetas INT NOT NULL,
  FOREIGN KEY (id_datosenvio) REFERENCES DatosEnvio(id),
  FOREIGN KEY (id_orden) REFERENCES ordenes(id),
  FOREIGN KEY (id_tarjetas) REFERENCES tarjetas(id)

);

CREATE TABLE IF NOT EXISTS ingredientes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    descripcion VARCHAR(255) NOT NULL

);

CREATE TABLE IF NOT EXISTS instrucciones (
    id INT AUTO_INCREMENT PRIMARY KEY,
    descripcion VARCHAR(255) NOT NULL

);

