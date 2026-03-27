package vm

import "fmt"

// ExInfo is an exception value carrying a message and data map.
// Equivalent to Clojure's ExceptionInfo.
type ExInfo struct {
	message string
	data    *PersistentMap
	cause   error
}

func NewExInfo(message string, data *PersistentMap, cause error) *ExInfo {
	return &ExInfo{message: message, data: data, cause: cause}
}

func (e *ExInfo) Type() ValueType    { return ExInfoType }
func (e *ExInfo) Unbox() interface{} { return e }
func (e *ExInfo) String() string {
	return fmt.Sprintf("#error {:message %q, :data %s}", e.message, e.data.String())
}
func (e *ExInfo) Error() string        { return e.message }
func (e *ExInfo) Message() string      { return e.message }
func (e *ExInfo) Data() *PersistentMap { return e.data }
func (e *ExInfo) Cause() error         { return e.cause }

type theExInfoType struct{}

func (t *theExInfoType) String() string     { return t.Name() }
func (t *theExInfoType) Type() ValueType    { return TypeType }
func (t *theExInfoType) Unbox() interface{} { return nil }
func (t *theExInfoType) Name() string       { return "let-go.lang.ExceptionInfo" }
func (t *theExInfoType) Box(bare interface{}) (Value, error) {
	return NIL, NewTypeError(bare, "can't be boxed as", t)
}

var ExInfoType *theExInfoType = &theExInfoType{}
