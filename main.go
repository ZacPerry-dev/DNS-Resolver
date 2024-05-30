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
	ID    uint16
	FLAGS uint16
	// QR      uint16
	// OPCODE  uint16
	// AA      uint16
	// TC      uint16
	// RD      uint16
	// RA      uint16
	// Z       uint16
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
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, message.Header.ID)
	binary.Write(buf, binary.BigEndian, message.Header.FLAGS)
	binary.Write(buf, binary.BigEndian, message.Header.RCODE)
	binary.Write(buf, binary.BigEndian, message.Header.QDCOUNT)
	binary.Write(buf, binary.BigEndian, message.Header.ANCOUNT)
	binary.Write(buf, binary.BigEndian, message.Header.NSCOUNT)
	binary.Write(buf, binary.BigEndian, message.Header.ARCOUNT)

	buf.Write(message.Question.QNAME)
	binary.Write(buf, binary.BigEndian, message.Question.QTYPE)
	binary.Write(buf, binary.BigEndian, message.Question.QCLASS)

	return buf.Bytes()
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

	_, err = conn.Write(dnsMessageBytes)
	if err != nil {
		fmt.Println("Error writing to the socket connection")
		os.Exit(1)
	}

	buf := make([]byte, 512)
	_, err = conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading response into buffer")
		os.Exit(1)
	}

	fmt.Printf("Response in hex: %x\n", buf)
	responseId := binary.BigEndian.Uint16(buf[:2])
	fmt.Println(responseId)

	// responseHeader := DNSHeader{
	// 	ID:    binary.BigEndian.Uint16(buf[0:2]),
	// 	FLAGS: binary.BigEndian.Uint16(buf[2:4]),
	// }
}
