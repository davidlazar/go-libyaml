package yaml

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func compareToJSON(t *testing.T, yamlPath string) {
	yaml, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		t.Fatal(err)
	}

	vy, err := Load(yaml)
	if err != nil {
		t.Fatal(err)
	}

	jsonPath := strings.TrimSuffix(yamlPath, ".yaml") + ".json"
	jsonFile, err := os.Open(jsonPath)
	if err != nil {
		t.Fatal(err)
	}

	decoder := json.NewDecoder(jsonFile)
	var vj interface{}
	decoder.Decode(&vj)

	if !reflect.DeepEqual(vj, vy) {
		t.Fatalf("%s:\n\texpect %#v\n\tgot    %#v", yamlPath, vj, vy)
	}
}

func TestLoad(t *testing.T) {
	paths, err := filepath.Glob("tests/*.yaml")
	if err != nil {
		t.Fatal(err)
	}

	if len(paths) < 5 {
		t.Fatal("wrong tests/ directory?")
	}

	for _, path := range paths {
		compareToJSON(t, path)
	}
}

func TestLoadSimple(t *testing.T) {
	o := func(b []byte, expected interface{}) {
		actual, err := Load(b)
		if err != nil {
			t.Fatalf("Load(%q): %v", b, err)
		}
		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("Load(%q):\n\texpect %#v\n\tgot    %#v", b, expected, actual)
		}
	}

	o(nil, nil)
	o([]byte("a"), "a")
	o([]byte("[a,b,c]"), []interface{}{"a", "b", "c"})
	o([]byte("a: x\nb: y\n"), map[string]interface{}{"a": "x", "b": "y"})
	o([]byte("a: [x, b: y]"), map[string]interface{}{"a": []interface{}{"x", map[string]interface{}{"b": "y"}}})
	o([]byte("hello: 世界\n😂: smile"), map[string]interface{}{"hello": "世界", "😂": "smile"})
}
