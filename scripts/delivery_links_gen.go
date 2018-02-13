package main

import (
	"fmt"
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
			Name:      "Alandreia",
			Addresses: ``,
		},
		nameAndAddresses{
			Name:      "David",
			Addresses: ``,
		},
		nameAndAddresses{
			Name:      "Kimberly",
			Addresses: ``,
		},
		nameAndAddresses{
			Name:      "Tim",
			Addresses: ``,
		},
		nameAndAddresses{
			Name:      "Traci",
			Addresses: ``,
		},
		// Founders
		nameAndAddresses{
			Name:      "Piyush",
			Addresses: ``,
		},
		nameAndAddresses{
			Name:      "Enis",
			Addresses: ``,
		},
		nameAndAddresses{
			Name:      "Chris",
			Addresses: ``,
		},
		nameAndAddresses{
			Name:      "Atish",
			Addresses: ``,
		},
	}
)

func main() {
	fmt.Printf("Doc link:\n%s\n----------\n", docSheet)
	startAddress := []string{"166+Chesapeake+Harbor+Blvd,+Hendersonville"}
	del := "\n"
	replacer := strings.NewReplacer("\t", "", " ", "+")
	for _, i := range input {
		if len(i.Addresses) > 0 {
			addresses := strings.Split(replacer.Replace(i.Addresses), del)
			addresses = append(startAddress, addresses...)
			printMapLinks(i.Name, addresses)
		}
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
			addressesString = addressesString + "/" + address
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
