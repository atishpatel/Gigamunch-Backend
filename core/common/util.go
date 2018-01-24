package common

// CommaSeperate seperates a list of strings with comma space. []string{"1","2","3"} would return "1, 2, 3".
func CommaSeperate(s []string) string {
	var v string
	for i := range s {
		if i == len(s)-1 {
			v += s[i]
			break
		}
		v += s[i] + ", "
	}
	return v
}
