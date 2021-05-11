package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"os"
	"os/signal"
	"sms_server/common"
	"syscall"
	"time"
)

type RedisCache struct {
	pool      *redis.Pool
	prefix    string
	marshal   func(v interface{}) ([]byte, error)
	unmarshal func(data []byte, v interface{}) error
}

type Options struct {
	Network     string
	Addr        string
	Password    string
	Db          int
	MaxActive   int
	MaxIdle     int
	IdleTimeout int                                    // 空闲连接的超时时间，超过该时间则关闭连接。单位为秒。默认值是5分钟。值为0时表示不关闭空闲连接。此值应该总是大于redis服务的超时时间。
	Prefix      string                                 // 键名前缀
	Marshal     func(v interface{}) ([]byte, error)    // 数据序列化方法，默认使用json.Marshal序列化
	Unmarshal   func(data []byte, v interface{}) error // 数据反序列化方法，默认使用json.Unmarshal序列化
}

//New 根据配置参数创建redis工具实例
func New(options Options) (*RedisCache, error) {
	rc := &RedisCache{}
	err := rc.StartAndGC(options)
	return rc, err
}

//初始化redis，并在进程退出时关闭连接池。
func (rc *RedisCache) StartAndGC(options interface{}) error {
	switch opts := options.(type) {
	case Options:
		if opts.Network == "" {
			opts.Network = "tcp"
		}
		if opts.Addr == "" {
			opts.Addr = "127.0.0.1:6379"
		}
		if opts.MaxIdle == 0 {
			opts.MaxIdle = 3
		}
		if opts.Prefix == "" {
			rc.prefix = "default:"
		} else {
			rc.prefix = opts.Prefix
		}
		if opts.IdleTimeout == 0 {
			opts.IdleTimeout = 300
		}
		if opts.Marshal == nil {
			rc.marshal = json.Marshal
		}
		if opts.Unmarshal == nil {
			rc.unmarshal = json.Unmarshal
		}
		pool := &redis.Pool{
			MaxActive:   opts.MaxActive,
			MaxIdle:     opts.MaxIdle,
			IdleTimeout: time.Duration(opts.IdleTimeout) * time.Second,

			Dial: func() (redis.Conn, error) {
				conn, err := redis.Dial(opts.Network, opts.Addr)
				if err != nil {
					return nil, err
				}
				if opts.Password != "" {
					if _, err := conn.Do("auth", opts.Password); err != nil {
						conn.Close()
						return nil, err
					}
				}
				if _, err := conn.Do("select", opts.Db); err != nil {
					conn.Close()
					return nil, err
				}
				return conn, err
			},
			TestOnBorrow: func(conn redis.Conn, t time.Time) error {
				_, err := conn.Do("ping")
				return err
			},
		}
		rc.pool = pool
		rc.closePool()
		return nil
	default:
		return errors.New("unsupported options")
	}
}

//程序进程退出时关闭连接池
func (rc *RedisCache) closePool() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, syscall.SIGTERM)
	signal.Notify(ch, syscall.SIGKILL)
	go func() {
		<-ch //如果没有信号进来则一直阻塞
		rc.pool.Close()
		os.Exit(0)
	}()
}

func (rc *RedisCache) getKey(key string) string {
	return rc.prefix + key
}

//序列化要保存的值
func (rc *RedisCache) encode(val interface{}) (interface{}, error) {
	var value interface{}
	switch v := val.(type) {
	case string, int, uint, int8, int16, int32, int64, float32, float64, bool:
		value = v
	default:
		json, err := rc.marshal(v)
		if err != nil {
			return nil, err
		}
		value = string(json)
	}
	return value, nil
}

func (rc *RedisCache) decode(reply interface{}, err error, val interface{}) error {
	str, err := common.String(reply, err)
	if err != nil {
		return err
	}
	return rc.unmarshal([]byte(str), val)
}

//执行 redis 命令接口
func (rc *RedisCache) Do(command string, args ...interface{}) (reply interface{}, err error) {
	conn := rc.pool.Get()
	defer conn.Close()
	return conn.Do(command, args...)
}

//获取键值，一般不直接使用该值，而是配合下面的工具类方法获取具体类型，
//或者直接使用github.com/gomodule/redigo/redis包的工具方法。
func (rc *RedisCache) get(key string) (interface{}, error) {
	return rc.Do("get", rc.getKey(key))
}

