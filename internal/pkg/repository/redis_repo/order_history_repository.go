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

type redisOrderHistoryRepoImpl struct {
	redisClient *redis.Client
	logger      *zerolog.Logger
}

type RedisOrderHistoryRepository interface {
	GetOrderHistoryCtx(ctx context.Context, orderhistoryid, userid int) (*dtos.ResDataOrderHistoryItem, error)
	SetOrderHistoryCtx(ctx context.Context, data *dtos.ResDataOrderHistoryItem, userid int) error
	DeleteOrderHistoryCtx(ctx context.Context, orderhistoryid, userid int) error
}

func NewRedisRepoOrderHistory(redisClient *redis.Client, logger *zerolog.Logger) RedisOrderHistoryRepository {
	return &redisOrderHistoryRepoImpl{
		redisClient: redisClient,
		logger:      logger,
	}
}

func keyOrderHistoryGenerator(orderhistory, userid int) string {
	var key = "orderhistory:"
	return fmt.Sprint(key, userid, ":", orderhistory)
}

var expireTimeOrderHistory = time.Minute * time.Duration(10)

func (roi *redisOrderHistoryRepoImpl) GetOrderHistoryCtx(ctx context.Context, orderhistoryid, userid int) (*dtos.ResDataOrderHistoryItem, error) {

	// Get data from Redis
	realKey := keyOrderHistoryGenerator(orderhistoryid, userid)
	log.Info().Msg(fmt.Sprintf("Get keys %s from redis\n", realKey))
	result, err := roi.redisClient.HGet(ctx, realKey, "").Result()
	if err != nil {
		return nil, errors.Wrap(err, "orderHistoryRedisRepo.GetOrderHistoryCtx.redisClient.Get")
	}

	if len(result) == 0 {
		return nil, errors.Wrap(err, fmt.Sprintf("%s not found", realKey))
	}

	newBase := &dtos.ResDataOrderHistoryItem{}

	if err := newBase.UnmarshalBinary([]byte(result)); err != nil {
		return nil, errors.Wrap(err, "orderHistoryRedisRepo.GetOrderHistoryCtx.mapstructure.Decode")
	}

	log.Info().Msg(fmt.Sprintf("Succedd get keys %s from redis\n", realKey))
	return newBase, nil
}

func (roi *redisOrderHistoryRepoImpl) SetOrderHistoryCtx(ctx context.Context, data *dtos.ResDataOrderHistoryItem, userid int) error {
	newBytes, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "orderHistoryRedisRepo.SetOrderHistoryCtx.json.Marshal")
	}

	log.Info().Msg("set key")
	realKey := keyOrderHistoryGenerator(int(data.ID), userid)
	if err = roi.redisClient.HSet(ctx, realKey, "", newBytes).Err(); err != nil {
		return errors.Wrap(err, "orderHistoryRedisRepo.SetOrderHistoryCtx.redisClient.set")
	}

	if err = roi.redisClient.Expire(ctx, realKey, expireTimeOrderHistory).Err(); err != nil {
		return errors.Wrap(err, "orderHistoryRedisRepo.SetOrderHistoryCtx.redisClient.set")
	}

	log.Info().Msg(fmt.Sprintf("Set keys %s to redis\n", realKey))
	return nil
}

func (roi *redisOrderHistoryRepoImpl) DeleteOrderHistoryCtx(ctx context.Context, orderhistoryid, userid int) error {
	realKey := keyOrderHistoryGenerator(orderhistoryid, userid)
	log.Info().Msg(fmt.Sprintf("Delete keys %s from redis\n", realKey))
	if err := roi.redisClient.Del(ctx, realKey).Err(); err != nil {
		return errors.Wrap(err, "orderHistoryRedisRepo.DeleteOrderHistoryCtx.redisClient.Del")
	}

	log.Info().Msg(fmt.Sprintf("Delete keys %s from redis succeed\n", realKey))
	return nil
}
