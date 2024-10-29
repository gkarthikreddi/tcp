package cmdparser

import "fmt"

const MAX_CHILDREN int = 16

// function type
type callback func(*Param, *SerBuff) bool
type validation func(string) bool

type ParamMode int

const (
	CMD ParamMode = iota
	LEAF
)

type Type int

const (
	INT Type = iota
	STRING
	IPV4
	FLOAT
	IPV6
	BOOLEAN
	INVALID
)

type Leaf struct {
	leafType Type
	value    string
	fn       validation
	id       string
}

type Cmd struct {
	name string
}

type ParamType struct {
	cmd  *Cmd
	leaf *Leaf
}

type Param struct {
	mode     ParamMode
	kind     ParamType
	fn       callback
	help     string
	children [MAX_CHILDREN]*Param
	parent   *Param
	code     int
}

var root Param
var show Param
var config Param
var run Param

func GetLeaf(param *Param) *Leaf {
	return param.kind.leaf
}

func InitParam(param *Param,
	mode ParamMode,
	name string,
	run callback,
	fn validation,
	leafType Type,
	id string,
	help string) {

	if mode == CMD {
		param.mode = CMD
        param.kind = ParamType{cmd: &Cmd{name: name}}
	}

	if mode == LEAF {
		param.mode = LEAF
        param.kind = ParamType{leaf: &Leaf{leafType: leafType, fn: fn, id: id}}
	}

	param.parent = nil
	param.fn = run
	param.help = help
	param.code = -1
}

func LibcliRegisterParam(parent *Param, child *Param) error {
	if parent == nil {
		parent = &root
	}

	for i := 0; i < MAX_CHILDREN; i++ {
		if parent.children[i] == nil {
			parent.children[i] = child
			child.parent = parent
			return nil
		}
	}
	return fmt.Errorf("Can't register child to parent")
}

func SetParamCmdCode(param *Param, code int) {
    param.code = code
}

func InitLibcli() {
	InitParam(&root, CMD, "root", nil, nil, INVALID, "", "root command")

	InitParam(&show, CMD, "show", nil, nil, INVALID, "", "show command")
	LibcliRegisterParam(&root, &show)

	InitParam(&config, CMD, "config", nil, nil, INVALID, "", "config command")
	LibcliRegisterParam(&root, &config)

	InitParam(&run, CMD, "run", nil, nil, INVALID, "", "run command")
	LibcliRegisterParam(&root, &run)
}

func GetRootHook() *Param {
	return &root
}
func GetConfigHook() *Param {
	return &config
}
func GetRunHook() *Param {
	return &run
}
func GetShowHook() *Param {
	return &show
}
