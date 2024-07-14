package parser

import "github.com/frederik-jatzkowski/havel/internal/tooling/errors"

type Program struct {
	Packages []Package
}

func (program *Program) ResolveNames(errorsCollector *errors.Collector) {

}
