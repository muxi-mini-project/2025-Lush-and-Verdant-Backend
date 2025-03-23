package dao

import (
	"2025-Lush-and-Verdant-Backend/api/request"
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
)

type LikeDAO interface {
	IncrementLike(to string)
	DecrementLike(to string)
	SaveLike(like *request.ForestLikeReq)
	SaveUnlike(like *request.ForestLikeReq)
	Check(like *request.ForestLikeReq) bool
	GetLikes(to string) (string, error)
}

type LikeDAOImpl struct {
	rdb *redis.Client
}

func NewLikeDAOImpl(rdb *redis.Client) *LikeDAOImpl {
	return &LikeDAOImpl{
		rdb: rdb,
	}
}

func (dao *LikeDAOImpl) IncrementLike(to string) {
	dao.rdb.Incr(context.Background(), "like_count:"+to)
}
func (dao *LikeDAOImpl) DecrementLike(to string) {
	dao.rdb.Decr(context.Background(), "like_count:"+to)
}
func (dao *LikeDAOImpl) SaveLike(like *request.ForestLikeReq) {
	dao.rdb.SAdd(context.TODO(), "like_from_to:"+like.From, like.To)
}
func (dao *LikeDAOImpl) SaveUnlike(like *request.ForestLikeReq) {
	dao.rdb.SRem(context.TODO(), "like_from_to:"+like.From, like.To)
}
func (dao *LikeDAOImpl) Check(like *request.ForestLikeReq) bool {
	return dao.rdb.SIsMember(context.TODO(), "like_from_to:"+like.From, like.To).Val()
}

func (dao *LikeDAOImpl) GetLikes(to string) (string, error) {
	numStr, err := dao.rdb.Get(context.TODO(), "like_count:"+to).Result()
	if errors.Is(err, redis.Nil) {
		return "0", nil
	} else if err != nil {
		return "", err
	}
	return numStr, nil
}
