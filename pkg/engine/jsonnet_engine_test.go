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
			{Name: "templates/test.jsonnet", Data: []byte(`function(Values, Chart)Chart.name`)},
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

	expected := "moby"

	if expected != out["moby/templates/test.jsonnet"] {
		t.Errorf("Expected '%s', got %s", expected, out["moby/templates/test.jsonnet"])
	}
}
