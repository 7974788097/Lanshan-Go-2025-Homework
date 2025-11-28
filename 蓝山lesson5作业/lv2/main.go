package main

// 测试板块一,证明并发安全
/*
func a(input *int) {

	*input += 1
}
func main() {
	sum := 0
	tasks := 100000
	c := &handout.Control{}
	c.Begin(10, tasks)
	for range tasks {
		b := func() {
			a(&sum)
		}
		c.Input(b)
	}
	c.Gorun()
	c.Close()
	fmt.Println(sum)
}
*/

// 测试模块二，证明多线程运行（输出顺序混乱）
/*
func a(input int) {
	fmt.Println(input)
}
func main() {
	tasks := 100
	c := handout.Control{}
	c.Begin(10, tasks)
	for i := range tasks {
		j := i
		b := func() {
			a(j)
		}
		c.Input(b)
	}
	c.Gorun()
	c.Close()
}
*/
