package maputil_test

import (
	"encoding/json"
	"log"
	"reflect"
	"strings"
	"testing"

	. "github.com/dylannz-sailthru/maputil"
)

func TestStringValue(t *testing.T) {
	input := map[string]interface{}{
		"foo": "bar",
	}
	tr := NewMapTraverser(input)

	expected := "bar"
	result, ok := tr.Value("foo", "")
	if !ok {
		t.Error("expected ok to be true, was false")
	}
	if expected != result {
		t.Errorf("expected value %#v, got: %#v", expected, result)
	}
}

func TestMissingStringValue(t *testing.T) {
	input := map[string]interface{}{
		"foo": "bar",
	}
	tr := NewMapTraverser(input)

	result, ok := tr.Value("non_existent_key", "")
	if ok {
		t.Error("expected ok to be false")
	}
	if nil != result {
		t.Errorf("expected nil, got: %#v", result)
	}
}

func TestStringValueWithWrongType(t *testing.T) {
	input := map[string]interface{}{
		"foo": 123,
	}
	tr := NewMapTraverser(input)

	result, ok := tr.Value("foo", "")
	if ok {
		t.Error("expected ok to be false")
	}
	if nil != result {
		t.Errorf("expected nil, got: %#v", result)
	}
}

func TestValidChild(t *testing.T) {
	input := map[string]interface{}{
		"parent": map[string]interface{}{
			"foo": "bar",
		},
	}
	tr := NewMapTraverser(input)

	expected := NewMapTraverser(input["parent"].(map[string]interface{}))
	result := tr.Child("parent")
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("expected value %#v, got: %#v", expected, result)
	}
}

func TestInvalidChild(t *testing.T) {
	input := map[string]interface{}{
		"parent": map[string]interface{}{
			"foo": "bar",
		},
	}
	tr := NewMapTraverser(input)

	expected := NewMapTraverser(nil)
	result := tr.Child("invalid")
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("expected value %#v, got: %#v", expected, result)
	}
}

func TestMultipleInvalidChildren(t *testing.T) {
	input := map[string]interface{}{
		"parent": map[string]interface{}{
			"foo": "bar",
		},
	}
	tr := NewMapTraverser(input)

	expected := NewMapTraverser(nil)
	result := tr.Child("parent", "invalid", "child")
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("expected value %#v, got: %#v", expected, result)
	}
}

func TestNestedStringValue(t *testing.T) {
	input := map[string]interface{}{
		"parent": map[string]interface{}{
			"foo": "bar",
		},
	}
	tr := NewMapTraverser(input)

	expected := "bar"
	result, ok := tr.Child("parent").Value("foo", "")
	if !ok {
		t.Error("expected ok to be true, was false")
	}
	if expected != result {
		t.Errorf("expected value %#v, got: %#v", expected, result)
	}
}

func TestNestedIntValueWithWrongType(t *testing.T) {
	input := map[string]interface{}{
		"parent": map[string]interface{}{
			"foo": "bar",
		},
	}
	tr := NewMapTraverser(input)

	result, ok := tr.Child("parent").Value("foo", 123)
	if ok {
		t.Error("expected ok to be false")
	}
	if nil != result {
		t.Errorf("expected nil, got: %#v", result)
	}
}

func TestFindAllWithKey(t *testing.T) {
	input := map[string]interface{}{
		"parent": map[string]interface{}{
			"foo": "bar1",
			"another_child": map[string]interface{}{
				"foo": map[string]interface{}{
					"key": "value",
					"foo": "bar2",
				},
			},
			"some_other_thing": map[string]interface{}{"key": "value"},
		},
		"array": []interface{}{
			map[string]interface{}{
				"foo": "bar3",
			},
		},
	}
	tr := NewMapTraverser(input)

	expected := []MapTraverser{
		NewMapTraverser(input["parent"].(map[string]interface{})),
		NewMapTraverser(input["parent"].(map[string]interface{})["another_child"].(map[string]interface{})),
		NewMapTraverser(input["parent"].(map[string]interface{})["another_child"].(map[string]interface{})["foo"].(map[string]interface{})),
		NewMapTraverser(input["array"].([]interface{})[0].(map[string]interface{})),
	}
	result := tr.FindAllWithKey("foo")
	if len(expected) != len(result) {
		t.Errorf("expected MapTraversers with length %d, got: %d", len(expected), len(result))
	}

	for _, v := range expected {
		found := false
		for _, r := range result {
			if reflect.DeepEqual(r, v) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected result to contain MapTraverser with value: %#v, got: %#v", v, result)
		}
	}
}

func TestSetStringValue(t *testing.T) {
	input := map[string]interface{}{
		"foo": "world",
	}
	tr := NewMapTraverser(input)

	tr.Set("foo", "bar")

	expected := "bar"
	result, ok := tr.Value("foo", "")
	if !ok {
		t.Error("expected ok to be true, was false")
	}
	if expected != result {
		t.Errorf("expected value %#v, got: %#v", expected, result)
	}
}

func TestInvalidSetStringValue(t *testing.T) {
	input := map[string]interface{}{
		"foo": "world",
	}
	tr := NewMapTraverser(input)

	tr.Child("invalid").Set("foo", "bar")

	expected := "world"
	result, ok := tr.Value("foo", "")
	if !ok {
		t.Error("expected ok to be true, was false")
	}
	if expected != result {
		t.Errorf("expected value %#v, got: %#v", expected, result)
	}
}

func TestDeleteStringValue(t *testing.T) {
	input := map[string]interface{}{
		"foo": "world",
	}
	tr := NewMapTraverser(input)

	tr.Delete("foo")

	result, ok := tr.Value("foo", "")
	if ok {
		t.Error("expected ok to be false")
	}
	if nil != result {
		t.Errorf("expected nil, got: %#v", result)
	}
}

func TestInvalidDeleteStringValue(t *testing.T) {
	input := map[string]interface{}{
		"foo": "world",
	}
	tr := NewMapTraverser(input)

	tr.Child("invalid").Delete("foo")

	expected := "world"
	result, ok := tr.Value("foo", "")
	if !ok {
		t.Error("expected ok to be true, was false")
	}
	if expected != result {
		t.Errorf("expected value %#v, got: %#v", expected, result)
	}
}

func unmarshalJSON(j string) map[string]interface{} {
	o := map[string]interface{}{}
	dec := json.NewDecoder(strings.NewReader(j))
	dec.UseNumber()
	err := dec.Decode(&o)
	if err != nil {
		log.Fatal(err)
	}
	return o
}

func TestSetAll(t *testing.T) {
	result := unmarshalJSON(`{
		"foo": 111,
		"bar": 222.0,
		"nested": {
			"foo": 333,
			"bar": 444.5
		},
		"array": [
			555,
			666.6,
			{
				"foo": 777,
				"bar": 888.8
			}
		]
	}`)

	tr := NewMapTraverser(result)
	tr.SetAll(func(k Key, value interface{}) (interface{}, bool) {
		switch t := value.(type) {
		case json.Number:
			i, err := t.Int64()
			if err != nil {
				f, err := t.Float64()
				if err != nil {
					return t.String(), true
				}
				return f, true
			}
			return int(i), true
		default:
			return t, false
		}
	})

	expected := map[string]interface{}{
		"foo": 111,
		"bar": 222.0,
		"nested": map[string]interface{}{
			"foo": 333,
			"bar": 444.5,
		},
		"array": []interface{}{
			555,
			666.6,
			map[string]interface{}{
				"foo": 777,
				"bar": 888.8,
			},
		},
	}
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("expected value %#v, got: %#v", expected, result)
	}
}
