/**@constructor*/function Foo(){
/**@type{!Number}*/this.bar = new Number(5);
}
/**@type{!Foo}*/var f = new Foo();
/**@type{!String}*/var n = f.bar;
