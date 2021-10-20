package terr

import (
	"fmt"
	"runtime"
)

// tagged error
type tagerror interface {
	error
	GetTraces() (trace string)
	AddTrace() *TagError
	AddTag(tag string) *TagError
	IfTagged(tags []string) bool
	IfNotTagged(tags []string) bool
	Wrap(tagerror)
	Wrapped() tagerror
}

type TagError struct {
	Tags    map[string]*struct{}
	String  string
	Trace   []string
	wrapped tagerror
}

func NewError(tags []string, format string, a ...interface{}) (e *TagError) {
	e = &TagError{
		Tags:   map[string]*struct{}{},
		String: fmt.Sprintf(format, a...),
	}
	for _, t := range tags {
		e.Tags[t] = nil
	}
	pc, fn, line, _ := runtime.Caller(1)
	e.Trace = append(e.Trace, fmt.Sprintf("%s[%s:%d]", runtime.FuncForPC(pc).Name(), fn, line))
	tagstr := e.GetTagsString()
	fmt.Printf("%v %v", tagstr, e.String)
	return e
}

func NewTaggedErrorFrom(tags []string, err error) (e *TagError) {
	e = &TagError{
		Tags:   map[string]*struct{}{},
		String: err.Error(),
	}
	for _, t := range tags {
		e.Tags[t] = nil
	}
	pc, fn, line, _ := runtime.Caller(1)
	e.Trace = append(e.Trace, fmt.Sprintf("%s[%s:%d]", runtime.FuncForPC(pc).Name(), fn, line))
	tagstr := e.GetTagsString()
	fmt.Printf("%v %v", tagstr, e.String)
	return e
}

func (e *TagError) Error() string {
	tagstr := e.GetTagsString()
	return fmt.Sprintf("%v %v", tagstr, e.String)
}

func (e *TagError) GetTraces() (trace string) {
	trace = "TRACES:\n"
	for i := len(e.Trace) - 1; i >= 0; i-- {
		trace += e.Trace[i] + "\n"
	}
	return
}

func (e *TagError) GetTagsString() (tags string) {
	tags = "["
	for t, _ := range e.Tags {
		tags += t + ", "
	}
	tags = tags[:len(tags)-2] + "]"
	return
}

func (e *TagError) AddTag(tag string) *TagError {
	e.Tags[tag] = nil
	return e
}

func (e *TagError) AddTrace() *TagError {
	pc, fn, line, _ := runtime.Caller(1)
	e.Trace = append(e.Trace, fmt.Sprintf("%s[%s:%d]", runtime.FuncForPC(pc).Name(), fn, line))
	return e
}

func (e *TagError) IfTagged(tags []string) bool {
	for _, t := range tags {
		if _, ok := e.Tags[t]; !ok {
			return false
		}
	}
	return true
}

func (e *TagError) IfNotTagged(tags []string) bool {
	for _, t := range tags {
		if _, ok := e.Tags[t]; ok {
			return false
		}
	}
	return true
}

func (e *TagError) Wrap(tagerr tagerror) {
	e.wrapped = tagerr
}

func (e *TagError) Wrapped() tagerror {
	return e.wrapped
}

