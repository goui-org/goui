package goui

// import "syscall/js"

// type stack struct {
// 	items []js.Value
// 	len   int
// }

// func newJsValueStack() *stack {
// 	return &stack{}
// }

// func (q *stack) Push(item js.Value) {
// 	if q.len < len(q.items) {
// 		q.items[q.len] = item
// 	} else {
// 		q.items = append(q.items, item)
// 	}
// 	q.len++
// }

// func (q *stack) Pop() js.Value {
// 	if q.len == 0 {
// 		return js.Null()
// 	}
// 	q.len--
// 	return q.items[q.len]
// }

// func (q *Stack) Peek() T {
// 	q.mu.Lock()
// 	defer q.mu.Unlock()
// 	if q.len == 0 {
// 		var t T
// 		return t
// 	}
// 	return q.items[q.len-1]
// }

// func (q *stack) Len() int {
// 	return q.len
// }

// func (q *Stack) Slice() []T {
// 	q.mu.Lock()
// 	defer q.mu.Unlock()
// 	return q.items[:q.len]
// }

// func (q *Stack) Clear() {
// 	q.mu.Lock()
// 	defer q.mu.Unlock()
// 	q.items = nil
// 	q.len = 0
// }

// func (q *Stack) String() string {
// 	q.mu.Lock()
// 	defer q.mu.Unlock()
// 	return fmt.Sprintf("%v", q.items[:q.len])
// }
