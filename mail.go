package mailparser

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"mime"
	"net/mail"
	"strings"
)

type ParsedEmail struct {
	message mail.Message
	parts   []emailPart
}

type Attachment struct {
	Data            []byte
	MimeFileType    string
	FileName        string
	Size            int
	CreationDate    string
	ModificatioDate string
	ReadDate        string
	Encoding        string // UTF-8, 8bit, Base64
}

func ParseEmail(data []byte) (*ParsedEmail, error) {
	// create a buffer from the message bytes
	buffer := bytes.NewBuffer(data)

	// create a reader
	reader := bufio.NewReader(buffer)

	// parse the mail message
	message, err := mail.ReadMessage(reader)
	if err != nil {
		return nil, err
	}

	parts, err := parseEmailParts(*message)
	if err != nil {
		log.Println("ERROR Parsing email parts", err)
		return nil, err
	}

	parsedMessage := new(ParsedEmail)
	parsedMessage.message = *message
	parsedMessage.parts = parts
	return parsedMessage, nil

}

func (em *ParsedEmail) Header(name string) string {
	header := em.message.Header.Get(name)

	// if the header is empty then return an empty string
	if header == "" {
		return ""
	}

	decoder := new(mime.WordDecoder)
	decoded, err := decoder.Decode(header)
	if err != nil {
		return header
	}

	return decoded
}

func (em *ParsedEmail) Body() string {

	if strings.Contains(em.Header("Content-Type"), "multipart") {
		if len(em.parts) > 0 {
			for _, part := range em.parts {
				if !part.IsAttachment {
					return string(part.DataAsBytes())
				}
			}

		}
	}

	// otherwise it'a normal email
	data, err := ioutil.ReadAll(em.message.Body)
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(data)
}

func (em *ParsedEmail) Attachments() []Attachment {
	var attachments []Attachment
	for _, part := range em.parts {
		if part.IsAttachment {
			attachments = append(attachments, Attachment{
				Data:            part.DataAsBytes(),
				MimeFileType:    part.MimeFileType,
				FileName:        part.FileName,
				Size:            part.Size,
				CreationDate:    part.CreationDate,
				ModificatioDate: part.ModificatioDate,
				ReadDate:        part.ReadDate,
				Encoding:        part.Encoding,
			})
		}
	}
	return attachments
}
