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

type processAffectedOptions struct {
	*Options
	wg  *sync.WaitGroup
	sem chan bool
}

func (l *List) Affected() error {
	l.printAffected()

	// init process options
	sem := make(chan bool, l.Options.MaxConcurrency+1)
	wg := sync.WaitGroup{}
	pOpts := &processAffectedOptions{
		Options: l.Options,
		wg:      &wg,
		sem:     sem,
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

	isAffected := p.isAffected(opts.Target)

	if isAffected && opts.DryRun {
		printer.Printf("%s %s %s\n", color.BlueString(p.Name), color.YellowString("=>"), p.Path)
	}

	if opts.Target == TargetBuild && !opts.DryRun && isAffected {
		if _, err := p.buildProject(); err != nil {
			printer.Errorf("failed to build: %s\n", err.Error())
			return
		}
		if err := p.writeCache(opts.Target); err != nil {
			printer.Errorf("%v\n", err)
			return
		}
	}
}

func (l *List) countAffected() (n int) {
	for _, p := range l.Projects {
		if p.isAffected(l.Options.Target) {
			n++
		}
	}
	return n
}

func (p *Project) isAffected(target Target) bool {
	// TODO: make FileExist optional w/ cmd cli --CheckBinary
	if p.Cache.Targets[target] == nil || !p.Cache.Targets[target].isMetadataMatch(p.Metadata) || !utils.FileExist(path.Join(p.Path, p.Name)) {
		return true
	}
	return false
}
