package assets

import (
  "bytes"
)

type Asset struct {
  content bytes.Buffer
}

func (a *Asset) Append(content string) {
  a.content.WriteString(content)
}

func (a *Asset) Content() string {
  return a.content.String()
}

