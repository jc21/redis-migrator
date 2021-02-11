package config

import (
	"redismigrator/pkg/model"

	"github.com/JeremyLoy/config"
	"github.com/alexflint/go-arg"
)

var appArguments model.ArgConfig

// GetConfig returns the ArgConfig
func GetConfig() model.ArgConfig {
	config.FromEnv().To(&appArguments)
	arg.MustParse(&appArguments)
	return appArguments
}
