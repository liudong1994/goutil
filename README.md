## goutil

go语言编写的一些小组件

#### rediscli

------

- Binary2string：二进制数据 转换为 redis-cli中显示的字符串，可以直接redis-cli中set；

- String2binary：redis-cli显示的字符 转换为 二进制数据。


#### redisconvert

------

- Rediscli2pb2json：redis-cli显示字符串 转换为 pb结构 转换为 json数据；

    - 用于将redis-cli中显示字符串 转换为 json可视化数据，方便查看。

- Json2pb2rediscli：json数据 转换为 pb结构 转换为 redis-cli字符串。

    - 用于将mock的json数据转换为redis-cli字符串，set到redis中，方便造假数据使用。
