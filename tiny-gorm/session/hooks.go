package session

import (
	"reflect"
	"tinygorm/log"
)

/*
Hooks 主要思想是提前在可能增加功能的地方预设好一个钩子
当我们需要重新修改或者增加这个地方的逻辑的时候，把扩展的类或者方法挂载到这个点 即可
例如 github action 中，当出现 git push 的时候，自动进行一下构建
还有前端的 hot reload 功能，写一点之后自己编译打包
vue 中也一样
*/

// Hooks constants
const (
	BeforeQuery  = "BeforeQuery"
	AfterQuery   = "AfterQuery"
	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"
	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"
	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
)

// CallMethod calls the registered hooks
func (s *Session) CallMethod(method string, value interface{}) {
	fm := reflect.ValueOf(s.RefTable().Model).MethodByName(method)
	// 同样使用反射实现
	// s.RefTable().Model 或者 value 就是当前会话正在操作的对象
	// 使用 methodByName 方法反射得到该对象的方法

	if value != nil {
		fm = reflect.ValueOf(value).MethodByName(method)
	}
	param := []reflect.Value{reflect.ValueOf(s)}
	if fm.IsValid() {
		if v := fm.Call(param); len(v) > 0 {
			if err, ok := v[0].Interface().(error); ok {
				log.Error(err)
			}
		}
	}
	return
}

// 可以改成接口，判断是否实现了某个接口，实现了就调用
// Hooks constants
// const (
// 	BeforeQuery  = "BeforeQuery"
// 	AfterQuery   = "AfterQuery"
// 	BeforeUpdate = "BeforeUpdate"
// 	AfterUpdate  = "AfterUpdate"
// 	BeforeDelete = "BeforeDelete"
// 	AfterDelete  = "AfterDelete"
// 	BeforeInsert = "BeforeInsert"
// 	AfterInsert  = "AfterInsert"
// )
//
// type IAfterQuery interface {
// 	AfterQuery(s *Session) error
// }
//
// type IBeforeInsert interface {
// 	BeforeInsert(s *Session) error
// }
//
// // CallMethod calls the registered hooks
// func (s *Session) CallMethod(method string, value interface{}) {
// 	param := reflect.ValueOf(value)
// 	switch method {
// 	case AfterQuery:
// 		if i, ok := param.Interface().(IAfterQuery); ok {
// 			fmt.Println("after query successful")
// 			i.AfterQuery(s)
// 		}
// 	case BeforeInsert:
// 		if i, ok := param.Interface().(IBeforeInsert); ok {
// 			i.BeforeInsert(s)
// 		}
// 	default:
// 		panic("unsupported hook method")
// 	}
// 	return
// }
