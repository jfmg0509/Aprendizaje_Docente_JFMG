# Sistema de Gestión de Libros Electrónicos – POO en Go

## Estudiante
Juan Francisco Morán G.  
GitHub: https://github.com/jfmg0509

## Fecha
18 de diciembre de 2025

## Objetivo del programa
Desarrollar un sistema de software utilizando Programación Orientada a Objetos en el lenguaje Go, que permita gestionar libros electrónicos mediante una API REST, integrando base de datos MySQL, concurrencia, servicios web y una interfaz web básica en HTML.

## Tecnologías utilizadas
- Lenguaje Go
- Gorilla Mux
- MySQL (Laragon)
- godotenv
- HTML con templates
- Goroutines y Channels
- Testing en Go

## Funcionalidades principales
- Gestión de libros (crear, listar, actualizar y eliminar)
- Gestión de usuarios
- Registro de accesos concurrentes
- Servicios web REST
- Frontend HTML
- Manejo de errores
- Tests automatizados

## Servicios Web Implementados
- GET /books
- GET /books/{id}
- POST /books
- PUT /books/{id}
- DELETE /books/{id}
- GET /users
- POST /users
- POST /access

## Ejecución del sistema
1. Configurar MySQL en Laragon
2. Crear la base de datos usando el script schema.sql
3. Configurar el archivo .env
4. Ejecutar el proyecto con:
   go run ./cmd/api

## Observaciones
El proyecto cumple con los principios de Programación Orientada a Objetos, manejo de concurrencia y buenas prácticas de desarrollo de software.
