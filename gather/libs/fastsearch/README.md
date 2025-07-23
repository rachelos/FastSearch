# fastsearch Node.js 客户端

## 安装服务端

在使用fastsearch之前，需要安装服务端，服务端分源码版和二进制文件版，二进制文件版可以通过下载安装，源码版需要先编译安装。


## 安装

```shell
  npm install fastsearch
```

## 使用

```javascript
let Client = require('fastsearch');

let client = new Client("http://127.0.0.1:5678/api");


//添加索引
client.addDocument(100, "小红今年在上一年级", {
    "name": "test", "age": "18",
}).then(r => {
    console.log("索引返回：", r.data)
})

//查询
client.query("小红")
    .then(res => {
        console.log("查询返回：", res.data)
        // 遍历输出文档，
        // 这里为什么有这么多个data，
        // 原因是因为fastsearch响应的带data字段，
        // 而axios这个库也带data字段，就是这么巧。
        if (res.data.state && res.data.data.documents) {
            res.data.data.documents.forEach(item => {
                console.log(item)
            })
        }
    })
    .catch(err => {
        console.log(err)
    })


//删除索引
client.removeDocument(100).then(r => {
    console.log("删除索引返回：", r.data)
})
```
