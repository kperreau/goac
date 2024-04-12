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
	sem := make(chan bool, l.Options.MaxConcurrency+1)
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
		printer.Printf("%s %s %s\n", color.BlueString(p.Name), color.YellowString("=>"), p.Path)
	}

	if p.CMDOptions.Target == TargetBuild && !p.CMDOptions.DryRun && isAffected {
		if _, err := p.buildProject(); err != nil {
			printer.Errorf("failed to build: %s\n", err.Error())
			return
		}
		if err := p.writeCache(p.CMDOptions.Target); err != nil {
			printer.Errorf("%v\n", err)
			return
		}
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

	if p.Cache.Targets[p.CMDOptions.Target] == nil || !p.Cache.Targets[p.CMDOptions.Target].isMetadataMatch(p.Metadata) {
		return true
	}

	if p.CMDOptions.BinaryCheck && !utils.FileExist(path.Join(p.Path, p.Name)) {
		return true
	}

	return false
}
