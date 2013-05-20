package tmpl

import (
	"github.com/murz/go-handlebars/handlebars"
	"net/http"
	"fmt"
	"strings"
)

type HandlebarsTemplateEngine struct {
	TemplateEngine
}

type HandlebarsTemplate struct {
	tmpl *handlebars.Template
	CompiledTemplate
}

func (te *HandlebarsTemplateEngine) Compile(t string) CompiledTemplate {
	tmpl, err := handlebars.ParseFile(t)
	if err != nil {
		fmt.Printf("err: %v", err)
		panic(err)
	}
	return HandlebarsTemplate{tmpl:tmpl}
}

func (tpl HandlebarsTemplate) Execute(ctx map[string]interface{}) (out *string, err error) {
	str := tpl.tmpl.Render(ctx)
	return &str, nil
}

func (tpl HandlebarsTemplate) ExecuteRW(w http.ResponseWriter, ctx map[string]interface{}) error {
	str := strings.Trim(tpl.tmpl.Render(ctx), "\n")
	w.Write([]byte(str))
	return nil
}