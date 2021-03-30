package v2

var ColorReset string = "\033[0m"

var FgColorRed string = "\033[31m"
var FgColorGreen string = "\033[32m"
var FgColorYellow string = "\033[33m"
var FgColorBlue string = "\033[34m"
var FgColorMagenta string = "\033[35m"
var FgColorCyan string = "\033[36m"

var BgColorRed string = "\033[41m"
var BgColorGreen string = "\033[41m"
var BgColorYellow string = "\033[41m"
var BgColorBlue string = "\033[41m"
var BgColorMagenta string = "\033[41m"
var BgColorCyan string = "\033[41m"

func FgRed(s string) string {
	return FgColorRed + s + ColorReset
}

func FgGreen(s string) string {
	return FgColorGreen + s + ColorReset
}

func FgBlue(s string) string {
	return FgColorBlue + s + ColorReset
}

func BgRed(s string) string {
	return FgColorRed + s + ColorReset
}

func BgGreen(s string) string {
	return FgColorGreen + s + ColorReset
}

func BgBlue(s string) string {
	return BgColorBlue + s + ColorReset
}
