package util

import (
	"bytes"
	"html/template"
)

func TemplateString(templateString string, parameters Parameters) (string, error) {
	templatevar, err := template.New("template").Parse(templateString)
	if err != nil {
		return "", err
	}
	var byteBuffer bytes.Buffer
	err = templatevar.Execute(&byteBuffer, parameters)
	if err != nil {
		return "", err
	}
	return byteBuffer.String(), nil
}
