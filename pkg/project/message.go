package project

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/kperreau/goac/pkg/printer"
)

func (l *List) printAffected() {
	affectedCounter := l.countAffected()
	affected := color.HiBlueString("%d", affectedCounter)
	if affectedCounter == 0 {
		affected = color.HiBlackString("%d", affectedCounter)
	}
	printer.Printf("Affected: %s/%s\n", affected, color.HiBlueString("%d", len(l.Projects)))
}

func MessageFound(n int) string {
	return fmt.Sprintf("Found %s projects", color.YellowString("%d", n))
}

func MessageName(n string) string {
	return color.HiBlueString("%s", n)
}
