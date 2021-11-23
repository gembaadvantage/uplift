## [{{.Tag}}] - {{.Date}}
{{range $chg := .Changes }}
{{.AbbrevHash}} {{.Message}}
{{- end}}
