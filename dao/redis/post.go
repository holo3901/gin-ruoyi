package redis

//func getIDsFormKey(key string, page, size int64) ([]string, error) {
//	start := (page - 1) * size
//	end := start + size - 1
//	//3.ZREVRANGE按分数从大到小的顺序查询指定数量的元素
//	return client.ZRevRange(key, start, end).Result()
//}
//func GetPostIDsOrder(p *models.ParamPostList) ([]string, error) {
//	//从redis获取id
//	//1.根据用户请求中携带的order参数确定要查询的redis key
//	orderKey := getRedisKey(KeyPostTimeZSet)
//	if p.Order == models.OrderScore {
//		orderKey = getRedisKey(KeyPostScoreZSet)
//	}
//	//2.确定查询的索引起始点
//	return getIDsFormKey(orderKey, p.Page, p.Size)
//}
//
//// GetPostVoteData 根据IDS查询每篇帖子的投赞成票的数据
//func GetPostVoteData(ids []string) (data []int64, err error) {
//	/*data =make([]int64,0,len(ids))
//	for _,id :=range ids{
//		key:=getRedisKey(keyPostVotedZSetPF+id)
//		//查找key中分数是1的元素的数量->统记每篇帖子的赞成票数量 ，统计反对票投票数量把min和max改为-1
//		v:=client.ZCount(key,"1","1").Val()
//		data=append(data,v)
//	}*/ //每次都要通过client.ZCount得到数据,需要优化
//	//使用pipeline一次发送多条命令，减少RTT
//	pipeline := client.Pipeline() //通过pipeline避免重复查询redis
//	for _, id := range ids {
//		key := getRedisKey(keyPostVotedZSetPF + id)
//		pipeline.ZCount(key, "1", "1")
//	}
//	cmders, err := pipeline.Exec()
//	if err != nil {
//		return nil, err
//	}
//	data = make([]int64, 0, len(cmders))
//	for _, cmder := range cmders {
//		v := cmder.(*redis.IntCmd).Val()
//		data = append(data, v)
//	}
//	return
//}
//
//// GetCommunityPostIDsInOrder 按社区查询ids
//func GetCommunityPostIDsInOrder(P *models.ParamPostList) ([]string, error) {
//	//使用zintterstore把分区的帖子set与帖子分数的zset生成一个新的zset
//	//针对新的zset按之前的逻辑取数据
//	orderKey := getRedisKey(KeyPostTimeZSet)
//	if P.Order == models.OrderScore {
//		orderKey = getRedisKey(KeyPostScoreZSet)
//	}
//	//社区的key
//	cKey := getRedisKey(keyCommunitySetPF + strconv.Itoa(int(P.CommunityID)))
//	//利用缓存key减少zinterstore执行的次数,当运行过getIDsFormKey之后的60分钟内，直接根据key查询ids
//	key := orderKey + strconv.Itoa(int(P.CommunityID))
//	if client.Exists(key).Val() < 1 {
//		//不存在，需要计算
//		pipeline := client.Pipeline()
//		pipeline.ZInterStore(key, redis.ZStore{
//			Aggregate: "MAX",
//		}, cKey, orderKey) //zinterstore 计算
//		pipeline.Expire(key, 60*time.Second) //设置超时时间
//		_, err := pipeline.Exec()
//		if err != nil {
//			return nil, err
//		}
//	}
//	//存在的话就直接根据key查询ids
//	return getIDsFormKey(key, P.Page, P.Size)
//
//}
//
