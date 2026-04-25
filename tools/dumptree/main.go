package main
import ("fmt";"os";"strings";"github.com/paullepoulpe/jsgo/internal/rhino"; "github.com/paullepoulpe/jsgo/internal/parser")
func dump(n *rhino.Node, depth int){if n==nil{return};fmt.Printf("%s%s str=%q\n",strings.Repeat("  ",depth),n.Token().String(),n.String_());for c:=n.FirstChild();c!=nil;c=c.Next(){dump(c,depth+1)}}
func main(){f:=os.Args[1];b,_:=os.ReadFile(f);n,_:=parser.Parse(string(b));dump(n,0)}
