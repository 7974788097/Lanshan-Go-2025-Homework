# 项目概述
这是一个能够直接打包成docker容器的，基于go语言的轻量级社区交流后端，支持用户注册登录，文章发表，发送评论，同时可以关注其他用户，权限系统设置了普通用户，VIP，管理员三级权限，其中管理员能够为其他用户设置禁言时间和禁言原因。项目整体采用Gin框架，GORM作为ORM，使用mysql储存数据，redis缓存用户的登录状态和账号基本信息，使用JWT生成Token进行身份认证，Snowflake生成唯一的评论或者文章id，同时还使用了"golang.org/x/time/rate"进行流量限制。  
## 项目结构展示  
```
a_pile_of_shit/
├── .idea/                          # IDEA 配置文件
│   ├── a_pile_of_shit.iml
│   ├── dataSources.local.xml       # 本地数据源配置
│   ├── dataSources.xml             # 数据源配置
│   ├── modules.xml                 # 模块配置
│   ├── workspace.xml               # 工作区配置
│   └── inspectionProfiles/
│       └── Project_Default.xml     # 检查配置文件
│
├── app/                            # 应用主目录
│   ├── api/                        # API层
│   │   └── api.go                  # 路由和接口初始化
│   │
│   ├── dao/                        # 数据访问层
│   │   ├── dbManage.go             # 数据库连接管理
│   │   ├── mysql.go                 # MySQL配置
│   │   └── redis.go                 # Redis配置
│   │   └── dbmanage/                # 数据库操作
│   │       ├── datamanage/          # 数据管理
│   │       │   ├── datamanage.go    # 数据操作接口
│   │       │   ├── commentmanage/   # 评论管理
│   │       │   │   ├── commentmanage.go
│   │       │   │   └── mysql.go
│   │       │   └── messagemanage/   # 文章管理
│   │       │       ├── messagemanage.go
│   │       │       └── mysql.go
│   │       └── usermanage/          # 用户管理
│   │           ├── usermanage.go    # 用户操作接口
│   │           ├── userinfomanage/  # 用户信息管理
│   │           │   ├── userinfomanage.go
│   │           │   ├── mysql.go
│   │           │   └── redis.go
│   │           ├── userloginmanage/ # 用户登录管理
│   │           │   ├── userloginmanage.go
│   │           │   ├── mysql.go
│   │           │   └── redis.go
│   │           ├── usermutemanage/  # 用户禁言管理
│   │           │   ├── usermutemanage.go
│   │           │   └── mysql.go
│   │           └── userfollowermanage/ # 用户关注管理
│   │               ├── userfollowermanage.go
│   │               └── mysql.go
│   │
│   ├── middleware/                  # 中间件
│   │   ├── getJSONInfo.go           # JSON解析
│   │   ├── limiter.go               # 限流器
│   │   ├── loginState.go            # 登录状态
│   │   ├── permission.go             # 权限控制
│   │   └── token.go                  # JWT令牌
│   │   └── service/                  # 服务层中间件
│   │       ├── CommentManage.go      # 评论服务
│   │       ├── MessageManage.go      # 文章服务
│   │       ├── muteState.go          # 禁言状态
│   │       └── UserManage.go         # 用户服务
│   │
│   ├── model/                        # 数据模型
│   │   ├── dataModel.go              # 评论和文章模型
│   │   ├── dbModel.go                # 数据库模型
│   │   ├── timeModel.go              # 时间配置
│   │   └── userModel.go              # 用户模型
│   │
│   └── other/                        # 工具包
│       └── tool.go                    # 通用工具函数
│
├── confirm/                          # 确认模块(空)
│
├── go.mod                            # Go模块定义
├── go.sum                            # Go依赖校验
├── main.go                           # 程序入口                        # 项目说明
└── README.md                         # 项目说明
```


## 创建镜像
在确保终端是当前文件夹时，使用`docker compose up -d --build`构建可直接运行的镜像，占用8080端口  
## 接口调用  

