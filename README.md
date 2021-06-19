# EnvManager

## Concept

The envManager knows **profiles** and **storages**. A **storage** is a secure location where secrets are stored (e.g.
keepass). A **profile** is a set of environment variables which should be set / unset together. The environment variables
can either have a constant value or get their value from a storage. A profile can depend on other profiles, which means
that loading the first profile will load all dependencies as well. Even circular dependencies are allowed.

## Setup

Add the function in `wrapper.sh` to your shell (e.g. `. ./wrapper.sh`). The file contains the function `envManager`
which in turn calls `envManager-bin` (assuming it is in your PATH). Should your `envManager-bin` not be in your PATH,
replace `envManager-bin` with an absolute path to the binary (e.g. `/home/john.doe/envManager/envManager-bin`).

## Usage

Create an initial config with `envManager config init`. By default, the application creates a `.envManager.yml` in your
home directory.
You can now add storages and profiles by hand or let the application do it for you. Run `envManager config add storage`
to add your first storage and `envManager config add profile` to add your first profile. After that, you can easily
copy&paste the profile and storage configurations.

## Available storage adapters

Currently, only the `keepass` adapter is supported.

### Keepass / KeepassX / KeepassXC

This adapter can read keepass2 files (.kdbx). Its config contains one key, `path` which contains an absolute path to the
kdbx file.

**Example config**
```yaml
storages:
  myStorageName:
    type: keepass
    config:
      path: /home/john.doe/myKeepassFile.kdbx
```

## FAQ 

### Can I use multiple storages for one profile?

No, but you can create one profile which depends on multiple profiles. If you load the "main" profile, the dependencies
will be loaded automatically.

### Can I have dependencies across multiple storages?

Yes.

### I want to have two profiles depending on each other, will I run into an endless loop?

No, circular dependencies are resolved without looping endlessly.

## Test data

In the `/testData` directory is a dummy `keepass.kdbx` containing the following entries. The password for this database is `1234`.

- `entry1` with `user1` and `pass1` as well as the additional attribute `advanced1` with `advanced1-value`
- `group1/g1e1` with `g1e1-user` and `g1e1-pass` and the additional attribute `g1e1-advanced1` with `g1e1-advanced1-value`

## Used libraries

- [github.com/josa42/go-prompt](https://github.com/josa42/go-prompt)
- [github.com/manifoldco/promptui](https://github.com/manifoldco/promptui)
- [github.com/spf13/cobra](https://github.com/spf13/cobra)
- [github.com/tobischo/gokeepasslib/v3](https://github.com/tobischo/gokeepasslib)
- [gopkg.in/errgo.v2](https://gopkg.in/errgo.v2)
- [gopkg.in/yaml.v2 v2.4.0](https://gopkg.in/yaml.v2)