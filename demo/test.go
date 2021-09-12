package main

import "fmt"

/*
@author RandySun
@create 2021-09-12-21:19
*/
func main() {
	s := "123你好"
	fmt.Println(s, len(s))
	//for i := 0; i < len(s); i++ { //byte
	//	fmt.Printf("%v(%c) ", s[i], s[i])
	//}
	//123博客 9
	//[49 50 51 21338 23458] 5 6
	//fmt.Println(rune(s))
	s1 := []rune(s)
	fmt.Println(s1, len(s1), cap(s1))


}