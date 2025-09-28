package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func main() {
	ctx := context.Background()

	fmt.Println("测试 eino 框架基本功能...")

	// 创建一个简单的 Chain
	chain := compose.NewChain[map[string]any, *schema.Message]()

	// 编译 Chain
	runnable, err := chain.Compile(ctx)
	if err != nil {
		log.Fatalf("编译 Chain 失败: %v", err)
	}

	// 执行 Chain
	input := map[string]any{
		"query": "Hello, eino!",
	}

	output, err := runnable.Invoke(ctx, input)
	if err != nil {
		log.Fatalf("执行 Chain 失败: %v", err)
	}

	fmt.Printf("Chain 输出: %+v\n", output)
	fmt.Println("eino 框架测试完成！")
}
