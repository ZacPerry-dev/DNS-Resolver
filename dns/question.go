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
func DecodeDNSQuestion(response []byte) (DNSQuestion, int) {
	// Inital offset is the byte in which the header ends.
	offset := 12

	qname, newOffset := DecodeName(response, offset)
	testOffset := newOffset - offset
	offset = offset + testOffset

	qtype := binary.BigEndian.Uint16(response[offset : offset+2])
	qclass := binary.BigEndian.Uint16(response[offset+2 : offset+4])

	offset += 4

	return DNSQuestion{
		QNAME:  qname,
		QTYPE:  qtype,
		QCLASS: qclass,
	}, offset
}

// Error Checking (TODO)
