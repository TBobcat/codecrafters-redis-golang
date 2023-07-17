## July 16th
Every Golang program has 1 main package, and 1 main function

- socket stage:
    - open a socket on a port and write response to it

```go
/*
## init go module, and download modules for import in code
go mod init codecrafters-redis-go/app
go mod tidy
*/

func (variable_name variable_data_type) function_name() [return_type]{
   /* function body*/
}

it goes receiver, function name, function return type,  in the header

## example
package main

import (
	"fmt"
	"math"
)


func main() {
	v := Vertex{3, 4}
	fmt.Println(v.Abs())
	fmt.Println(v.anyfunc())
}

type Vertex struct {
	X, Y float64
}

func (v Vertex) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// any instance of Vertex v, can call this function with v.anyfunc()
func (v Vertex) anyfunc() float64 {
	return v.X
}
```
ref:
Methods and interfaces: https://go.dev/tour/methods/4

https://www.golinuxcloud.com/golang-tcp-server-client/#Building_a_Simple_Golang_TCP_Server
