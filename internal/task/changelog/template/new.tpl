# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

## [{{.Tag}}] - {{.Date}}
{{range $chg := .Changes }}
{{.AbbrevHash}} {{.Message}}
{{- end}}
