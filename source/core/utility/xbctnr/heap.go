package xbctnr

import (
	"iter"
	"math"
	"slices"

	"github.com/starryck/x-lib-go/source/core/base/xbtype"
	"github.com/starryck/x-lib-go/source/core/toolkit/xbvalue"
)

const heapNodeIndexFirst = 0

type Heap[T HeapValue[T]] struct {
	values []T
}

type HeapValue[T any] xbtype.Comparable[T]

func (heap *Heap[T]) Size() int {
	return len(heap.values)
}

func (heap *Heap[T]) Push(value T) {
	node := newHeapNode(heap, heap.Size())
	heap.push(value)
	heap.swim(node)
}

func (heap *Heap[T]) Pull() (T, bool) {
	if heap.Size() == 0 {
		return xbvalue.Zero[T](), false
	}
	if heap.Size() == 1 {
		return heap.pop(), true
	}
	node := newHeapNode(heap, heapNodeIndexFirst)
	return heap.drop(node), true
}

func (heap *Heap[T]) Peek() (T, bool) {
	if heap.Size() == 0 {
		return xbvalue.Zero[T](), false
	}
	return heap.values[heapNodeIndexFirst], true
}

func (heap *Heap[T]) Slice() []T {
	return slices.Clone(heap.values)
}

func (heap *Heap[T]) Sequence() iter.Seq2[int, T] {
	return slices.All(heap.values)
}

func (heap *Heap[T]) Iterator() *HeapIterator[T] {
	return &HeapIterator[T]{heap: heap}
}

func (heap *Heap[T]) Clear() {
	heap.values = []T{}
}

func (heap *Heap[T]) push(value T) {
	heap.values = append(heap.values, value)
}

func (heap *Heap[T]) pop() T {
	value := heap.values[heap.Size()-1]
	heap.values = heap.values[:heap.Size()-1]
	return value
}

func (heap *Heap[T]) drop(node *heapNode[T]) T {
	dest := node.toLast()
	dest.exchange(node)
	value := heap.pop()
	heap.sink(dest)
	return value
}

func (heap *Heap[T]) swim(node *heapNode[T]) {
	next := node.toParent()
	for next != nil {
		if node.isValueGte(next) {
			break
		}
		node.exchange(next)
		next = node.toParent()
	}
}

func (heap *Heap[T]) sink(node *heapNode[T]) {
	next := node.toLeftChild()
	for next != nil {
		if dest := node.toRightChild(); dest != nil && dest.isValueLt(next) {
			next = dest
		}
		if node.isValueLte(next) {
			break
		}
		node.exchange(next)
		next = node.toLeftChild()
	}
}

func newHeapNode[T HeapValue[T]](heap *Heap[T], index int) *heapNode[T] {
	return &heapNode[T]{heap: heap, index: index}
}

type heapNode[T HeapValue[T]] struct {
	heap  *Heap[T]
	index int
}

func (node *heapNode[T]) getValue() T {
	return node.heap.values[node.index]
}

func (node *heapNode[T]) exchange(dest *heapNode[T]) {
	*node, *dest = *dest, *node
	values := node.heap.values
	values[node.index], values[dest.index] = values[dest.index], values[node.index]
}

func (node *heapNode[T]) compare(dest *heapNode[T]) int {
	return node.getValue().Compare(dest.getValue())
}

func (node *heapNode[T]) isValueGt(dest *heapNode[T]) bool {
	return node.compare(dest) > 0
}

func (node *heapNode[T]) isValueGte(dest *heapNode[T]) bool {
	return node.compare(dest) >= 0
}

func (node *heapNode[T]) isValueLt(dest *heapNode[T]) bool {
	return node.compare(dest) < 0
}

func (node *heapNode[T]) isValueLte(dest *heapNode[T]) bool {
	return node.compare(dest) <= 0
}

func (node *heapNode[T]) toLast() *heapNode[T] {
	if index := node.makeLastIndex(); index >= 0 {
		return newHeapNode(node.heap, index)
	}
	return nil
}

func (node *heapNode[T]) toParent() *heapNode[T] {
	if index := node.makeParentIndex(); index >= 0 {
		return newHeapNode(node.heap, index)
	}
	return nil
}

func (node *heapNode[T]) toLeftChild() *heapNode[T] {
	if index := node.makeLeftChildIndex(); index < node.heap.Size() {
		return newHeapNode(node.heap, index)
	}
	return nil
}

func (node *heapNode[T]) toRightChild() *heapNode[T] {
	if index := node.makeRightChildIndex(); index < node.heap.Size() {
		return newHeapNode(node.heap, index)
	}
	return nil
}

func (node *heapNode[T]) makeLastIndex() int {
	return node.heap.Size() - 1
}

func (node *heapNode[T]) makeParentIndex() int {
	return (node.index+1)/2 - 1
}

