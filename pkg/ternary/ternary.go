package ternary

func String(cond bool, trueVal, falseVal string) string {
	if cond {
		return trueVal
	}
	return falseVal
}

func Int(cond bool, trueVal, falseVal int) int {
	if cond {
		return trueVal
	}
	return falseVal
}

func Float(cond bool, trueVal, falseVal float32) float32 {
	return Float32(cond, trueVal, falseVal)
}

func Float32(cond bool, trueVal, falseVal float32) float32 {
	if cond {
		return trueVal
	}
	return falseVal
}

func Float64(cond bool, trueVal, falseVal float64) float64 {
	if cond {
		return trueVal
	}
	return falseVal
}

func Int64(cond bool, trueVal, falseVal int64) int64 {
	if cond {
		return trueVal
	}
	return falseVal
}

func Int32(cond bool, trueVal, falseVal int32) int32 {
	if cond {
		return trueVal
	}
	return falseVal
}

func Int16(cond bool, trueVal, falseVal int16) int16 {
	if cond {
		return trueVal
	}
	return falseVal
}

func Int8(cond bool, trueVal, falseVal int8) int8 {
	if cond {
		return trueVal
	}
	return falseVal
}

func Uint64(cond bool, trueVal, falseVal uint64) uint64 {
	if cond {
		return trueVal
	}
	return falseVal
}

func Uint32(cond bool, trueVal, falseVal uint32) uint32 {
	if cond {
		return trueVal
	}
	return falseVal
}

func Uint16(cond bool, trueVal, falseVal uint16) uint16 {
	if cond {
		return trueVal
	}
	return falseVal
}

func Uint8(cond bool, trueVal, falseVal uint8) uint8 {
	if cond {
		return trueVal
	}
	return falseVal
}

func Uint(cond bool, trueVal, falseVal uint) uint {
	if cond {
		return trueVal
	}
	return falseVal
}

func Byte(cond bool, trueVal, falseVal byte) byte {
	if cond {
		return trueVal
	}
	return falseVal
}

func Rune(cond bool, trueVal, falseVal rune) rune {
	if cond {
		return trueVal
	}
	return falseVal
}
