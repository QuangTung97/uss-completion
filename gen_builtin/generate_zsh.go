package gen_builtin

import (
	_ "embed"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed generate_zsh.tmpl
var generateZshTemplate string

func GenerateZshBuiltinCompletionScript(
	output string,
	cmdList ...Command,
) {
	tmpl := template.Must(template.New("completion").Parse(generateZshTemplate))

	if err := os.MkdirAll(filepath.Dir(output), 0755); err != nil {
		panic(err)
	}

	file, err := os.Create(output)
	if err != nil {
		panic(err)
	}

	data := templateData{}
	for _, cmd := range cmdList {
		configOption := ""
		setOption := "compopt -o nospace"
		if cmd.NoSpace {
			configOption = " -o nospace"
			setOption = ""
		}

		data.Commands = append(data.Commands, commandData{
			Name:         cmd.Name,
			NameUpper:    strings.ToUpper(cmd.Name),
			SetOption:    setOption,
			ConfigOption: configOption,
		})
	}

	if err := tmpl.Execute(file, data); err != nil {
		panic(err)
	}

	if err := file.Close(); err != nil {
		panic(err)
	}
}
