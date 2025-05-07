package usecases

import (
	"context"
	"encoding/json"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/logger"
)

// CreateResource tries to fetch a resource from the cache using the idempotencyKey as a cache key, which was obtained from
// the request header of the same name. If it doesn't find anything, it creates a new resource on the database and stores
// it in the cache for ttl seconds.
//
// R indicates the type of the resource being created. This is useful for marshalling and unmarshalling the data to and
// from the cache.
//
// persistToDB is a function parameter that creates a new resource on the database. It's only called if there's no
// cached resource.
func CreateResource[R models.Resource](ctx context.Context, cache ResourceCacheManager, idempotenceKey string, ttl int64, persistToDB func() (R, error)) (R, error) {
	cachedResourceStr, err := cache.GetResource(ctx, idempotenceKey)
	if err != nil {
		logger.Error("get_resource_from_cache_failed", err, models.Any("idempotency_key", idempotenceKey))
	}

	if cachedResourceStr != "" {
		logger.AddToContext("idempotency_key", idempotenceKey)
		logger.AddToContext("cached_resource_string", cachedResourceStr)

		var cachedResource R
		err = json.Unmarshal([]byte(cachedResourceStr), &cachedResource)
		if err != nil {
			logger.Error("unmarshal_cached_resource_failed", err)
		}

		if cachedResource != nil {
			logger.Info("returning_cached_resource", models.Any("cached_resource", cachedResource))
			return cachedResource, nil
		}
	}

	createdResource, err := persistToDB()
	if err != nil {
		return nil, err
	}

	err = addToCache(ctx, cache, idempotenceKey, createdResource, ttl)
	if err != nil {
		logger.Error("add_created_resource_to_cache_failed", err, models.Any("idempotency_key", idempotenceKey),
			models.Any("created_resource", createdResource))
	}

	return createdResource, nil
}

func addToCache[R models.Resource](ctx context.Context, cache ResourceCacheManager, idempotenceKey string, resource R, ttl int64) error {
	createdResourceData, err := json.Marshal(resource)
	if err != nil {
		return err
	}

	err = cache.AddResource(ctx, idempotenceKey, createdResourceData, ttl)
	if err != nil {
		return err
	}

	return nil
}
