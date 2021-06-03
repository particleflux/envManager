package secretsStorage

import (
	"envManager/environment"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
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
	Storage   string            `yaml:"storage"`
	Path      string            `yaml:"path"`
	ConstEnv  map[string]string `yaml:"constEnv"`
	Env       map[string]string `yaml:"env"`
	DependsOn []string          `yaml:"dependsOn"`
	visited   bool
}

func LoadConfigurationFromFile(path string) (*Configuration, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config Configuration

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
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
//given environment.Environment instance. After that, it will load all profiles
//which are listed as dependency of this profile.
func (p *Profile) AddToEnvironment(env *environment.Environment) error {
	if p.visited {
		// skip if already visited
		return nil
	}
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
	p.visited = true
	for i := 0; i < len(p.DependsOn); i++ {
		profileName := p.DependsOn[i]
		profile, _ := GetRegistry().GetProfile(profileName)
		err := profile.AddToEnvironment(env)
		if err != nil {
			return err
		}
	}
	return nil
}

//RemoveFromEnvironment removes the environment variables defined by this profile
//from the given environment.Environment instance.
func (p Profile) RemoveFromEnvironment(env *environment.Environment) error {
	if p.visited {
		// skip if already visited
		return nil
	}
	// unload constEnv
	for key, _ := range p.ConstEnv {
		err := env.Unset(key)
		if err != nil {
			return err
		}
	}

	// unload env from storage
	if len(p.Env) > 0 {
		for key, _ := range p.Env {
			err := env.Unset(key)
			if err != nil {
				return err
			}
		}
	}
	p.visited = true
	return nil
}
