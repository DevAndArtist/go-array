## Go-Array

A type-safe and thread-safe array type in `Go`. If you are tired of `slices` and the `Go list.list` is not powerful enough for you, than you are at the right place.

Custom types should be always registred first.

**Registry functions:**
```go
GetTypeName(value interface{}) string 
```
```go
IsTypeRegistered(typeName string) bool
```
```go
RegisterType(value interface{})
```
```go
RegisteredTypes() []string
```

**Types:**
```go
type Element interface{}

type Array struct {
  // contains filtered or unexported fields
}
```
  
**Array constructor:**
```go
ArrayOfType(typeName string) *Array
```
**Array methods:**
```go
SetType(typeName string) *Array
```
```go
Type() string
```
```go
Append(newElement Element)
```
```go
InsertAtIndex(newElement Element, index int)
```
```go
RemoveAtIndex(index int) Element
```
```go
Remove(element Element)
```
```go
RemoveFirst() Element
```
```go
RemoveLast() Element
```
```go
RemoveAll()
```
```go
Count() int
```
```go
IsEmpty() bool
```
```go
ContainsElement(element Element) bool
```
```go
IndexForElement(element Element) int
```
```go
ElementAtIndex(index int) Element 
```
```go
FirstElement() Element
```
```go
LastElement() Element 
```
```go
SetAtIndex(element Element, index int)
```
```go
String() string
```

**TODO-List:** 
  * more examples
  * deeper explanation
  * delegation callbacks (to make the array observable)
  * code enhancement

## License

BSD license. See the `LICENSE` file for details.