func (node *heapNode[T]) makeLeftChildIndex() int {
	return 2*node.index + 1
}

func (node *heapNode[T]) makeRightChildIndex() int {
	return 2*node.index + 2
}

type HeapIterator[T HeapValue[T]] struct {
	heap  *Heap[T]
	index int
}

func (iterator *HeapIterator[T]) Next() (T, bool) {
	heap := iterator.heap
	if iterator.index >= heap.Size() {
		return xbvalue.Zero[T](), false
	}
	value := heap.values[iterator.index]
	iterator.index++
	return value, true
}

const (
	deapNodeIndexMin = 0
	deapNodeIndexMax = 1
	deapNodeClassMin = 0
	deapNodeClassMax = 1
)

type Deap[T DeapValue[T]] struct {
	values []T
}

type DeapValue[T any] HeapValue[T]

func (deap *Deap[T]) Size() int {
	return len(deap.values)
}

func (deap *Deap[T]) Push(value T) {
	node := newDeapNode(deap, deap.Size())
	deap.push(value)
	deap.swap(node)
	deap.swim(node)
}

func (deap *Deap[T]) PullMin() (T, bool) {
	if deap.Size() == 0 {
		return xbvalue.Zero[T](), false
	}
	if deap.Size() == 1 {
		return deap.pop(), true
	}
	node := newDeapNode(deap, deapNodeIndexMin)
	return deap.drop(node), true
}

func (deap *Deap[T]) PullMax() (T, bool) {
	if deap.Size() == 0 {
		return xbvalue.Zero[T](), false
	}
	if deap.Size() <= 2 {
		return deap.pop(), true
	}
	node := newDeapNode(deap, deapNodeIndexMax)
	return deap.drop(node), true
}

func (deap *Deap[T]) PeekMin() (T, bool) {
	if deap.Size() == 0 {
		return xbvalue.Zero[T](), false
	}
	return deap.values[deapNodeIndexMin], true
}

func (deap *Deap[T]) PeekMax() (T, bool) {
	if deap.Size() == 0 {
		return xbvalue.Zero[T](), false
	}
	if deap.Size() == 1 {
		return deap.values[deapNodeIndexMin], true
	}
	return deap.values[deapNodeIndexMax], true
}

func (deap *Deap[T]) Slice() []T {
	return slices.Clone(deap.values)
}

func (deap *Deap[T]) Sequence() iter.Seq2[int, T] {
	return slices.All(deap.values)
}

func (deap *Deap[T]) Iterator() *DeapIterator[T] {
	return &DeapIterator[T]{deap: deap}
}

func (deap *Deap[T]) Clear() {
	deap.values = []T{}
}

func (deap *Deap[T]) push(value T) {
	deap.values = append(deap.values, value)
}

func (deap *Deap[T]) pop() T {
	value := deap.values[deap.Size()-1]
	deap.values = deap.values[:deap.Size()-1]
	return value
}

func (deap *Deap[T]) drop(node *deapNode[T]) T {
	dest := node.toLast()
	dest.exchange(node)
	value := deap.pop()
	deap.sink(dest)
	if next, ok := deap.swap(dest); ok {
		deap.swim(dest)
		deap.sort(next)
	} else {
		deap.sort(dest)
	}
	return value
}

func (deap *Deap[T]) swim(node *deapNode[T]) {
	if node.isInMinHeap() {
		deap.swimInMinHeap(node)
	} else {
		deap.swimInMaxHeap(node)
	}
}

func (deap *Deap[T]) swimInMinHeap(node *deapNode[T]) {
	next := node.toParent()
	for next != nil {
		if node.isValueGte(next) {
			break
		}
		node.exchange(next)
		next = node.toParent()
	}
}

func (deap *Deap[T]) swimInMaxHeap(node *deapNode[T]) {
	next := node.toParent()
	for next != nil {
		if node.isValueLte(next) {
			break
		}
		node.exchange(next)
		next = node.toParent()
	}
}

func (deap *Deap[T]) sink(node *deapNode[T]) {
	if node.isInMinHeap() {
		deap.sinkInMinHeap(node)
	} else {
		deap.sinkInMaxHeap(node)
	}
}

func (deap *Deap[T]) sinkInMinHeap(node *deapNode[T]) {
	next := node.toLeftChild()
	for next != nil {
		if dest := node.toRightChild(); dest != nil && dest.isValueLt(next) {
			next = dest
		}
		if node.isValueLte(next) {
			break
		}
		node.exchange(next)
		next = node.toLeftChild()
	}
}

func (deap *Deap[T]) sinkInMaxHeap(node *deapNode[T]) {
	next := node.toLeftChild()
	for next != nil {
		if dest := node.toRightChild(); dest != nil && dest.isValueGt(next) {
			next = dest
		}
		if node.isValueGte(next) {
			break
		}
		node.exchange(next)
		next = node.toLeftChild()
	}
}

