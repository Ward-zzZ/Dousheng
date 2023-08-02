# README

# Tiktok-demo 简易版抖音

# 一、项目介绍

**基于 Hertz HTTP框架 + Kitex 微服务框架实现的一个简易版抖音服务端**

**项目具体接口文档**：[简易版抖音项目方案说明](https://bytedance.feishu.cn/docs/doccnKrCsU5Iac6eftnFBdsXTof#)，[API文档](https://apifox.com/apidoc/shared-09d88f32-0b6c-4157-9d07-a36d32d7a75c/api-50707521)

**项目App端下载与使用说明**：[客户端APP](https://bytedance.feishu.cn/docs/doccnM9KkBAdyDhg8qaeGlIz7S7)

# 二、项目实现

## 2.1 项目架构

### 调用关系

![call_relations.png](https://s2.loli.net/2023/07/16/AN1q3rwsKR4FWbQ.png)​

### 系统架构

**考虑到功能组件升级、用户规模变动以及各功能间的负载差异不同，项目采用微服务架构**，可根据用户规模选择**单机部署**或者**集群部署**。

![image-20230723092211619](https://s2.loli.net/2023/07/23/tv8UHjYWF3Mmpo2.png)​

### 服务关系

**一个api服务，六个rpc服务**

1. `http-api`: 作为微服务网关对外提供http服务，接受客户端的请求并解析参数，通过JWT鉴权，并调用相应的RPC服务，将数据返回给客户端
2. `user`: 负责用户的注册、登录和用户信息获取
3. `relation`: 负责用户的关注、粉丝和好友的操作和查询等功能
4. `message`​：负责用户好友的聊天功能
5. `comment`​：负责视频评论信息的操作和查询
6. `favorite`: 负责视频点赞信息的操作和查询
7. `video`: 负责视频的发布、已发布列表，获取视频流等功能

![image-20230723092250122](https://s2.loli.net/2023/07/23/3S1pP7xCQJ5DouX.png)​

## 2.2 技术选型

| **功能**               | **实现**                       |
| ---------------------- | ------------------------------ |
| **HTTP 框架**          | **Hertz**                      |
| **RPC 框架**           | **Kitex**                      |
| **数据库**             | **MySQL、Redis**               |
| **ORM框架**            | **GORM**                       |
| **身份鉴权**           | **JWT**                        |
| **服务发现与配置中心** | **Consul**                     |
| **消息队列**           | **RabbitMQ**                   |
| **服务治理**           | **OpenTelemetry**              |
| **链路追踪**           | **Jaeger、prometheus、logrus** |
| **限流熔断**           | **Sentinel**                   |
| **对象存储**           | **Minio**                      |
| **反向代理与负载均衡** | **Nginx**                      |
| **服务部署**           | **Docker**                     |

## 2.3 数据库设计

### 2.3.1 Mysql表设计

项目采用了**垂直分表**的方法，不同的服务使用不同的数据库，提高数据库的性能和维护性，**外键约束转移至业务层实现**。

#### user表

索引设置：用户主键id使用雪花算法生成，并设定为索引

#### relation表

索引设置：

1. 关注关系id的自增主键索引
2. 用户id和关注的用户id的联合UNIQUE索引，用于在数据库中查询两个用户之间的关注关系

#### message表

索引设置：

1. 消息id的自增主键索引
2. 发送用户id和接收用户id的联合索引，用于在数据库中查询两个用户之间的消息
3. 接收用户id索引，用于查询数据库中指定的用户接收的消息

#### favorite表

索引设置：

1. 点赞关系id的自增主键索引
2. 用户id和视频id的联合索引，用于在数据库中查询用户是否点赞视频

#### comment表

索引设置：

1. 评论id的自增主键索引
2. 视频id的索引，用于在数据库中查询某条视频对应的评论内容
3. 不会真的删除，标记delete_time

#### video表

索引设置：

1. 视频唯一id的索引，由雪花算法生成
2. 发布时间的索引，用户在数据库中查询指定时间范围的视频
3. 作者id的索引，用于查询指定作者的视频列表

![image-20230723084755340](https://s2.loli.net/2023/07/23/NDfTv39O5pmARC6.png)​

### 2.3.2 Redis缓存设计

### user模块

**需求分析：**

在用户模块中，用户信息是高频读的数据，需要做缓存优化。用一个**哈希表(Hash)缓存用户信息**，包括用户名，id，关注数，被关注数（粉丝）。缓存没命中时，利用协程异步插入缓存。

**方案：**

写优化：当发生关注/取关时，引入MQ削峰填谷，**删除缓存**，并通过mq异步更新数据库。

![2023_07_24_23_55_23](https://s2.loli.net/2023/08/03/usiIDXQd9wL8Ak2.png)​​

### relation模块

**需求分析：**

热点事件会导致关注/取关的流量暴增，导致数据库压力增加，因此加入消息队列 RabbitMQ 异步处理数据库写入请求。同时为了保证缓存命中率和一致性，先更新缓存（因为这里更新缓存比较简单，不是来自多张底层数据表的聚合，如果更新了一个字段，那么就要更新整个数据库，还要关联的去查询和汇总各个周边业务系统的数据，这个操作会非常耗时，这种情况删除更好），再异步更新数据库，消息队列会保证一致性。

**方案：**

1、读优化：用户信息的is_follow是一个高频读的数据，将**用户关注列表**存到集合Set中。同时，粉丝列表读次数相对较少，我们将粉丝数大于一定阈值的**用户的粉丝列表**存到一个集合Set中。命中缓存直接返回，否则读取数据库后更新缓存

2、写优化：首先调用user模块更新关注和被关注数，然后更新缓存+异步更新数据库。Write behind （异步缓存写入）

```js
Set: key: user_id, Member: follow_id
Set: key: user_id, value: follower_id
```

### favorite和comment模块

**需求分析：**

因为点赞信息是高频读数据，且对**数据一致性要求较低**，我们更关心**缓存的命中率**。常用的**Cache aside**策略先更新数据库的数据然后删除缓存中的数据，通过缓存数据的删除能加速数据一致性的同步，但是降低了缓存命中率，不太适合这个场景。由于我们**缓存的数据比较简单**，不需要复杂的计算和读取多个表的数据，因此用更新缓存的策略可能更好，同时设计一个较小的过期时间，减少如果不一致对业务的影响。

综上，我们使用**更新缓存，通过消息队列更新数据库**，并设置缓存过期时间兜底。

**方案:**

Set保存用户点赞视频列表，String保存视频点赞数

ZSet保存视频的评论id，时间戳排序，hash保存一条评论的用户id和内容

更新策略：redis缓存+消息队列异步推送mysql

comment使用cache aside（读多写少）

**优点:**

1. 降低对数据库的操作，起到削峰填谷的作用
2. 提高点赞的效率

**缺点:**

1. redis挂掉，或者mq延迟使数据库数据与redis数据不一致

   解决方案:定时同步redis与数据库数据

2. 高频写时会频繁更新缓存会降低性能

3. MQ挂掉，丢失数据

**KV设计：**

```js
// favorite
Set: key: user_id, Member: video_id
String: key: video_id, value: count
// comment
ZSet: key: video_id, Member: comment_id, Score: timestamp
Hash: key: comment_id, field: user_id,content
```

### video和message模块

这两个模块没有使用到缓存，但是这里做一下分析

**需求分析：**

视频信息是高频读的数据，而且由于使用了分表设计，video的数据在插入后就不会改变，因此缓存设计比较简单，不太需要考虑更新的一致性问题。

**方案：**

ZSet保存每个用户的发布视频列表，发布时间排序；

List保存一对好友间的有限条聊天信息；

更新策略：写入视频信息后，然后存在缓存则更新，并设置过期时间

**KV设计：**

```js
// video
ZSet: key: user_id, Member: video_info_json, score: create_time
// message
List：key: user1_id&user2_id, Member: message_info_json
```

‍

## 2.4 项目优化

### 场景和安全优化

1. **限流熔断**：利用sentinel流量治理组件，保障服务的稳定性

2. **身份鉴权**：利用JWT在Api处进行鉴权，对于发布视频、点赞评论等需要用户权限的操作，JWT进行token查验，获取用户Id。对于获取视频、查看评论等操作可以不带token，默认Id为0

3. **数据脱敏：**

   * **不明文存储密码，将用户密码结合**MD5加盐，生成的数据摘要和盐保存起来 ，提高安全性![2023_07_16_17_08_46](https://s2.loli.net/2023/07/16/EibYBw93UQqfI7a.png)​

   * **好友间的聊天消息同样不明文存储，使用**AES对称加密对消息内容加密

     ![2023_07_29_12_30_48](https://s2.loli.net/2023/08/03/s5AfYociZn4NX2E.png)​

4. **雪花算法生成用户和视频的主键Id**：确保 ID 是唯一的、全局的、递增的、可扩展的和高性能的。

5. **用户内容检查**：基于前缀树实现评论内容和视频标题的敏感词过滤，维护友好的交流氛围。识别用户上传文件的MIME类型，判断是否是视频文件。不信任客户端发送的任何文件元数据，例如文件名，文件类型等。所有这些信息都在服务端进行重建。

   ![截屏2023-07-26 下午3.44.49](https://s2.loli.net/2023/08/03/AsUjkSKvZoYx9pL.png)​

6. **MQ实现流量削峰和异步持久化**：对于热点视频的点赞和热点人物的关注取关造成的流量高峰，使用**消息队列**进行流量削锋，减轻数据库压力。当出现取关或者关注动作时，发送请求给对应的消息队列，异步完成数据库的读写

7. **高频读优化**：对于用户信息，是否关注，是否点赞等高频读的信息，利用Redis做缓存优化，提高查询效率，并设置随机缓存时间，防止缓存雪崩。

8. **数据一致性**：为了解决缓存带来的数据一致性问题，在更新数据时使用**Write behind、cache aside**等策略，并通过缓存超时兜底。同时，对于一些写操作比如关注，使用GORM的**事务机制**，保证”增加用户的关注数&增加被关注用户的粉丝数“，保证一致性。

9. **大流量的负载均衡**：使用minio的分布式部署，并通过nginx网关实现反向代理和负载均衡

10. **配置中心**：利用consul实现配置中心，相比文件配置，配置中心能够实现动态更新，分布式支持以及更高的安全性，防止配置数据如服务的地址与端口直接暴露在GitHub仓库中![2023_07_16_20_43_53](https://s2.loli.net/2023/07/16/sMkXqihYayI6WJD.png)​

11. **ffmpeg-go截取视频封面**：采用了ffmpeg-go工具库，获取视频的封面。同时，为了提高响应数据，开启协程上传视频封面

12. SQL注入问题：针对可能的SQL注入问题，使用GORM的参数占位符来构造SQL语句，不使用字符串拼接。

### 实现优化

1. **并行请求**：对于一些请求，比如查询用户信息，获取视频流等信息，需要调用多个RPC模块。这里采用了go协程同步进行，提高效率。
2. **Redis事务管道传输：**
3. **快速上传**：为了提高上传视频接口的响应速度，接口第一时间返回成功，然后使用两个协程，一个上传文件到minio，并通过channel通知另一个协程根据上传结果写入数据库。将接口耗时从原来串行的2s优化到100ms
4. **错误处理**：实现了完整的错误码（share/errno），对业务出现的错误进行分类和提示
5. **Docker一键部署环境**：使用了Docker Compose对需要的环境进行部署，使得服务的管理更为简单便捷，并具有可移植性。

# 三、快速部署

1. **生成环境**：`docker-compose up -d`

2. **拉取依赖**：`go mod tidy`

3. **填写consul的K/V键值，模板在config.yaml**

4. **启动服务**：

   ```
   chmod +x start.sh
   ./start.sh
   ```

# 四、测试结果

## 4.1功能测试

‍见doc

![eaf3daac62700e8c15e8a6e2142cb83](https://s2.loli.net/2023/07/16/8RXbK7nNcOWEyiA.jpg)​

![03d5c228cea65bf2262f979c03ffa9e](https://s2.loli.net/2023/07/16/egoj7xSFAMz2NLV.jpg)​

## 4.2压力测试

**todo**

# 五、项目总结与反思

## 5.1 目前仍存在的问题

1. **尚未对接口进行单元测试和基准测试**
2. **缓存和消息队列的使用不成熟**
3. **项目使用一个数据库，可以做分库分表设计**
4. **日志不完善**
5. 关注列表和粉丝列表时，传入要查询的用户id，但是无法查询返回用户列表和token用户的关注关系。（当前改为只能查看自己的信息）

## 5.2 已识别出的优化项

1. ~~**大V和热门视频的缓存设计使用延时双删会降低缓存命中率，可以使用更新缓存并定时同步数据库的方案**~~
2. ~~**对于点赞数等高频读数据，并不要求非常强的数据一致性，缓存命中率可能更重要。可以采用更新数据库+更新缓存的方法。同时，固定过期时间可以改为随机时间，避免缓存雪崩等情况**~~
3. **优化代码架构，封装公共代码**
4. **可以进行压力测试，找到项目的优化目标**

‍

# 六、开发手册

## IDL

首先根据我们需要的功能编写对应的IDL文件，可以见`idl`​文件夹内的示例

## 代码生成

### Kitex

* 服务依赖代码：在shared文件夹内执行下面的命令生成服务依赖代码，如果服务存在会根据idl更新，否则会新建服务

  ```sh
  kitex -module tiktok-demo -I  ./../idl ./../idl/MessageServer.proto
  ```

* 服务端代码：上面的代码只是定义结构体和方法，我们还需要生成带有脚手架的代码，如下命令。在 `handler.go`​ 的接口中填充业务代码后，执行 `main.go`​ 的主函数即可快速启动 Kitex Server。

  ```bash
  kitex -service MessageService -module tiktok-demo -use tiktok-demo/shared/kitex_gen -I ./../../idl ./../../idl/MessageServer.proto
  ```

## Hertz

在cmd的api目录生成Hertz的网关代码

* 新建：`hz new -idl ./../../idl/ApiServer.proto -mod tiktok-demo/cmd/api`​
* 更新：`hz update -idl ./../../idl/ApiServer.proto -mod tiktok-demo/cmd/api`​

更多Hertz和Kitex的使用方法可以查看[官方文档](https://www.cloudwego.io/zh/ "https://www.cloudwego.io/zh/")，有很多示例可供快速上手

## 业务开发

可以参考现有的服务进行快速开发，项目的目录介绍如下：

```
tiktok-demo
├─ cmd #项目的主要应用程序
│  ├─ api #微服务网关
│  │  ├─ biz
│  │  │  ├─ handler
│  │  │  │  ├─ ApiServer # 业务代码
│  │  │  │  └─ ping.go
│  │  │  ├─ middleware   # 中间件
│  │  │  ├─ model        # 数据模型，Hertz生成
│  │  │  └─ router
│  │  │     ├─ ApiServer # 路由以及对应的中间件添加
│  │  ├─ cert            # TLS证书
│  │  ├─ config          # 配置结构体定义
│  │  ├─ config.yaml     # 配置文件，是consul的地址和端口
│  │  ├─ initialize
│  │  │  ├─ config.go    # 配置初始化，读取consul对应的 K/V 键值
│  │  │  ├─ cors.go
│  │  │  ├─ jwt.go       # JWT鉴权
│  │  │  ├─ logger.go    # 日志初始化
│  │  │  ├─ registry.go  # 服务注册
│  │  │  ├─ rpc          # 需要的RPC客户端初始化
│  │  │  ├─ sentinel.go  # 限流熔断
│  │  │  └─ tls.go       # TLS证书读取
│  │  ├─ main.go         # Api网关主函数
│  │  ├─ pkg             # 一些公共的包，主要是响应体的封装
│  ├─ comment #评论服务
│  │  ├─ build.sh        #编译脚本
│  │  ├─ config          #配置定义
│  │  ├─ config.yaml     #配置文件
│  │  ├─ handler.go      #业务处理代码
│  │  ├─ initialize
│  │  │  ├─ config.go    #配置初始化
│  │  │  ├─ db.go        #数据库连接初始化
│  │  │  ├─ flag.go      #初始化RPC的ip和端口
│  │  │  ├─ logger.go    #日志初始化
│  │  │  ├─ mq.go        #消息队列初始化
│  │  │  ├─ redis.go     #redis初始化
│  │  │  ├─ registry.go  #服务注册
│  │  │  └─ user_service.go #需要的RPC客户端初始化
│  │  ├─ main.go         #主函数
│  │  ├─ pkg             #数据操作的封装，以及其他工具包
│  │  │  ├─ mq
│  │  │  ├─ mysql
│  │  │  ├─ pack
│  │  │  └─ redis
│  ├─ favorite
      ......             # 同上
│  ├─ message
      ......             # 同上
│  ├─ relation
      ......             # 同上
│  ├─ user
      ......             # 同上
│  └─ video
      ......             # 同上
├─ data   # Docker相关配置和数据
│  ├─ configs           # Docker启动时的配置文件
│  ├─ data              # 数据库数据
│  └─ minio             # 对象存储数据
├─ idl                  # IDL文件
├─ shared # 公共代码
│  ├─ consts
│  │  └─ consts.go      # 一些常量
│  ├─ errno             # 全局定义的错误码
│  ├─ kitex_gen         # Kitex生成的代码
│  ├─ middleware        # 全局中间件
│  │  ├─ client_log.go
│  │  ├─ common_log.go
│  │  ├─ recovery.go    # 用于恢复panic
│  │  └─ server_log.go  # 用于打印RPC调用信息
│  └─ tools
│     ├─ cover.go      # 用于生成视频封面
│     ├─ ip.go
│     ├─ port.go       # 用于生成随机端口
│     └─ resp.go       # 用于封装基本响应体
```

‍
