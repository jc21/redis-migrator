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

	logger.Info("Migration running, each dot is 1,000 keys, each row is 50,000 keys")

	var cursor uint64
	var n int

	counter := 0
	skipped := 0
	dots := 0

	for {
		var keys []string
		var scanErr error
		keys, cursor, scanErr = sourceClient.Scan(helpers.Ctx, cursor, keyFilter, 1000).Result()
		if scanErr != nil {
			logger.Error("Scan Error: %+v", scanErr)
			os.Exit(1)
		}

		for _, sourceKey := range keys {
			destinationKey := keyPrefix + sourceKey

			// read from source
			typeCmd := sourceClient.Type(helpers.Ctx, sourceKey)
			keyType := typeCmd.Val()
			ttl := sourceClient.TTL(helpers.Ctx, sourceKey).Val()

			logger.Trace("Key '%s' type '%s' ttl: %v destination: '%s'", sourceKey, keyType, ttl, destinationKey)

			switch keyType {
			case "string":
				copyString(sourceClient, destinationClient, sourceKey, destinationKey, ttl)
			case "hash":
				copyHash(sourceClient, destinationClient, sourceKey, destinationKey)
			case "list":
				copyList(sourceClient, destinationClient, sourceKey, destinationKey)
			case "none":
				// Key does not exist, or at least not anymore.
				skipped++
			default:
				logger.Error("Key type not yet supported: %s", keyType)
				logger.Error("Migration was NOT completed!")
				os.Exit(1)
			}

			counter++
		}

		fmt.Print(".")
		dots++
		n += len(keys)
		if cursor == 0 {
			break
		}
		if dots%50 == 0 {
			fmt.Print("\n")
		}
	}

	fmt.Print("\n")
	logger.Info("Migration completed with %d keys, %d skipped :)", counter, skipped)
	os.Exit(0)
}

func copyString(sourceClient, destinationClient *redis.Client, sourceKey, destinationKey string, ttl time.Duration) {
	val, err := sourceClient.Get(helpers.Ctx, sourceKey).Result()
	logger.Trace("Key '%s' val '%v'", sourceKey, val)
	if err == nil {
		// Small hack for situations where duration is negative, causing syntax error
		if ttl < 1 {
			ttl = 0
		}
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
