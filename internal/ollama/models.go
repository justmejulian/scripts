package ollama

type Model struct{ name string }

func (m Model) String() string { return m.name }

var (
	ModelQwen3_8B = Model{name: "qwen3:8b"}
	ModelQwen3_4B = Model{name: "qwen3:4b"}
)
