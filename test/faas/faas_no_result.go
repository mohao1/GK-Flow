package faas

import (
	"context"
	"fmt"
	"kis-folw/kis"
)

// type FaaS func(context.Context, Flow) error

func NoResultFuncHandler(ctx context.Context, flow kis.Flow) error {
	fmt.Println("---> Call NoResultFuncHandler ----")

	for _, row := range flow.Input() {
		str := fmt.Sprintf("In FuncName = %s, FuncId = %s, row = %s", flow.GetThisFuncConf().FName, flow.GetThisFunction().GetId(), row)
		fmt.Println(str)
	}

	//return flow.Next()

	return flow.Next(kis.ActionForceEntryNext)
}
