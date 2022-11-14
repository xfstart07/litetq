# 单机轻量的任务队列

可以将任务放在 kv 存储，内存等存储方式的任务队列。

图:

![](./litetq.excalidraw)

## 内存版

内存版本，任务消息存储在内存中。不能进行持久化。

* [x] - 内存存储

## 功能列表

* [ ] - 多种 kv 存储的方式(leveldb, bbolt)
* [x] - 任务重试
* [ ] - 支持 HTTP，GRPC 方式提交任务

## 场景描述

TODO
