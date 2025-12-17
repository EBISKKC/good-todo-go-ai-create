package pkg

import (
	"github.com/google/uuid"
)

type IUUIDGenerator interface {
	Generate() string
}

type UUIDGenerator struct{}

func NewUUIDGenerator() IUUIDGenerator {
	return &UUIDGenerator{}
}

func (g *UUIDGenerator) Generate() string {
	return uuid.New().String()
}
