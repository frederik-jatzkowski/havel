package scope

import "github.com/frederik-jatzkowski/havel/pkg/tool"

type Object interface {
	tool.NodeLike
	Identifier() string
}
