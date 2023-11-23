package redis

//redis key

//redis key尽量使用命名空间方式,方便查询和拆分

const (
	Prefix             = "ruoyi-go:"
	KeyPostTimeZSet    = "post:time"   //zset;帖子及发帖时间          //命名方法，目的是在多个公司合作时，区别不同的redis
	KeyPostScoreZSet   = "post:score"  //zset;帖子及投票时间
	keyPostVotedZSetPF = "post:voted:" //zset;记录用户及投票类型,参数是post id
	keyLoginZsetPF     = "login:"
	keyCommunitySetPF  = "community:" //保存每个分区下帖子的id
	KeyEmailZsetPF     = "email:"     // 保存邮箱的id
	KeyYanZengZsetPF   = "yanzeng:"
)

// 给redis key加上前缀
func getRedisKey(key string) string {
	return Prefix + key
}
