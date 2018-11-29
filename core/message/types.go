package message

import "regexp"

var (
	// EmployeeNumbers are the numbers of the founders
	EmployeeNumbers = &employeeNumbers{
		chris:  "615-545-4989",
		enis:   "615-397-5516",
		atish:  "931-644-5311",
		piyush: "931-644-6755",
	}
)

type employeeNumbers struct {
	chris  string
	enis   string
	atish  string
	piyush string
}

func (f *employeeNumbers) Chris() string {
	return f.chris
}

func (f *employeeNumbers) Enis() string {
	return f.enis
}

func (f *employeeNumbers) Atish() string {
	return f.atish
}

func (f *employeeNumbers) Piyush() string {
	return f.piyush
}

func (f *employeeNumbers) CustomerSupport() string {
	return f.chris
}

func (f *employeeNumbers) OnCallDeveloper() string {
	return f.atish
}

func (f *employeeNumbers) IsEmployee(number string) bool {
	cleanNumber := GetCleanPhoneNumber(number)
	if cleanNumber == f.chris || cleanNumber == f.atish || cleanNumber == f.enis || cleanNumber == f.piyush {
		return true
	}
	return false
}

// GetCleanPhoneNumber takes a raw phone number and formats it to clean phone number.
func GetCleanPhoneNumber(rawNumber string) string {
	reg := regexp.MustCompile("[^0-9]+")
	cleanNumber := reg.ReplaceAllString(rawNumber, "")
	if len(cleanNumber) < 10 {
		return cleanNumber
	}
	cleanNumber = cleanNumber[len(cleanNumber)-10:]
	cleanNumber = cleanNumber[:3] + "-" + cleanNumber[3:6] + "-" + cleanNumber[6:]
	return cleanNumber
}
