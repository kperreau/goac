package project

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"golang.org/x/mod/modfile"

	"github.com/kperreau/goac/pkg/hasher"
	"github.com/kperreau/goac/pkg/scan"
	"github.com/kperreau/goac/pkg/utils"
	"gopkg.in/yaml.v3"

	"github.com/kperreau/goac/pkg/printer"
)

type Env struct {
	Key   string `yaml:",omitempty"`
	Value string `yaml:",omitempty"`
}

type Exec struct {
	CMD    string
	Params []string `yaml:",omitempty"`
}

type TargetConfig struct {
	Envs []Env `yaml:",omitempty"`
	Exec *Exec
}

type Project struct {
	Version    string
	Name       string
	Path       string `yaml:",omitempty"`
	CleanPath  string `yaml:",omitempty"`
	Target     map[Target]*TargetConfig
	HashPath   string     `yaml:",omitempty"`
	Module     *Module    `yaml:",omitempty"`
	HashPool   *sync.Pool `yaml:",omitempty"`
	Metadata   *Metadata  `yaml:",omitempty"`
	Cache      *Cache     `yaml:",omitempty"`
	Rule       *scan.Rule `yaml:",omitempty"`
	CMDOptions *Options   `yaml:",omitempty"`
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
	Target         Target
	DryRun         bool
	MaxConcurrency int
	BinaryCheck    bool
	Force          bool
	DockerIgnore   bool
	ProjectsName   []string
	Debug          []string
	PrintStdout    bool
}

var RootPath = "."

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

func loadConfig(file string, opts *processProjectOptions) (*Project, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("error opening project config: %w", err)
	}

	var project Project
	if err = yaml.Unmarshal(data, &project); err != nil {
		printer.Errorf("failed to unmarshal project config: %s", err.Error())
		return nil, err
	}
	project.CleanPath = utils.CleanPath(file, configFileName)
	project.Path = utils.AddCurrentDirPrefix(project.CleanPath)
	project.HashPool = opts.hashPool
	project.CMDOptions = opts.Options

	hashPath, err := hasher.WithPool(opts.hashPool, project.CleanPath)
	if err != nil {
		return nil, fmt.Errorf("error hashing files: %w", err)
	}
	project.HashPath = hashPath

	return &project, nil
}

type processProjectOptions struct {
	*Options
	gomod     *modfile.File
	projectCh chan *Project
	errorsCh  chan error
	hashPool  *sync.Pool
	wg        *sync.WaitGroup
	sem       chan bool
}

func getProjects(opt *Options) (projects []*Project, err error) {
	if opt.MaxConcurrency < 1 {
		return nil, fmt.Errorf("max concurrency can't be less than 1")
	}

	projectsFiles, err := find(RootPath, configFileName)
	if err != nil {
		return nil, err
	}

	// preload go mod file dependencies
	gomod, err := loadGOModFile(RootPath)
	if err != nil {
		return nil, err
	}

	// init process options
	sem := make(chan bool, opt.MaxConcurrency)
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
		gomod:     gomod,
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
			if project != nil {
				projects = append(projects, project)
			}
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
	project, err := loadConfig(projectFile, opt)
	if err != nil {
		go func() { opt.errorsCh <- fmt.Errorf("error loading config: %w", err) }()
		return
	}

	// Skip if option cli cmd --projects is set and not match this project name
	if len(opt.Options.ProjectsName) > 0 && !slices.Contains(opt.Options.ProjectsName, project.Name) {
		go func() { opt.projectCh <- nil }()
		return
	}

	// load go modules with go list cmd cli (list imports and dependencies)
	if err := project.LoadGOModules(opt.gomod); err != nil {
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

	if opt.DockerIgnore {
		// load includes/excludes rule
		project.LoadRule(opt.Target)
	}

	// load hashs
	if err := project.LoadHashs(); err != nil {
		go func() { opt.errorsCh <- err }()
		return
	}

	go func() { opt.projectCh <- project }()
}
