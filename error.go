package ehclient

import (
	"fmt"
	"strings"
)

type ParserNotFoundError struct {
	Module string
	Node   string
	Attr   string
}

func (e *ParserNotFoundError) Error() string {
	var sb strings.Builder
	if e.Module != "" {
		sb.WriteString(e.Module)
		sb.WriteString(" ")
	}
	sb.WriteString("Node ")
	sb.WriteString(e.Node)
	if e.Attr != "" {
		sb.WriteString(" Attr ")
		sb.WriteString(e.Attr)
	}
	sb.WriteString(" not found")
	return sb.String()
}

func newNodeNotFoundError(module string, node string) *ParserNotFoundError {
	return &ParserNotFoundError{
		Module: module,
		Node:   node,
	}
}

func newAttrNotFoundError(module string, node string, attr string) *ParserNotFoundError {
	return &ParserNotFoundError{
		Module: module,
		Node:   node,
		Attr:   attr,
	}
}

type ParserParseError struct {
	Module   string
	Field    string
	RawValue string
	Reason   error
}

func (e *ParserParseError) Error() string {
	return fmt.Sprintf("%s Field %s parse error: %v (value: %q)", e.Module, e.Field, e.Reason, e.RawValue)
}

func newParserParseError(module string, field string, rawValue string, reason error) *ParserParseError {
	return &ParserParseError{
		Module:   module,
		Field:    field,
		RawValue: rawValue,
		Reason:   reason,
	}
}

type HttpError struct {
	StatusCode int
	Status     string
	Body       string
	Err        error
}

func (e *HttpError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("HTTP error: %d %s %v", e.StatusCode, e.Status, e.Err)
	} else {
		return fmt.Sprintf("HTTP error: %d %s %s", e.StatusCode, e.Status, e.Body)
	}
}

func newHttpError(statusCode int, status string, body string, err error) *HttpError {
	return &HttpError{
		StatusCode: statusCode,
		Status:     status,
		Body:       body,
		Err:        err,
	}
}
