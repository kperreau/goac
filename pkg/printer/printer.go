package printer

import (
	"fmt"

	"github.com/fatih/color"
)

func Sprintf(format string, a ...any) string {
	return fmt.Sprintf(format, a...)
}

func Printf(format string, a ...any) {
	fmt.Printf(format, a...)
}

func Errorf(format string, a ...any) {
	red := color.New(color.FgRed).SprintFunc()
	fmt.Printf("%s", red(fmt.Sprintf(format, a...)))
}

func Warnf(format string, a ...any) {
	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Printf("%s", yellow(fmt.Sprintf(format, a...)))
}

func BoldGreen(s string) string {
	c := color.New(color.Bold).Add(color.FgGreen).SprintFunc()
	return c(s)
}
