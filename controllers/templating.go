package controllers

import (
	"bytes"
	"html/template"
)

func templateString(templateparam string, parameters Parameters) (string, error) {
	templatevar, err := template.New("template").Parse(templateparam)
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
