package tmpl

import(
	"github.com/hoisie/mustache"
	"log"
	"strings"
	"io/ioutil"
	"path"
)

type TemplateManager struct {
	pkgName string
	templates map[string]*mustache.Template
}

func (tm *TemplateManager) Parse(filename string) {
	t, err := mustache.ParseFile(filename)
	if err != nil {
		panic(err)
	}
	key := filename
	if tm.pkgName != "" {
		key = strings.Replace(key, tm.pkgName + "/app/views/", "", -1)
		log.Printf("KEY: %v", key)
	}
	log.Printf("parsing: %s\n", key)
	tm.templates[key] = t
}

func (tm *TemplateManager) ParseDir(dirname string) {
	dirlist, err := ioutil.ReadDir(tm.pkgName + dirname)
	if err != nil {
		log.Fatalf("Error reading %s: %s\n", dirname, err)
	}
	for _, f := range dirlist {
		filename := path.Join(tm.pkgName, dirname, f.Name())
		if f.IsDir() {
			tm.ParseDir(path.Join(dirname, f.Name()))
		} else {
			tm.Parse(filename)
		}
	}
}

func (tm *TemplateManager) Render(key string, context interface{}) string {
	if (key[0:1] == "/") {
		key = key[1:len(key)]
	}
	return tm.templates[key].Render(context)
}

func (tm *TemplateManager) RenderInLayout(layoutKey string, key string, context interface{}) string {
	log.Printf("layoutKey: %v", "layouts/"+layoutKey)
	return tm.templates[key].RenderInLayout(tm.templates["layouts/"+layoutKey], context)
}