package utils

type LimitStackMap struct {
	Order []string
	Map   map[string]interface{}
	Size  int
	Limit int
}

func NewLimitStackMap(limit int) *LimitStackMap {
	return &LimitStackMap{
		Order: make([]string, 0, limit),
		Map:   make(map[string]interface{}),
		Size:  0,
		Limit: limit,
	}
}

func (m *LimitStackMap) Peek() interface{} {
	if m.Size > 0 {
		return m.Map[m.Order[m.Size-1]]
	}
	return nil
}

func (m *LimitStackMap) Find(k string) interface{} {
	return m.Map[k]
}

func (m *LimitStackMap) Push(k string, i interface{}) {
	if m.Size >= m.Limit {
		ancient := m.Order[0]
		m.Order = m.Order[1:]

		delete(m.Map, ancient)
		m.Size--
	}

	m.Order = append(m.Order, k)
	m.Map[k] = i
	m.Size++
}
