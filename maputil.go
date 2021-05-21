package maputil

import "reflect"

type MapTraverser struct {
	m map[string]interface{}
}

type MapTraversers []MapTraverser

func NewMapTraverser(m map[string]interface{}) MapTraverser {
	return MapTraverser{m: m}
}

func (t MapTraverser) child(key string) MapTraverser {
	if t.m == nil {
		return MapTraverser{}
	}

	v, ok := t.Value(key, map[string]interface{}{})
	if !ok {
		return MapTraverser{}
	}

	return MapTraverser{m: v.(map[string]interface{})}
}

func (t MapTraverser) Child(keys ...string) MapTraverser {
	for _, k := range keys {
		t = t.child(k)
	}
	return t
}

func (t MapTraverser) FindAllWithKey(key string) MapTraversers {
	if t.m == nil {
		return nil
	}

	o := MapTraversers{}

	for k, v := range t.m {
		if k == key {
			o = append(o, t)
		}

		switch m := v.(type) {
		case map[string]interface{}:
			o = append(o, MapTraverser{m}.FindAllWithKey(key)...)
		}
	}

	return o
}

func (t MapTraverser) Value(key string, i interface{}) (interface{}, bool) {
	if t.m == nil {
		return nil, false
	}

	v, ok := t.m[key]
	if !ok {
		return nil, false
	}

	if reflect.TypeOf(i) == reflect.TypeOf(v) {
		return v, true
	}

	return nil, false
}

func (t MapTraverser) Set(key string, i interface{}) bool {
	if t.m == nil {
		return false
	}

	t.m[key] = i
	return true
}

func (t MapTraverser) Delete(key string) bool {
	if t.m == nil {
		return false
	}

	delete(t.m, key)
	return true
}
