//
// Copyright (c) 2016, Adrian Zubarev (alias DevAndArtist)
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// * Redistributions of source code must retain the above copyright notice, this
// list of conditions and the following disclaimer.
//
// * Redistributions in binary form must reproduce the above copyright notice,
// this list of conditions and the following disclaimer in the documentation
// and/or other materials provided with the distribution.
//
// * Neither the name of the copyright holder nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//

package array

import "reflect"
import "strings"
import "sync"
import "log"
import "fmt"

//==============--------------------------------------------==============//
//==============------- registry type container/guard ------==============//
//==============--------------------------------------------==============//
var types = []string{"bool", "uint8", "complex128", "complex64", "float32", "float64", "int", "int16",
	"int32", "int64", "int8", "int32", "string", "uint", "uint16", "uint32", "uint64", "uint8"}
var typeGuard = new(sync.Mutex)

//==============--------------------------------------------==============//
//==============--------- type registry functions ----------==============//
//==============--------------------------------------------==============//
func GetTypeName(value interface{}) string {

	// parse type names to support pointers (package name will be ignored)
	var typeName = reflect.TypeOf(value).String()
	var slitedType = strings.SplitAfter(typeName, ".")

	if len(slitedType) == 2 {

		var typeName = ""
		for i := 0; reflect.DeepEqual(slitedType[0][i], "*"[0]); i++ {

			typeName += "*"
		}
		return typeName + slitedType[1]
	}

	return typeName
}

func IsTypeRegistered(typeName string) bool {

	typeGuard.Lock()
	for _, registeredType := range types {

		if equalTypes(registeredType, typeName) {

			typeGuard.Unlock()
			return true
		}
	}
	typeGuard.Unlock()
	return false
}

func RegisterType(value interface{}) {

	var typeName = GetTypeName(value)

	if !IsTypeRegistered(typeName) {

		typeGuard.Lock()
		types = append(types, typeName)
		typeGuard.Unlock()
	}
}

func RegisteredTypes() []string {

	return append([]string{}, types...) // do not just return types, because then it could be modified
}

//==============--------------------------------------------==============//
//==============------------------ types -------------------==============//
//==============--------------------------------------------==============//
type Element interface{}

type Array struct {
	elements    []Element  // Element container
	elementType string     // Element type
	guard       sync.Mutex // Cuncurrency guard
}

//==============--------------------------------------------==============//
//==============------------- array constructor ------------==============//
//==============--------------------------------------------==============//
func ArrayOfType(typeName string) *Array {

	var array = new(Array)
	return array.SetType(typeName)
}

//==============--------------------------------------------==============//
//==============--------- array type getter/setter ---------==============//
//==============--------------------------------------------==============//
func (array *Array) SetType(typeName string) *Array {

	array.guard.Lock()
	if equalTypes(array.elementType, "") {

		if !IsTypeRegistered(typeName) {

			array.guard.Unlock() // unlock before throw
			log.Fatalf("FATAL: type <%s> is not registered in module array yet\nUSAGE: RegisterType(instanceOfNotRegisteredType)\n", typeName)
		}
		array.elementType = typeName
		array.guard.Unlock()

	} else {

		array.guard.Unlock() // unlock before throw
		log.Fatalf("FATAL: type of array <%p> is already set to <%s>\n", array, array.elementType)
	}
	return array
}

func (array *Array) Type() string {

	array.guard.Lock()
	if equalTypes(array.elementType, "") {

		array.guard.Unlock() // unlock before throw
		log.Fatalf("FATAL: type for array <%p> is not set yet\n", array)
	}
	var elementType = array.elementType
	array.guard.Unlock()
	return elementType
}

//==============--------------------------------------------==============//
//==============-------------- elements adder --------------==============//
//==============--------------------------------------------==============//
func (array *Array) Append(newElement Element) {

	array.typeCheck(newElement)

	array.guard.Lock()
	array.elements = append(array.elements, newElement)
	array.guard.Unlock()
}

func (array *Array) InsertAtIndex(newElement Element, index int) {

	array.typeCheck(newElement)

	array.guard.Lock()
	var anIndex = 0

	if index > anIndex {

		anIndex = index
	}

	if anIndex > len(array.elements) {

		array.elements = append(array.elements, newElement)

	} else {

		var lhs = array.elements[:anIndex]
		var rhs = append([]Element{newElement}, array.elements[anIndex:]...)
		array.elements = append(lhs, rhs...)
	}
	array.guard.Unlock()
}

