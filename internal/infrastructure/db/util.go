package db

import "strings"

func boolToTiny(b bool) int {
    if b { return 1 }
    return 0
}

func isMySQLDuplicate(err error) bool {
    // go-sql-driver/mysql expone *mysql.MySQLError, pero evitamos import directo aquí.
    // Heurística por mensaje/código textual.
    if err == nil { return false }
    msg := strings.ToLower(err.Error())
    return strings.Contains(msg, "duplicate") || strings.Contains(msg, "1062")
}
