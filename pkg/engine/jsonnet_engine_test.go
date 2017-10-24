package engine

import (
	"testing"

	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

func TestJsonnetRender(t *testing.T) {

	c := &chart.Chart{
		Metadata: &chart.Metadata{
			Name:    "moby",
			Version: "1.2.3",
			Engine:  "jsonnet",
		},
		Templates: []*chart.Template{
			{Name: "templates/test.jsonnet", Data: []byte(`function(Values, Release, Chart) Values.name`)},
		},
		Values: &chart.Config{
			Raw: "name: Test",
		},
	}

	vals := &chart.Config{}

	jsonnetEngine := NewJsonnetEngine()

	v, err := chartutil.CoalesceValues(c, vals)
	if err != nil {
		t.Fatalf("Failed to coalesce values: %s", err)
	}
	out, err := jsonnetEngine.Render(c, v)
	if err != nil {
		t.Errorf("Failed to render templates: %s", err)
	}

	for name, tpl := range out {
		println(name)
		println(tpl)
	}

}

func TestJsonnetWithDependencies(t *testing.T) {
	e := NewJsonnetEngine()
	deptpl := `{Person: {Name: "test"}}`
	toptpl := `local lib = import "innerchart/templates/_lib.jsonnet"; {Person: lib.Person{Name:  "Test4"}}`
	ch := &chart.Chart{
		Metadata: &chart.Metadata{Name: "outerchart"},
		Templates: []*chart.Template{
			{Name: "templates/outer", Data: []byte(toptpl)},
		},
		Dependencies: []*chart.Chart{
			{
				Metadata: &chart.Metadata{Name: "innerchart"},
				Templates: []*chart.Template{
					{Name: "templates/_lib.jsonnet", Data: []byte(deptpl)},
				},
			},
		},
	}

	out, err := e.Render(ch, map[string]interface{}{})

	if err != nil {
		t.Fatalf("failed to render chart: %s", err)
	}

	if len(out) != 2 {
		t.Errorf("Expected 2, got %d", len(out))
	}

	expect := "Hello World"
	if out["outerchart/templates/outer"] != expect {
		t.Errorf("Expected %q, got %q", expect, out["outer"])
	}

}
