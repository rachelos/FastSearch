#config.yaml参考配置
```
#监听地址
addr: 0.0.0.0:5679

#数据目录
data: ./data/db

# 负面词目录
negative_data: ./data/neg_data/

#词典目录
dictionary: ./data/dictionary.txt

#是否启用admin
enableAdmin: true

#是否允许删除
allow_drop: true

# 最大线程数
gomaxprocs: 2

# admin 用户名和密码
# auth: 帐号:密码

# 接口是否开启压缩
enableGzip: true

# 数据库关闭超时时间
timeout: 600

# 分片数量
shard: 10

# 分片缓冲数量
bufferNum: 1000

# 是否开启调试模式
debug: false

#通知配置
notice:
    enable: true
    #钉钉通知
    DingTalk:
        enable: true
        web_hook: https://oapi.dingtalk.com/robot/send?access_token=

    # 企业微信通知
    QyWeixin:
        enable: false
        corp_id: ""
        corp_secret: ""

```