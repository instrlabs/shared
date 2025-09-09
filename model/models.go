package model

import (
	"time"
)

type Product struct {
	ID          string    `json:"id,omitempty"`
	Key         string    `json:"key"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	Description string    `json:"description,omitempty"`
	Image       string    `json:"image,omitempty"`
	ProductType string    `json:"productType,omitempty"`
	UserID      string    `json:"userId,omitempty"`
	Active      bool      `json:"active"`
	IsFree      bool      `json:"isFree"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type InstructionStatus string

const (
	InstructionStatusPending    InstructionStatus = "PENDING"
	InstructionStatusProcessing InstructionStatus = "PROCESSING"
	InstructionStatusCompleted  InstructionStatus = "COMPLETED"
	InstructionStatusFailed     InstructionStatus = "FAILED"
)

type Instruction struct {
	ID        string            `json:"id"`
	UserID    string            `json:"user_id"`
	ProductID string            `json:"product_id"`
	Status    InstructionStatus `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}
