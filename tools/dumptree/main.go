package main
import ("fmt";"os";"strings";"github.com/paullepoulpe/jsgo/internal/parser";"github.com/paullepoulpe/jsgo/internal/rhino")
func dump(n *rhino.Node, d int){fmt.Printf("%s%s @%d:%d str=%q\n",strings.Repeat("  ",d),n.Token().String(),n.Lineno(),n.Charno(),n.String_());for c:=n.FirstChild();c!=nil;c=c.Next(){dump(c,d+1)}}
func main(){f:=os.Args[1];b,_:=os.ReadFile(f);n,_:=parser.Parse(string(b));dump(n,0)}
