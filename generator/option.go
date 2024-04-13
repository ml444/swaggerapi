package generator

type Option func(gen *Generator)

// DescriptorFilePath: @usage: generator.DescriptorFilePath("myapp.proto", "pkg/mygroup/myapp.descriptor.pb")
func DescriptorFilePath(protoName, descriptorFile string) Option {
	return func(gen *Generator) {
		gen.registry.AddDescriptorMap(protoName, descriptorFile)
	}
}

func PreserveRPCOrder() Option {
	return func(gen *Generator) {
		gen.registry.SetPreserveRPCOrder(true)
	}
}

func AllowDeleteBody(b bool) Option {
	return func(gen *Generator) {
		gen.registry.SetAllowDeleteBody(b)
	}
}

// UseJSONNamesForFields if disabled, the original proto name will be used for generating OpenAPI definitions
func UseJSONNamesForFields(b bool) Option {
	return func(gen *Generator) {
		gen.registry.SetUseJSONNamesForFields(b)
	}
}

// RecursiveDepth maximum recursion count allowed for a field type
func RecursiveDepth(depth int) Option {
	return func(gen *Generator) {
		gen.registry.SetRecursiveDepth(depth)
	}
}

// EnumsAsInts whether to render enum values as integers, as opposed to string values
func EnumsAsInts(b bool) Option {
	return func(gen *Generator) {
		gen.registry.SetEnumsAsInts(b)
	}
}

// MergeFileName target OpenAPI file name prefix after merge
func MergeFileName(name string) Option {
	return func(gen *Generator) {
		gen.registry.SetMergeFileName(name)
	}
}

// DisableDefaultErrors if set, disables generation of default errors. This is useful if you have defined custom error handling
func DisableDefaultErrors(b bool) Option {
	return func(gen *Generator) {
		gen.registry.SetDisableDefaultErrors(b)
	}
}
