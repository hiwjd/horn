# horn doc

## 消息接口

### 地址

```
POST http://app.horn.com/api/message

消息体
```

### 消息格式

#### 普通对话消息

```
{
    "mid": "消息ID 发送时该字段省略",
    "type": "text",
    "from": {
        "id": "userID",
        "name": "userName"
    },
    "chat": {
        "id": "chatID"
    },
    "text": "你好"
}
```

#### 文件消息

```
{
    "mid": "消息ID 发送时该字段省略",
    "type": "file",
    "from": {
        "id": "userID",
        "name": "userName"
    },
    "chat": {
        "id": "chatID"
    },
    "file": {
        "name": "文件名",
        "src": "文件链接地址",
        "size" 0 // 大小
    }
}
```

#### 图片消息

```
{
    "mid": "消息ID 发送时该字段省略",
    "type": "image",
    "from": {
        "id": "userID",
        "name": "userName"
    },
    "chat": {
        "id": "chatID"
    },
    "image": {
        "src": "文件链接地址",
        "width": 1, // 宽
        "height": 1, // 高
        "size" 0 // 大小
    }
}
```

#### 请求对话

```
{
    "mid": "消息ID 发送时该字段省略",
    "type": "event",
    "from": {
        "id": "userID",
        "name": "userName"
    },
    "event": {
        "chat": { // 发送时该字段省略 对话ID由服务端在请求对话时就创建好
            "id": "对话ID"
        },
        "uids": ["userID"]
    }
}
```

#### 加入对话

```
{
    "mid": "消息ID 发送时该字段省略",
    "type": "event",
    "from": {
        "id": "userID",
        "name": "userName"
    },
    "event": {
        "chat": {
            "id": "对话ID"
        }
    }
}
```