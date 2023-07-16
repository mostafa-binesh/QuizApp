package utils

func ConvertBoolToUint(boolValue bool) uint {
	var uintValue uint
	if boolValue {
		uintValue = 1
	} else {
		uintValue = 0
	}
	return uintValue
}