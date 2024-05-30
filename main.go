package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strings"
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

/* Encodes the host name into a byte array with each substrings length prepended to each substring (i.e. 3dns6google3com) */
func encodeHostName(hostName string) []byte {
	hostNameParts := strings.Split(hostName, ".")
	var formattedHostName []byte

	for _, sub := range hostNameParts {
		length := byte(len(sub))
		formattedHostName = append(formattedHostName, length)
		formattedHostName = append(formattedHostName, sub...)
	}

	formattedHostName = append(formattedHostName, 0)
	return formattedHostName
}

/* Converts the message to a byte string, where all of it's fields are 2 bytes each and appended to one another */
func convertToByteString(message DNSMessage) []byte {
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

	return dnsMessageBytes
}

func main() {
	hostNameArg := os.Args[1]
	encodedHostName := encodeHostName(hostNameArg)

	// Create the DNS message
	message := createDNSMessage()
	message.Question.QNAME = encodedHostName
	fmt.Println("DNS Message:", message)

	dnsMessageBytes := convertToByteString(message)

	// Print the byte string in hex
	fmt.Printf("DNS Message in hex: %x\n", dnsMessageBytes)

	// In progress: sending message to google DNS server

	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		fmt.Println("Error connecting to the socket")
		os.Exit(1)
	}
	defer conn.Close()

	n, err := conn.Write(dnsMessageBytes)
	if err != nil {
		fmt.Println("Error writing to the socket connection")
		os.Exit(1)
	}

	buf := make([]byte, n)
	_, err = conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading response into buffer")
		os.Exit(1)
	}

	fmt.Printf("Response in hex: %x\n", buf)
	response := binary.BigEndian.Uint16(buf[:2])
	fmt.Println(response)
}
