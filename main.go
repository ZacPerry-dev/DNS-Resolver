package main

import (
	"bytes"
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
// QR, OPCODE, AA, TC, RD, RA, Z, and RCODE are packed into a "FLAGS" field (2bytes)(16bits)
type DNSHeader struct {
	ID      uint16
	FLAGS   uint16
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
	Header     DNSHeader
	Question   DNSQuestion
	Answer     []byte
	Authority  []byte
	Additional []byte
}

func createDNSMessage() DNSMessage {
	header := DNSHeader{
		ID:      22, //uint16(rand.Intn(65535)), make random at some point. 22 for now
		FLAGS:   0x0100,
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

// Encodes the host name into a byte array with each substrings length prepended to each substring (i.e. 3dns6google3com)
func encodeQName(hostName string) []byte {
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

// Encode the header section of the DNS message
func encodeDNSHeader(header DNSHeader) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, header.ID)
	binary.Write(buf, binary.BigEndian, header.FLAGS)
	binary.Write(buf, binary.BigEndian, header.QDCOUNT)
	binary.Write(buf, binary.BigEndian, header.ANCOUNT)
	binary.Write(buf, binary.BigEndian, header.NSCOUNT)
	binary.Write(buf, binary.BigEndian, header.ARCOUNT)

	return buf.Bytes()
}

// Encode the question section of the DNS message
func encodeDNSQuestion(question DNSQuestion) []byte {
	buf := new(bytes.Buffer)
	buf.Write(question.QNAME)
	binary.Write(buf, binary.BigEndian, question.QTYPE)
	binary.Write(buf, binary.BigEndian, question.QCLASS)

	return buf.Bytes()
}

// Encode the DNS Message to be sent
func encodeDNSMessage(message DNSMessage) []byte {
	var query []byte

	encodedHeader := encodeDNSHeader(message.Header)
	encodedQuestion := encodeDNSQuestion(message.Question)

	query = append(query, encodedHeader...)
	query = append(query, encodedQuestion...)

	return query
}

// Send request to the name server
func sendRequest(message []byte) []byte {
	// TODO: Don't hard code the IP and port, change later
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		fmt.Println("Error connecting to the socket")
		os.Exit(1)
	}
	defer conn.Close()

	_, err = conn.Write(message)
	if err != nil {
		fmt.Println("Error writing to the socket connection")
		os.Exit(1)
	}

	// read the response into a buffer
	buf := make([]byte, 128)
	_, err = conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading response into buffer")
		os.Exit(1)
	}

	return buf
}

// Get the header from the response
func extractResponseHeader(response []byte) DNSHeader {
	responseHeader := DNSHeader{
		ID:      binary.BigEndian.Uint16(response[0:2]),
		FLAGS:   binary.BigEndian.Uint16(response[2:4]),
		QDCOUNT: binary.BigEndian.Uint16(response[4:6]),
		ANCOUNT: binary.BigEndian.Uint16(response[6:8]),
		NSCOUNT: binary.BigEndian.Uint16(response[8:10]),
		ARCOUNT: binary.BigEndian.Uint16(response[10:12]),
	}

	// Extracting the flags
	QR := (responseHeader.FLAGS >> 15) & 0x1
	RCODE := responseHeader.FLAGS & 0xF

	// Printing the flags
	fmt.Println("QR:", QR)
	fmt.Println("RCODE:", RCODE)

	// Check this is a response
	if QR != 1 {
		fmt.Println("Error with the response: QR does not indicate this message is a response (1)...")
		os.Exit(1)
	}

	// check for any errors returned within the header
	switch RCODE {
	case 1:
		fmt.Println("RCODE ERROR: 1 (Format Error), Name server was unable to interpret the query...")
		os.Exit(1)
	case 2:
		fmt.Println("RCODE ERROR: 2 (Server Error), Name server was unable to process the query due to a server error...")
		os.Exit(1)
	case 3:
		fmt.Println("RCODE ERROR: 3 (Name Error), Domain referenced in the query does not exist...")
		os.Exit(1)
	case 4:
		fmt.Println("RCODE ERROR: 4 (Not Implemented), The name server does not support this kind of query...")
		os.Exit(1)
	case 5:
		fmt.Println("RCODE ERROR: 5 (Refused), The name server refuses to perform this operation for policy reasons...")
		os.Exit(1)
	default:
	}

	return responseHeader
}

func decodeQName(data []byte, offset int) (string, int) {
	var qnamePieces []string
	saveOffset := offset
	// loop until we hit the last null byte (0)
	for {
		length := int(data[offset])
		if length == 0 {
			break
		}
		offset++
		qnamePieces = append(qnamePieces, string(data[offset:offset+length]))

		offset += length
	}
	fullQname := strings.Join(qnamePieces, ".")
	finalOffset := offset - saveOffset + 1

	return fullQname, finalOffset
}

func extractResponseQuestion(response []byte) DNSQuestion {

	// Inital offset is the byte in which the header ends.
	responseQuestion := DNSQuestion{}
	offset := 12

	qname, newOffset := decodeQName(response, offset)
	offset = offset + newOffset

	qtype := binary.BigEndian.Uint16(response[offset : offset+2])
	qclass := binary.BigEndian.Uint16(response[offset+2 : offset+4])

	fmt.Printf("Decoded QNAME: %s -> Offset: %d\n", qname, newOffset)
	fmt.Printf("Deocded QTYPE: %d\n", qtype)
	fmt.Printf("Deocded QCLASS: %d\n", qclass)
	return responseQuestion
}

func main() {
	hostNameArg := os.Args[1]
	encodedHostName := encodeQName(hostNameArg)

	// Create the DNS message
	message := createDNSMessage()
	message.Question.QNAME = encodedHostName

	// Encode the Message
	dnsMessageBytes := encodeDNSMessage(message)
	fmt.Printf("DNS Message in hex: %x\n", dnsMessageBytes)

	// Send request
	response := sendRequest(dnsMessageBytes)
	fmt.Printf("Response in hex: %x\n", response)

	// decode response header
	responseHeader := extractResponseHeader(response)
	fmt.Println("RESPONSE HEADER: ", responseHeader)

	// decode the response question
	responseQuestion := extractResponseQuestion(response)
	fmt.Println("RESPONSE QUESTION", responseQuestion)

	// TODO: Parse the response, get the answer, etc.
}
