package models

import "time"

const (
	ExchangeTypeDirect ExchangeType = iota
	ExchangeTypeFanout
	ExchangeTypeTopic
)

// ExchangeType represents a
type ExchangeType uint8

type ExchangeBindingRuleOperator uint8

// Exchange represents an input that accepts tasks from the publishers and routes them to the appropriate queues.
type Exchange struct {
	Id        string // unique ID
	Name      string // unique name
	Type      ExchangeType
	CreatedAt time.Time
}