//==============--------------------------------------------==============//
//==============-------------- element remover -------------==============//
//==============--------------------------------------------==============//
func (array *Array) RemoveAtIndex(index int) Element {

	array.guard.Lock()
	if index < 0 || index >= len(array.elements) {

		array.guard.Unlock()
		return nil
	}

	var element = array.elements[index]

	array.elements = append(array.elements[:index], array.elements[index+1:]...)
	array.guard.Unlock()

	return element
}

func (array *Array) Remove(element Element) {

	array.typeCheck(element)

	array.guard.Lock()
	var index = -1
	for anIndex, anElement := range array.elements {

		if reflect.DeepEqual(element, anElement) {

			index = anIndex
			break // we found an equal element
		}
	}

	if index >= 0 && index < len(array.elements) {

		array.elements = append(array.elements[:index], array.elements[index+1:]...)
	}
	array.guard.Unlock()
}

func (array *Array) RemoveFirst() Element {

	return array.RemoveAtIndex(0)
}

func (array *Array) RemoveLast() Element {

	return array.RemoveAtIndex(array.Count() - 1)
}

func (array *Array) RemoveAll() {

	array.guard.Lock()
	array.elements = []Element{}
	array.guard.Unlock()
}

//==============--------------------------------------------==============//
//==============-------------- element counter -------------==============//
//==============--------------------------------------------==============//
func (array *Array) Count() int {

	array.guard.Lock()
	var count = len(array.elements)
	array.guard.Unlock()

	return count
}

func (array *Array) IsEmpty() bool {

	return array.Count() == 0
}

//==============--------------------------------------------==============//
//==============---------- element/index searcher ----------==============//
//==============--------------------------------------------==============//
func (array *Array) ContainsElement(element Element) bool {

	array.typeCheck(element)

	array.guard.Lock()
	for _, anElement := range array.elements {

		if reflect.DeepEqual(anElement, element) {

			array.guard.Unlock()
			return true
		}
	}
	array.guard.Unlock()
	return false
}

func (array *Array) IndexForElement(element Element) int {

	array.typeCheck(element)

	array.guard.Lock()
	for index, anElement := range array.elements {

		if reflect.DeepEqual(element, anElement) {

			array.guard.Unlock()
			return index
		}
	}
	array.guard.Unlock()

	return -1
}

//==============--------------------------------------------==============//
//==============-------------- element getter --------------==============//
//==============--------------------------------------------==============//
func (array *Array) ElementAtIndex(index int) Element {

	array.guard.Lock()
	var count = len(array.elements)

	if index < 0 || index >= count {

		array.guard.Unlock()
		return nil
	}

	var element = array.elements[index]
	array.guard.Unlock()

	return element
}

func (array *Array) FirstElement() Element {

	return array.ElementAtIndex(0)
}

func (array *Array) LastElement() Element {

	return array.ElementAtIndex(array.Count() - 1)
}

//==============--------------------------------------------==============//
//==============-------------- element setter --------------==============//
//==============--------------------------------------------==============//
func (array *Array) SetAtIndex(element Element, index int) {

	array.guard.Lock()

	if index < 0 || index >= len(array.elements) {

		array.guard.Unlock()
		log.Fatalf("FATAL: index %d out of range for array <%p> with elements: %v\n", index, array, array.elements)
	}

	array.elements[index] = element

	array.guard.Unlock()
}

//==============--------------------------------------------==============//
//==============-------------- array printer ---------------==============//
//==============--------------------------------------------==============//
func (array *Array) String() string {

	return fmt.Sprintf("Array <%p> of type <%s> with elements: %v", array, array.elementType, array.elements)
}

//==============--------------------------------------------==============//
//==============-------------- private helper --------------==============//
//==============--------------------------------------------==============//
func equalTypes(type1 string, type2 string) bool {

	return strings.Compare(type1, type2) == 0
}

func (array *Array) typeCheck(element Element) {

	var elementType = GetTypeName(element)

	array.guard.Lock()
	if equalTypes(array.elementType, "") {

		array.guard.Unlock() // unlock before throw
		log.Fatalf("FATAL: type for array <%p> is not set yet\n", array)

	} else if !equalTypes(elementType, array.elementType) {

		array.guard.Unlock() // unlock before throw
		log.Fatalf("FATAL: array <%p> of type <%s> can not procced with an element of type <%s>\n", array, array.elementType, elementType)
	}
	array.guard.Unlock()
}
