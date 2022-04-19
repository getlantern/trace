package trace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTracing(t *testing.T) {
	tracer := NewTracer("test")

	requireCurrentEquals := func(expected context.Context) {
		current := CurrentContext()
		currentWrapper, ok := current.(*contextWrapper)
		require.True(t, ok)
		require.Equal(t, expected, currentWrapper.Context)
	}

	topCtx, top := tracer.Start(context.Background(), "top")
	requireCurrentEquals(topCtx)

	twoCtx, two := tracer.Start(topCtx, "two")
	requireCurrentEquals(twoCtx)

	threeCtx, three := tracer.Start(CurrentContext(), "three")
	requireCurrentEquals(threeCtx)

	fourCtx, four := tracer.Continue("four")
	requireCurrentEquals(fourCtx)

	four.End()
	requireCurrentEquals(threeCtx)

	three.End()
	requireCurrentEquals(twoCtx)

	two.End()
	requireCurrentEquals(topCtx)

	top.End()
	require.Equal(t, context.Background(), CurrentContext())
}
