package usecases

import (
	"context"
	"encoding/json"
	"github.com/JoelD7/money/backend/models"
	"github.com/JoelD7/money/backend/shared/logger"
)

// IdempotencyManager manages the creation and caching of resources.
type IdempotencyManager struct {
	cache ResourceCacheManager
}

func NewIdempotencyManager(cache ResourceCacheManager) *IdempotencyManager {
	return &IdempotencyManager{
		cache: cache,
	}
}

// CreateResource tries to fetch a resource from the cache using the idempotencyKey as a cache key, which obtained from
// the request header of the same name. If it doesn't find anything, it creates a new resource on the database and stores
// it in the cache for ttl seconds.
//
// Every caller of this funcion must:
// 1. Provide a function that saves the resource to the database via the persistToDB parameter.
// 2. Type assert the returned resource to the type of the resource that is being created.
func (i *IdempotencyManager) CreateResource(ctx context.Context, idempotenceKey string, ttl int64, persistToDB func() (any, error)) (any, error) {
	cachedResourceStr, err := i.cache.GetResource(ctx, idempotenceKey)
	if err != nil {
		logger.Error("get_resource_from_cache_failed", err, models.Any("idempotency_key", idempotenceKey))
	}

	if cachedResourceStr != "" {
		var cachedResource interface{}
		err = json.Unmarshal([]byte(cachedResourceStr), &cachedResource)
		if err != nil {
			logger.Error("unmarshal_cached_resource_failed", err, models.Any("idempotency_key", idempotenceKey),
				models.Any("cached_resource", cachedResourceStr))
		}

		if cachedResource != nil {
			logger.Info("returning_cached_resource", models.Any("idempotency_key", idempotenceKey))
			return cachedResource, nil
		}
	}

	createdResource, err := persistToDB()
	if err != nil {
		return nil, err
	}

	err = i.addToCache(ctx, idempotenceKey, createdResource, ttl)
	if err != nil {
		logger.Error("add_created_resource_to_cache_failed", err, models.Any("idempotency_key", idempotenceKey),
			models.Any("created_resource", createdResource))
	}

	return createdResource, nil
}

func (i *IdempotencyManager) addToCache(ctx context.Context, idempotenceKey string, resource any, ttl int64) error {
	createdResourceData, err := json.Marshal(resource)
	if err != nil {
		return err
	}

	err = i.cache.AddResource(ctx, idempotenceKey, createdResourceData, ttl)
	if err != nil {
		return err
	}

	return nil
}
