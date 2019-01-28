package util

import (
	"bytes"
	"html/template"
)

func RenderTemplate(tplString string, object interface{}) (string, error) {
	t := template.New(tplString)
	t, err := t.Parse(tplString)
	if err != nil {
		return "", nil
	}
	var tpl bytes.Buffer
	if err = t.Execute(&tpl, object); err != nil {
		return "", err
	}
	return tpl.String(), err
}
