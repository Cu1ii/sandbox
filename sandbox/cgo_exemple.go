package sandbox

/*
 #include <stdio.h>

char*  say_hello() {


	return "hello world";
}
*/
import "C"
import "fmt"

func sayHello() {
	s := C.say_hello()
	fmt.Printf("%v, %T", C.GoString(s), s)
}
