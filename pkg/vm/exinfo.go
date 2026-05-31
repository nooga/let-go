package vm

import "fmt"

// ExInfo is an exception value carrying a message and data map.
// Equivalent to Clojure's ExceptionInfo.
type ExInfo struct {
	message string
	data    *PersistentMap
	cause   error
	meta    Value
}

func NewExInfo(message string, data *PersistentMap, cause error) *ExInfo {
	return &ExInfo{message: message, data: data, cause: cause}
}

func (e *ExInfo) Type() ValueType { return ExInfoType }
func (e *ExInfo) Unbox() any      { return e }
func (e *ExInfo) String() string {
	return fmt.Sprintf("#error {:message %q, :data %s}", e.message, e.data.String())
}
func (e *ExInfo) Error() string        { return e.message }
func (e *ExInfo) Message() string      { return e.message }
func (e *ExInfo) Data() *PersistentMap { return e.data }
func (e *ExInfo) Cause() error         { return e.cause }

// Meta / WithMeta implement IMeta. Type hints compile to runtime metadata in
// let-go, so medley's (.getMessage ^Throwable ex) reaches (with-meta ex ...);
// supporting metadata lets that hinted call work. WithMeta returns a copy that
// is still *ExInfo, preserving Type()==ExInfoType (so (instance? Throwable ex)
// stays true) and InvokeMethod (so .getMessage/.getCause still dispatch).
func (e *ExInfo) Meta() Value {
	if e.meta == nil {
		return NIL
	}
	return e.meta
}

func (e *ExInfo) WithMeta(m Value) Value {
	cp := *e
	cp.meta = m
	return &cp
}

// InvokeMethod implements Receiver so the JVM-interop method forms
// (.getMessage ex) / (.getCause ex) dispatch to an ex-info value. This lets
// libraries that branch on (instance? Throwable ex) and then call .getMessage
// — e.g. medley's ex-message/ex-cause on their :clj branch — work against
// let-go ex-info values. Mirrors the semantics of core's ex-message/ex-cause.
func (e *ExInfo) InvokeMethod(name Symbol, args []Value) (Value, error) {
	switch name {
	case "getMessage":
		return String(e.message), nil
	case "getCause":
		if cev, ok := e.cause.(*ExInfo); ok {
			return cev, nil
		}
		return NIL, nil
	}
	return NIL, fmt.Errorf("ExceptionInfo has no method %s", name)
}

type theExInfoType struct{}

func (t *theExInfoType) String() string  { return t.Name() }
func (t *theExInfoType) Type() ValueType { return TypeType }
func (t *theExInfoType) Unbox() any      { return nil }
func (t *theExInfoType) Name() string    { return "let-go.lang.ExceptionInfo" }
func (t *theExInfoType) Box(bare any) (Value, error) {
	return NIL, NewTypeError(bare, "can't be boxed as", t)
}

var ExInfoType *theExInfoType = &theExInfoType{}
