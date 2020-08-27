 # rc 

Remember the last command of the terminal and its frequency 

### 配置

```
	historyPath = "/Users/edz/.local/share/fish/fish_history"
	cachePath   = "/var/rc/cache.log"
```
### 依赖

`/usr/local/bin/fish`

### changelog

#### 2020-08-11

* 记录命令以及其记录频率，按照频率大小从高到低显示命令
* 增加命令更新时间,以table形式显示命令列表

#### 2020-08-12

* 记录命令时候增加备注信息
* 增加记录删除命令

#### 2020-08-15

* 增加历史命令执行功能
* 增加编辑缓存命令
* 增加将命令赋值到粘贴板

#### 2020-08-27
* 修改命令用bash执行

### todo

* [ ] 配置文件
