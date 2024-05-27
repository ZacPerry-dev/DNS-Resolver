package main

import (
	"fmt"
)

/*
Information for structs was found here:
	https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.1
*/
// Struct for DNS Header
type DNSHeader struct {
	ID      uint16
	QR      uint16
	OPCODE  uint16
	AA      uint16
	TC      uint16
	RD      uint16
	RA      uint16
	Z       uint16
	RCODE   uint16
	QDCOUNT uint16
	ANCOUNT uint16
	NSCOUNT uint16
	ARCOUNT uint16
}

// Struct for DNS Question
type DNSQuestion struct {
	QNAME  []byte
	QTYPE  uint16
	QCLASS uint16
}

// Struct for the DNS Message
type DNSMessage struct {
	Header   DNSHeader
	Question DNSQuestion
}

func createDNSMessage() DNSMessage {
	header := DNSHeader{
		ID:      22, //uint16(rand.Intn(65535)), make random at some point. 22 for now
		QR:      1,
		OPCODE:  0,
		AA:      0,
		TC:      0,
		RD:      1,
		RA:      0,
		Z:       0,
		RCODE:   0,
		QDCOUNT: 1,
		ANCOUNT: 0,
		NSCOUNT: 0,
		ARCOUNT: 0,
	}

	question := DNSQuestion{
		QNAME:  []byte{},
		QTYPE:  1,
		QCLASS: 1,
	}

	return DNSMessage{
		Header:   header,
		Question: question,
	}
}

func main() {
	// Create the DNS message
	message := createDNSMessage()

	// Abstract this out to a function, enable users to input domain name on the command line
	message.Question.QNAME = []byte{3, 'w', 'w', 'w', 5, 'g', 'o', 'o', 'g', 'l', 'e', 3, 'c', 'o', 'm', 0}
	fmt.Println("DNS Message:", message)

	// Create the byte string
	dnsMessageBytes := []byte{}

	// Add the header fields to the byte string
	dnsMessageBytes = append(dnsMessageBytes, byte(message.Header.ID>>8), byte(message.Header.ID))
	dnsMessageBytes = append(dnsMessageBytes, byte(message.Header.QR>>8), byte(message.Header.QR))
	dnsMessageBytes = append(dnsMessageBytes, byte(message.Header.OPCODE>>8), byte(message.Header.OPCODE))
	dnsMessageBytes = append(dnsMessageBytes, byte(message.Header.AA>>8), byte(message.Header.AA))
	dnsMessageBytes = append(dnsMessageBytes, byte(message.Header.TC>>8), byte(message.Header.TC))
	dnsMessageBytes = append(dnsMessageBytes, byte(message.Header.RD>>8), byte(message.Header.RD))
	dnsMessageBytes = append(dnsMessageBytes, byte(message.Header.RA>>8), byte(message.Header.RA))
	dnsMessageBytes = append(dnsMessageBytes, byte(message.Header.Z>>8), byte(message.Header.Z))
	dnsMessageBytes = append(dnsMessageBytes, byte(message.Header.RCODE>>8), byte(message.Header.RCODE))
	dnsMessageBytes = append(dnsMessageBytes, byte(message.Header.QDCOUNT>>8), byte(message.Header.QDCOUNT))
	dnsMessageBytes = append(dnsMessageBytes, byte(message.Header.ANCOUNT>>8), byte(message.Header.ANCOUNT))
	dnsMessageBytes = append(dnsMessageBytes, byte(message.Header.NSCOUNT>>8), byte(message.Header.NSCOUNT))
	dnsMessageBytes = append(dnsMessageBytes, byte(message.Header.ARCOUNT>>8), byte(message.Header.ARCOUNT))

	// Add the question fields to the byte string
	dnsMessageBytes = append(dnsMessageBytes, message.Question.QNAME...)
	dnsMessageBytes = append(dnsMessageBytes, byte(message.Question.QTYPE>>8), byte(message.Question.QTYPE))
	dnsMessageBytes = append(dnsMessageBytes, byte(message.Question.QCLASS>>8), byte(message.Question.QCLASS))

	// Print the byte string in hex
	fmt.Printf("DNS Message in hex: %x\n", dnsMessageBytes)
}
