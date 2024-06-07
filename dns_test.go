package main

import (
	"bytes"
	"testing"
)

func TestEncodeHostName(t *testing.T) {
	hostName := "dns.google.com"
	targetArray := []byte{3, 100, 110, 115, 6, 103, 111, 111, 103, 108, 101, 3, 99, 111, 109, 0}

	encodedHostName := encodeHostName(hostName)

	if !bytes.Equal(encodedHostName, targetArray) {
		t.Errorf("Expected 3dns6google3com0. Got: %s \n", string(encodedHostName))
	}
}

func TestDecodeHostName(t *testing.T) {

}

func TestEncodeDNSHeader(t *testing.T) {

}

func TestEncodeDNSQuestion(t *testing.T) {

}

func TestEncodeDNSMessage(t *testing.T) {

}
