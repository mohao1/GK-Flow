package test

import (
	"context"
	"kis-folw/common"
	"kis-folw/config"
	"kis-folw/file"
	"kis-folw/flow"
	"kis-folw/kis"
	"kis-folw/test/faas"
	"kis-folw/test/load_conf/proto"
	"testing"
)

func TestAutoInjectParamWithConfig(t *testing.T) {
	ctx := context.Background()

	// 0. 注册Function 回调业务
	kis.Pool().FaaS("AvgStuScore", faas.AvgStuScore)
	kis.Pool().FaaS("PrintStuAvgScore", faas.PrintStuAvgScore)

	// 0. 注册ConnectorInit 和 Connector 回调业务
	//kis.Pool().CaaSInit("ConnName1", caas.InitConnDemo1)
	//kis.Pool().CaaS("ConnName1", "AvgStuScore", common.S, caas.CaasDemoHanler1)

	// 1. 加载配置文件并构建Flow
	if err := file.ConfigImportYaml("D:/kis-flow/kis-flow/test/load_conf/"); err != nil {
		panic(err)
	}

	// 2. 获取Flow
	flow1 := kis.Pool().GetFlow("StuAvg")

	// 3. 提交原始数据
	_ = flow1.CommitRow(&faas.AvgStuScoreIn{
		StuScores: proto.StuScores{
			StuId:  100,
			Score1: 1,
			Score2: 2,
			Score3: 3,
		},
	})
	_ = flow1.CommitRow(`{"stu_id":101}`)
	_ = flow1.CommitRow(faas.AvgStuScoreIn{
		StuScores: proto.StuScores{
			StuId:  100,
			Score1: 1,
			Score2: 2,
			Score3: 3,
		},
	})

	// 4. 执行flow1
	if err := flow1.Run(ctx); err != nil {
		panic(err)
	}
}

func TestAutoInjectParam(t *testing.T) {
	ctx := context.Background()

	kis.Pool().FaaS("AvgStuScore", faas.AvgStuScore)
	kis.Pool().FaaS("PrintStuAvgScore", faas.PrintStuAvgScore)

	source1 := config.KisSource{
		Name: "Test",
		Must: []string{},
	}

	avgStuScoreConfig := config.NewFuncConfig("AvgStuScore", common.C, &source1, nil)
	if avgStuScoreConfig == nil {
		panic("AvgStuScore is nil")
	}

	printStuAvgScoreConfig := config.NewFuncConfig("PrintStuAvgScore", common.C, &source1, nil)
	if printStuAvgScoreConfig == nil {
		panic("printStuAvgScoreConfig is nil")
	}

	myFlowConfig1 := config.NewFlowConfig("cal_stu_avg_score", common.FlowEnable)

	flow1 := flow.NewKisFlow(myFlowConfig1)

	// 4. 拼接Functioin 到 Flow 上
	if err := flow1.Link(avgStuScoreConfig, nil); err != nil {
		panic(err)
	}
	if err := flow1.Link(printStuAvgScoreConfig, nil); err != nil {
		panic(err)
	}

	// 3. 提交原始数据
	_ = flow1.CommitRow(&faas.AvgStuScoreIn{
		StuScores: proto.StuScores{
			StuId:  100,
			Score1: 1,
			Score2: 2,
			Score3: 3,
		},
	})
	_ = flow1.CommitRow(`{"stu_id":101}`)
	_ = flow1.CommitRow(faas.AvgStuScoreIn{
		StuScores: proto.StuScores{
			StuId:  100,
			Score1: 1,
			Score2: 2,
			Score3: 3,
		},
	})

	// 4. 执行flow1
	if err := flow1.Run(ctx); err != nil {
		panic(err)
	}
}
