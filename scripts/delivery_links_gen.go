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
	input    = ``
)

func main() {
	fmt.Printf("Doc link:\n%s\n----------\n", docSheet)
	startAddress := []string{"1001+Thompson+Pl,+Nashville,+TN+37217"}
	del := "\n"
	replacer := strings.NewReplacer("\t", "", ".", "", " ", "+")
	parsedInput := getNameAndAddresses(input)
	for _, i := range parsedInput {
		if len(i.Addresses) > 0 {
			addresses := strings.Split(replacer.Replace(i.Addresses), del)
			addresses = append(startAddress, addresses...)
			printMapLinks(i.Name, addresses)
		}
	}

}

func getNameAndAddresses(s string) []nameAndAddresses {
	rows := strings.Split(s, "\n")
	i := 0
	driverIndex := 0
	addressIndex := 4
	if strings.Contains(rows[0], "Driver") {
		headers := strings.Split(rows[0], "\t")
		for j, header := range headers {
			if strings.ToUpper(header) == "DRIVER" {
				driverIndex = j
			}
			if strings.ToUpper(header) == "ADDRESS" {
				addressIndex = j
			}
		}
		i++
	}
	deliveries := []nameAndAddresses{}
	delivery := nameAndAddresses{}
	for ; i < len(rows); i++ {
		row := strings.Split(rows[i], "\t")
		if delivery.Name != row[driverIndex] {
			deliveries = append(deliveries, delivery)
			delivery = nameAndAddresses{
				Name: row[driverIndex],
			}
		}
		if delivery.Addresses == "" {
			delivery.Addresses = row[addressIndex]
		} else {
			delivery.Addresses += "\n" + row[addressIndex]
		}
	}
	deliveries = append(deliveries, delivery)
	return deliveries[1:]
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
