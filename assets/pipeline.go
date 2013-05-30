package assets

import (
)

type AssetPipeline struct {
  inputDirs []string
  outputDir string
  processor *DirectiveProcessor
}

func NewAssetPipeline() *AssetPipeline {
  return &AssetPipeline{
    inputDirs: make([]string, 0),
    outputDir: "./public",
    processor: NewDirectiveProcessor(),
  }
}

func (p *AssetPipeline) CompileAll() {

}

func (p *AssetPipeline) CompileAsset(name string) {
  directives := ProcessFile(name)
  for _, directive := range directives {
    switch directive.Name {
      case "require":
        
      case "require_self":
        
    }
  }
}

func (p *AssetPipeline) AddInputDir(dir string) {
  p.inputDirs = append(p.inputDirs, dir)
}

func (p *AssetPipeline) SetOutputDir(dir string) {
  p.outputDir = dir
}
