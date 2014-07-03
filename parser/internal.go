package parser

import (
	"io"
	"io/ioutil"

	"github.com/zimmski/container/list/linkedlist"

	"github.com/zimmski/tavor/log"
	"github.com/zimmski/tavor/token"
)

func ParseInternal(root token.Token, src io.Reader) (token.Token, error) {
	log.Debug("Start internal parsing")

	if root == nil {
		return nil, &token.ParserError{
			Message: "Root token is nil",
			Type:    token.ParseErrorRootIsNil,
		}
	}

	p := &token.InternalParser{}

	if d, err := ioutil.ReadAll(src); err == nil {
		p.Data = string(d)
		p.DataLen = len(p.Data)
	} else {
		panic(err)
	}

	queue := linkedlist.New()

	nex, err := root.Parse(p, &token.ParserList{})
	if err != nil {
		return nil, err
	}

	for i := len(nex) - 1; i > -1; i-- {
		queue.Unshift(nex[i])
	}

	for !queue.Empty() {
		v, _ := queue.Shift()
		l, _ := v.(token.ParserList)

		for i := len(l.Tokens) - 1; i > -1; i-- {
			if l.Tokens[i].Index != l.Tokens[i].MaxIndex {
				nex, err = l.Tokens[i].Token.Parse(p, &l)

				for i := len(nex) - 1; i > -1; i-- {
					queue.Unshift(nex[i])
				}
			}
		}

		if l.Index == p.DataLen {
			log.Debugf("Finished internal parsing with token %v", l.Tokens[0].Token)

			return l.Tokens[0].Token, nil
		}
	}

	log.Debugf("Internal parsing failed %v", err)

	return nil, err
}
