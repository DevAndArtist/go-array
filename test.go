package main

import . "fmt"
import . "./array"

type A struct{}

func init() {

	RegisterType(A{})
	RegisterType(&A{})
}

func main() {

	var array = ArrayOfType("int")
	array.Append(1)
	array.Append(2)
	array.Append(3)
	Println(array)
	array.InsertAtIndex(42, 0)
	array.InsertAtIndex(99, 10)
	Println(array)
	Println(array.RemoveFirst())
	Println(array.RemoveLast())
	Println(array)
	array.RemoveAll()
	array.RemoveFirst()
	array.RemoveLast()
	Println(array)
	array.Append(44)
	Println(array)
	array.SetAtIndex(55, 0)
	Println(array)
	Println(array.Count())
	Println(array.IsEmpty())
	Println(array.ElementAtIndex(0))
	Println(array.ElementAtIndex(-1))
	Println(array.ElementAtIndex(1))
	Println(array.IndexForElement(55))
	Println(array.IndexForElement(-55))
	array.Remove(55)
	Println(array)

	Println()
	Println(IsTypeRegistered("typeName"))
	Println(IsTypeRegistered("A"))
	var pointer = new(A)
	Println(GetTypeName(pointer))
	RegisteredTypes()[0] = "xxx"
	Println(RegisteredTypes())
	Println(IsTypeRegistered("*A"))
	Println()

	var newArray = new(Array)
	newArray.SetType("string")
	newArray.Append("newElement")
	Println(newArray)
	Println(newArray.ContainsElement("newElement"))
	Println(newArray.IndexForElement("newElement"))
	newArray.Remove("newElement")
	Println(newArray)
	newArray.InsertAtIndex("test", 0)
	newArray.InsertAtIndex("foo", -100)
	newArray.InsertAtIndex("boo", 0)
	Println(newArray)
}
