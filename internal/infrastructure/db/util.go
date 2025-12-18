package db // Infraestructura DB: funciones utilitarias específicas de MySQL

import "strings"

// boolToTiny convierte un bool de Go a un entero compatible con MySQL.
// MySQL suele representar booleanos como TINYINT(1):
// - true  -> 1
// - false -> 0
func boolToTiny(b bool) int {
	if b {
		return 1
	}
	return 0
}

// isMySQLDuplicate detecta errores de clave duplicada en MySQL.
//
// NOTA:
// El driver go-sql-driver/mysql expone el tipo *mysql.MySQLError,
// pero aquí evitamos depender directamente de ese tipo para:
//
// - reducir acoplamiento al driver
// - simplificar el código
// - mantener la infraestructura intercambiable
//
// Se usa una heurística segura basada en el mensaje de error:
// - "duplicate"
// - código "1062" (duplicate entry)
func isMySQLDuplicate(err error) bool {

	// Si no hay error, no hay duplicado
	if err == nil {
		return false
	}

	// Normaliza el mensaje a minúsculas para comparar
	msg := strings.ToLower(err.Error())

	// Detecta error de clave duplicada
	return strings.Contains(msg, "duplicate") || strings.Contains(msg, "1062")
}
