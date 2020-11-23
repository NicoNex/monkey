package evaluator

import "github.com/NicoNex/monkey/obj"

var builtins = map[string]*obj.Builtin{
	"len": &obj.Builtin{
		Fn: func(args ...obj.Object) obj.Object {
			if l := len(args); l != 1 {
				return newError("wrong number of arguments: got %d, want 1", l)
			}

			switch arg := args[0].(type) {

			case *obj.String:
				return &obj.Integer{Value: int64(len(arg.Value))}

			default:
				return newError("type not supported: got %s", arg.Type())
			}
		},
	},
}
