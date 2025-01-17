//go:build !yq_nojson

package yqlib

import (
	"bytes"
	"io"

	"github.com/goccy/go-json"
)

type jsonEncoder struct {
	indentString string
	colorise     bool
	UnwrapScalar bool
}

func NewJSONEncoder(indent int, colorise bool, unwrapScalar bool) Encoder {
	var indentString = ""

	for index := 0; index < indent; index++ {
		indentString = indentString + " "
	}

	return &jsonEncoder{indentString, colorise, unwrapScalar}
}

func (je *jsonEncoder) CanHandleAliases() bool {
	return false
}

func (je *jsonEncoder) PrintDocumentSeparator(_ io.Writer) error {
	return nil
}

func (je *jsonEncoder) PrintLeadingContent(_ io.Writer, _ string) error {
	return nil
}

func (je *jsonEncoder) Encode(writer io.Writer, node *CandidateNode) error {
	log.Debugf("I need to encode %v", NodeToString(node))
	log.Debugf("kids %v", len(node.Content))

	if node.Kind == ScalarNode && je.UnwrapScalar {
		return writeString(writer, node.Value+"\n")
	}

	destination := writer
	tempBuffer := bytes.NewBuffer(nil)
	if je.colorise {
		destination = tempBuffer
	}

	var encoder = json.NewEncoder(destination)
	encoder.SetEscapeHTML(false) // do not escape html chars e.g. &, <, >
	encoder.SetIndent("", je.indentString)

	err := encoder.Encode(node)
	if err != nil {
		return err
	}
	if je.colorise {
		return colorizeAndPrint(tempBuffer.Bytes(), writer)
	}
	return nil
}