//查询string类型的键值
func (rc *RedisCache) GetString(key string) (string, error) {
	return common.String(rc.get(key))
}

//获取int类型的键值
func (rc *RedisCache) GetInt(key string) (int, error) {
	return common.Int(rc.get(key))
}

//获取int64类型的键值
func (rc *RedisCache) GetInt64(key string) (int64, error) {
	return common.Int64(rc.get(key))
}

//获取bool类型的键值
func (rc *RedisCache) GetBool(key string) (bool, error) {
	return common.Bool(rc.get(key))
}

//获取非基本类型struct的键值。
//在实现上，使用json的Marshal和Unmarshal做序列化存取。
func (rc *RedisCache) GetObject(key string, val interface{}) error {
	reply, err := rc.get(key)
	return rc.decode(reply, err, val)
}

// 存储并设置有效时长。时间单位为秒
// 基础类型直接保存，其他用json.Marshal 转成string 保存
func (rc *RedisCache) Set(key string, val interface{}, expire int64) error {
	json, err := rc.encode(val)
	if err != nil {
		return err
	}
	if expire > 0 {
		_, err := rc.Do("setex", rc.getKey(key), expire, json)
		return err
	}
	_, err = rc.Do("set", rc.getKey(key), json)
	return err
}

// 检查键是否存在
func (rc *RedisCache) Exists(key string) (bool, error) {
	return common.Bool(rc.Do("exists", rc.getKey(key)))
}

// 删除键
func (rc *RedisCache) Del(key string) error {
	_, err := rc.Do("del", rc.getKey(key))
	return err
}

// 清空当前数据库中所有key
func (rc *RedisCache) Flush() error {
	_, err := rc.Do("flushdb")
	return err
}

//ttl 以秒为单位，当key不存在时，返回 -2 。当key存在但没有设置剩余生存时间时，返回 -1
func (rc *RedisCache) TTL(key string) (ttl int64, err error) {
	return common.Int64(rc.Do("ttl", rc.getKey(key)))
}

//设置键过期时间，expire的单位为秒
func (rc *RedisCache) Expire(key string, expire int64) error {
	_, err := common.Bool(rc.Do("expire", rc.getKey(key), expire))
	return err
}

//incr 自增
func (rc *RedisCache) Incr(key string) (val int64, err error) {
	return common.Int64(rc.Do("incr", rc.getKey(key)))
}

// IncrBy 将 key 所储存的值加上给定的增量值（increment）。
func (rc *RedisCache) IncrBy(key string, amount int64) (val int64, err error) {
	return common.Int64(rc.Do("incrby", rc.getKey(key), amount))
}

// Decr 将 key 中储存的数字值减一。
func (rc *RedisCache) Decr(key string) (val int64, err error) {
	return common.Int64(rc.Do("decr", rc.getKey(key)))
}

// DecrBy key 所储存的值减去给定的减量值（decrement）。
func (rc *RedisCache) DecrBy(key string, amount int64) (val int64, err error) {
	return common.Int64(rc.Do("decrby", rc.getKey(key), amount))
}

//hmset 将一个map存到redis hash，同时设置有效期，单位:秒
func (rc *RedisCache) HMSet(key string, val interface{}, expire int) (err error) {
	conn := rc.pool.Get()
	defer conn.Close()
	err = conn.Send("hmset", redis.Args{}.Add(rc.getKey(key)).AddFlat(val)...)
	if err != nil {
		return err
	}
	if expire > 0 {
		err = conn.Send("expire", rc.getKey(key), int64(expire))
	}
	if err != nil {
		return err
	}
	conn.Flush()
	_, err = conn.Receive()
	return err
}

//redis hash 是一个string类型的field和value的映射表，hash特别适合用于存储对象。
// hset 将哈希表key中的字段field的值设为val
func (rc *RedisCache) Hset(key, field string, val interface{}) (interface{}, error) {
	json, err := rc.encode(val)
	if err != nil {
		return nil, err
	}
	return rc.Do("hset", rc.getKey(key), field, json)
}

//hget 获取存储在哈希表中指定字段的值
func (rc *RedisCache) HGet(key, field string) (reply interface{}, err error) {
	reply, err = rc.Do("hget", rc.getKey(key), field)
	return
}

// HGetString HGet的工具方法，当字段值为字符串类型时使用
func (rc *RedisCache) HGetString(key, field string) (reply string, err error) {
	reply, err = common.String(rc.HGet(key, field))
	return
}

