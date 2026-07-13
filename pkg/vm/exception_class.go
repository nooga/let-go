package vm

// ExceptionClass is a ValueType representing a JVM-style exception class.
// Instances form the java.lang.* hierarchy through the parent pointer; the
// walk in pkg/rt's directTypeAncestors follows it for instance? and typed
// catch dispatch. Exception VALUES are *ExInfo tagged with one of these
// classes; there is one Go struct for all exception values, not one per
// class.
type ExceptionClass struct {
	name   string
	parent *ExceptionClass
}

func (c *ExceptionClass) String() string  { return c.name }
func (c *ExceptionClass) Type() ValueType { return TypeType }
func (c *ExceptionClass) Unbox() any      { return nil }
func (c *ExceptionClass) Name() string    { return c.name }
func (c *ExceptionClass) Box(bare any) (Value, error) {
	return NIL, NewTypeError(bare, "can't be boxed as", c)
}
func (c *ExceptionClass) Parent() *ExceptionClass { return c.parent }

// ExceptionClasses lists every canonical class in registration order so
// pkg/rt can Def names and constructors by iteration.
var ExceptionClasses []*ExceptionClass

func newExceptionClass(name string, parent *ExceptionClass) *ExceptionClass {
	c := &ExceptionClass{name: name, parent: parent}
	ExceptionClasses = append(ExceptionClasses, c)
	return c
}

var (
	ClassThrowable            = newExceptionClass("java.lang.Throwable", nil)
	ClassError                = newExceptionClass("java.lang.Error", ClassThrowable)
	ClassAssertionError       = newExceptionClass("java.lang.AssertionError", ClassError)
	ClassException            = newExceptionClass("java.lang.Exception", ClassThrowable)
	ClassIOException          = newExceptionClass("java.io.IOException", ClassException)
	ClassInterruptedException = newExceptionClass("java.lang.InterruptedException", ClassException)
	ClassRuntimeException     = newExceptionClass("java.lang.RuntimeException", ClassException)
	ClassArithmeticException  = newExceptionClass("java.lang.ArithmeticException", ClassRuntimeException)
	ClassClassCastException   = newExceptionClass("java.lang.ClassCastException", ClassRuntimeException)
	ClassIllegalArgument      = newExceptionClass("java.lang.IllegalArgumentException", ClassRuntimeException)
	ClassNumberFormat         = newExceptionClass("java.lang.NumberFormatException", ClassIllegalArgument)
	ClassIllegalState         = newExceptionClass("java.lang.IllegalStateException", ClassRuntimeException)
	ClassIndexOutOfBounds     = newExceptionClass("java.lang.IndexOutOfBoundsException", ClassRuntimeException)
	ClassNullPointer          = newExceptionClass("java.lang.NullPointerException", ClassRuntimeException)
	ClassUnsupportedOperation = newExceptionClass("java.lang.UnsupportedOperationException", ClassRuntimeException)
)
