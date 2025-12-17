package mock

type MockUUIDGenerator struct {
	UUIDs []string
	index int
}

func NewMockUUIDGenerator(uuids ...string) *MockUUIDGenerator {
	return &MockUUIDGenerator{UUIDs: uuids}
}

func (g *MockUUIDGenerator) Generate() string {
	if g.index >= len(g.UUIDs) {
		return "mock-uuid"
	}
	id := g.UUIDs[g.index]
	g.index++
	return id
}
