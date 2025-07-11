module github.com/ahboujelben/go-crossword/cli

go 1.24

require (
	github.com/ahboujelben/go-crossword/modules v0.0.0
	github.com/charmbracelet/lipgloss v1.0.0
	golang.org/x/term v0.29.0
)

require (
	github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
	github.com/charmbracelet/x/ansi v0.4.2 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/muesli/termenv v0.15.2 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	golang.org/x/sys v0.30.0 // indirect
)

replace github.com/ahboujelben/go-crossword/modules => ../modules
