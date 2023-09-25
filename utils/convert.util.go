package utils

// boolArgument == true ? 1 : 0
func ConvertBoolToUint(boolValue bool) uint {
	var uintValue uint
	if boolValue {
		uintValue = 1
	} else {
		uintValue = 0
	}
	return uintValue
}
