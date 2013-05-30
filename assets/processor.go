package assets

import (
  "io/ioutil"
//  "strings"
  "fmt"
  "regexp"
)

type Directive struct {
  Name string
  Path string
}

type DirectiveProcessor struct {
  directiveRegexp *regexp.Regexp
}

func NewDirectiveProcessor() *DirectiveProcessor {
  r, _ := regexp.Compile(`(//|\*)= ([\w_]+) ?([A-Za-z\-\.]+)?`)
  return &DirectiveProcessor{
    directiveRegexp: r,
  }
}

func (dp *DirectiveProcessor) ProcessFile(filename string) ([]Directive, error) {
  content, err := ioutil.ReadFile(filename)
  if err != nil {
    return nil, err
  }
  return dp.ProcessString(string(content)), nil
}

func (dp *DirectiveProcessor) ProcessString(content string) ([]Directive, error) {
  // Parse directives
  matches := dp.directiveRegexp.FindAllStringSubmatch(content, -1)
  directives := make([]Directive, 0)
  for _, match := range matches {
    directives = append(directives, Directive{
      Name: match[2],
      Path: match[3],
    })
  }
  return directives, nil
}


