package sqvm

type PrintFunc func(vm *VM, format string, args ...any)

type VM struct {
}

func Open(initialStackSize uint) *VM {
	return &VM{}
}

func (vm *VM) Close() {
}

func (vm *VM) SetPrintFunc(onPrint, onError PrintFunc) {
}

func (vm *VM) SetErrorHandler() {
}

func (vm *VM) Push(idx int) {
}

func (vm *VM) Pop(n int) {
}

func (vm *VM) Remove(idx int) {
}

func (vm *VM) GetTop() int {
	panic("not implemented")
}

func (vm *VM) SetTop(top int) {
}

func (vm *VM) PushString(value string) {
}

func (vm *VM) PushFloat(value float64) {
}

func (vm *VM) PushInteger(value int64) {
}

func (vm *VM) PushUserPointer(value any) {
}

func (vm *VM) PushBool(value bool) {
}

func (vm *VM) PushNull() {
}

func (vm *VM) GetType(idx int) ObjectType {
	return TypeNull
}

func (vm *VM) GetString(idx int) string {
	panic("not implemented")
}

func (vm *VM) GetInteger(idx int) int64 {
	panic("not implemented")
}

func (vm *VM) GetFloat(idx int) float64 {
	panic("not implemented")
}

func (vm *VM) GetUserPointer(idx int) any {
	panic("not implemented")
}

func (vm *VM) GetUserData(idx int) any {
	panic("not implemented")
}

func (vm *VM) GetBool(idx int) bool {
	panic("not implemented")
}

func (vm *VM) Cmp() int64 {
	panic("not implemented")
}

func (vm *VM) PushRootTable() {
}

func (vm *VM) Call(nArgs int, pushResult, raiseError bool) error {
	return nil
}
