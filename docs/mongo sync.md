## mongo 同步到 starrocks

### 使用说明
#### 环境准备
```
1. GO构建版本 v1.18.10
2. mongo版本 4.2及以上；使用change stream监听方式，理论支持3.6及以上
```
#### 1. 新增配置文件
mongo-to-starrocks.toml
```toml
# name 必填，多实例运行时保证全局唯一
name = "mongo2starrocks"

[input]
# 指定初次监听开始时间点，当_xxx-pos.info点位文件内容存在时，此选项不生效
start-position = 2023-03-27T11:00:00 # Without double quotes, it means the local time zone
# 转换document Field从camelCase到snakeCase，默认false；例如 userName（mongo） -> user_name（starrocks）
# 设置为true，还会处理主键Field命名：_id（mongo） -> id（starrocks）
convert-snake-case = false

[mongo]
uri = "192.168.0.1:3717,192.168.0.2:3717,192.168.0.3:3717/admin?replicaSet=mgset-xxxxx"
username = "root"
password = ""

[starrocks]
host = "127.0.0.1"
port = 8040
username = "root"
password = ""

[sync-param]
# 同步chan队列最大值，达到会进行flush，最小100
channel-size = 10240
# 同步延迟秒数，达到会进行flush，最小1
flush-delay-second = 10

[[rule]]
source-schema = "mongo_test"
source-table = "coll1"
target-schema = "starrocks_test"
target-table = "coll1"

[[rule]]
source-schema = "mongo_test"
source-table = "coll2"
target-schema = "starrocks_test"
target-table = "coll2"
```
#### 2. 启动
```shell
[sr@ ~]$ ./go-mongo-starrocks-linux-xxxxxx -config mongo-to-starrocks.toml
```
#### 3. 查看日志
默认输出到控制台

指定log-file参数运行
```shell
[sr@ ~]$ ./go-mongo-starrocks-linux-xxxxxx -config mongo-to-starrocks.toml -log-file mongo2starrocks.log
[sr@ ~]$ tail -f mongo2starrocks.log
```

#### 4. 查看帮助
```shell
[sr@ ~]$ ./go-mongo-starrocks-linux-xxxxxx -h
```
#### 5. 后台运行
```shell
[sr@ ~]$ (nohup ./go-mongo-starrocks-linux-xxxxxx -config mongo-to-starrocks.toml -log-file mongo2starrocks.log -level info &)
```
