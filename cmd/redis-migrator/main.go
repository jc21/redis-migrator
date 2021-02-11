package main

import (
	"fmt"
	"os"

	"redismigrator/pkg/config"
	"redismigrator/pkg/helpers"
	"redismigrator/pkg/logger"
	"redismigrator/pkg/migrator"
)

func main() {
	argConfig := config.GetConfig()
	logger.Init(argConfig)
	logger.Trace("Args: %+v", argConfig)

	argConfig.Print()

	// Ensure configurations are not identical
	if argConfig.IsIdenticalServers() {
		logger.Error("Source and Destination configuration cannot be identical")
		os.Exit(1)
	}

	source := argConfig.GetSource()
	if err := source.Check(); err != nil {
		logger.Error(fmt.Sprintf("SOURCE: %s", err.Error()))
		os.Exit(1)
	}

	destination := argConfig.GetDestination()
	if err := destination.Check(); err != nil {
		logger.Error(fmt.Sprintf("DESTINATION: %s", err.Error()))
		os.Exit(1)
	}

	// Check connectivity to servers
	sourceClient := helpers.NewRedisClient(source)
	destinationClient := helpers.NewRedisClient(destination)

	size1, err1 := helpers.GetDBSize(sourceClient, "*")
	if err1 != nil {
		logger.Error(fmt.Sprintf("SOURCE: %v", err1))
		os.Exit(1)
	}
	logger.Info("Source has %d keys", size1)

	if size1 == 0 {
		logger.Warn("Migration was aborted as SOURCE has no keys")
		os.Exit(0)
	}

	size2, err2 := helpers.GetDBSize(destinationClient, "*")
	if err2 != nil {
		logger.Error(fmt.Sprintf("DESTINATION: %v", err2))
		os.Exit(1)
	}
	logger.Info("Destination has %d keys", size2)
	migrator.DoMigration(sourceClient, destinationClient, argConfig.KeyFilter, argConfig.KeyPrefix)
}
