package dns

import (
	"encoding/binary"
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

// NOTE: There is issues here and idk what yet
func DecodeName(data []byte, offset int) ([]byte, int) {
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
