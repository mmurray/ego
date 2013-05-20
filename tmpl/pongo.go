package tmpl

import (
	"github.com/flosch/pongo"
	"net/http"
)

type PongoTemplateEngine struct {
	TemplateEngine
}

type PongoTemplate struct {
	tmpl *pongo.Template
	CompiledTemplate
}

func (te *PongoTemplateEngine) Compile(t string) CompiledTemplate {
	return PongoTemplate{
		tmpl: pongo.Must(pongo.FromFile(t, nil)),
	}
}

func (tpl PongoTemplate) Execute(ctx map[string]interface{}) (out *string, err error) {
	return tpl.tmpl.Execute(&pongo.Context{})
}

func (tpl PongoTemplate) ExecuteRW(w http.ResponseWriter, ctx map[string]interface{}) error {
	pctx := pongo.Context(ctx)
	return tpl.tmpl.ExecuteRW(w, &pctx)
}