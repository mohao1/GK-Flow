package faas

import (
	"context"
	"kis-folw/kis"
	"kis-folw/serialize"
	"kis-folw/test/load_conf/proto"
)

type AvgStuScoreIn struct {
	serialize.DefaultSerialize
	proto.StuScores
}

type AvgStuScoreOut struct {
	serialize.DefaultSerialize
	proto.StuAvgScore
}

// AvgStuScore(FaaS) 计算学生平均分
func AvgStuScore(ctx context.Context, flow kis.Flow, rows []*AvgStuScoreIn) error {
	for _, row := range rows {
		avgScore := proto.StuAvgScore{
			StuId:    row.StuId,
			AvgScore: float64(row.Score1+row.Score2+row.Score3) / 3,
		}
		// 提交结果数据
		_ = flow.CommitRow(AvgStuScoreOut{StuAvgScore: avgScore})
	}

	return nil
}
