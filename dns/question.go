package dns

import (
	"bytes"
	"encoding/binary"
)

type DNSQuestion struct {
	QNAME  []byte
	QTYPE  uint16
	QCLASS uint16
}

// Encode the question section of the DNS message
func EncodeDNSQuestion(question DNSQuestion) []byte {
	buf := new(bytes.Buffer)
	buf.Write(question.QNAME)
	binary.Write(buf, binary.BigEndian, question.QTYPE)
	binary.Write(buf, binary.BigEndian, question.QCLASS)

	return buf.Bytes()
}

// decoding
func DecodeDNSQuestion(reader *bytes.Reader) *DNSQuestion {
	var responseQuestion DNSQuestion
	responseQuestion.QNAME = []byte(DecodeDomainName(reader))
	binary.Read(reader, binary.BigEndian, &responseQuestion.QTYPE)
	binary.Read(reader, binary.BigEndian, &responseQuestion.QCLASS)

	return &responseQuestion
}
