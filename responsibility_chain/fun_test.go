package responsibility_chain

import (
	"context"
	"log"
	"testing"
)

/*
责任链模式有两种写法：
一种是闭包写法【难度大，但是更加灵活，可以自由组合，实际框架中使用较多】
一种是结构体写法【看起来直观，但是逻辑相对固定】
*/

type Message struct {
	Name string
}

type Option struct {
}

type MessageHandler func(context.Context, Message) error

//责任链模式的闭包写法，输入一个方法，返回一个方法
type handlerWrapper func(MessageHandler, *Option) MessageHandler

// Note that this is an init-time only list.
var handlerChain []handlerWrapper

func init() {
	handlerChain = []handlerWrapper{
		wrapShadow,
		wrapFlowControl,
	}
}

//闭包责任链从后往前加载，执行的时候就是从前往后执行
func newMessageHandler(handler MessageHandler, opt *Option) MessageHandler {
	for i := len(handlerChain) - 1; i >= 0; i-- {
		handler = handlerChain[i](handler, opt)
	}
	return handler
}

func wrapShadow(next MessageHandler, opt *Option) MessageHandler {
	return func(ctx context.Context, msg Message) error {
		log.Println("wrapShadow...")
		return next(ctx, msg)
	}
}

func wrapFlowControl(next MessageHandler, opt *Option) MessageHandler {
	targetTopics := "var wrapFlowControl"
	return func(ctx context.Context, msg Message) error {
		//判断是否需要进行流控
		log.Println("exec wrapFlowControl,use", targetTopics)
		ok := true
		if !ok {
			return next(ctx, msg)
		}
		return nil
	}
}

func TestFunc(t *testing.T) {
	selfhandler := func(ctx context.Context, msg Message) error {
		log.Println(msg)
		return nil
	}
	handler := newMessageHandler(selfhandler, nil)
	handler(context.Background(), Message{Name: "yang"})
}
