package main

import (
	"log"
	"strings"
	"text/template"
)

var (
	tFactories *template.Template
	fns        = template.FuncMap{
		"last": func(x int, y []any) bool {
			return x == len(y)-1
		},
		"nlast": func(x int, y []Field) bool {
			return x != len(y)-1
		},
		"toupper": func(str string) string {
			return strings.ToUpper(str[:1]) + str[1:]
		},
	}
)

func init() {
	var err error
	tFactories, err = template.New("factories").Funcs(fns).Parse(FactoriesTemplateFile)
	if err != nil {
		log.Fatal(err)
	}
}

const FactoriesTemplateFile = `// Code generated by proto-factories. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.6
// - proto-factories    v0.1.0

package {{.GoPackageName}}

{{range $message:= .Messages}}{{range $key, $field := $message.Oneof.Fields}}{{if ne $field.Name "err" }}func Fac{{$message.Name}}{{toupper $field.Name}}({{$field.Name}} *{{$message.Name}}_{{$field.Type}})*{{$message.Name}}{
	return &{{$message.Name}}{
		{{toupper $message.Oneof.Name}}: &{{$message.Name}}_{{toupper $field.Type}}_{
			{{$field.Type}}:{{$field.Name}},
		},
	}
}
{{end}}
{{else}}
{{end}}
{{range $embeddedMessage := $message.EmbeddedMessages}}{{ if eq $embeddedMessage.Name "Error" }}{{range $key, $field := $embeddedMessage.Oneof.Fields}}func Fac{{$message.Name}}{{toupper $field.Name}}()*{{$message.Name}}{
	return &{{$message.Name}}{
		{{toupper $embeddedMessage.MessageOneOfName}}: &{{$message.Name}}_{{toupper $embeddedMessage.Name}}_{
			{{toupper $embeddedMessage.Name}}: &{{$message.Name}}_{{toupper $embeddedMessage.Name}}{
				Error: &{{$message.Name}}_Error_{{toupper $field.Name}}{
					{{toupper $field.Name}}: &{{toupper $field.Name}}{},
				},
			},
		},
	}
}
{{end}}
{{end}}
{{end}}
{{ else }}
{{ end }}
`