func (deap *Deap[T]) swap(node *deapNode[T]) (*deapNode[T], bool) {
	dest := node.toContrast()
	if dest == nil {
		return nil, false
	}
	if (node.isInMinHeap() && node.isValueGt(dest)) || (node.isInMaxHeap() && node.isValueLt(dest)) {
		node.exchange(dest)
		return dest, true
	}
	return nil, false
}

func (deap *Deap[T]) sort(node *deapNode[T]) {
	dest := node.toContrast()
	if dest == nil || dest.isInMaxHeap() {
		return
	}
	next := dest.toLeftChild()
	if next == nil {
		return
	}
	if more := dest.toRightChild(); more != nil && more.isValueGt(next) {
		next = more
	}
	if node.isValueLt(next) {
		node.exchange(next)
	}
}

func newDeapNode[T DeapValue[T]](deap *Deap[T], index int) *deapNode[T] {
	node := &deapNode[T]{deap: deap, index: index}
	node.initialize()
	return node
}

type deapNode[T DeapValue[T]] struct {
	deap  *Deap[T]
	index int
	level int
	class int
}

func (node *deapNode[T]) getValue() T {
	return node.deap.values[node.index]
}

func (node *deapNode[T]) initialize() {
	node.level = node.makeLevel()
	node.class = node.makeClass()
}

func (node *deapNode[T]) exchange(dest *deapNode[T]) {
	*node, *dest = *dest, *node
	values := node.deap.values
	values[node.index], values[dest.index] = values[dest.index], values[node.index]
}

func (node *deapNode[T]) compare(dest *deapNode[T]) int {
	return node.getValue().Compare(dest.getValue())
}

func (node *deapNode[T]) isValueGt(dest *deapNode[T]) bool {
	return node.compare(dest) > 0
}

func (node *deapNode[T]) isValueGte(dest *deapNode[T]) bool {
	return node.compare(dest) >= 0
}

func (node *deapNode[T]) isValueLt(dest *deapNode[T]) bool {
	return node.compare(dest) < 0
}

func (node *deapNode[T]) isValueLte(dest *deapNode[T]) bool {
	return node.compare(dest) <= 0
}

func (node *deapNode[T]) isInMinHeap() bool {
	return node.class == deapNodeClassMin
}

func (node *deapNode[T]) isInMaxHeap() bool {
	return node.class == deapNodeClassMax
}

func (node *deapNode[T]) toLast() *deapNode[T] {
	if index := node.makeLastIndex(); index >= 0 {
		return newDeapNode(node.deap, index)
	}
	return nil
}

func (node *deapNode[T]) toContrast() *deapNode[T] {
	if index := node.makeContrastIndex(); index >= 0 {
		return newDeapNode(node.deap, index)
	}
	return nil
}

func (node *deapNode[T]) toParent() *deapNode[T] {
	if index := node.makeParentIndex(); index >= 0 {
		return newDeapNode(node.deap, index)
	}
	return nil
}

func (node *deapNode[T]) toLeftChild() *deapNode[T] {
	if index := node.makeLeftChildIndex(); index < node.deap.Size() {
		return newDeapNode(node.deap, index)
	}
	return nil
}

func (node *deapNode[T]) toRightChild() *deapNode[T] {
	if index := node.makeRightChildIndex(); index < node.deap.Size() {
		return newDeapNode(node.deap, index)
	}
	return nil
}

func (node *deapNode[T]) makeLevel() int {
	return int(math.Log2(float64(node.index)/2 + 1))
}

func (node *deapNode[T]) makeClass() int {
	if node.index < int(3*math.Pow(2, float64(node.level))-2) {
		return deapNodeClassMin
	} else {
		return deapNodeClassMax
	}
}

func (node *deapNode[T]) makeLastIndex() int {
	return node.deap.Size() - 1
}

func (node *deapNode[T]) makeContrastIndex() int {
	half := int(math.Pow(2, float64(node.level)))
	if node.isInMaxHeap() {
		return node.index - half
	}
	index := node.index + half
	if index < node.deap.Size() {
		return index
	}
	return index/2 - 1
}

func (node *deapNode[T]) makeParentIndex() int {
	return node.index/2 - 1
}

func (node *deapNode[T]) makeLeftChildIndex() int {
	return 2 * (node.index + 1)
}

func (node *deapNode[T]) makeRightChildIndex() int {
	return 2*(node.index+1) + 1
}

type DeapIterator[T DeapValue[T]] struct {
	deap  *Deap[T]
	index int
}

func (iterator *DeapIterator[T]) Next() (T, bool) {
	deap := iterator.deap
	if iterator.index >= deap.Size() {
		return xbvalue.Zero[T](), false
	}
	value := deap.values[iterator.index]
	iterator.index++
	return value, true
}
