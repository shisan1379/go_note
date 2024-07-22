package controllers

import (
	cache "gin_gorm_redis_demo/catch"
	"gin_gorm_redis_demo/model"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

type PlayerController struct{}

var rK string = "ranking:"

func (p PlayerController) GetPlayers(c *gin.Context) {
	aidStr := c.DefaultPostForm("aid", "0")
	sortStr := c.DefaultPostForm("sort", "id")
	aid, _ := strconv.Atoi(aidStr)

	players, err := model.GetPlayers(aid, sortStr)
	if err != nil {
		ReturnError(c, 4004, err.Error())
		return
	}

	ReturnSuccess(c, 0, "success", players, int64(len(players)))

}

func (p PlayerController) GetRanking(c *gin.Context) {
	aidStr := c.DefaultPostForm("aid", "0")
	aid, _ := strconv.Atoi(aidStr)

	//定义rediskey
	var redisKey string
	redisKey = "ranking:" + aidStr
	//获取redis信息
	rs, err := cache.Rdb.ZRevRange(cache.Rctx, redisKey, 0, -1).Result()

	if err == nil && len(rs) > 0 {
		idArr, _ := stringSliceToIntSlice(rs)
		players, _ := model.GetPlayerInfoByIds(idArr)
		ReturnSuccess(c, 0, "success-redis", players, 1)
		return
	}

	rsDb, errDb := model.GetPlayers(aid, "score desc")
	if errDb == nil {
		// 注意这里存储的是ID
		// 一般存储的都是ID，这种标识性的值
		// 如果要存储其它类型的值，如结构体，一般会先序列化，然后存储，使用时在反序列化
		// 但是，这通常不是推荐的做法，因为它会使Redis中的数据变得难以查询和维护。相反，最好只存储必要的唯一标识符（如ID）和分数，并在需要时从其他数据源检索完整的数据。
		for _, value := range rsDb {
			cache.Rdb.ZAdd(cache.Rctx, redisKey, cache.Zscore(value.Id, value.Score)).Err()
		}
		//遍历完成以后为rediskey设置过期时间
		cache.Rdb.Expire(cache.Rctx, redisKey, 24*time.Hour)
		ReturnSuccess(c, 0, "success", rsDb, 1)
		return
	}

	ReturnError(c, 4004, "没有相关信息")
}

func stringSliceToIntSlice(s []string) ([]int, error) {
	var result []int
	for _, str := range s {
		// 尝试将字符串转换为整数
		num, err := strconv.Atoi(str)
		if err != nil {
			// 如果转换失败，返回错误
			return nil, err
		}
		// 如果转换成功，将整数添加到结果切片中
		result = append(result, num)
	}
	return result, nil
}