### 用户管理  
用户权限分普通用户，vip用户和管理员，对应权限等级为1，2，3 ，用户权限存放在Token中 
禁言状态会返回"endTime""muteReason"

| 方法   | 接口                   | 用处     | 传入参数                                      | 补充                                                  |
|------|----------------------|--------|-------------------------------------------|-----------------------------------------------------|
| POST | /user/register       | 注册用户信息 | username,password                         | username长度介于3-18，password长度介于6-18，默认用户权限为1          |
| POST | /user/login          | 用户登录   | username,password                         | 登陆时间为一小时                                            |
| POST | /user/logout         | 用户登出   | Token                                     |                                                     |
| POST | /user/changePassword | 更改密码   | Token,old_password,new_password           | 新旧密码不能相同                                            |
| GET  | /user/get            | 获取用户信息 | Token                                     | 返回followed_names,follow_name,permission,其他信息单独调接口查询 |
| POST | /user/setPermission  | 设置用户权限 | Token,goal_username,goal_permission       | 为目标用户设置权限等级                                         |
| POST | /user/follow         | 关注其他用户 | Token,followed_username                   |                                                     |
| POST | /user/unFollow       | 取消关注   | Token,followed_username                   |                                                     |
| POST | /user/setMute        | 设置禁言   | Token,mute_username,mute_reason,mute_time |                                                     |
| POST | /user/releaseMute    | 取消禁言   | Token,mute_username                       |                                                     |                                                         |

### 文章管理
文章id的末尾为0  
文章的权限即阅读文章所需要的用户权限，1为普通文章，2为vip文章  
文章状态1为草稿状态，search无法获取该状态的文章，2为已发布状态，正常处理  

| 方法   | 接口                     | 用处        | 传入参数                                         | 补充                                                                            |                                         
|------|------------------------|-----------|----------------------------------------------|-------------------------------------------------------------------------------|
| POST | /message/set           | 上传文章      | Token,message_name,message_content           | 文章默认是草稿状态                                                                     |
| GET  | /message/get           | 获取文章内容    | Token,get_way,message_id                     | 当get_way==message_id时，message_id非空，get_way==writer时，从Token中获取username作为writer |
| POST | /message/send          | 发送文章      | Token,message_id                             | 将文章变为已发布状态                                                                    |
| POST | message/update         | 更新文章标题或正文 | Token,message_content,message_id             | message_content,message_id不能同时为空                                              |
| POST | /message/delete        | 删除文章      | Token,message_id                             | 同时也会删除相关联的评论                                                                  |
| GET  | /message/search        | 搜索文章      | Token,message_id,writer,message_name,get_way | 根据get_way中的字段判断搜索方式，但是只能搜索已发布的文章                                              |
| POST | /message/setPermission | 设置文章阅读权限  | Token,message_id,message_permission          |                                                                               |

### 评论管理  
评论id的末尾为1  
通过树结构储存
通过将祖先节点的id转为16进制后以id1/id2/...的形式包装成path,用于快速、高效的完成删除操作，避免了常规树结构调用递归函数的栈溢出风险、数据库操作中断递归不利于回滚以及反复查询开销大的问题。  
同时每个评论节点储存了子节点的id，如果中断等问题可以根据comment_id重新往下查询。  

| 方法   | 接口              | 用处     | 传入参数                                             | 补充                                          |
|------|-----------------|--------|--------------------------------------------------|---------------------------------------------|
| POST | /comment/set    | 添加评论   | Token,comment_content,message_id,parent_node_id  | parent_node_id和message_id只能且必须传入一个作为父节点     |
| GET  | /comment/get    | 获取评论信息 | Token，get_way,comment_id,message_id,conment_path | 根据get_way选择获取评论的方式                          |
| POST | /comment/delete | 删除评论   | Token,comment_path                               | 将所删除的评论对应的path传回，如果删除的文章，将文章id转16进制作为path即可 |