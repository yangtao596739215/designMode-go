package closure

import (
	"context"
	"log"
	"testing"
)

type Message struct {
	Name string
}

type MessageHandler func(context.Context, *Message) error

// Consume方法有自己的逻辑，其中的一部分逻辑就是调用MessageHandler去处理，但是MessageHandler的逻辑是用户定义的，这时候就需要把方法作为函数的参数专递进来
//【MessageHandler是Consume的一个子逻辑，但是这个子逻辑是可变的】这样理解就好一些了
// 使用函数指针的好处在于，可以将实现同一功能的多个模块统一起来标识，这样一来更容易后期的维护，系统结构更加清晰。
// 或者归纳为：便于分层设计、利于系统抽象、降低耦合度以及使接口与实现分开。【实现的效果和面向对象中的多态是一样的】
//【即使新增handler，consume也不需要改变，代码复用性强】
func Consume(ctx context.Context, handler MessageHandler) error {
	log.Println("consume logic")
	err := handler(ctx, &Message{})
	log.Println(err)
	return nil
}

func TestConsume(t *testing.T) {
	userHandler := func(context.Context, *Message) error {
		log.Println("user handler")
		return nil
	}

	sysHandler := func(context.Context, *Message) error {
		log.Println("sys handler")
		return nil
	}
	ctx := context.Background()

	//通过函数参数实现了多态
	Consume(ctx, userHandler)
	Consume(ctx, sysHandler)
}
