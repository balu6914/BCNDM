package auth

// Action represents any action.
type Action int64

// Predefined list of actions.
const (
	Any Action = iota
	Create
	CreateBulk
	Read
	Update
	Delete
	Write
	List
	Buy
	Withdraw
	Sign
)

const id = "id"

// Rule represents Policy rule.
type Rule struct {
	Action    Action
	Type      string
	Condition Condition
}

// Condition represents Rule condition.
type Condition interface {
	Validate(subject, object Resource) bool
}

// SimpleCondition represents a simple key:value condition.
type SimpleCondition struct {
	Key   string
	Value string
}

// Validate implements Condition interface.
// If there is no value, value is evaluated
// against subject with the same attribute.
func (sc SimpleCondition) Validate(subject, object Resource) bool {
	attrs := subject.Attributes()
	if sc.Value == "" {
		return object.Attributes()[sc.Key] == subject.Attributes()[id]
	}

	for k, v := range attrs {
		if k == sc.Key && v == sc.Value {
			return true
		}
	}
	return false
}
