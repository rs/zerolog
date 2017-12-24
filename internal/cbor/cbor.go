package cbor

const (
	majorOffset   = 5
	additionalMax = 23
	//Non Values
	additionalTypeBoolFalse byte = 20
	additionalTypeBoolTrue  byte = 21
	additionalTypeNull      byte = 22
	//Integer (+ve and -ve) Sub-types
	additionalTypeIntUint8  byte = 24
	additionalTypeIntUint16 byte = 25
	additionalTypeIntUint32 byte = 26
	additionalTypeIntUint64 byte = 27
	//Float Sub-types
	additionalTypeFloat16 byte = 25
	additionalTypeFloat32 byte = 26
	additionalTypeFloat64 byte = 27
	additionalTypeBreak   byte = 31
	//Tag Sub-types
	additionalTypeTimestamp byte = 01
	//Unspecified number of elements
	additionalTypeInfiniteCount byte = 31
)
const (
	majorTypeUnsignedInt    byte = iota << majorOffset // Major type 0
	majorTypeNegativeInt                               // Major type 1
	majorTypeByteString                                // Major type 2
	majorTypeUtf8String                                // Major type 3
	majorTypeArray                                     // Major type 4
	majorTypeMap                                       // Major type 5
	majorTypeTags                                      // Major type 6
	majorTypeSimpleAndFloat                            // Major type 7
)

func appendCborTypePrefix(dst []byte, major byte, number uint64) []byte {
	byteCount := 8
	var minor byte
	switch {
	case number < 256:
		byteCount = 1
		minor = additionalTypeIntUint8

	case number < 65536:
		byteCount = 2
		minor = additionalTypeIntUint16

	case number < 4294967296:
		byteCount = 4
		minor = additionalTypeIntUint32

	default:
		byteCount = 8
		minor = additionalTypeIntUint64

	}
	dst = append(dst, byte(major|minor))
	byteCount--
	for ; byteCount >= 0; byteCount-- {
		dst = append(dst, byte(number>>(uint(byteCount)*8)))
	}
	return dst
}
