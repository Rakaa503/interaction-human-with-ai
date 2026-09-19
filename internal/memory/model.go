package memory

type MemoryType string

const (
	MemoryName     MemoryType = "name"
	MemoryActivity MemoryType = "activity"
)

type Memory struct {
	Type       MemoryType
	Value      string
	Confidence float64
}
