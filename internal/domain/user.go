package domain

import (
    "fmt"
    "strings"
    "time"
)

// Struct + encapsulación (campos privados)
type User struct {
    id        uint64
    name      string
    email     string
    role      Role
    active    bool
    createdAt time.Time
    updatedAt time.Time
}

// Constructor
func NewUser(name, email string, role Role) (*User, error) {
    u := &User{active: true}
    if err := u.SetName(name); err != nil {
        return nil, err
    }
    if err := u.SetEmail(email); err != nil {
        return nil, err
    }
    if err := u.SetRole(role); err != nil {
        return nil, err
    }
    u.createdAt = time.Now()
    return u, nil
}

// Factory para cargar desde DB
func HydrateUser(id uint64, name, email string, role Role, active bool, createdAt, updatedAt time.Time) (*User, error) {
    u, err := NewUser(name, email, role)
    if err != nil {
        return nil, err
    }
    u.id = id
    u.active = active
    u.createdAt = createdAt
    u.updatedAt = updatedAt
    return u, nil
}

// Getters
func (u *User) ID() uint64        { return u.id }
func (u *User) Name() string      { return u.name }
func (u *User) Email() string     { return u.email }
func (u *User) Role() Role        { return u.role }
func (u *User) Active() bool      { return u.active }
func (u *User) CreatedAt() time.Time { return u.createdAt }
func (u *User) UpdatedAt() time.Time { return u.updatedAt }

// Setters (encapsulación)
func (u *User) SetName(name string) error {
    name = strings.TrimSpace(name)
    if len(name) < 2 {
        return fmt.Errorf("%w: name must have at least 2 characters", ErrValidation)
    }
    u.name = name
    u.updatedAt = time.Now()
    return nil
}

func (u *User) SetEmail(email string) error {
    email = strings.TrimSpace(strings.ToLower(email))
    if !strings.Contains(email, "@") || len(email) < 5 {
        return fmt.Errorf("%w: invalid email", ErrValidation)
    }
    u.email = email
    u.updatedAt = time.Now()
    return nil
}

func (u *User) SetRole(role Role) error {
    role = Role(strings.ToUpper(string(role)))
    if !role.IsValid() {
        return fmt.Errorf("%w: %s", ErrInvalidRole, role)
    }
    u.role = role
    u.updatedAt = time.Now()
    return nil
}

func (u *User) Deactivate() {
    u.active = false
    u.updatedAt = time.Now()
}

func (u *User) Activate() {
    u.active = true
    u.updatedAt = time.Now()
}
