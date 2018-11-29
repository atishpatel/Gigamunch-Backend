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
	date     = ""
	input    = ``
)

func main() {
	fmt.Printf("Doc link:\n%s\n----------\n", docSheet)
	startAddress := []string{"1001+Thompson+Pl,+Nashville,+TN+37217"}
	del := "\n"
	replacer := strings.NewReplacer("\t", "", ".", "", " ", "+")
	parsedInput := getNameAndAddresses(input)
	driverOutput := []string{}
	for _, i := range parsedInput {
		if len(i.Addresses) > 0 {
			addresses := strings.Split(replacer.Replace(i.Addresses), del)
			addresses = append(startAddress, addresses...)
			output := printMapLinks(i.Name, addresses)
			driverOutput = append(driverOutput, output)
		}
	}
	for _, output := range driverOutput {
		printMailToLinks(output)
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

func printMapLinks(driverName string, addresses []string) string {
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
	output := ""
	output += fmt.Sprintf("%s:\n", driverName)
	for _, u := range urls {
		output += fmt.Sprintf("%s\n\n", u)
	}
	output += "\n----------\n"
	fmt.Print(output)
	return output
}

func printMailToLinks(driverOutput string) {
	docLink := fmt.Sprintf("Doc link:\n%s\n----------\n", docSheet)
	split := strings.Split(driverOutput, ":")
	fmt.Printf("<a href=\"mailto:john@joydriv.com?subject=%s&body=%s\">mail to: %s</a><br>\n", url.QueryEscape(date+" - Gigamunch Deliveries"), url.QueryEscape(docLink+driverOutput), split[0])
}
