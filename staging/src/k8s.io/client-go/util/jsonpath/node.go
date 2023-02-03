/*
Copyright 2015 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package jsonpath

import "fmt"

// NodeType identifies the type of a parse tree node.
type NodeType int

// Type returns itself and provides an easy default implementation
func (t NodeType) Type() NodeType {
	return t
}

func (t NodeType) String() string {
	return NodeTypeName[t]
}

const (
	NodeText NodeType = iota
	NodeArray
	NodeList
	NodeField
	NodeParent
	NodeIdentifier
	NodeFilter
	NodeInt
	NodeFloat
	NodeWildcard
	NodeRecursive
	NodeUnion
	NodeBool
)

var NodeTypeName = map[NodeType]string{
	NodeText:       "NodeText",
	NodeArray:      "NodeArray",
	NodeList:       "NodeList",
	NodeField:      "NodeField",
	NodeParent:     "NodeParent",
	NodeIdentifier: "NodeIdentifier",
	NodeFilter:     "NodeFilter",
	NodeInt:        "NodeInt",
	NodeFloat:      "NodeFloat",
	NodeWildcard:   "NodeWildcard",
	NodeRecursive:  "NodeRecursive",
	NodeUnion:      "NodeUnion",
	NodeBool:       "NodeBool",
}

type Node interface {
	Type() NodeType
	String() string
}

// ListNode holds a sequence of nodes.
type ListNode struct {
	NodeType
	Nodes  []Node // The element nodes in lexical order.
	Parent interface{}
}

func newList(parent interface{}) *ListNode {
	return &ListNode{
		NodeType: NodeList,
		Parent:   parent,
	}
}

func (l *ListNode) append(n Node) {
	l.Nodes = append(l.Nodes, n)
}

func (l *ListNode) String() string {
	return l.Type().String()
}

// TextNode holds plain text.
type TextNode struct {
	NodeType
	Text   string // The text; may span newlines.
	Parent interface{}
}

func newText(text string, parent interface{}) *TextNode {
	return &TextNode{
		NodeType: NodeText,
		Text:     text,
		Parent:   parent,
	}
}

func (t *TextNode) String() string {
	return fmt.Sprintf("%s: %s", t.Type(), t.Text)
}

// FieldNode holds field of struct
type FieldNode struct {
	NodeType
	Value  string
	Parent interface{}
}

func newField(value string, parent interface{}) *FieldNode {
	return &FieldNode{
		NodeType: NodeField,
		Value:    value,
		Parent:   parent,
	}
}

func (f *FieldNode) String() string {
	return fmt.Sprintf("%s: %s", f.Type(), f.Value)
}

// IdentifierNode holds an identifier
type IdentifierNode struct {
	NodeType
	Name   string
	Parent interface{}
}

func newIdentifier(value string, parent interface{}) *IdentifierNode {
	return &IdentifierNode{
		NodeType: NodeIdentifier,
		Name:     value,
		Parent:   parent,
	}
}

func (f *IdentifierNode) String() string {
	return fmt.Sprintf("%s: %s", f.Type(), f.Name)
}

// ParamsEntry holds param information for ArrayNode
type ParamsEntry struct {
	Value   int
	Known   bool // whether the value is known when parse it
	Derived bool
}

// ArrayNode holds start, end, step information for array index selection
type ArrayNode struct {
	NodeType
	Params [3]ParamsEntry // start, end, step
	Parent interface{}
}

func newArray(params [3]ParamsEntry, parent interface{}) *ArrayNode {
	return &ArrayNode{
		NodeType: NodeArray,
		Params:   params,
		Parent:   parent,
	}
}

func (a *ArrayNode) String() string {
	return fmt.Sprintf("%s: %v", a.Type(), a.Params)
}

// FilterNode holds operand and operator information for filter
type FilterNode struct {
	NodeType
	Left     *ListNode
	Right    *ListNode
	Operator string
	Parent   interface{}
}

func newFilter(left, right *ListNode, operator string, parent interface{}) *FilterNode {
	return &FilterNode{
		NodeType: NodeFilter,
		Left:     left,
		Right:    right,
		Operator: operator,
		Parent:   parent,
	}
}

func (f *FilterNode) String() string {
	return fmt.Sprintf("%s: %s %s %s", f.Type(), f.Left, f.Operator, f.Right)
}

// IntNode holds integer value
type IntNode struct {
	NodeType
	Value  int
	Parent interface{}
}

func newInt(num int, parent interface{}) *IntNode {
	return &IntNode{
		NodeType: NodeInt,
		Value:    num,
		Parent:   parent,
	}
}

func (i *IntNode) String() string {
	return fmt.Sprintf("%s: %d", i.Type(), i.Value)
}

// FloatNode holds float value
type FloatNode struct {
	NodeType
	Value  float64
	Parent interface{}
}

func newFloat(num float64, parent interface{}) *FloatNode {
	return &FloatNode{
		NodeType: NodeFloat,
		Value:    num,
		Parent:   parent,
	}
}

func (i *FloatNode) String() string {
	return fmt.Sprintf("%s: %f", i.Type(), i.Value)
}

// WildcardNode means a wildcard
type WildcardNode struct {
	NodeType
	Parent interface{}
}

func newWildcard(parent interface{}) *WildcardNode {
	return &WildcardNode{
		NodeType: NodeWildcard,
		Parent:   parent,
	}
}

func (i *WildcardNode) String() string {
	return i.Type().String()
}

// RecursiveNode means a recursive descent operator
type RecursiveNode struct {
	NodeType
	Parent interface{}
}

func newRecursive(parent interface{}) *RecursiveNode {
	return &RecursiveNode{
		NodeType: NodeRecursive,
		Parent:   parent,
	}
}

func (r *RecursiveNode) String() string {
	return r.Type().String()
}

// UnionNode is union of ListNode
type UnionNode struct {
	NodeType
	Nodes  []*ListNode
	Parent interface{}
}

func newUnion(nodes []*ListNode, parent interface{}) *UnionNode {
	return &UnionNode{
		NodeType: NodeUnion,
		Nodes:    nodes,
		Parent:   parent,
	}
}

func (u *UnionNode) String() string {
	return u.Type().String()
}

// BoolNode holds bool value
type BoolNode struct {
	NodeType
	Value  bool
	Parent interface{}
}

func newBool(value bool, parent interface{}) *BoolNode {
	return &BoolNode{
		NodeType: NodeBool,
		Value:    value,
		Parent:   parent,
	}
}

func (b *BoolNode) String() string {
	return fmt.Sprintf("%s: %t", b.Type(), b.Value)
}

// ParentNode
type ParentNode struct {
	NodeType
	Parent interface{}
}

func newParent(parent interface{}) *ParentNode {
	return &ParentNode{
		NodeType: NodeParent,
		Parent:   parent,
	}
}

func (p *ParentNode) String() string {
	return fmt.Sprintf("%s: %t", p.Type(), p.Parent)
}
