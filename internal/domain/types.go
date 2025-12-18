package domain

import (
    "strings"
)

type Role string

const (
    RoleAdmin  Role = "ADMIN"
    RoleReader Role = "READER"
    RoleConsultor Role = "CONSULTOR"
)

// Array (requerimiento): set fijo de roles permitidos
var AllowedRoles = [3]Role{RoleAdmin, RoleReader, RoleConsultor}

func (r Role) IsValid() bool {
    rr := strings.ToUpper(string(r))
    for _, allowed := range AllowedRoles {
        if rr == string(allowed) {
            return true
        }
    }
    return false
}

type AccessType string

const (
    AccessApertura AccessType = "APERTURA"
    AccessLectura  AccessType = "LECTURA"
    AccessDescarga AccessType = "DESCARGA"
)

var AllowedAccessTypes = [3]AccessType{AccessApertura, AccessLectura, AccessDescarga}

func (a AccessType) IsValid() bool {
    aa := strings.ToUpper(string(a))
    for _, allowed := range AllowedAccessTypes {
        if aa == string(allowed) {
            return true
        }
    }
    return false
}
