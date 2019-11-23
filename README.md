# GolangConfig
golang config

golang读取配置，很多go配置都需要预定义结构体去读取，很不方便配置热更吧！
所以这里把配置和代码分开，代码比较乱。
## 用法
LoadPathConfig(filepath)
## 测试配置
// 单行注释
/* 多行注释
    abc
    *b/
*/
include("./test1.config");  //包含另一个本地配置，路径为当前相对路径

logpath = "log11"
logpath1 = 
logger = loglevel + "/skynet.log"	// 引擎日志输出文件
log_dailyrotate = true				// 是否按天切分日志
loglevel = "debug"					// 日志级别: debug/trace/info/warn/fatal
daemon = logpath + "/skynet.pid"
start = "app/main"					// 启动脚本

start = "fabdd\
dd\"fd+sf\
哈哈"
