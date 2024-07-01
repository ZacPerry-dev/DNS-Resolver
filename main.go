package main

import (
	"dns-resolver/dns"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strings"
)

/*
  TODO:
    [] Get to work for any host or domain name
      - Includes deocding different data types (IPv6, etc)
      - Updating to run in a loop in case it needs multiple contact points to resolve
      - Need to figure out how to query any name
    [] Try and read the response into the DNSResponse struct
    [] review but if valid, move functions out to different files for sanity
    [] Update tests (they are commented out for now)
    [] Refactor structs & functions
*/

/*
Information for structs was found here:
	https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.1
*/

// Represents the Answer, Auth, and Additionals
type DNSRecord struct {
	NAME     []byte
	TYPE     uint16
	CLASS    uint16
	TTL      uint32
	RDLENGTH uint16
	RDATA    []byte
}

// Struct for the DNS Message
type DNSMessage struct {
	Header   dns.DNSHeader
	Question dns.DNSQuestion
}

type DNSResponse struct {
	Header     dns.DNSHeader
	Question   dns.DNSQuestion
	Answer     DNSRecord
	Authority  DNSRecord
	Additional DNSRecord
}

func createDNSMessage() DNSMessage {
	header := dns.DNSHeader{
		ID:      22, //uint16(rand.Intn(65535)), make random at some point. 22 for now
		FLAGS:   0x0000,
		QDCOUNT: 1,
		ANCOUNT: 0,
		NSCOUNT: 0,
		ARCOUNT: 0,
	}

	question := dns.DNSQuestion{
		QNAME:  []byte{},
		QTYPE:  1,
		QCLASS: 1,
	}

	return DNSMessage{
		Header:   header,
		Question: question,
	}
}

// Encode the DNS Message to be sent
func encodeDNSMessage(message DNSMessage) []byte {
	var query []byte

	encodedHeader := dns.EncodeDNSHeader(message.Header)
	encodedQuestion := dns.EncodeDNSQuestion(message.Question)

	query = append(query, encodedHeader...)
	query = append(query, encodedQuestion...)

	return query
}

// Send request to the name server
func sendRequest(message []byte) []byte {
	// TODO: Don't hard code the IP and port, change later
	// 8.8.8.8:53 -> for google testing, set recursion to 1 (FLAGS = 0x0100)
	conn, err := net.Dial("udp", "192.203.230.10:53")
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

func decodeName(data []byte, offset int) ([]byte, int) {
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
		fmt.Println("---------TYPE NS")

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
	name, newOffset := decodeName(response, offset)
	offset = newOffset
	qtype := binary.BigEndian.Uint16(response[offset : offset+2])
	qclass := binary.BigEndian.Uint16(response[offset+2 : offset+4])
	ttl := binary.BigEndian.Uint32(response[offset+4 : offset+8])
	rdlength := binary.BigEndian.Uint16(response[offset+8 : offset+10])
	rdata := response[offset+10 : offset+10+int(rdlength)]
	offset += 10 + int(rdlength)

	//decodedRdata := decodeAnswerRData(rdata, qtype, rdlength)

	return DNSRecord{
		NAME:     name,
		TYPE:     qtype,
		CLASS:    qclass,
		TTL:      ttl,
		RDLENGTH: rdlength,
		RDATA:    rdata,
	}, offset
}

func main() {
	// NOTE: may need to change the arg number later depending on how I wanna do this...
	hostNameArg := os.Args[1]
	encodedHostName := dns.EncodeName(hostNameArg)

	// Create the DNS message
	message := createDNSMessage()
	message.Question.QNAME = encodedHostName

	// Encode the Message
	dnsMessageBytes := encodeDNSMessage(message)

	// Send request
	response := sendRequest(dnsMessageBytes)

	// NOTE: Consider decoding into a DNSMessage struct instance
	// decode response header
	responseHeader := dns.DecodeDNSHeader(response)
	fmt.Println("RESPONSE HEADER: ", responseHeader)

	// decode the response question
	responseQuestion, offset := dns.DecodeDNSQuestion(response)
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

	auth, offset := extractResponseAnswer(response, offset)
	fmt.Println("AUTH : ", auth)
	// Auth and additionals
	/*
		auth := make([]DNSRecord, 0)
		for range responseHeader.NSCOUNT {
			responseAuth, newOffset := extractResponseAnswer(response, offset)
			auth = append(auth, responseAuth)
			offset = newOffset
		}
		fmt.Println("Response AUTH: ", auth)
		fmt.Println("NEW OFFSET (AUTH): ", offset) */

	additionals, offset := extractResponseAnswer(response, offset)
	fmt.Println("RESPONSE ADDITIONALS: ", additionals)
	fmt.Println("FINAL OFFSET: ", offset)
}
