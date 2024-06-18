Title: {{.Title}}
{{range .Story}}
{{.}}
{{end}}

Options:
{{range $index, $element := .Options}}
[{{$index}}] ({{$element.Arc}}) {{$element.Text}}
{{end}}