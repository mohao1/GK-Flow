package function

import (
	"context"
	"kis-folw/kis"
	"kis-folw/log"
)

type KisFunctionC struct {
	BaseFunction
}

func (f *KisFunctionC) Call(ctx context.Context, flow kis.Flow) error {
	log.Logger().InfoF("KisFunctionC, flow = %+v\n", flow)

	// 通过KisPool 路由到具体的执行计算Function中
	if err := kis.Pool().CallFunction(ctx, f.Config.FName, flow); err != nil {
		log.Logger().ErrorFX(ctx, "Function Called Error err = %s\n", err)
		return err
	}

	return nil
}

func NewKisFunctionC() kis.Function {
	f := new(KisFunctionC)

	// 初始化metaData
	f.metaData = make(map[string]interface{})

	return f
}
