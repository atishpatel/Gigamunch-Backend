package main

import (
	"fmt"
	"net/url"
	"strings"
)

type nameAndAddresses struct {
	Name      string
	Addresses string
}

var (
	docSheet = ""
	input    = []nameAndAddresses{
		nameAndAddresses{
			Name:      "Ali",
			Addresses: ``,
		},
		nameAndAddresses{
			Name:      "Haseen",
			Addresses: ``,
		},
		nameAndAddresses{
			Name:      "Atish",
			Addresses: ``,
		},
		nameAndAddresses{
			Name:      "Chris",
			Addresses: ``,
		},
		nameAndAddresses{
			Name:      "Enis",
			Addresses: ``,
		},
		nameAndAddresses{
			Name:      "Piyush",
			Addresses: ``,
		},
	}
)

func main() {
	fmt.Printf("Doc link:\n%s\n----------\n", docSheet)
	startAddress := []string{"166 Chesapeake Harbor Blvd, Hendersonville"}
	del := "\n"
	for _, i := range input {
		addresses := strings.Split(strings.Replace(i.Addresses, "\t", "", -1), del)
		addresses = append(startAddress, addresses...)
		printMapLinks(i.Name, addresses)
	}

}

func printMapLinks(driverName string, addresses []string) {
	var urls []string
	i := 0
	for i < len(addresses) {
		if i != 0 {
			i--
		}
		n := i + 9
		if n > len(addresses) {
			n = len(addresses)
		}
		var addressesString string
		for _, address := range addresses[i:n] {
			addressesString = addressesString + "/" + url.PathEscape(address)
		}
		urls = append(urls, fmt.Sprintf("https://www.google.com/maps/dir%s", addressesString))
		i = n
	}
	fmt.Printf("%s:\n", driverName)
	for _, u := range urls {
		fmt.Printf("%s\n\n", u)
	}
	fmt.Print("\n----------\n")
}
