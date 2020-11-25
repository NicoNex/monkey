package evaluator

import "github.com/NicoNex/monkey/obj"

var builtins = map[string]*obj.Builtin{
	"len": &obj.Builtin{
		Fn: func(args ...obj.Object) obj.Object {
			if l := len(args); l != 1 {
				return newError("len: wrong number of arguments: got %d, want 1", l)
			}

			switch arg := args[0].(type) {

			case *obj.String:
				return &obj.Integer{Value: int64(len(arg.Value))}

			case *obj.Array:
				return &obj.Integer{Value: int64(len(arg.Elements))}

			default:
				return newError("len: type not supported, got %s", arg.Type())
			}
		},
	},
	"append": &obj.Builtin{
		Fn: func(args ...obj.Object) obj.Object {
			if len(args) == 0 {
				return newError("append: no arguments provided")
			}

			arr, ok := args[0].(*obj.Array)
			if !ok {
				return newError("append: first argument must be an array")
			}

			if len(args) > 1 {
				for _, a := range args[1:] {
					arr.Elements = append(arr.Elements, a)
				}
			}
			return arr
		},
	},
}
