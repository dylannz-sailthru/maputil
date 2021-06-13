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

type SetFn func(k string, value interface{}) (interface{}, bool)

func (t MapTraverser) SetAll(fn SetFn) int {
	if t.m == nil {
		return 0
	}

	return setAll(t.m, fn)
}

func setAll(i interface{}, fn SetFn) int {
	o := 0
	switch t := i.(type) {
	case map[string]interface{}:
		for k, v := range t {
			res, changed := fn(k, v)
			if changed {
				t[k] = res
				o++
			}

			o += setAll(v, fn)
		}
	case []interface{}:
		for k := range t {
			switch t[k].(type) {
			case map[string]interface{}:
				o += setAll(t[k], fn)
			default:
				res, changed := fn("", t[k])
				if changed {
					t[k] = res
					o++
				}
			}
		}
	}
	return o
}

type FindFn func(k string, value interface{}) bool

func (t MapTraverser) FindAll(fn FindFn) MapTraversers {
	if t.m == nil {
		return nil
	}

	return findAll(t.m, fn)
}

func findAll(i interface{}, fn FindFn) MapTraversers {
	o := MapTraversers{}
	switch t := i.(type) {
	case map[string]interface{}:
		for k, v := range t {
			if fn(k, v) {
				o = append(o, MapTraverser{t})
			}

			o = append(o, findAll(v, fn)...)
		}
	case []interface{}:
		for _, v := range t {
			o = append(o, findAll(v, fn)...)
		}
	}
	return o
}

func (t MapTraverser) FindAllWithKey(key string) MapTraversers {
	return t.FindAll(func(k string, value interface{}) bool {
		return k == key
	})
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
