package engine

import (
	"encoding/json"
	"strings"

	"github.com/ghodss/yaml"
	jsonnet "github.com/google/go-jsonnet"
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

	vm := jsonnet.MakeVM()
	for name, renderable := range templateMap {
		println("name: " + name)
		if strings.Contains(name, "_") {
			continue
		}

		chartName := name

		b, _ := json.MarshalIndent(renderable.vals.AsMap(), "", "  ")
		println(string(b))

		jsonValues, err := json.Marshal(renderable.vals)
		if err != nil {
			return nil, err
		}
		jsonChart, err := json.Marshal(chrt.Metadata)
		if err != nil {
			return nil, err
		}
		println(string(jsonChart))
		/*jsonRelease, err := json.Marshal(renderable.vals["Release"])
		if err != nil {
			return nil, err
		}*/
		vm.TLACode("Values", string(jsonValues[:]))
		vm.TLACode("Chart", string(jsonChart[:]))
		//vm.TLACode("Release", string(jsonRelease[:]))

		//vm.Importer(chartBasedImportCallback(templateMap))
		println(string(jsonValues[:]))
		result, err := vm.EvaluateSnippet(chartName, renderable.tpl)
		if err != nil {
			return nil, err
		}
		yamlBytes, err := yaml.JSONToYAML([]byte(result))
		if err != nil {
			return nil, err
		}
		result = strings.TrimSuffix(string(yamlBytes)[:], "\n")
		renderResults[chartName] = result
	}
	return renderResults, nil
}

/*func chartBasedImportCallback(templateMap map[string]renderable) jsonnet.Importer {
	return &FileImporter{
		templateMap: templateMap,
	}
}

type FileImporter struct {
	templateMap map[string]renderable
}

func (i *FileImporter) Import(codeDir string, importedPath string) jsonnet.ImportedData {

	for templateName, renderable := range i.templateMap {
		if strings.HasSuffix(templateName, importedPath) {
			return jsonnet.ImportedData{
				content:   renderable.tpl,
				foundHere: importedPath,
			}
		}
	}

	return jsonnet.ImportedData{err: fmt.Errorf("Import not available %v", importedPath)}

}
*/
