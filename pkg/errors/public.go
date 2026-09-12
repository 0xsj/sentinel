package errors

import "maps"

// PublicInfo is a transport-neutral public snapshot with no diagnostic data.
// Kind has its stable name through String; this struct is not a wire envelope.
type PublicInfo struct {
	Kind    Kind
	Message string
	Type    string
	Fields  map[string]string
}

// Public projects one complete failure frame. The boolean reports presence,
// not classification: only a literal nil error interface returns false.
func Public(err error) (PublicInfo, bool) {
	view, present := publicFrame(err)
	view.Fields = maps.Clone(view.Fields)
	return view, present
}

// publicFrame applies disclosure once. Its fields still belong to the selected
// error; exported accessors must copy them before returning them to a caller.
func publicFrame(err error) (PublicInfo, bool) {
	if err == nil {
		return PublicInfo{}, false
	}
	e, kind, classified := inspect(err)
	if !classified || kind == Internal {
		return PublicInfo{Kind: Internal, Message: "internal error"}, true
	}
	view := PublicInfo{Kind: kind, Message: "request failed"}
	if e != nil {
		if e.message != "" {
			view.Message = e.message
		}
		view.Type = e.typ
		view.Fields = e.fields
	}
	return view, true
}

// Message returns only the selected frame's safe public message.
func Message(err error) string {
	view, _ := publicFrame(err)
	return view.Message
}

// TypeOf never inherits public metadata from a cause.
func TypeOf(err error) string {
	view, _ := publicFrame(err)
	return view.Type
}

// FieldsOf returns an independent public projection.
func FieldsOf(err error) map[string]string {
	view, _ := publicFrame(err)
	return maps.Clone(view.Fields)
}