// HGetInt HGet的工具方法，当字段值为int类型时使用
func (rc *RedisCache) HGetInt(key, field string) (reply int, err error) {
	reply, err = common.Int(rc.HGet(key, field))
	return
}

// HGetInt64 HGet的工具方法，当字段值为int64类型时使用
func (rc *RedisCache) HGetInt64(key, field string) (reply int64, err error) {
	reply, err = common.Int64(rc.HGet(key, field))
	return
}

// HGetBool HGet的工具方法，当字段值为bool类型时使用
func (rc *RedisCache) HGetBool(key, field string) (reply bool, err error) {
	reply, err = common.Bool(rc.HGet(key, field))
	return
}

// HGetObject HGet的工具方法，当字段值为非基本类型的stuct时使用
func (rc *RedisCache) HGetObject(key, field string, val interface{}) error {
	reply, err := rc.HGet(key, field)
	return rc.decode(reply, err, val)
}

// HGetAll HGetAll("key", &val)
func (rc *RedisCache) HGetAll(key string, val interface{}) error {
	v, err := redis.Values(rc.Do("hgetall", rc.getKey(key)))
	if err != nil {
		return err
	}
	if err := redis.ScanStruct(v, val); err != nil {
		fmt.Println(err)
	}
	return err
}

//blpop 它是 lpop命令的阻塞版本，当给定列表内没有任何元素可供弹出的时候，连接将被blpop命令阻塞，直到等待超时或发现可弹出元素为止。
//超时参数 timeout 接受一个以秒为单位的数字作为值。超时参数为设置0时，表示阻塞时间可以无限期延长(block indefinitely)
func (rc *RedisCache) BLPop(key string, timeout int) (interface{}, error) {
	values, err := redis.Values(rc.Do("blpop", rc.getKey(key), timeout))
	if err != nil {
		return nil, err
	}
	if len(values) != 2 {
		return nil, fmt.Errorf("redisgo: unexpected number of values,got %d", len(values))
	}
	return values[1], err
}

// BLPopInt BLPop的工具方法，元素类型为int时
func (rc *RedisCache) BLPopInt(key string, timeout int) (int, error) {
	return common.Int(rc.BLPop(key, timeout))
}

// BLPopInt64 BLPop的工具方法，元素类型为int64时
func (rc *RedisCache) BLPopInt64(key string, timeout int) (int64, error) {
	return common.Int64(rc.BLPop(key, timeout))
}

// BLPopString BLPop的工具方法，元素类型为string时
func (rc *RedisCache) BLPopString(key string, timeout int) (string, error) {
	return common.String(rc.BLPop(key, timeout))
}

// BLPopBool BLPop的工具方法，元素类型为bool时
func (rc *RedisCache) BLPopBool(key string, timeout int) (bool, error) {
	return common.Bool(rc.BLPop(key, timeout))
}

// BLPopObject BLPop的工具方法，元素类型为object时
func (rc *RedisCache) BLPopObject(key string, timeout int, val interface{}) error {
	reply, err := rc.BLPop(key, timeout)
	return rc.decode(reply, err, val)
}

// LPop 移出并获取列表中的第一个元素（表头，左边）
func (rc *RedisCache) LPop(key string) (interface{}, error) {
	return rc.Do("lpop", rc.getKey(key))
}

// LPopInt 移出并获取列表中的第一个元素（表头，左边），元素类型为int
func (rc *RedisCache) LPopInt(key string) (int, error) {
	return common.Int(rc.LPop(key))
}

// LPopInt64 移出并获取列表中的第一个元素（表头，左边），元素类型为int64
func (rc *RedisCache) LPopInt64(key string) (int64, error) {
	return common.Int64(rc.LPop(key))
}

// LPopString 移出并获取列表中的第一个元素（表头，左边），元素类型为string
func (rc *RedisCache) LPopString(key string) (string, error) {
	return common.String(rc.LPop(key))
}

// LPopBool 移出并获取列表中的第一个元素（表头，左边），元素类型为bool
func (rc *RedisCache) LPopBool(key string) (bool, error) {
	return common.Bool(rc.LPop(key))
}

// LPopObject 移出并获取列表中的第一个元素（表头，左边），元素类型为非基本类型的struct
func (rc *RedisCache) LPopObject(key string, val interface{}) error {
	reply, err := rc.LPop(key)
	return rc.decode(reply, err, val)
}

