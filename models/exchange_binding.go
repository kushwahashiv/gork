package models

// ExchangeBinding represents a link between exchange and queue.
type ExchangeBinding struct {
	Id         string // unique ID
	ExchangeId string // related exchange ID
	QueueId    string // related queue ID
	Rules      []*ExchangeBindingRule
}

// ExchangeBindingRule represents a rule that should be checked.
type ExchangeBindingRule struct {
	Key      string
	Operator ExchangeBindingRuleOperator
	Value    string
}
