package log

import (
	"context"
	"fmt"
)

// KisDefaultLog 默认提供的日志对象
type KisDefaultLog struct {
}

func (log *KisDefaultLog) InfoFX(ctx context.Context, str string, v ...interface{}) {
	fmt.Println(ctx)
	fmt.Printf(str, v...)
}

func (log *KisDefaultLog) ErrorFX(ctx context.Context, str string, v ...interface{}) {
	fmt.Println(ctx)
	fmt.Printf(str, v...)
}

func (log *KisDefaultLog) DebugFX(ctx context.Context, str string, v ...interface{}) {
	fmt.Println(ctx)
	fmt.Printf(str, v...)
}

func (log *KisDefaultLog) InfoF(str string, v ...interface{}) {
	fmt.Printf(str, v...)
}

func (log *KisDefaultLog) ErrorF(str string, v ...interface{}) {
	fmt.Printf(str, v...)
}

func (log *KisDefaultLog) DebugF(str string, v ...interface{}) {
	fmt.Printf(str, v...)
}

func init() {
	// 如果没有设置Logger, 则启动时使用默认的kisDefaultLog对象
	if kisLog == nil {
		SetLogger(&KisDefaultLog{})
	}

}
