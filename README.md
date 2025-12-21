Sistema de Gestión de Libros

Descripción General

El Sistema de Gestión de Libros es una aplicación web desarrollada como proyecto académico para la asignatura Programación Orientada a Objetos, cuyo objetivo es demostrar la correcta aplicación de principios de diseño orientado a objetos, arquitectura en capas, separación de responsabilidades y desarrollo de interfaces web.

El sistema permite la gestión básica de usuarios y libros, integrando un frontend web y un backend estructurado, con persistencia de datos en una base de datos relacional.

Objetivo del Proyecto
Objetivo General

Desarrollar una aplicación web funcional que implemente los principios de la Programación Orientada a Objetos, integrando backend, frontend y base de datos, bajo una arquitectura limpia y mantenible.

Objetivos Específicos

Aplicar correctamente conceptos de encapsulación, abstracción y modularidad.

Implementar una arquitectura en capas (dominio, casos de uso, infraestructura y presentación).

Diseñar interfaces web claras, funcionales y orientadas al usuario final.

Justificación del Proyecto

Este proyecto fue seleccionado debido a que representa un caso realista y escalable de aplicación de la Programación Orientada a Objetos en sistemas modernos.
La gestión de información (usuarios y libros) es un problema común en múltiples organizaciones, lo que permite visualizar su proyección futura hacia sistemas más complejos como bibliotecas digitales, sistemas académicos o plataformas educativas.

Además, el proyecto integra:

Backend estructurado

Frontend web

Base de datos relacional

Buenas prácticas de desarrollo

Arquitectura del Sistema

El sistema está organizado bajo una arquitectura en capas, lo que facilita la mantenibilidad, escalabilidad y pruebas.

Capas del sistema:

Dominio

Entidades (User, Book)

Reglas de negocio

Casos de Uso (Usecase)

Lógica de aplicación

Coordinación entre dominio y repositorios

Infraestructura

Persistencia en base de datos (MySQL)

Repositorios

Presentación

Handlers HTTP

Templates HTML (frontend)

Tecnologías Utilizadas

Lenguaje: Go (Golang)

Framework HTTP: Gorilla Mux

Base de Datos: MySQL

Frontend: HTML + CSS

Arquitectura: MVC / Clean Architecture

Control de versiones: Git / GitHub

Descripción de las Pantallas del Sistema
Página de Inicio

Pantalla de bienvenida

Texto centrado

Enlaces a:

Usuarios

Libros

Buscar libros

Gestión de Usuarios

Crear usuarios

Listar usuarios existentes

Campos: nombre, correo, rol, estado

Gestión de Libros

Crear libros

Listar libros

Acceder al detalle de cada libro

Búsqueda de Libros

Búsqueda por:

Texto

Autor

Categoría

Resultados dinámicos

Detalle del Libro

Información completa del libro

Acceso desde el listado o búsqueda

Nota: La sección de estadísticas fue eliminada para simplificar la estabilidad del sistema.

Pruebas Realizadas
Pruebas Unitarias

Validación de entidades del dominio

Pruebas de servicios (casos de uso)

Verificación de reglas de negocio

Pruebas de Integración

Conexión con la base de datos

Flujo completo frontend–backend

Pruebas de QA

Navegación entre pantallas

Validación de formularios

Manejo de errores

Resultado:
El sistema respondió de forma estable, cumpliendo con los requisitos funcionales planteados.

Ejecución del Proyecto
1. Clonar el repositorio
git clone https://github.com/jfmg0509/sistema_libros_funcional_go.git

2. Configurar variables de entorno

Crear archivo .env:

DB_USER=root
DB_PASS=
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=libros_poo
APP_ADDR=:8081

3. Ejecutar la aplicación
go run main.go

4. Acceder desde el navegador
http://localhost:8081

Estructura del Proyecto
/internal
 ├── domain
 ├── usecase
 ├── infrastructure
 │    └── db
 └── transport
      └── http
/web
 └── templates
main.go
README.md

🔮 Visualización del Futuro

El sistema puede evolucionar hacia:

Autenticación y roles avanzados

API REST pública

Reportes y estadísticas

Migración a microservicios

Integración con sistemas académicos

Conclusiones

El desarrollo del Sistema de Gestión de Libros permitió aplicar de manera práctica los principios fundamentales de la Programación Orientada a Objetos, integrando backend, frontend y base de datos bajo una arquitectura clara y mantenible.

El proyecto cumple con los objetivos académicos propuestos y sienta las bases para futuras ampliaciones funcionales y tecnológicas.

Autor

Juan Francisco Morán Gortaire
Proyecto Académico – Programación Orientada a Objetos
