package project

import (
	"fmt"
	"path"

	"github.com/fatih/color"
	"github.com/kperreau/goac/pkg/printer"
	"github.com/kperreau/goac/pkg/utils"
	"golang.org/x/sync/errgroup"
)

type Target string

const (
	TargetNone       Target = "none"
	TargetBuild      Target = "build"
	TargetBuildImage Target = "build-image"
)

func (t Target) String() string { return string(t) }

func (l *List) Affected() error {
	l.printAffected()

	eg := errgroup.Group{}
	eg.SetLimit(l.Options.MaxConcurrency)
	for _, p := range l.Projects {
		eg.Go(func() error {
			return processAffected(p)
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}

func processAffected(p *Project) error {
	isAffected := p.isAffected()

	if isAffected && p.CMDOptions.DryRun {
		printer.Printf("%s %s %s\n", color.BlueString(p.Name), color.YellowString("=>"), p.CleanPath)
	}

	if p.CMDOptions.DryRun || !isAffected {
		return nil
	}

	if err := p.build(); err != nil {
		return fmt.Errorf("error building: %s", err.Error())
	}

	if err := p.writeCache(); err != nil {
		return err
	}

	return nil
}

func (l *List) countAffected() (n int) {
	for _, p := range l.Projects {
		if p.isAffected() {
			n++
		}
	}
	return n
}

func (p *Project) isAffected() bool {
	if p.CMDOptions.Force {
		return true
	}

	if p.Cache.Target[p.CMDOptions.Target] == nil || !p.Cache.Target[p.CMDOptions.Target].isMetadataMatch(p.Metadata) {
		return true
	}

	if p.CMDOptions.BinaryCheck && !utils.FileExist(path.Join(p.CleanPath, p.Name)) {
		return true
	}

	return false
}

func StringToTarget(s string) Target {
	switch s {
	case TargetBuild.String():
		return TargetBuild
	case TargetBuildImage.String():
		return TargetBuildImage
	}
	return TargetNone
}

func (l *List) printAffected() {
	affectedCounter := l.countAffected()
	affected := color.HiBlueString("%d", affectedCounter)
	if affectedCounter == 0 {
		affected = color.HiBlackString("%d", affectedCounter)
	}
	printer.Printf("Affected: %s/%s\n", affected, color.HiBlueString("%d", len(l.Projects)))
}
