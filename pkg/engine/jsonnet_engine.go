package engine

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ghodss/yaml"
	jsonnet "github.com/strickyak/jsonnet_cgo"
	"k8s.io/helm/pkg/chartutil"
	chart "k8s.io/helm/pkg/proto/hapi/chart"
)

type JsonnetEngine struct {
}

func NewJsonnetEngine() *JsonnetEngine {
	return &JsonnetEngine{}
}

func (e *JsonnetEngine) Render(chrt *chart.Chart, values chartutil.Values) (map[string]string, error) {
	templateMap := allTemplates(chrt, values)

	renderResults := make(map[string]string)

	vm := jsonnet.Make()
	for name, renderable := range templateMap {
		println("name: " + name)
		if strings.Contains(name, "_") {
			continue
		}
		chartName := renderable.basePath + "/" + name

		jsonValues, err := json.Marshal(renderable.vals["Values"])
		if err != nil {
			return nil, err
		}
		jsonChart, err := json.Marshal(renderable.vals["Chart"])
		if err != nil {
			return nil, err
		}
		jsonRelease, err := json.Marshal(renderable.vals["Release"])
		if err != nil {
			return nil, err
		}
		vm.TlaCode("Values", string(jsonValues[:]))
		vm.TlaCode("Chart", string(jsonChart[:]))
		vm.TlaCode("Release", string(jsonRelease[:]))

		vm.ImportCallback(chartBasedImportCallback(templateMap))
		println(string(jsonValues[:]))
		result, err := vm.EvaluateSnippet(chartName, renderable.tpl)
		if err != nil {
			return nil, err
		}
		yamlBytes, err := yaml.JSONToYAML([]byte(result))
		if err != nil {
			return nil, err
		}
		result = string(yamlBytes)[:]

		println("Result: " + result)
		renderResults[chartName] = result
	}

	vm.Destroy()
	return renderResults, nil
}

func chartBasedImportCallback(templateMap map[string]renderable) func(base, rel string) (result string, path string, err error) {
	return func(base, rel string) (result string, path string, err error) {
		for templateName, renderable := range templateMap {
			if strings.HasSuffix(templateName, rel) {
				return renderable.tpl, rel, nil
			}
		}
		return "", "", fmt.Errorf("Import library %s not found", rel)
	}
}
