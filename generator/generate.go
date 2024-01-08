package generator

import (
	"github.com/ml444/swaggerapi/generator/descriptor"
	"github.com/ml444/swaggerapi/generator/genopenapi"
	"google.golang.org/protobuf/types/pluginpb"
)

type Generator struct {
	registry *descriptor.Registry
}

func NewGenerator(options ...Option) *Generator {
	gen := &Generator{
		registry: descriptor.NewRegistry(),
	}
	gen.registry.SetUseJSONNamesForFields(true)
	gen.registry.SetMergeFileName("apidocs")
	gen.registry.SetDisableDefaultErrors(true)
	for _, o := range options {
		o(gen)
	}
	return gen
}

func (g *Generator) Gen(req *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error) {
	reg := g.registry
	if reg == nil {
		reg = NewGenerator().registry
	}
	if err := reg.SetRepeatedPathParamSeparator("csv"); err != nil {
		return nil, err
	}

	gen := genopenapi.New(reg, genopenapi.FormatJSON)

	if err := genopenapi.AddErrorDefs(reg); err != nil {
		return nil, err
	}

	if err := reg.Load(req); err != nil {
		return nil, err
	}
	var targets []*descriptor.File
	for _, target := range req.FileToGenerate {
		f, err := reg.LookupFile(target)
		if err != nil {
			return nil, err
		}
		targets = append(targets, f)
	}

	out, err := gen.Generate(targets)
	if err != nil {
		return nil, err
	}
	return emitFiles(out), nil
}
func emitFiles(out []*descriptor.ResponseFile) *pluginpb.CodeGeneratorResponse {
	files := make([]*pluginpb.CodeGeneratorResponse_File, len(out))
	for idx, item := range out {
		files[idx] = item.CodeGeneratorResponse_File
	}
	resp := &pluginpb.CodeGeneratorResponse{File: files}
	sf := uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
	resp.SupportedFeatures = &sf
	return resp
}