// RPop 移出并获取列表中的最后一个元素（表尾，右边）
func (rc *RedisCache) RPop(key string) (interface{}, error) {
	return rc.Do("rpop", rc.getKey(key))
}

// RPopInt 移出并获取列表中的最后一个元素（表尾，右边），元素类型为int
func (rc *RedisCache) RPopInt(key string) (int, error) {
	return common.Int(rc.RPop(key))
}

// RPopInt64 移出并获取列表中的最后一个元素（表尾，右边），元素类型为int64
func (rc *RedisCache) RPopInt64(key string) (int64, error) {
	return common.Int64(rc.RPop(key))
}

// RPopString 移出并获取列表中的最后一个元素（表尾，右边），元素类型为string
func (rc *RedisCache) RPopString(key string) (string, error) {
	return common.String(rc.RPop(key))
}

// RPopBool 移出并获取列表中的最后一个元素（表尾，右边），元素类型为bool
func (rc *RedisCache) RPopBool(key string) (bool, error) {
	return common.Bool(rc.RPop(key))
}

// RPopObject 移出并获取列表中的最后一个元素（表尾，右边），元素类型为非基本类型的struct
func (rc *RedisCache) RPopObject(key string, val interface{}) error {
	reply, err := rc.RPop(key)
	return rc.decode(reply, err, val)
}

// LPush 将一个值插入到列表头部
func (rc *RedisCache) LPush(key string, member interface{}) error {
	value, err := rc.encode(member)
	if err != nil {
		return err
	}
	_, err = rc.Do("lpush", rc.getKey(key), value)
	return err
}

// RPush 将一个值插入到列表尾部
func (rc *RedisCache) RPush(key string, member interface{}) error {
	value, err := rc.encode(member)
	if err != nil {
		return err
	}
	_, err = rc.Do("rpush", rc.getKey(key), value)
	return err
}

// LREM 根据参数 count 的值，移除列表中与参数 member 相等的元素。
// count 的值可以是以下几种：
// count > 0 : 从表头开始向表尾搜索，移除与 member 相等的元素，数量为 count 。
// count < 0 : 从表尾开始向表头搜索，移除与 member 相等的元素，数量为 count 的绝对值。
// count = 0 : 移除表中所有与 member 相等的值。
// 返回值：被移除元素的数量。
func (rc *RedisCache) LREM(key string, count int, member interface{}) (int, error) {
	return common.Int(rc.Do("lrem", rc.getKey(key), count, member))
}

// LLen 获取列表的长度
func (rc *RedisCache) LLen(key string) (int64, error) {
	return common.Int64(rc.Do("rpop", rc.getKey(key)))
}

// LRange 返回列表 key 中指定区间内的元素，区间以偏移量 start 和 stop 指定。
// 下标(index)参数 start 和 stop 都以 0 为底，也就是说，以 0 表示列表的第一个元素，以 1 表示列表的第二个元素，以此类推。
// 你也可以使用负数下标，以 -1 表示列表的最后一个元素， -2 表示列表的倒数第二个元素，以此类推。
// 和编程语言区间函数的区别：end 下标也在 LRANGE 命令的取值范围之内(闭区间)。
func (rc *RedisCache) LRange(key string, start, end int) (interface{}, error) {
	return rc.Do("lrange", rc.getKey(key), start, end)
}

/**
Redis 有序集合和集合一样也是string类型元素的集合,且不允许重复的成员。
不同的是每个元素都会关联一个double类型的分数。redis正是通过分数来为集合中的成员进行从小到大的排序。
有序集合的成员是唯一的,但分数(score)却可以重复。
集合是通过哈希表实现的，所以添加，删除，查找的复杂度都是O(1)。
**/

// ZAdd 将一个 member 元素及其 score 值加入到有序集 key 当中。
func (rc *RedisCache) ZAdd(key string, score int64, member string) (reply interface{}, err error) {
	return rc.Do("zadd", rc.getKey(key), score, member)
}

// ZRem 移除有序集 key 中的一个成员，不存在的成员将被忽略。
func (rc *RedisCache) ZRem(key string, member string) (reply interface{}, err error) {
	return rc.Do("zrem", rc.getKey(key), member)
}

// ZScore 返回有序集 key 中，成员 member 的 score 值。 如果 member 元素不是有序集 key 的成员，或 key 不存在，返回 nil 。
func (rc *RedisCache) ZScore(key string, member string) (int64, error) {
	return common.Int64(rc.Do("zscore", rc.getKey(key), member))
}

