package dns

import (
	"encoding/binary"
	"fmt"
	"strings"
)

// Represents the Answer, Auth, and Additionals
type DNSRecord struct {
	NAME     []byte
	TYPE     uint16
	CLASS    uint16
	TTL      uint32
	RDLENGTH uint16
	RDATA    []byte
}

// Decode the Rdata based on the qtype
func DecodeRData(rdata []byte, qtype uint16, rdlength uint16) string {
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

// Extract the record from the response
// ALSO there is issues here and idk what yet
func DecodeDNSRecord(response []byte, offset int) (DNSRecord, int) {
	name, newOffset := DecodeName(response, offset)

	fmt.Println("\n\n\nUSING THIS OFFSET IN DECODE DNS RECORD: ", newOffset)
	qtype := binary.BigEndian.Uint16(response[newOffset : newOffset+2])
	qclass := binary.BigEndian.Uint16(response[newOffset+2 : newOffset+4])
	ttl := binary.BigEndian.Uint32(response[newOffset+4 : newOffset+8])
	rdlength := binary.BigEndian.Uint16(response[newOffset+8 : newOffset+10])
	rdata := response[newOffset+10 : newOffset+10+int(rdlength)]

	newOffset += 10 + int(rdlength)

	return DNSRecord{
		NAME:     name,
		TYPE:     qtype,
		CLASS:    qclass,
		TTL:      ttl,
		RDLENGTH: rdlength,
		RDATA:    rdata,
	}, newOffset
}
