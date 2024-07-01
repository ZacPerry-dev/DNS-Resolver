package dns

import (
	"fmt"
	"net"
	"os"
)

type DNSMessage struct {
	Header   DNSHeader
	Question DNSQuestion
}

type DNSResponse struct {
	Header     DNSHeader
	Question   DNSQuestion
	Answer     DNSRecord
	Authority  DNSRecord
	Additional DNSRecord
}

func CreateDNSMessage(domainName string) DNSMessage {
	header := DNSHeader{
		ID:      22, //uint16(rand.Intn(65535)), make random at some point. 22 for now
		FLAGS:   0x0000,
		QDCOUNT: 1,
		ANCOUNT: 0,
		NSCOUNT: 0,
		ARCOUNT: 0,
	}

	question := DNSQuestion{
		QNAME:  EncodeName(domainName),
		QTYPE:  1,
		QCLASS: 1,
	}

	return DNSMessage{
		Header:   header,
		Question: question,
	}
}

// Encode the DNS Message to be sent
func EncodeDNSMessage(message DNSMessage) []byte {
	var query []byte

	encodedHeader := EncodeDNSHeader(message.Header)
	encodedQuestion := EncodeDNSQuestion(message.Question)

	query = append(query, encodedHeader...)
	query = append(query, encodedQuestion...)

	return query
}

// Send request to the name server
func ResolveDNSRequest(domainName string) string {

	// Create and encode the DNS Message
	dnsMessage := CreateDNSMessage(domainName)
	encodedMessage := EncodeDNSMessage(dnsMessage)

	// Send the message / request
	// 8.8.8.8:53 -> for google testing, set recursion to 1 (FLAGS = 0x0100)
  // TODO: THIS WILL NEED TO BE IN A LOOP & LOOP THROUGH WHATEVER IPs I HAVE SAVED currently
	conn, err := net.Dial("udp", "192.203.230.10:53")
	if err != nil {
		fmt.Println("Error connecting to the socket")
		os.Exit(1)
	}
	defer conn.Close()

	_, err = conn.Write(encodedMessage)
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
  
  // HEADER // 
	responseHeader := DecodeDNSHeader(buf)
	fmt.Println("   --- HEADER: ", responseHeader)
  
  // QUESTION // 
	responseQuestion, offset := DecodeDNSQuestion(buf)
	fmt.Println("   --- QUESTION: ", responseQuestion)
	fmt.Println("         CURR. OFFSET: ", offset)
  
  // ANSWER (MORE TODO HERE) //
	responseAnswers := make([]DNSRecord, 0)
	for range responseHeader.ANCOUNT {
		anw, newOffset := DecodeDNSRecord(buf, offset)
		responseAnswers = append(responseAnswers, anw)
		offset = newOffset
	}
	fmt.Println("   --- RESPONSE ANSWER: ", responseAnswers)
	fmt.Println("         CURR. OFFSET: ", offset)
  
  // AUTH //


  // ADDITIONALS//


	return ""
}
