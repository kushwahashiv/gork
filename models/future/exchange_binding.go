package models

import "time"

// ExchangeBindingRuleOperator represents a logical operator that is used in exchange binding rules.
type ExchangeBindingRuleOperator uint8

// ExchangeBinding represents a link between exchange and queue.
type ExchangeBinding struct {
	Id         string                 // unique ID
	ExchangeId string                 // related exchange ID
	QueueId    string                 // related queue ID
	Rules      []*ExchangeBindingRule // routing rules
	CreatedAt  time.Time              // creation time
}

// ExchangeBindingRule represents a rule that should be checked.
type ExchangeBindingRule struct {
	Key      string                      // key of the header param
	Operator ExchangeBindingRuleOperator // comparison operator
	Value    string                      // comparison value
}
