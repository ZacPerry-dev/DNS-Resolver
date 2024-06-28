package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strings"
)

// TODO (not including making this resolve any host or domain name, thats later :D )
/*
  [ ] Refactor function returns: they are inconsistent
  [ ] Decode the RDATA in the response Answer (based on the TYPE)
      [x] Ensure I am getting all of the IP addresses
        - Had to loop and decode however many answers are returned
  [ ] Simplify logic for decoding the Answer, along with the the following sections
      - Lots of this can be used for the following sections as well
  [ ] Decode the rest of the response sections
*/

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

// Represents the Answer, Auth, and Additionals
type DNSRecord struct {
	NAME     []byte
	TYPE     uint16
	CLASS    uint16
	TTL      uint32
	RDLENGTH uint16
	RDATA    string //TODO: maybe change to string later, need to decode this first depending on it's format
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

// TODO
// Need to both rename this and refactor to handle any "name", both for the question and answer from the response.
// SUpport compression
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

func extractResponseQuestion(response []byte) (DNSQuestion, int) {

	// Inital offset is the byte in which the header ends.
	offset := 12

	qname, newOffset := decodeQName(response, offset)
	offset = offset + newOffset

	qtype := binary.BigEndian.Uint16(response[offset : offset+2])
	qclass := binary.BigEndian.Uint16(response[offset+2 : offset+4])

	offset += 4

	// NOTE: Re-encoding the QNAME into a byte array to represent the string.This is kind of weird and I may need to change later
	return DNSQuestion{
		QNAME:  encodeQName(qname),
		QTYPE:  qtype,
		QCLASS: qclass,
	}, offset
}

// Refactor this, just testing with it for now, will clean later
func testAnswerDecode(data []byte, offset int) ([]byte, int) {
	var qnamePieces []byte
	saveOffset := offset
	jumped := false

	for {
		length := int(data[offset])

		if length == 0 {
			qnamePieces = append(qnamePieces, 0)
			offset++
			break
		}

		// Check if length indicates a pointer (first two bits are 11)
		if length&0xC0 == 0xC0 {
			if !jumped {
				saveOffset = offset + 2 // Save offset before jumping
			}
			jumped = true
			offset = int(binary.BigEndian.Uint16(data[offset:offset+2])) & 0x3FFF
			continue
		}

		offset++
		qnamePieces = append(qnamePieces, data[offset:offset+length]...)
		offset += length
	}

	if jumped {
		return qnamePieces, saveOffset
	}
	return qnamePieces, offset
}

// Decode the RData into a string based on the qtype
func decodeAnswerRData(rdata []byte, qtype uint16, rdlength uint16) string {
	// check for the different qtype values to determine how to decode this
	// TODO: Abstract these out to functions later
	var addressPieces []string

	switch qtype {

	// TYPE A - IPV4 (host address)
	case 1:
		fmt.Println("TYPE A - aka IPV4 -- Decoding...")
		for i := range rdlength {
			addressPieces = append(addressPieces, fmt.Sprintf("%d", rdata[i]))
		}
		fmt.Println("IPV4 DECODED -> ", strings.Join(addressPieces, "."))
		return strings.Join(addressPieces, ".")

		// TYPE NS - Nameserver (TODO Later)
	case 2:
		fmt.Println("TYPE NS")

	// TYPE AAAA - IPV6
	case 28:
		fmt.Println("TYPE AAAA")
		// Assuming you just loop and insert colons instead of periods

	// TYPE CNAME - Cannonical name -> points to a domain name that points to the IP address or another CNAME
	case 5:
		fmt.Println("Type CNAME")

	default:
		fmt.Println("NOT SUPPORTED")
	}

	return ""
}

// Extract the answer from the response
func extractResponseAnswer(response []byte, offset int) (DNSRecord, int) {
	name, newOffset := testAnswerDecode(response, offset)
	offset = newOffset
	qtype := binary.BigEndian.Uint16(response[offset : offset+2])
	qclass := binary.BigEndian.Uint16(response[offset+2 : offset+4])
	ttl := binary.BigEndian.Uint32(response[offset+4 : offset+8])
	rdlength := binary.BigEndian.Uint16(response[offset+8 : offset+10])
	rdata := response[offset+10 : offset+10+int(rdlength)]
	offset += 10 + int(rdlength)

	decodedRdata := decodeAnswerRData(rdata, qtype, rdlength)

	return DNSRecord{
		NAME:     name,
		TYPE:     qtype,
		CLASS:    qclass,
		TTL:      ttl,
		RDLENGTH: rdlength,
		RDATA:    decodedRdata,
	}, offset
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

	// NOTE: Consider decoding into a DNSMessage struct instance
	// decode response header
	responseHeader := extractResponseHeader(response)
	fmt.Println("RESPONSE HEADER: ", responseHeader)

	// decode the response question
	responseQuestion, offset := extractResponseQuestion(response)
	fmt.Println("RESPONSE QUESTION: ", responseQuestion)

	fmt.Println("Response Header ANCOUNT: ", responseHeader.ANCOUNT)

	// Need to loop here in order to get all of the Answers, since there can be more than one (based on the header ANCOUNT)
	answers := make([]DNSRecord, 0)
	for range responseHeader.ANCOUNT {
		responseAnswer, newOffset := extractResponseAnswer(response, offset)
		answers = append(answers, responseAnswer)
		offset = newOffset
	}
	fmt.Println("RESPONSE ANSWER: ", answers)
	fmt.Println("NEW OFFSET: ", offset)
	// TODO: Next, get the auth and additional settings
	auth, offset := extractResponseAnswer(response, offset)
	fmt.Println("Auth: ", auth)
	fmt.Print("NEW OFFSEWRT AFTER AUTHL ", offset)

	additionals, offset := extractResponseAnswer(response, offset)
	fmt.Println("additionals: ", additionals)
	fmt.Print("NEW OFFSEWRT AFTER additionals ", offset)
}
