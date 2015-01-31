package parser

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/zimmski/tavor/log"
	"github.com/zimmski/tavor/token"
)

// ParseInternal reads and parses an input modelled after the given token graph.
// The errors return argument is not nil if an error is encountered during reading or parsing the input e.g. if the input does not match the given token graph.
func ParseInternal(root token.Token, src io.Reader) []error {
	log.Debug("start internal parsing")

	if root == nil {
		return []error{&token.ParserError{
			Message: "root token is nil",
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

	if len(errs) > 0 {
		log.Debugf("internal parsing failed %v", errs)

		return errs
	} else if nex != p.DataLen {
		i := p.DataLen - nex
		msg := ""

		if i > 5 {
			msg = fmt.Sprintf("Expected EOF but still %q and more left", p.Data[nex:nex+5])
		} else {
			msg = fmt.Sprintf("Expected EOF but still %q left", p.Data[nex:nex+i])
		}

		return []error{&token.ParserError{
			Message: msg,
			Type:    token.ParseErrorExpectedEOF,

			Position: p.GetPosition(nex),
		}}
	}

	log.Debugf("finished internal parsing")

	return nil
}
