package main

import (
	"bytes"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/yoheimuta/go-protoparser"
	"github.com/yoheimuta/go-protoparser/parser"
)

func main() {
	out := flag.String("out", "", "the output folder")
	flag.Parse()
	protosPath := flag.Args()
	if *out == "" {
		log.Fatal("out can't by empty")
	}
	if len(protosPath) == 0 {
		log.Fatal("need to provide one or more proto file path")
	}
	protoFile := new(ProtoFile)
	for _, protoPath := range protosPath {
		body, err := os.ReadFile(protoPath)
		if err != nil {
			log.Fatal(err)
		}
		got, err := protoparser.Parse(bytes.NewBuffer(body))
		if err != nil {
			log.Fatal(err)
		}
		for _, visitee := range got.ProtoBody {
			switch visitee := visitee.(type) {
			case *parser.Option:
				if visitee.OptionName == "go_package" {
					visitee.Constant = strings.ReplaceAll(visitee.Constant, "/", "")
					visitee.Constant = strings.ReplaceAll(visitee.Constant, ".", "")
					visitee.Constant = strings.ReplaceAll(visitee.Constant, "\"", "")
					protoFile.GoPackageName = visitee.Constant
				}
			case *parser.Message:
				if strings.Contains(visitee.MessageName, "Response") {
					message := parseMessage(visitee)
					protoFile.Messages = append(
						protoFile.Messages,
						message,
					)
				}
			}
		}
	}
	factoriesFile := *out + protoFile.GoPackageName + "/factories.go"
	os.Remove(factoriesFile)
	file, err := os.OpenFile(factoriesFile, os.O_RDWR|os.O_CREATE, 0775)
	if err != nil {
		log.Fatal(err)
	}
	err = tFactories.Execute(file, protoFile)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
}

func parseMessage(visitee *parser.Message) Message {
	message := Message{
		Name: visitee.MessageName,
	}
	for _, messageVisitee := range visitee.MessageBody {
		switch messageVisitee := messageVisitee.(type) {
		case *parser.Message:
			embeddedMessage := parseMessage(messageVisitee)
			message.EmbeddedMessages = append(
				message.EmbeddedMessages,
				EmbeddedMessage{
					Message: embeddedMessage,
				},
			)
		case *parser.Field:
			if messageVisitee.Type == "string" ||
				messageVisitee.Type == "bool" {
			} else if messageVisitee.Type == "bytes" {
				messageVisitee.Type = "[]byte"
			} else if messageVisitee.Type == "double" {
				messageVisitee.Type = "float64"
			} else {
				messageVisitee.Type = "*" + messageVisitee.Type
			}
			if messageVisitee.IsRepeated {
				messageVisitee.Type = "[]" + messageVisitee.Type
			}
			if messageVisitee.FieldName == "error" {
				messageVisitee.FieldName = "err"
			}
			message.Fields = append(
				message.Fields,
				Field{
					Name: messageVisitee.FieldName,
					Type: messageVisitee.Type,
				},
			)
		case *parser.Oneof:
			oneOfFields := []Field{}
			for _, oneOfField := range messageVisitee.OneofFields {
				if oneOfField.FieldName == "error" {
					oneOfField.FieldName = "err"
				}
				oneOfFields = append(
					oneOfFields,
					Field{
						Name: oneOfField.FieldName,
						Type: oneOfField.Type,
					},
				)
			}
			message.Oneof =
				Oneof{
					Name:   messageVisitee.OneofName,
					Fields: oneOfFields,
				}
		}
	}
	for i := range message.EmbeddedMessages {
		message.EmbeddedMessages[i].MessageOneOfName = message.Oneof.Name
	}
	return message
}
