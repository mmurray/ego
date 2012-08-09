package tmpl

import (
	"github.com/hoisie/mustache"
)

var tm = &TemplateManager {
	templates: make(map[string]*mustache.Template),
}

func Render(key string, context interface{}) string {
	return tm.Render(key, context)
}

func RenderInLayout(layoutKey string, key string, context interface{}) string {
	return tm.RenderInLayout(layoutKey, key, context)
}

func Parse(filename string) {
	tm.Parse(filename)
}

func ParseDir(dirname string) {
	tm.ParseDir(dirname)
}

func SetPackageName(pkgName string) {
	tm.pkgName = pkgName
}