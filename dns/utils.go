package dns

import (
	"bytes"
	"io"
	"strings"
)

// encode name
// Encodes the host name into a byte array with each substrings length prepended to each substring (i.e. 3dns6google3com)
func EncodeName(hostName string) []byte {
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

// TODO: Handler Pointer Name
func DecodeDomainName(reader *bytes.Reader) string {
	var domainName bytes.Buffer

	for {
		// length for next piece of the domain name
		length, _ := reader.ReadByte()

		//check for pointer

		// If we are at the end of the domain name, break
		if length == 0 {
			break
		}

		// If a normal part of the name, just decode as usual (reading into byte array and appending a "."
		piece := make([]byte, length)
		io.ReadFull(reader, piece)
		domainName.Write(piece)
		domainName.WriteByte('.')
	}

	domainNameString := domainName.String()
	finalIndex := len(domainNameString) - 1

	return domainNameString[:finalIndex]
}

func GetDomainPointer(reader *bytes.Reader, length byte) string {

	return ""
}
