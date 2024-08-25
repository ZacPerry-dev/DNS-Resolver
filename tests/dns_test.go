package main

import ()

/*
func TestEncodeName(t *testing.T) {
	hostName := "dns.google.com"
	targetArray := []byte{3, 100, 110, 115, 6, 103, 111, 111, 103, 108, 101, 3, 99, 111, 109, 0}

	encodedHostName := dns.EncodeName(hostName)

	if !bytes.Equal(encodedHostName, targetArray) {
		t.Errorf("Expected 3dns6google3com0. Got: %s \n", string(encodedHostName))
	}
}

// May need to update this as this function is getting refactored
func TestDecodeName(t *testing.T) {
	response := []byte{
		0,
		22,
		129,
		128,
		0,
		1,
		0,
		2,
		0,
		0,
		0,
		0,
		3,
		100,
		110,
		115,
		6,
		103,
		111,
		111,
		103,
		108,
		101,
		3,
		99,
		111,
		109,
		0,
		0,
		1,
		0,
		1,
		192,
		12,
		0,
		1,
		0,
		1,
		0,
		0,
		3,
		132,
		0,
		4,
		8,
		8,
		4,
		4,
		192,
		12,
		0,
		1,
		0,
		1,
		0,
		0,
		3,
		132,
		0,
		4,
		8,
		8,
		8,
		8,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
	}

	expected := []byte{100, 110, 115, 103, 111, 111, 103, 108, 101, 99, 111, 109}

	qname, _ := dns.DecodeHostName(response, 12)

	if !bytes.Equal(qname, expected) {
		fmt.Println("Actual qname: ", qname)
		t.Errorf("Error decoding the QNAME!")
	}
}

func TestEncodeDNSHeader(t *testing.T) {
	expectedEncodedHeader := []byte{0, 22, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0}
	dnsMessage := dns.CreateDNSMessage("dns.google.com")

	encodedHeader := dns.EncodeDNSHeader(dnsMessage.Header)

	if !bytes.Equal(encodedHeader, expectedEncodedHeader) {
		t.Errorf("Expected the encoded DNS Header bytes to be %v. Instead, found %v\n", expectedEncodedHeader, encodedHeader)
	}
}

func TestEncodeDNSQuestion(t *testing.T) {
	expectedEncodedQuestion := []byte{3, 100, 110, 115, 6, 103, 111, 111, 103, 108, 101, 3, 99, 111, 109, 0, 0, 1, 0, 1}
	hostName := "dns.google.com"

	dnsMessage := dns.CreateDNSMessage(hostName)
	encodedQuestion := dns.EncodeDNSQuestion(dnsMessage.Question)

	if !bytes.Equal(encodedQuestion, expectedEncodedQuestion) {
		t.Errorf("Expected the encoded DNS Question bytes to be %v. Instead, found %v\n", expectedEncodedQuestion, encodedQuestion)
	}
}

func TestEncodeDNSMessage(t *testing.T) {
	expectedEncodedMessage := []byte{0, 22, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 3, 100, 110, 115, 6, 103, 111, 111, 103, 108, 101, 3, 99, 111, 109, 0, 0, 1, 0, 1}
	hostName := "dns.google.com"

	dnsMessage := dns.CreateDNSMessage(hostName)
	encodedMessage := dns.EncodeDNSMessage(dnsMessage)

	if !bytes.Equal(encodedMessage, expectedEncodedMessage) {
		t.Errorf("Expected the encoded DNS Message bytes to be %v. Instead, found %v\n", expectedEncodedMessage, encodedMessage)
	}
}

// TODO
//func TestExtractResponseHeader() { }

// TODO
// test to extract the different parts of the header
// For the answer / authority / additional sections, just write a seperate test for each even though they use the same functions
*/
