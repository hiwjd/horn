# horn doc

## 消息接口

### 地址

```
POST http://app.horn.com/api/message

消息体<见下方消息格式>
```

### 消息格式

#### 普通对话消息

```
{
    "mid": "消息ID 发送时该字段省略",
    "type": "text",
    "from": {
        "id": "userID",
        "name": "userName",
        "role": "user" // user:访客 staff:客服
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
        "name": "userName",
        "role": "user" // user:访客 staff:客服
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
        "name": "userName",
        "role": "user" // user:访客 staff:客服
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
        "name": "userName",
        "role": "user" // user:访客 staff:客服
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
        "name": "userName",
        "role": "user" // user:访客 staff:客服
    },
    "event": {
        "chat": {
            "id": "对话ID"
        }
    }
}
```

### 响应

```
{
    "code": 0, // 0:成功
    "mid": "b0fhhegeibm0m3kpetgg", // 消息ID
    "msg": "ok"
}
```

## 初始化接口

> 在接入推送服务前，需要调用初始化接口来获取之前的状态数据和推送服务的地址
> 后续收到的事件类型消息，通过判断消息ID是否大于状态版本version来决定是否响应

### 地址

```
GET http://app.horn.com/api/state/init?uid=UID&fp=FINGERPRINT
```

### 响应

```
{
    "code":0,
    "uid":"kefu001",
    "addr":"p1.horn.com:9001",
    "track_id":"2016110623570800kefu001aHc0IdFuWj4NJ6iB",
    "state":{
        "chats":[{
            "chat_id":"cqaM7Ja8i13eW8wiGVc620Gik",
            "gid":"",
            "creator":"FWaBA3y0dwwrn6mmgvjW6aQ",
            "kfid":"0",
            "user_num":"3",
            "state":"request",
            "created_at":"2016-11-05 17:05:01",
            "ended_at":"2016-11-05 17:05:01",
            "id":"cqaM7Ja8i13eW8wiGVc620Gik"
        }],
        "version":"b0f18ageibm0m3kpeqe0" // 状态版本
    }
}
```

## 上发访问信息

### 地址

```
POST http://app.horn.com/api/track

消息体<见下方>
```

### 消息体

```
{
    "uid": "用户ID",
    "fp": "指纹",
    "gid": "公司ID",
    "url": "当前访问地址",
    "title": "当前访问页面的标题",
    "referer": "来源页地址",
    "os": "操作系统信息",
    "browser": "浏览器信息",
    "ip": "IP"
}
```

### 响应

```
{
    "code": 0,
    "msg": "ok",
    "track_id": ""
}
```

## 消息推送

### 地址

#### longpolling
```
GET http://<推送服务地址>/pull?uid=UID&track_id=TRACKID
```

#### websocket
```
GET http://<推送服务地址>/ws?uid=UID&track_id=TRACKID
```

### 响应

> data字段内的消息格式请参看`消息接口`

```
{
    "code":0,
    "msg":"",
    "data":[{
        "type":"text",
        "t":{
            "t0":1478477956,
            "t1":1478477956
        },
        "mid":"b0fsh10eibm0m3kpetpg",
        "from":{
            "id":"FWaBA3y0dwwrn6mmgvjW6aQ",
            "name":"",
            "role":"user"
        },
        "chat":{
            "id":"YAbJu1BTEmAIS22SdIezLlWaA"
        },
        "text":"qqq"
    }]
}
```
