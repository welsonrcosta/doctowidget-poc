package assets

import (
	"embed"
	_ "embed"
)

//go:embed style.css
var Style string

//go:embed *.svg
var Images embed.FS

//go:embed ui.glade
var Ui string
