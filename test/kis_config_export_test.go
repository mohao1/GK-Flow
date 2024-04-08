package test

import (
	"fmt"
	"kis-folw/common"
	"kis-folw/file"
	"kis-folw/kis"
	"kis-folw/test/caas"
	"kis-folw/test/faas"
	"testing"
)

func TestConfigExportYaml(*testing.T) {
	// 0. 注册Function 回调业务
	kis.Pool().FaaS("funcName1", faas.FuncDemo1Handler)
	kis.Pool().FaaS("funcName2", faas.FuncDemo2Handler)
	kis.Pool().FaaS("funcName3", faas.FuncDemo3Handler)

	// 0. 注册ConnectorInit 和 Connector 回调业务
	kis.Pool().CaaSInit("ConnName1", caas.InitConnDemo1)
	kis.Pool().CaaS("ConnName1", "funcName2", common.S, caas.CaasDemoHanler1)

	// 1. 加载配置文件并构建Flow
	if err := file.ConfigImportYaml("D:/kis-flow/kis-flow/test/load_conf/"); err != nil {
		panic(err)
	}

	fmt.Println("-------")

	// 2. 讲构建的内存KisFlow结构配置导出的文件当中
	flows := kis.Pool().GetFlows()
	for _, flow := range flows {
		if err := file.ConfigExportYaml(flow, "D:/kis-flow/kis-flow/test/export_conf/"); err != nil {
			panic(err)
		}
	}
}
