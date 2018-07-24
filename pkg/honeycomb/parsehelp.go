package honeycomb

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

type Renderer interface {
	RenderTo(m KV)
}

func (i *item) RenderTo(m KV) {
	m[i.key] = i.value
}

func (g *group) RenderTo(m KV) {
	sub := make(KV)
	for _, k := range g.kids {
		k.RenderTo(sub)
	}
	m[g.key] = sub
}

func toIfaceSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	return v.([]interface{})
}
