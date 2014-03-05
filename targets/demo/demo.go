package demo

type Demo struct {
    name string
}

type Command struct {
    Name string
    command string
    Args string
    Help string
}

func Setup() *Demo {
    return &Demo{}
}

func (self *Demo) GetCommands() {
    one = &Command{name: "DoFoo", command: "foo", args: nil, help: "Do Foo"}
    two = &Command{name: "DoBar", command: "bar", args: nil, help: "Do Bar"}

}
