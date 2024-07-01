package main

import (
	"dns-resolver/dns"
	"os"
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

func main() {
	dns.ResolveDNSRequest(os.Args[1])
}
