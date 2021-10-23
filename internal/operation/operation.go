package operation

type Type uint8

const (
	Get Type = iota
	Set
	Delete
	Keys
)

type Operation struct {
	Type  Type
	Key   string
	Value interface{}
}