// ZRank 返回有序集中指定成员的排名。其中有序集成员按分数值递增(从小到大)顺序排列。score 值最小的成员排名为 0
func (rc *RedisCache) ZRank(key, member string) (int64, error) {
	return common.Int64(rc.Do("zrank", rc.getKey(key), member))
}

// ZRevrank 返回有序集中成员的排名。其中有序集成员按分数值递减(从大到小)排序。分数值最大的成员排名为 0 。
func (rc *RedisCache) ZRevrank(key, member string) (int64, error) {
	return common.Int64(rc.Do("zrevrank", rc.getKey(key), member))
}

// ZRange 返回有序集中，指定区间内的成员。其中成员的位置按分数值递增(从小到大)来排序。具有相同分数值的成员按字典序(lexicographical order )来排列。
// 以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推。或 以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推。
func (rc *RedisCache) ZRange(key string, from, to int64) (map[string]int64, error) {
	return redis.Int64Map(rc.Do("zrange", rc.getKey(key), from, to, "withscores"))
}

// ZRevrange 返回有序集中，指定区间内的成员。其中成员的位置按分数值递减(从大到小)来排列。具有相同分数值的成员按字典序(lexicographical order )来排列。
// 以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推。或 以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推。
func (rc *RedisCache) ZRevrange(key string, from, to int64) (map[string]int64, error) {
	return redis.Int64Map(rc.Do("zrevrange", rc.getKey(key), from, to, "withscores"))
}

// ZRangeByScore 返回有序集合中指定分数区间的成员列表。有序集成员按分数值递增(从小到大)次序排列。
// 具有相同分数值的成员按字典序来排列
func (rc *RedisCache) ZRangeByScore(key string, from, to, offset int64, count int) (map[string]int64, error) {
	return redis.Int64Map(rc.Do("zrangebyscore", rc.getKey(key), from, to, "withscores", "limit", offset, count))
}

// ZRevrangeByScore 返回有序集中指定分数区间内的所有的成员。有序集成员按分数值递减(从大到小)的次序排列。
// 具有相同分数值的成员按字典序来排列
func (rc *RedisCache) ZRevrangeByScore(key string, from, to, offset int64, count int) (map[string]int64, error) {
	return redis.Int64Map(rc.Do("zrevrangebyscore", rc.getKey(key), from, to, "withscores", "limit", offset, count))
}

/**
Redis 发布订阅(pub/sub)是一种消息通信模式：发送者(pub)发送消息，订阅者(sub)接收消息。
Redis 客户端可以订阅任意数量的频道。
当有新消息通过 PUBLISH 命令发送给频道 channel 时， 这个消息就会被发送给订阅它的所有客户端。
**/

// Publish 将信息发送到指定的频道，返回接收到信息的订阅者数量
func (rc *RedisCache) Publish(channel, message string) (int, error) {
	return common.Int(rc.Do("publish", channel, message))
}

// subscribe 订阅给定的一个或多个频道的信息。
// 支持redis服务停止或网络异常等情况时，自动重新订阅。
// 一般的程序都是启动后开启一些固定channel的订阅，也不会动态的取消订阅，这种场景下可以使用本地方法。
// 复杂场景的使用可以直接参考 https://godoc.org/github.com/gomodule/redigo/redis#hdr-Publish_and_Subscribe
func (rc *RedisCache) Subscribe(onMsg func(channel string, data []byte) error, channels ...string) error {
	conn := rc.pool.Get()
	psc := redis.PubSubConn{Conn: conn}
	err := psc.Subscribe(redis.Args{}.AddFlat(channels)...)
	//如果订阅失败，休息1秒后重新订阅(比如当redis服务停止或网络异常)
	if err != nil {
		fmt.Println(err)
		time.Sleep(time.Second)
		return rc.Subscribe(onMsg, channels...)
	}
	quit := make(chan int, 1)
	// 处理信息
	go func() {
		for {
			switch v := psc.Receive().(type) {
			case redis.Message:
				go onMsg(v.Channel, v.Data)
			case redis.Subscription:
				fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
			case error:
				quit <- 1
				fmt.Println(err)
				return
			}
		}
	}()
	//异常情况下重新订阅
	go func() {
		<-quit
		time.Sleep(time.Second)
		psc.Close()
		rc.Subscribe(onMsg, channels...)
	}()
	return err
}

// init 注册到cache 具体使用方法需要google
//func init() {
//	cache.Register("redis", &RedisCache{})
//}
