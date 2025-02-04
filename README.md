
# Proyecto Labora

Este proyecto se encarga de procesar una base de datos de emails. Realiza las siguientes funciones:

1. Recorre un directorio con emails en formato texto.

2. Guarda los emails en una base de datos MySQL.

3. Indexa los emails en ZincSearch.

4. Proporciona una pequeña API para consultar los datos almacenados.


## Requisitos Previos

Antes de comenzar, asegúrese de tener lo siguiente instalado en su sistema:

Go v1.19 o superior.

Docker para ejecutar ZincSearch.

MySQL para la base de datos.
## Configuración del proyecto
Paso 1: Clonar el repositorio

```bash
  git clone https://github.com/tomygp97/proyecto-labora.git
cd proyecto-labora
```

Paso 2: Configurar la base de datos MySQL:
1. Cree una base de datos llamada emails_db:
```bash
CREATE DATABASE emails_db;
```

2. Cree la tabla emails ejecutando el siguiente script SQL en la base de datos:
```bash
CREATE TABLE emails (
    id INT AUTO_INCREMENT PRIMARY KEY,
    message_id VARCHAR(255) NOT NULL,
    sender VARCHAR(255),
    receiver VARCHAR(255),
    subject TEXT,
    mime_version VARCHAR(50),
    content_type VARCHAR(255),
    encoding VARCHAR(50),
    folder VARCHAR(255),
    body TEXT,
    date DATETIME
);
```

Paso 3: Levantar ZincSearch con Docker

Ejecute el siguiente comando para iniciar ZincSearch en un contenedor Docker:
```bash
docker run -d -v zinc-data:/data -p 4080:4080 \
    --name zincpublic2 \
    -e ZINC_FIRST_ADMIN_USER=admin \
    -e ZINC_FIRST_ADMIN_PASSWORD=ComplexPassword123 \
    public.ecr.aws/zinclabs/zinc:latest
```

Paso 4: Configurar las variables de entorno`

1. Copie el archivo .env.example y renómbrelo a .env:

```bash
cp .env.example .env
```

2. Actualice el archivo .env con las credenciales de su base de datos MySQL y ZincSearch si es necesario.

Paso 5: Ejecutar la aplicación

Ejecute el proyecto:
```bash
go run cmd/main.go
```

Esto iniciará el procesamiento de los emails.
## Uso de la API


#### Endpoints disponibles

| Método | endpoint     | Descripción                |
| :-------- | :------- | :------------------------- |
| `GET` | `/` | Bienvenida a la API |
| `GET` | `/emails` | Lista de todos los emails. |
| `GET` | `/emails/:id` | Devuelve un email especifico por su ID |
| `GET` | `/emails/search` | Realiza busquedas en los emails indexados |


La colección de Postman se encuentra en el repositorio en la ruta:


```
/docs/Labora APi.postman_collection.json
```

1. Abra Postman.

2. Importe la colección desde la ubicación indicada.

3. Configure la URL base de la API si es necesario.

## Autor

Tomas Gutierrez
- [@tomygp97](https://github.com/tomygp97)

