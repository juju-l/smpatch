package main

import (
	"fmt"
	"github.com/your-org/yamlpatch"
	"gopkg.in/yaml.v3"
)

func main() {
	// 加载 YAML
	var obj map[string]interface{}
	yaml.Unmarshal([]byte(`
spec:
  components:
    - name: web
      image: nginx:1.20
  tags:
    - dev
    - legacy
`), &obj)

	// 构造 Patch
	patches := []yamlpatch.PatchOp{
		{
			Ope:     "merge",
			PathKey: "/spec/components",
			ByKey:   "name",
			Value: []interface{}{
				map[string]interface{}{
					"name":  "web",
					"image": "nginx:1.25",
				},
			},
		},
		{
			Ope:     "merge",
			PathKey: "/spec/tags",
			ItemOps: "replace",
			Old:     "legacy",
			Value:   "production",
		},
	}

	// dry-run
	after, _ := yamlpatch.DryRun(obj, patches)
	diffs := yamlpatch.DiffMaps(obj, after)
	for _, d := range diffs {
		fmt.Printf("%s: -%v +%v\n", d.Path, d.Old, d.New)
	}
}