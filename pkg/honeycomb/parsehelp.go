package honeycomb

//go:generate pigeon -o tm_msg_parser.go tm_msg_parser.peggo

// KV is just a Key-Value map
type KV map[string]interface{}

type item struct {
	key   string
	value interface{}
}

func newItem(k string, v interface{}) *item {
	return &item{key: k, value: v}
}

type group struct {
	key  string
	kids []Renderer
	hash string
}

func newGroup(k string, its []Renderer, h string) *group {
	return &group{key: k, kids: its, hash: h}
}

// Renderer is an interface we use to render items and groups recursively
// to a KV map.
type Renderer interface {
	RenderTo(m KV)
}

// RenderTo implements Renderer for item
func (i *item) RenderTo(m KV) {
	m[i.key] = i.value
}

// RenderTo implements renderer for group
func (g *group) RenderTo(m KV) {
	sub := make(KV)
	for _, k := range g.kids {
		k.RenderTo(sub)
	}
	m[g.key] = sub
}

// toIfaceSlice converts a generic interface to an array of interfaces
func toIfaceSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	return v.([]interface{})
}
