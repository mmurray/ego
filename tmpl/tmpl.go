package tmpl

import (
	"html/template"
	"io/ioutil"
	"io"
	"strings"
	"log"
	"path"
)

var pkgName string
var templates = make(map[string]*template.Template)
var helpers = make([]*Helper, 0)
var partials = make([]*Partial, 0)

// TODO: parse partials into a structure first,
// then use them in Parse()

func Parse(filename string) {
	basePath := pkgName + "/app/views/"
	dirlist, err := ioutil.ReadDir(basePath + "layouts")
	if err != nil {
		log.Fatalf("Error reading /app/views/layouts")
	}
	for _, f := range dirlist {
		filenames := []string {
			// basePath + "layouts/" + f.Name(),
			// TODO: add partials
			filename,
		}
		key := f.Name() + "||" + filename
		if pkgName != "" {
			key = strings.Replace(key, basePath, "", -1)
			log.Printf("KEY: %v", key)
		}
		t := template.New(key)
		for _, helper := range helpers {
			t.Funcs(template.FuncMap{
				helper.Name: helper.Execute,
	        })
		}
		for _, partial := range partials {
			partialTemplate := t.New(partial.Name)
			partialTemplate.Parse(partial.TemplateString)
		}
		file, _ := ioutil.ReadFile(basePath + "layouts/" + f.Name())
		t.Parse(string(file))
		t.ParseFiles(filenames...)
		if (err != nil) {
			log.Panic(err)
		}
		templates[key] = t
	}
}

func ParseDir(dirname string) {
	ParsePartials(dirname)
	dirlist, err := ioutil.ReadDir(pkgName + dirname)
	if err != nil {
		log.Fatalf("Error reading %s: %s\n", dirname, err)
	}
	for _, f := range dirlist {
		filename := path.Join(pkgName, dirname, f.Name())
		if f.Name() == "layouts" || f.Name()[0:1] == "_" {
			continue
		}
		if f.IsDir() {
			ParseDir(path.Join(dirname, f.Name()))
		} else {
			Parse(filename)
		}
	}
}

func ParsePartials(dirname string) {
	basePath := pkgName + dirname
	dirlist, err := ioutil.ReadDir(basePath)
	if err != nil {
		log.Fatalf("Error reading %s: %s\n", dirname, err)
	}
	for _, f := range dirlist {
		filename := path.Join(pkgName, dirname, f.Name())
		if f.Name()[0:1] == "_" {
			key := f.Name()[1:len(f.Name())]
			key = strings.Replace(key, ".html", "", -1)
			// t := template.New(key)
			file, _ := ioutil.ReadFile(filename)
			// t.Parse(string(file))
			partial := &Partial{
				Name: key,
				TemplateString: string(file),
			}
			log.Printf("adding partial: %v", key)
			partials = append(partials, partial)
		}
	}
}

func Render(wr io.Writer, key string, layout string, context interface{}) {
	log.Printf("looking up: %v", layout+"||"+key)
	t := templates[layout + "||" + key]
	log.Printf("found: %v", t)
	err := t.Execute(wr, context)
	if (err != nil) {
		panic(err)
	}
	// TODO: do something with err instead of _
}

func SetPackageName(name string) {
	pkgName = name
}

func RegisterHelper(helper *Helper) *Helper {
	log.Printf("## REGISTERED helper: %v", helper.Name)
	helpers = append(helpers, helper)
	return helper
}