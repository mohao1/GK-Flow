package faas

import (
	"context"
	"fmt"
	"kis-folw/kis"
)

// type FaaS func(context.Context, Flow) error

func FuncDemo3Handler(ctx context.Context, flow kis.Flow) error {
	fmt.Println("---> Call funcName3Handler ----")

	fmt.Printf("Params = %+v\n", flow.GetFuncParamAll())

	for _, row := range flow.Input() {
		str := fmt.Sprintf("In FuncName = %s, FuncId = %s, row = %s", flow.GetThisFuncConf().FName, flow.GetThisFunction().GetId(), row)
		fmt.Println(str)
	}

	return nil
}
