package dict

type Consumer func(key string, val interface{}) bool

// Dict 迭代项目更方便
type Dict interface {
	Get(key string) (val interface{}, exists bool)        // 获取
	Len() (result int)                                    // 字典数据长度
	Put(key string, val interface{}) (result int)         // 存放
	PutIfAbsent(key string, val interface{}) (result int) // 如果没有在存入
	PutIfExists(key string, val interface{}) (result int) // 修改操作
	Remove(key string) (result int)                       // 删除
	ForEach(consumer Consumer)                            // 遍历
	Keys() []string
	RandomKeys(limit int) []string         // 随机返回
	RandomDistinctKeys(limit int) []string // 返回多个不同的键
	Clear()
}
