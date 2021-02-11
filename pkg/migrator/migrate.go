package migrator

import (
	"fmt"
	"math"
	"os"
	"time"

	"redismigrator/pkg/helpers"
	"redismigrator/pkg/logger"

	redis "github.com/go-redis/redis/v8"
)

// DoMigration ...
func DoMigration(sourceClient, destinationClient *redis.Client, keyFilter, keyPrefix string) {
	size, err := helpers.GetDBSize(sourceClient, keyFilter)
	if err != nil {
		logger.Error(fmt.Sprintf("SOURCE: %v", err))
		os.Exit(1)
	}
	logger.Info("Found %d keys on SOURCE with Key filter: %s", size, keyFilter)

	if size == 0 {
		logger.Warn("Migration was aborted as SOURCE has no keys matching filter")
		os.Exit(1)
	}

	logger.Info("Migration running...")

	cmd := sourceClient.Keys(helpers.Ctx, keyFilter)
	if cmd.Err() != nil {
		logger.Error("SOURCE: %s", cmd.Err())
		os.Exit(1)
	}

	keys, err := cmd.Result()
	if err != nil {
		logger.Error("SOURCE: %s", err)
		os.Exit(1)
	}

	counter := 0
	for _, sourceKey := range keys {
		destinationKey := keyPrefix + sourceKey

		// read from source
		typeCmd := sourceClient.Type(helpers.Ctx, sourceKey)
		keyType := typeCmd.Val()

		ttl := sourceClient.TTL(helpers.Ctx, sourceKey).Val()

		switch keyType {
		case "string":
			copyString(sourceClient, destinationClient, sourceKey, destinationKey, ttl)
		case "hash":
			copyHash(sourceClient, destinationClient, sourceKey, destinationKey)
		case "list":
			copyList(sourceClient, destinationClient, sourceKey, destinationKey)
		default:
			logger.Error("Key type not yet sypported: %s", keyType)
			os.Exit(1)
		}

		counter++
	}

	logger.Info("Migration completed with %d keys :)", counter)
	os.Exit(0)
}

func copyString(sourceClient, destinationClient *redis.Client, sourceKey, destinationKey string, ttl time.Duration) {
	val, err := sourceClient.Get(helpers.Ctx, sourceKey).Result()
	if err == nil {
		if setErr := destinationClient.Set(helpers.Ctx, destinationKey, val, ttl).Err(); setErr != nil {
			logger.Error("Could not set '%s' on destination: %s", destinationKey, setErr.Error())
		}
	} else {
		logger.Trace("copyString Error: %+v", err)
	}
}

func copyHash(sourceClient, destinationClient *redis.Client, sourceKey, destinationKey string) {
	fieldCount, err := sourceClient.HLen(helpers.Ctx, sourceKey).Result()
	if err == nil {
		// Count the fields in the hash
		logger.Trace("Hash '%s' has %d fields", sourceKey, fieldCount)
		if fieldCount > 0 {
			var cursor uint64
			// Loop over fields in hash
			for {
				hashKeys, newCursor, hscanErr := sourceClient.HScan(helpers.Ctx, sourceKey, cursor, "", 0).Result()
				if hscanErr == nil {
					cursor = newCursor
					// The hashKeys is a slice of strings
					// where the "key" and "value" are slice siblings
					// instead of a map. ie:
					// [key1, value1, key2, value2, ...]

					// Sanity check to make sure we have an even number of slice items
					isEven := len(hashKeys)%2 == 0
					if !isEven {
						logger.Error("Hash fields are not even for '%s'", sourceKey)
						os.Exit(1)
					}

					destinationClient.HSet(helpers.Ctx, destinationKey, hashKeys)
				} else {
					logger.Trace("copyHash Error: %+v", hscanErr)
				}

				if cursor == 0 {
					break
				}
			}
		}
	} else {
		logger.Trace("copyHash Error: %+v", err)
	}
}

func copyList(sourceClient, destinationClient *redis.Client, sourceKey, destinationKey string) {
	// populate test
	/*
		for x := int64(0); x < 3333; x++ {
			sourceClient.RPush(helpers.Ctx, sourceKey, fmt.Sprintf("example %v", x)).Result()
		}
		logger.Trace("added lots of examples for '%s'", sourceKey)
	*/

	itemsPerPage := float64(1000)
	itemCount, err := sourceClient.LLen(helpers.Ctx, sourceKey).Result()
	if err == nil {
		// Count the items in the List
		logger.Trace("List '%s' has %d items", sourceKey, itemCount)

		if itemCount > 0 {
			pages := math.Ceil(float64(itemCount) / itemsPerPage)
			// Remove the key from destination, prevents contamination
			destinationClient.Del(helpers.Ctx, destinationKey).Result()

			// For each page of up to 1000 items:
			for i := float64(0); i < pages; i++ {
				start := i * itemsPerPage
				end := math.Min(start+itemsPerPage, float64(itemCount))
				items, lrangeErr := sourceClient.LRange(helpers.Ctx, sourceKey, int64(start), int64(end)).Result()
				if lrangeErr == nil {
					// Add to destination
					// Requires a []interface{} not a []string
					itemsInterface := make([]interface{}, len(items))
					for i := range items {
						itemsInterface[i] = items[i]
					}
					destinationClient.RPush(helpers.Ctx, destinationKey, itemsInterface...)
				} else {
					logger.Trace("copyList Error: %+v", lrangeErr)
				}
			}
		}
	} else {
		logger.Trace("copyHash Error: %+v", err)
	}
}
