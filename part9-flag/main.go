package main

import (
	"flag"
	"fmt"
)

//Go 提供了一个 flag 包，支持基本的命令行标志解析。我们将用这个包来实现我们的命令行程序示例。

func main() {
	//基本的标记声明仅支持字符串、整数和布尔值选项。这里我们声明一个默认值为 "foo" 的字符串标志 word并带有一个简短的描述。这里的 flag.String 函数返回一个字符串指针（不是一个字符串值），在下面我们会看到是如何使用这个指针的。
	wordPtr := flag.String("word", "foo", "a string")
	//使用和声明 word 标志相同的方法来声明 numb 和 fork 标志。
	numbPtr := flag.Int("numb", 42, "an int")
	boolPtr := flag.Bool("fork", false, "a bool")

	//所有标志都声明完成以后，调用 flag.Parse() 来执行命令行解析。
	flag.Parse()

	//这里我们将仅输出解析的选项以及后面的位置参数。注意，我们需要使用类似 *wordPtr 这样的语法来对指针解引用，从而得到选项的实际值。
	fmt.Println("word:", *wordPtr)
	fmt.Println("numb:", *numbPtr)
	fmt.Println("fork:", *boolPtr)
	fmt.Println("tail:", flag.Args())

	flag.Usage()
}
