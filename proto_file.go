package main

type ProtoFile struct {
	GoPackageName string
	Messages      []Message
}
type Message struct {
	Name             string
	Fields           []Field
	Oneof            Oneof
	EmbeddedMessages []EmbeddedMessage
}
type EmbeddedMessage struct {
	Message
	MessageOneOfName string
}
type Oneof struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name string
	Type string
}
