package project

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/kperreau/goac/pkg/hasher"
	"github.com/kperreau/goac/pkg/scan"
	"github.com/kperreau/goac/pkg/utils"
	"gopkg.in/yaml.v3"

	"github.com/kperreau/goac/pkg/printer"
)

type Project struct {
	Version  string
	Name     string
	Path     string
	GoPath   string
	HashPath string
	Module   *Module
	hashPool *sync.Pool
	Metadata *Metadata
	Cache    *Cache
	Rule     *scan.Rule
}

type IList interface {
	List()
	Affected() error
}

type List struct {
	Projects []*Project
	Options  *Options
}

type Options struct {
	Path           string
	Target         Target
	DryRun         bool
	MaxConcurrency int
}

func NewProjectsList(opt *Options) (IList, error) {
	projects, err := getProjects(opt)
	if err != nil {
		return nil, err
	}

	return &List{
		Projects: projects,
		Options:  opt,
	}, err
}

func find(path string, projectFileName string) (files []string, err error) {
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), projectFileName) {
			files = append(files, path)
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error finding config files: %w", err)
	}

	return files, nil
}

func loadConfig(file string, hashPool *sync.Pool) (*Project, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("error opening project config: %w", err)
	}

	var project Project
	if err = yaml.Unmarshal(data, &project); err != nil {
		printer.Errorf("failed to unmarshal project config: %s", err.Error())
		return nil, err
	}
	project.Path = utils.CleanPath(file, configFileName)
	project.GoPath = utils.AddCurrentDirPrefix(project.Path)
	project.hashPool = hashPool

	hashPath, err := hasher.WithPool(hashPool, project.Path)
	if err != nil {
		return nil, fmt.Errorf("error hashing files: %w", err)
	}
	project.HashPath = hashPath

	return &project, nil
}

type processProjectOptions struct {
	*Options
	projectCh chan *Project
	errorsCh  chan error
	hashPool  *sync.Pool
	wg        *sync.WaitGroup
	sem       chan bool
}

func getProjects(opt *Options) (projects []*Project, err error) {
	projectsFiles, err := find(opt.Path, configFileName)
	if err != nil {
		return nil, err
	}

	// init process options
	sem := make(chan bool, opt.MaxConcurrency+1)
	projectsCh := make(chan *Project)
	errorsCh := make(chan error)
	wg := sync.WaitGroup{}
	pOpts := &processProjectOptions{
		Options:   opt,
		projectCh: projectsCh,
		errorsCh:  errorsCh,
		hashPool:  hasher.NewPool(),
		wg:        &wg,
		sem:       sem,
	}

	for _, projectFile := range projectsFiles {
		sem <- true // acquire
		wg.Add(1)
		go processProject(pOpts, projectFile)
	}

	wg.Wait()
	for i := 0; i < len(projectsFiles); i++ {
		select {
		case project := <-projectsCh:
			projects = append(projects, project)
		case err := <-errorsCh:
			return nil, err
		}
	}

	return projects, nil
}

func processProject(opt *processProjectOptions, projectFile string) {
	defer opt.wg.Done()
	defer func() {
		<-opt.sem // release
	}()

	// load config file .goacproject.yaml
	project, err := loadConfig(projectFile, opt.hashPool)
	if err != nil {
		go func() { opt.errorsCh <- fmt.Errorf("error loading config: %w", err) }()
		return
	}

	// load go modules with go list cmd cli (list imports and dependencies)
	if err := project.LoadGOModules(); err != nil {
		go func() { opt.errorsCh <- fmt.Errorf("error loading modules: %w", err) }()
		return
	}

	// no need affected data, return project (list for example)
	if opt.Target == TargetNone {
		go func() { opt.projectCh <- project }()
		return
	}

	// load caches data
	if err := project.LoadCache(); err != nil {
		go func() { opt.errorsCh <- err }()
		return
	}

	// TODO: make it optional w/ cmd cli --Dockerignore
	// load includes/excludes rule
	project.LoadRule(opt.Target)

	// load hashs
	if err := project.LoadHashs(); err != nil {
		go func() { opt.errorsCh <- err }()
		return
	}

	go func() { opt.projectCh <- project }()
}
