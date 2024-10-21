# Loyalty Campaigns Project

Este proyecto implementa un sistema de fidelización de clientes con campañas personalizables para comercios.

## Requisitos previos

- Go 1.16 o superior
- Docker y Docker Compose
- Git

## Configuración del proyecto

1. Clonar el repositorio:
   ```
   git clone https://github.com/Williberto14/loyalty-campaigns.git
   cd loyalty-campaigns
   ```

2. Instalar dependencias:
   ```
   go mod tidy
   ```

## Levantar la base de datos

Utilizamos Docker para gestionar la base de datos PostgreSQL. Para levantarla:

1. Asegúrate de que Docker esté corriendo en tu sistema.
2. Desde la raíz del proyecto, ejecuta:
   ```
   docker-compose up -d
   ```

Esto iniciará un contenedor PostgreSQL con la configuración especificada en el archivo `docker-compose.yml`.

## Ejecutar migraciones

Para crear las tablas necesarias en la base de datos debe levantar el proyecto

## Ejecutar el proyecto

Para iniciar el servidor:

```
go run main.go
```

El servidor estará disponible en `http://localhost:7070`.

## Documentación de la API

La documentación de la API está disponible a través de Swagger UI. Para acceder a ella:

1. Asegúrate de que el servidor esté corriendo.
2. Abre un navegador y ve a `http://localhost:7070/swagger/index.html`.

## Ejecutar pruebas

Para ejecutar las pruebas unitarias:

```
go test ./...
```

## Detener el proyecto

Para detener el servidor, simplemente presiona `Ctrl+C` en la terminal donde está corriendo.

Para detener y eliminar el contenedor de la base de datos:

```
docker-compose down
```