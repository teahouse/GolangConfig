# GolangConfig
golang config

golang读取配置，很多go配置都需要预定义结构体去读取，很不方便配置增\减\更新！
所以这里把配置和代码分开，代码比较乱。
## 用法
LoadPathConfig(filepath)  
# return map[string]string
# 运行测试
go run .\main.go .\LoadConfig.go .\test.config
## 测试配置
// 单行注释
/* 多行注释
    abc
    *b/
*/

include("./test1.config");  //include作为关键字，包含另一个本地配置，路径为当前相对路径

logpath = "log11"

logpath1 = 

logger = loglevel + "/game.log"	// 引擎日志输出文件

log_dailyrotate = true				// 是否按天切分日志

loglevel = "debug"					// 日志级别: debug/trace/info/warn/fatal

daemon = logpath + "/game.pid"

start = "app/main"					// 启动脚本

//多行字符串，字符串只能用双引号，不能用单引号

start = "fabdd\
dd\"fd+sf\
哈哈"  


