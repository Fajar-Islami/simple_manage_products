package redis_repo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Fajar-Islami/simple_manage_products/internal/pkg/dtos"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type redisOrderItemsRepoImpl struct {
	redisClient *redis.Client
	logger      *zerolog.Logger
}

type RedisOrderItemsRepository interface {
	GetOrderItemsCtx(ctx context.Context, orderitemsid int) (*dtos.ResDataOrderItemsData, error)
	SetOrderItemsCtx(ctx context.Context, data *dtos.ResDataOrderItemsData) error
	DeleteOrderItemsCtx(ctx context.Context, orderitemsid int) error
}

func NewRedisRepoOrderItems(redisClient *redis.Client, logger *zerolog.Logger) RedisOrderItemsRepository {
	return &redisOrderItemsRepoImpl{
		redisClient: redisClient,
		logger:      logger,
	}
}

func keyOrderItemsGenerator(id int) string {
	var key = "orderitems:"
	return fmt.Sprint(key, id)
}

var expireTimeOrderItems = time.Minute * time.Duration(10)

func (roi *redisOrderItemsRepoImpl) GetOrderItemsCtx(ctx context.Context, orderitemsid int) (*dtos.ResDataOrderItemsData, error) {

	// Get data from Redis
	realKey := keyOrderItemsGenerator(orderitemsid)
	log.Info().Msg(fmt.Sprintf("Get keys %s from redis\n", realKey))
	result, err := roi.redisClient.HGet(ctx, realKey, "").Result()
	if err != nil {
		return nil, errors.Wrap(err, "orderItemRedisRepo.GetOrderItemsCtx.redisClient.Get")
	}

	if len(result) == 0 {
		return nil, errors.Wrap(err, fmt.Sprintf("%s not found", realKey))
	}

	newBase := &dtos.ResDataOrderItemsData{}

	if err := newBase.UnmarshalBinary([]byte(result)); err != nil {
		return nil, errors.Wrap(err, "orderItemRedisRepo.GetOrderItemsCtx.mapstructure.Decode")
	}

	log.Info().Msg(fmt.Sprintf("Succedd get keys %s from redis\n", realKey))
	return newBase, nil
}

func (roi *redisOrderItemsRepoImpl) SetOrderItemsCtx(ctx context.Context, data *dtos.ResDataOrderItemsData) error {
	newBytes, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "orderItemRedisRepo.SetOrderItemsCtx.json.Marshal")
	}

	log.Info().Msg("set key")
	realKey := keyOrderItemsGenerator(int(data.ID))
	if err = roi.redisClient.HSet(ctx, realKey, "", newBytes).Err(); err != nil {
		return errors.Wrap(err, "orderItemRedisRepo.SetOrderItemsCtx.redisClient.set")
	}

	if err = roi.redisClient.Expire(ctx, realKey, expireTimeOrderItems).Err(); err != nil {
		return errors.Wrap(err, "orderItemRedisRepo.SetOrderItemsCtx.redisClient.set")
	}

	log.Info().Msg(fmt.Sprintf("Set keys %s to redis\n", realKey))
	return nil
}

func (roi *redisOrderItemsRepoImpl) DeleteOrderItemsCtx(ctx context.Context, orderitemsid int) error {
	realKey := keyOrderItemsGenerator(orderitemsid)
	log.Info().Msg(fmt.Sprintf("Delete keys %s from redis\n", realKey))
	if err := roi.redisClient.Del(ctx, realKey).Err(); err != nil {
		return errors.Wrap(err, "orderItemRedisRepo.DeleteOrderItemsCtx.redisClient.Del")
	}

	log.Info().Msg(fmt.Sprintf("Delete keys %s from redis succeed\n", realKey))
	return nil
}
