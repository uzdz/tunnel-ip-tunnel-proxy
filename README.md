## 用户限流策略

1. 通过当前请求TCP连接的授权信息或IP，获取到当前用户的唯一ID，如果获取不到，则表示未授权，直接返回。
2. 如果获取到了用户，表示可以允许代理当前连接，那么获取当前用户的连接限制，判断是否触发了限流。
3. 如果未获取到限流策略，表示无限流，直接代理转发。否则获得到连接限制，通过SQLite获取当前用户的访问数据。
4. 如果未达到限制，则转发请求，否则拒绝断开连接。

## 注意事项

1. 当处于下载状态时，所有的命令操作均阻塞等待下载完成，下载事件完成后放行命令。目前队列容量为20个，也就是说，当队列内超过20个命令时，会拒绝其他命令。

## 等待时间

位置 | 任务 | 沉睡秒数
--- | --- | ---
run.go#34 | VPS上报数据失败 | 5
run.go#74 | 下载时阻塞命令 | 5
m.go#27 | 刷新GUI面板 | 3
re.go#49 | 重新拨号时间间隔 | 取决服务器下发
run.go#57 | 等待关闭旧客户端 | 1 
download.go#44 | 定时检测是否有新版本 | 10