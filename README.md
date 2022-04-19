Adds support for goroutine-based propagation of contexts to OpenTelemetry tracing.

This library is inspired by https://github.com/tylerb/gls and https://github.com/jtolds/gls.

## Example

```go
topCtx, top := tracer.Start(context.Background(), "top")

// two will be a child span of top
_, two := tracer.Start(topCtx, "two")

// three will be a chlid span of two
_, three := tracer.Start(CurrentContext(), "three")

// four will be a child span of three
_, four := tracer.Continue("four")
```