package domain

import (
    "fmt"
    "strings"
    "time"
)

type AccessEvent struct {
    id         uint64
    userID     uint64
    bookID     uint64
    accessType AccessType
    createdAt  time.Time
}

func NewAccessEvent(userID, bookID uint64, accessType AccessType) (*AccessEvent, error) {
    accessType = AccessType(strings.ToUpper(string(accessType)))
    if !accessType.IsValid() {
        return nil, fmt.Errorf("%w: %s", ErrInvalidAccess, accessType)
    }
    if userID == 0 || bookID == 0 {
        return nil, fmt.Errorf("%w: user_id and book_id are required", ErrValidation)
    }
    return &AccessEvent{
        userID:     userID,
        bookID:     bookID,
        accessType: accessType,
        createdAt:  time.Now(),
    }, nil
}

func HydrateAccessEvent(id, userID, bookID uint64, accessType AccessType, createdAt time.Time) (*AccessEvent, error) {
    e, err := NewAccessEvent(userID, bookID, accessType)
    if err != nil {
        return nil, err
    }
    e.id = id
    e.createdAt = createdAt
    return e, nil
}

func (e *AccessEvent) ID() uint64 { return e.id }
func (e *AccessEvent) UserID() uint64 { return e.userID }
func (e *AccessEvent) BookID() uint64 { return e.bookID }
func (e *AccessEvent) AccessType() AccessType { return e.accessType }
func (e *AccessEvent) CreatedAt() time.Time { return e.createdAt }
