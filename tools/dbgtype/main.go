package main
import ("fmt";"os";"strings";"github.com/paullepoulpe/jsgo/internal/parser";"github.com/paullepoulpe/jsgo/internal/rhino")
func dump(n *rhino.Node, depth int){if n==nil{return};fmt.Printf("%s%s str=%q\n",strings.Repeat("  ",depth),n.Token().String(),n.String_());if info:=n.JSDocInfo();info!=nil{if t:=info.Type();t!=nil{fmt.Printf("%s  jsdoc-type: %s\n",strings.Repeat("  ",depth),t.TypeName)};for _,p:=range info.ParameterNames(){if pt:=info.ParameterType(p);pt!=nil{fmt.Printf("%s  @param %s: %s\n",strings.Repeat("  ",depth),p,pt.TypeName)}}};for c:=n.FirstChild();c!=nil;c=c.Next(){dump(c,depth+1)}}
func main(){f:=os.Args[1];b,_:=os.ReadFile(f);n,_:=parser.Parse(string(b));dump(n,0)}
