package dns

import ()

// Represents the Answer, Auth, and Additionals
type DNSRecord struct {
	NAME     []byte
	TYPE     uint16
	CLASS    uint16
	TTL      uint32
	RDLENGTH uint16
	RDATA    []byte
}

// TODO: Refactor and actually get it working
func DecodeDNSRecord(response []byte, offset int) (DNSRecord, int) {
	/*
		name, newOffset := DecodeDomainName(response, offset)
		newOffset = offset + newOffset
		qtype := binary.BigEndian.Uint16(response[newOffset : newOffset+2])
		qclass := binary.BigEndian.Uint16(response[newOffset+2 : newOffset+4])
		ttl := binary.BigEndian.Uint32(response[newOffset+4 : newOffset+8])
		rdlength := binary.BigEndian.Uint16(response[newOffset+8 : newOffset+10])
		rdata := []byte{}
		/*
			if qtype == 2 && qclass == 1 {
				rdata = []byte(decodeNSrData(response, response[10+newOffset:10+newOffset+int(rdlength)]))
			} else {
				rdata = response[10+newOffset : 10+newOffset+int(rdlength)]
			}

			fmt.Println("RDATA: ", rdata)

			newOffset += 10 + int(rdlength)
	*/

	return DNSRecord{}, 0

}
