# Sistema de Gestión de Libros Electrónicos (POO en Go)

**Autor / Grupo:** Juan Francisco Morán Gortaire 
**Materia:** Programación Orientada a Objetos  
**Fecha:** 18 de diciembre de 2025

Este repositorio implementa una **plataforma web** para **gestionar una biblioteca digital de libros técnicos**, con:
- Gestión de usuarios y roles.
- Registro, clasificación, búsqueda avanzada de libros.
- Registro de accesos (APERTURA, LECTURA, DESCARGA) y estadísticas por libro (map[AccessType]int).

La arquitectura del proyecto está organizada por capas:
- `/internal/domain` (entidades + interfaces)
- `/internal/usecase` (lógica)
- `/internal/infrastructure/*` (DB/config)
- `/internal/transport/http` (handlers/rutas)
- `/cmd/api` (entrypoint)

---

## Objetivo del programa

Desarrollar una plataforma para administrar un repositorio de libros técnicos y permitir acceso controlado/registrado.

---

## Tecnologías / Librerías

- Go (net/http, html/templates, encoding/json, context, etc.)
- **MySQL** (Laragon) + driver: `github.com/go-sql-driver/mysql`
- **Router:** `github.com/gorilla/mux`
- **Variables de entorno:** `github.com/joho/godotenv`

---

## Cómo ejecutar (Laragon + MySQL)

1) Crea la base y tablas ejecutando el script:

```sql
-- scripts/schema.sql
```

2) Copia el ejemplo de variables de entorno:

```bash
cp .env.example .env
```

3) Ejecuta:

```bash
go mod tidy
go run ./cmd/api
```

4) Abre:
- Frontend: `http://localhost:8081/`
- Health: `http://localhost:8081/health`
- API: `http://localhost:8081/api/...`

---

## Servicios Web (REST) implementados (>= 8)

### Usuarios
1. `GET /api/users` → listar usuarios  
2. `POST /api/users` → crear usuario  
3. `GET /api/users/{id}` → obtener usuario  
4. `PUT /api/users/{id}` → actualizar usuario  
5. `DELETE /api/users/{id}` → eliminar usuario  

### Libros
6. `GET /api/books` → listar libros  
7. `POST /api/books` → crear libro  
8. `GET /api/books/search?q=&author=&category=` → búsqueda avanzada  
9. `GET /api/books/{id}` → obtener libro  
10. `PUT /api/books/{id}` → actualizar libro  
11. `DELETE /api/books/{id}` → eliminar libro  
12. `GET /api/books/{id}/stats` → estadísticas de accesos  

### Accesos
13. `POST /api/access` → registrar acceso (APERTURA/LECTURA/DESCARGA)

---

## Frontend (HTML Templates + HTTP)

- Templates en `web/templates/*.html`
- Endpoints UI:
  - `GET /ui/users` (listar + formulario)
  - `GET /ui/books` (listar + formulario)
  - `GET /ui/books/search`
  - `GET /ui/books/{id}` (detalle + registrar acceso)

Los formularios usan **HTTP POST** y redirecciones; el backend expone además API REST completa (GET/POST/PUT/DELETE).

---

## Concurrencia (Goroutines + Canales)

El registro de accesos se realiza con una **cola** (`AccessQueue`) basada en `chan *AccessEvent` y **workers** en goroutines:
- El handler encola el evento.
- Los workers insertan en MySQL asíncronamente.
- Si la cola está llena, se hace fallback a inserción síncrona.

---

## Pruebas (testing)

Ejecutar:

```bash
go test ./...
```

---

## Estructura del repositorio

```text
cmd/api                # main
internal/domain         # entidades (User/Book/AccessEvent) + interfaces
internal/usecase        # lógica negocio + cola concurrente
internal/infrastructure # DB MySQL + config (godotenv)
internal/transport/http # mux router + handlers (JSON + templates)
web/templates           # frontend HTML
scripts                 # SQL schema
```
