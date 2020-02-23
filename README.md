# undao
云岛

<br>

**核心方法**

- 其中 Update，Add，Del，方法调用时需要给一个 Context，这样在ACID事务操作中尤为重要

| 函数 | 功能描述 |
|---|---|
| undao.CheckOID() | 对 ObjectID 作合规和是否存在检查
| undao.Update() | 更新所有满足条件的数据
| undao.Add() | 添加一条数据
| undao.Del() | 删除所有满足条件的数据
| undao.Len() | 获取满足条件的数据长度
| undao.Get) | 获取单条数据
| undao.GetList() | 获取数据列表，带分页。封装了多集合关联查询

<br>

**聚合管道查询演示案例**

- [传送门](https://github.com/dev-means/undao/blob/master/example/aggregate_lookup.go)

<br>

**其他方法**

- **undao.NewStorageMemory**，创建和mongodb兼容的内存数据库
- **undao.NewStorageDatabase**，连接mongodb数据库服务器
- **undao.MongoIndexesCreateMany**，创建mongodb字段索引
- **undao.MongoIndexesCreateUnique**，创建mongodb字段索引，带有唯一属性
