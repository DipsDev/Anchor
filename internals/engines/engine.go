type EngineExecutionResult string

type AnchorEngine[T any] interface {
	Parse(service &ServiceConfig) T
	Execute() void
}