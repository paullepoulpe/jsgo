package main
import ("fmt";"os";"github.com/paullepoulpe/jsgo/legacy/internal/jscomp")
func main(){
  f := os.Args[1]
  b, _ := os.ReadFile(f)
  d, _ := jscomp.CheckSource(f, string(b))
  for _, x := range d {
    fmt.Printf("%d:%d %s %s\n", x.Line, x.Col, x.Code, x.Message)
  }
}
