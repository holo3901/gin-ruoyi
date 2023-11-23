package controllers

import (
	"github.com/gin-gonic/gin"
	"ruoyi/dao/redis"
)

func CacheHandler(ctx *gin.Context) {
	list, err := redis.GetKeyList()
	if err != nil {
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, list)
}

func GetCacheKeysHandler(ctx *gin.Context) {
	name := ctx.Param("cacheName")
	byname, err := redis.GetCacheKeysByname(name)
	if err != nil {
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, byname)
}

func GetCacheValueHandler(ctx *gin.Context) {
	cacheName := ctx.Param("cacheName")
	cacheKey := ctx.Param("cacheKey")
	value, err := redis.GetCacheValue(cacheName, cacheKey)
	if err != nil {
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, value)
}

func ClearCacheNameHandler(ctx *gin.Context) {
	cacheName := ctx.Param("cacheName")
	err := redis.ClearCacheName(cacheName)
	if err != nil {
		ResponseError(ctx, CodeInvalidParam)

		return
	}
	ResponseSuccess(ctx, "清除成功")

}

func ClearCacheKeyHandler(ctx *gin.Context) {
	cacheName := ctx.Param("cacheName")
	cacheKey := ctx.Param("cacheKey")
	err := redis.ClearCacheKey(cacheName, cacheKey)
	if err != nil {
		ResponseError(ctx, CodeInvalidParam)

		return
	}
	ResponseSuccess(ctx, "清除成功")
}

func ClearCacheAllHandler(ctx *gin.Context) {
	if err := redis.ClearCacheAll(); err != nil {
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	ResponseSuccess(ctx, "清除所有成功")
}
