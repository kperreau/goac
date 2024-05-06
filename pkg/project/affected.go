package project

import (
	"path"
	"sync"

	"github.com/fatih/color"
	"github.com/kperreau/goac/pkg/printer"
	"github.com/kperreau/goac/pkg/utils"
)

type Target string

const (
	TargetNone       Target = "none"
	TargetBuild      Target = "build"
	TargetBuildImage Target = "build-image"
)

func (t Target) String() string { return string(t) }

type processAffectedOptions struct {
	wg  *sync.WaitGroup
	sem chan bool
}

func (l *List) Affected() error {
	l.printAffected()

	// init process options
	sem := make(chan bool, l.Options.MaxConcurrency)
	wg := sync.WaitGroup{}
	pOpts := &processAffectedOptions{
		wg:  &wg,
		sem: sem,
	}

	for _, p := range l.Projects {
		sem <- true // acquire
		wg.Add(1)
		go processAffected(p, pOpts)
	}

	wg.Wait()

	return nil
}

func processAffected(p *Project, opts *processAffectedOptions) {
	defer opts.wg.Done()
	defer func() {
		<-opts.sem // release
	}()

	isAffected := p.isAffected()

	if isAffected && p.CMDOptions.DryRun {
		printer.Printf("%s %s %s\n", color.BlueString(p.Name), color.YellowString("=>"), p.CleanPath)
	}

	if p.CMDOptions.DryRun || !isAffected {
		return
	}

	if err := p.build(); err != nil {
		printer.Errorf("error building: %s\n", err.Error())
		return
	}

	if err := p.writeCache(); err != nil {
		printer.Errorf("%v\n", err)
		return
	}
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
