package main

import (
	"dns-resolver/dns"
	"os"
)

/*
Information for structs was found here:
	https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.1
*/

func main() {
	dns.ResolveDNSRequest(os.Args[1])
}
