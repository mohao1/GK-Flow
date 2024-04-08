package faas

import (
	"context"
	"fmt"
	"kis-folw/kis"
	"kis-folw/serialize"
	"kis-folw/test/load_conf/proto"
)

type PrintStuAvgScoreIn struct {
	serialize.DefaultSerialize
	proto.StuAvgScore
}

type PrintStuAvgScoreOut struct {
	serialize.DefaultSerialize
}

func PrintStuAvgScore(ctx context.Context, flow kis.Flow, rows []*PrintStuAvgScoreIn) error {

	for _, row := range rows {
		fmt.Printf("stuid: [%+v], avg score: [%+v]\n", row.StuId, row.AvgScore)
	}

	return nil
}
