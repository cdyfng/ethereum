package main

import "fmt"

//Help function
func Help(args []string) int {
	for _, arg := range args {
		fmt.Println(arg)
	}
	fmt.Println("help end")
	return 0
}
func getCommandHandles() map[string]func(args []string) int {
	return map[string](func([]string) int){
		"hello": Help,
		"h":     Help,
	}

}
func main() {
	fmt.Println("Command>")
	handles := getCommandHandles()
	fmt.Println(handles)
	println(handles["hello"])
	println(handles["h"])
	println(handles["1"])

	fun := handles["hello"]
	fun([]string{"hi", "hi2", "hi3"})
}
