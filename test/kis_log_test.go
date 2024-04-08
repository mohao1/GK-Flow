package test

import (
	"context"
	"kis-folw/log"
	"testing"
)

func TestKisLogger(t *testing.T) {
	ctx := context.Background()

	log.Logger().DebugF("TestKisLogger DebugF")
	log.Logger().ErrorF("TestKisLogger ErrorF")
	log.Logger().InfoF("TestKisLogger InfoF")

	log.Logger().DebugFX(ctx, "TestKisLogger DebugFX")
	log.Logger().ErrorFX(ctx, "TestKisLogger ErrorFX")
	log.Logger().InfoFX(ctx, "TestKisLogger InfoFX")

}
