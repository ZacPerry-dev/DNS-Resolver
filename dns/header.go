package dns

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

// QR, OPCODE, AA, TC, RD, RA, Z, and RCODE are packed into a "FLAGS" field (2bytes)(16bits)
type DNSHeader struct {
	ID      uint16
	FLAGS   uint16
	QDCOUNT uint16
	ANCOUNT uint16
	NSCOUNT uint16
	ARCOUNT uint16
}

func EncodeDNSHeader(header DNSHeader) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, header.ID)
	binary.Write(buf, binary.BigEndian, header.FLAGS)
	binary.Write(buf, binary.BigEndian, header.QDCOUNT)
	binary.Write(buf, binary.BigEndian, header.ANCOUNT)
	binary.Write(buf, binary.BigEndian, header.NSCOUNT)
	binary.Write(buf, binary.BigEndian, header.ARCOUNT)

	return buf.Bytes()
}

// decode the header
// Get the header from the response
func DecodeDNSHeader(response []byte) DNSHeader {
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

	// This decides how many name server resource records are in the authority records section
	fmt.Println("NSCOUNT: ", responseHeader.NSCOUNT)

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

// error check the header maybe
