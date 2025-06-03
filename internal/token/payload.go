package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Payload struct {
	SessionID uuid.UUID `json:"session_id"`
	ID  primitive.ObjectID    `json:"id"`
	ExpiredAt time.Time `json:"expired_at"`
	IssuedAt  time.Time `json:"issued_at"`
}

func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return errors.New("token has expired")
	}

	if p.SessionID == uuid.Nil {
		return errors.New("invalid session ID")
	}

	if p.ID.Hex() == "" {
		return errors.New("email cannot be empty")
	}

	return nil
}
func NewPayload(objectId primitive.ObjectID, duration time.Duration) (*Payload, error) {
	sessionID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		SessionID: sessionID,
		ID:  objectId,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	if err := payload.Valid(); err != nil {
		return nil, err
	}

	return payload, nil
}
