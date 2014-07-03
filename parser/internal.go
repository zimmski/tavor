package parser

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/zimmski/tavor/log"
	"github.com/zimmski/tavor/token"
)

func ParseInternal(root token.Token, src io.Reader) []error {
	log.Debug("Start internal parsing")

	if root == nil {
		return []error{&token.ParserError{
			Message: "Root token is nil",
			Type:    token.ParseErrorRootIsNil,
		}}
	}

	p := &token.InternalParser{}

	if d, err := ioutil.ReadAll(src); err == nil {
		p.Data = string(d)
		p.DataLen = len(p.Data)
	} else {
		panic(err)
	}

	nex, errs := root.Parse(p, 0)

	if len(errs) != 0 {
		log.Debugf("Internal parsing failed %v", errs)

		return errs
	} else if nex != p.DataLen {
		i := p.DataLen - nex
		msg := ""

		if i > 5 {
			msg = fmt.Sprintf("Expected EOF but still %q and more left", p.Data[nex:5])
		} else {
			msg = fmt.Sprintf("Expected EOF but still %q left", p.Data[nex:nex+i])
		}

		return []error{&token.ParserError{
			Message: msg,
			Type:    token.ParseErrorExpectedEOF,
		}}
	}

	log.Debugf("Finished internal parsing")

	return nil
}
