Adds support for goroutine-based propagation of contexts to OpenTelemetry tracing.

This library is inspired by https://github.com/tylerb/gls and https://github.com/jtolds/gls.

## Example

```go
topCtx, top := tracer.Start(context.Background(), "top")

// Two will be a child span of top
twoCtx, two := tracer.Start(CurrentContext(), "two")

// Three will be a child span of two
threeCtx, three := tracer.Continue("three")
```