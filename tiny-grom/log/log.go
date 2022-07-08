package log

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
)

/*

log 标准库没有日志分级，不打印文件和行号，这就意味着我们很难定位具体发生问题的位置

这个简易的 log 具备以下特性
- 支持日志分级 info error disabled 三级
- 不同层级日志显示的时候使用不同的颜色进行区分
- 显示打印日志代码对应的文件名和行号

*/

var (
	errorLog = log.New(os.Stdout, "\033[31m[error]\033[0m ", log.LstdFlags|log.Lshortfile) // 红色
	infoLog  = log.New(os.Stdout, "\033[34m[info ]\033[0m ", log.LstdFlags|log.Lshortfile) // 蓝色
	loggers  = []*log.Logger{errorLog, infoLog}
	mu       sync.Mutex
)

// log methods 暴露四个方法：
var (
	Error  = errorLog.Println
	Errorf = errorLog.Printf
	Info   = infoLog.Println
	Infof  = infoLog.Printf
)

// log levels 设置日志的层级
const (
	InfoLevel = iota
	ErrorLevel
	Disabled
)

// SetLevel controls log level
func SetLevel(level int) {
	mu.Lock()
	defer mu.Unlock()

	for _, logger := range loggers {
		logger.SetOutput(os.Stdout) // 通过 setOutput 来控制日志是否打印
	}

	if ErrorLevel < level {
		errorLog.SetOutput(ioutil.Discard)
	}
	if InfoLevel < level {
		infoLog.SetOutput(ioutil.Discard)
	}
}
