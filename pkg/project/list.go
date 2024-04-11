package project

import (
	"github.com/fatih/color"
	"github.com/kperreau/goac/pkg/printer"
)

const configFileName = ".goacproject.yaml"

func (l *List) List() {
	printer.Printf("Found %s projects\n", color.YellowString("%d", len(l.Projects)))
	for _, project := range l.Projects {
		printer.Printf("%s %s %s\n", color.BlueString(project.Name), color.YellowString("=>"), project.Path)
	}
}
