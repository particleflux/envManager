package secretsStorage

import (
	"envManager/helper"
	"gopkg.in/errgo.v2/fmt/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Configuration struct {
	Options          Options             `yaml:"options,omitempty"`
	Storages         map[string]Storage  `yaml:"storages"`
	Profiles         map[string]Profile  `yaml:"profiles"`
	DirectoryMapping map[string][]string `yaml:"directoryMapping"`
}

type Storage struct {
	StorageType string            `yaml:"type"`
	Config      map[string]string `yaml:"config"`
}

//Options controls the general behavior of envManager. It is only read from the config file in the home directory and
//cannot be overridden by other config files.
type Options struct {
	//DisableCollisionDetection allows overwriting all profiles, storages and mappings instead of returning an error
	//setting this value to true makes CollisionDetectionIgnore meaningless.
	DisableCollisionDetection bool `yaml:"disableCollisionDetection"`
	//CollisionDetectionIgnore allows for fine-grained control over which collisions you want to permit. This setting
	//is ignored when DisableCollisionDetection is set to true.
	CollisionDetectionIgnore CollisionDetectionIgnore `yaml:"collisionDetectionIgnore,omitempty"`
}

// CollisionDetectionIgnore allows overwriting for the listed storages / profiles / mappings only.
type CollisionDetectionIgnore struct {
	Storages []string `yaml:"storages,omitempty"`
	Profiles []string `yaml:"profiles,omitempty"`
	Mappings []string `yaml:"mappings,omitempty"`
}

//NewConfiguration creates a new, empty configuration object
func NewConfiguration() Configuration {
	return Configuration{
		Options:          Options{},
		Storages:         map[string]Storage{},
		Profiles:         map[string]Profile{},
		DirectoryMapping: map[string][]string{},
	}
}

//LoadFromFile loads the config file at the given path. Calling this method on an existing configuration results
//in undefined behavior. Call MergeConfigFile if you want to add another configuration to the existing one.
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

//MergeConfigFile merges the configuration of a file into an existing configuration. Will return an error if a storage,
//profile or mapping of the same name / path already exists and disableCollisionDetection is set to false.
func (c *Configuration) MergeConfigFile(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	fragment := Configuration{}

	err = yaml.Unmarshal(data, &fragment)
	if err != nil {
		return err
	}

	// merging storages
	for name, storageConfig := range fragment.Storages {
		if !c.Options.DisableCollisionDetection {
			_, found := c.Storages[name]
			if found && !helper.SliceStringContains(name, c.Options.CollisionDetectionIgnore.Storages) {
				return errors.Newf("Collision detected, storage name %s is duplicated", name)
			}
		}
		c.Storages[name] = storageConfig
	}

	// merging profiles
	for name, profileConfig := range fragment.Profiles {
		if !c.Options.DisableCollisionDetection {
			_, found := c.Profiles[name]
			if found && !helper.SliceStringContains(name, c.Options.CollisionDetectionIgnore.Profiles) {
				return errors.Newf("Collision detected, profile name %s is duplicated", name)
			}
		}
		c.Profiles[name] = profileConfig
	}

	// merging directory mappings
	for name, mappingConfig := range fragment.DirectoryMapping {
		if !c.Options.DisableCollisionDetection {
			_, found := c.DirectoryMapping[name]
			if found && !helper.SliceStringContains(name, c.Options.CollisionDetectionIgnore.Mappings) {
				return errors.Newf("Collision detected, mapping %s is duplicated", name)
			}
		}
		c.DirectoryMapping[name] = mappingConfig
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
