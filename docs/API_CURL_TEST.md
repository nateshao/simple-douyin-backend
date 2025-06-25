# simple-douyin-backend 接口 curl 验证文档

本文件收录了项目主链路常用接口的 curl 测试命令，便于开发者快速验证服务功能。

---

## 1. 用户注册

```sh
curl -X POST "http://127.0.0.1:8888/douyin/user/register/" \
  -d "username=testuser&password=testpass"
```
- 返回：`user_id`、`token`、`status_code`、`status_msg`

---

## 2. 用户登录

```sh
curl -X POST "http://127.0.0.1:8888/douyin/user/login/" \
  -d "username=testuser&password=testpass"
```
- 返回：`user_id`、`token`、`status_code`、`status_msg`

---

## 3. 用户信息查询

```sh
curl "http://127.0.0.1:8888/douyin/user/?user_id=1&token=YOUR_TOKEN"
```
- `user_id` 用注册/登录返回的
- `token` 用注册/登录返回的

---

## 4. 视频 Feed 流

```sh
curl "http://127.0.0.1:8888/douyin/feed/"
```
- 可加参数 `latest_time`、`token`（可选）

---

## 5. 发布视频

> 需要 multipart/form-data，假设有 test.mp4

```sh
curl -X POST "http://127.0.0.1:8888/douyin/publish/action/" \
  -F "token=YOUR_TOKEN" \
  -F "title=测试视频" \
  -F "data=@test.mp4"
```
- 返回：`status_code`、`status_msg`

---

## 6. 获取已发布视频列表

```sh
curl "http://127.0.0.1:8888/douyin/publish/list/?user_id=1&token=YOUR_TOKEN"
```

---

## 7. 点赞/取消点赞

```sh
curl -X POST "http://127.0.0.1:8888/douyin/favorite/action/" \
  -d "token=YOUR_TOKEN&video_id=1&action_type=1"
# action_type=1 点赞，action_type=2 取消点赞
```

---

## 8. 获取点赞列表

```sh
curl "http://127.0.0.1:8888/douyin/favorite/list/?user_id=1&token=YOUR_TOKEN"
```

---

## 9. 评论操作

```sh
curl -X POST "http://127.0.0.1:8888/douyin/comment/action/" \
  -d "token=YOUR_TOKEN&video_id=1&action_type=1&comment_text=你好"
# action_type=1 发表评论，action_type=2 删除评论
```

---

## 10. 获取评论列表

```sh
curl "http://127.0.0.1:8888/douyin/comment/list/?video_id=1&token=YOUR_TOKEN"
```

---

## 11. 关注/取关

```sh
curl -X POST "http://127.0.0.1:8888/douyin/relation/action/" \
  -d "token=YOUR_TOKEN&to_user_id=2&action_type=1"
# action_type=1 关注，action_type=2 取关
```

---

## 12. 获取关注列表

```sh
curl "http://127.0.0.1:8888/douyin/relation/follow/list/?user_id=1&token=YOUR_TOKEN"
```

---

## 13. 获取粉丝列表

```sh
curl "http://127.0.0.1:8888/douyin/relation/follower/list/?user_id=1&token=YOUR_TOKEN"
```

---

## 14. 获取好友列表

```sh
curl "http://127.0.0.1:8888/douyin/relation/friend/list/?user_id=1&token=YOUR_TOKEN"
```

---

## 15. 发送私信

```sh
curl -X POST "http://127.0.0.1:8888/douyin/message/action/" \
  -d "token=YOUR_TOKEN&to_user_id=2&action_type=1&content=你好"
```

---

## 16. 获取聊天记录

```sh
curl "http://127.0.0.1:8888/douyin/message/chat/?token=YOUR_TOKEN&to_user_id=2"
```

---

## 17. 健康检查 & 监控

```sh
curl "http://127.0.0.1:8888/health"
curl "http://127.0.0.1:8888/metrics"
```

---

> **注意事项：**
> - `YOUR_TOKEN` 替换为注册/登录返回的 token。
> - `user_id`、`video_id`、`to_user_id` 等参数请根据实际数据填写。
> - 发布视频需本地有测试视频文件。
> - 如遇 500 或错误响应，请查看服务日志。 