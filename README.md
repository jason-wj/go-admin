## go-admin
结合项目经验以及网上优秀的开源项目，整合出一套比较通用的后台开发框架，首要目标是满足个人需要
当前v1.0.0可以满足基本需要，整体能跑通了。

## 项目地址
[github go-admin](https://github.com/jason-wj/go-admin)  
[gitee go-admin](https://gitee.com/jason-wj/go-admin)

## 当前进度
```text
jwt认证   -- 【已完成】  
公司、岗位、角色、管理员  -- 【已完成】  
参数管理    -- 【已完成】  
字典管理     -- 【已完成】  
内容管理(文章、通知)     -- 【已完成】  
APP版本管理(链接或OSS存储app)     -- 【已完成】  
登录日志、操作日志   -- 【已完成】
服务监控    -- 【已完成】
仅支持Mysql  -- 【已完成】   可见的将来，都只考虑mysql
代码生成    -- 【改造中】
支持插件化   -- 【改造中】

数据库开放    -- 【未完成】
管理页面开放    -- 【未完成】  数据库暂时不开放，等页面开房了再一并放出

废弃swagger接口管理   -- 【已完成】  理由：个人不喜欢，宁愿用第三方接口管理工具，如：yapi
```

## 部署数据库
数据库都在根目录的`db.zip`，直接自己建库`go-admin`，并导入数据库即可

## web
项目根目录的`web`目录为前端页面，启动：
```shell
npm i --force --registry=https://registry.npm.taobao.org
npm run serve
```

### 编译

开始编译：
```shell
# 当前
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o go-admin-api main.go
```

### 启动
```shell
nohup ./admin server &
```

## 开发注意
1. 使用自定义配置启动命令：go-admin server -c=config/settings.dev.yml
2. 检测并将接口加入到sys_api命令：go-admin server -a=true


## 感谢
[go-admin-team](https://github.com/go-admin-team)  
[flipped-aurora](https://github.com/flipped-aurora)
