package secretsStorage

import (
	"envManager/environment"
	"envManager/helper"
	"fmt"
	"gopkg.in/errgo.v2/fmt/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Configuration struct {
	Storages map[string]Storage `yaml:"storages"`
	Profiles map[string]Profile `yaml:"profiles"`
}

type Storage struct {
	StorageType string            `yaml:"type"`
	Config      map[string]string `yaml:"config"`
}

type Profile struct {
	Name      string
	Storage   string            `yaml:"storage"`
	Path      string            `yaml:"path"`
	ConstEnv  map[string]string `yaml:"constEnv"`
	Env       map[string]string `yaml:"env"`
	DependsOn []string          `yaml:"dependsOn"`
}

//NewConfiguration creates a new, empty configuration object
func NewConfiguration() Configuration {
	return Configuration{
		Storages: map[string]Storage{},
		Profiles: map[string]Profile{},
	}
}

//LoadFromFile loads the config file at the given path. It will merge the current config with the config from the file.
func (c *Configuration) LoadFromFile(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return err
	}

	return nil
}

//WriteToFile writes the current config to given path. It will not overwrite an existing file
//except when replace is set to true.
func (c Configuration) WriteToFile(path string, replace bool) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	fileInfo, err := os.Stat(path)
	if fileInfo != nil && !replace {
		return errors.Newf("Will not overwrite %s without being explicitly told to do so.", path)
	}
	err = ioutil.WriteFile(path, data, 0600)

	return err
}

//Validate checks the validity of the profile. The storage and all profiles this
//profile depends on must be known to the registry.
func (p Profile) Validate() []string {
	var out []string
	registry := GetRegistry()
	if !registry.HasStorage(p.Storage) {
		out = append(out, fmt.Sprintf("references storage %s which is not defined", p.Storage))
	}
	for i := 0; i < len(p.DependsOn); i++ {
		if !registry.HasProfile(p.DependsOn[i]) {
			out = append(out, fmt.Sprintf("depends on %s which is not defined", p.DependsOn[i]))
		}
	}
	return out
}

//AddToEnvironment adds the environment variables defined by this profile to the
//given environment.Environment instance.
func (p *Profile) AddToEnvironment(env *environment.Environment) error {
	// load constEnv
	for key, value := range p.ConstEnv {
		err := env.Set(key, value)
		if err != nil {
			return err
		}
	}

	// load env from storage
	if len(p.Env) > 0 {
		storage, err := GetRegistry().GetStorage(p.Storage)
		if err != nil {
			return err
		}
		entry, err := (*storage).GetEntry(p.Path)
		if err != nil {
			return err
		}
		for key, attributeName := range p.Env {
			value, err := entry.GetAttribute(attributeName)
			if err != nil {
				return err
			}
			err = env.Set(key, *value)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

//RemoveFromEnvironment removes the environment variables defined by this profile
//from the given environment.Environment instance.
func (p Profile) RemoveFromEnvironment(env *environment.Environment) error {
	// unload constEnv
	for key := range p.ConstEnv {
		err := env.Unset(key)
		if err != nil {
			return err
		}
	}

	// unload env from storage
	if len(p.Env) > 0 {
		for key := range p.Env {
			err := env.Unset(key)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

//GetDependencies gets the dependencies of this profile and its dependencies.
func (p Profile) GetDependencies(alreadyVisited []string) ([]string, error) {
	var dependencies []string
	alreadyVisited = append(alreadyVisited, p.Name)
	dependencies = append(dependencies, p.DependsOn...)
	for _, name := range p.DependsOn {
		if helper.SliceStringContains(name, alreadyVisited) {
			// do not load dependencies of a profile we already visited
			continue
		}
		childProfile, err := GetRegistry().GetProfile(name)
		if err != nil {
			return nil, errors.Newf("Unknown dependency %s", name)
		}
		alreadyVisited = append(alreadyVisited, name)
		childDeps, err := childProfile.GetDependencies(alreadyVisited)
		if err != nil {
			return nil, err
		}
		dependencies = append(dependencies, childDeps...)
	}
	dependencies = helper.SliceStringRemove(p.Name, dependencies)
	return dependencies, nil
}
