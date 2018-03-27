package cbor

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"strconv"
	"time"
	"unicode/utf8"
)

var decodeTimeZone *time.Location

const hexTable = "0123456789abcdef"

func decodeIntAdditonalType(src []byte, minor byte) (int64, uint, error) {
	val := int64(0)
	bytesRead := 0
	if minor <= 23 {
		val = int64(minor)
		bytesRead = 0
	} else {
		switch minor {
		case additionalTypeIntUint8:
			bytesRead = 1
		case additionalTypeIntUint16:
			bytesRead = 2
		case additionalTypeIntUint32:
			bytesRead = 4
		case additionalTypeIntUint64:
			bytesRead = 8
		default:
			return 0, 0, fmt.Errorf("Invalid Additional Type: %d in decodeInteger (expected <28)", minor)
		}
		for i := 0; i < bytesRead; i++ {
			val = val * 256
			val += int64(src[i])
		}
	}
	return val, uint(bytesRead), nil
}

func decodeInteger(src []byte) (int64, uint, error) {
	major := src[0] & maskOutAdditionalType
	minor := src[0] & maskOutMajorType
	if major != majorTypeUnsignedInt && major != majorTypeNegativeInt {
		return 0, 0, fmt.Errorf("Major type is: %d in decodeInteger!! (expected 0 or 1)", major)
	}
	val, bytesRead, err := decodeIntAdditonalType(src[1:], minor)
	if err != nil {
		return 0, 0, err
	}
	if major == 0 {
		return val, 1 + bytesRead, nil
	}
	return (-1 - val), 1 + bytesRead, nil
}

func decodeFloat(src []byte) (float64, uint, error) {
	major := (src[0] & maskOutAdditionalType)
	minor := src[0] & maskOutMajorType
	if major != majorTypeSimpleAndFloat {
		return 0, 0, fmt.Errorf("Incorrect Major type is: %d in decodeFloat", major)
	}

	switch minor {
	case additionalTypeFloat16:
		return 0, 0, fmt.Errorf("float16 is not suppported in decodeFloat")
	case additionalTypeFloat32:
		switch string(src[1:5]) {
		case float32Nan:
			return math.NaN(), 5, nil
		case float32PosInfinity:
			return math.Inf(0), 5, nil
		case float32NegInfinity:
			return math.Inf(-1), 5, nil
		}
		n := uint32(0)
		for i := 0; i < 4; i++ {
			n = n * 256
			n += uint32(src[i+1])
		}
		val := math.Float32frombits(n)
		return float64(val), 5, nil
	case additionalTypeFloat64:
		switch string(src[1:9]) {
		case float64Nan:
			return math.NaN(), 9, nil
		case float64PosInfinity:
			return math.Inf(0), 9, nil
		case float64NegInfinity:
			return math.Inf(-1), 9, nil
		}
		n := uint64(0)
		for i := 0; i < 8; i++ {
			n = n * 256
			n += uint64(src[i+1])
		}
		val := math.Float64frombits(n)
		return val, 9, nil
	}
	return 0, 0, fmt.Errorf("Invalid Additional Type: %d in decodeFloat", minor)
}

func decodeStringComplex(dst []byte, s string, pos uint) []byte {
	i := int(pos)
	const hex = "0123456789abcdef"
	start := 0

	for i < len(s) {
		b := s[i]
		if b >= utf8.RuneSelf {
			r, size := utf8.DecodeRuneInString(s[i:])
			if r == utf8.RuneError && size == 1 {
				// In case of error, first append previous simple characters to
				// the byte slice if any and append a replacement character code
				// in place of the invalid sequence.
				if start < i {
					dst = append(dst, s[start:i]...)
				}
				dst = append(dst, `\ufffd`...)
				i += size
				start = i
				continue
			}
			i += size
			continue
		}
		if b >= 0x20 && b <= 0x7e && b != '\\' && b != '"' {
			i++
			continue
		}
		// We encountered a character that needs to be encoded.
		// Let's append the previous simple characters to the byte slice
		// and switch our operation to read and encode the remainder
		// characters byte-by-byte.
		if start < i {
			dst = append(dst, s[start:i]...)
		}
		switch b {
		case '"', '\\':
			dst = append(dst, '\\', b)
		case '\b':
			dst = append(dst, '\\', 'b')
		case '\f':
			dst = append(dst, '\\', 'f')
		case '\n':
			dst = append(dst, '\\', 'n')
		case '\r':
			dst = append(dst, '\\', 'r')
		case '\t':
			dst = append(dst, '\\', 't')
		default:
			dst = append(dst, '\\', 'u', '0', '0', hex[b>>4], hex[b&0xF])
		}
		i++
		start = i
	}
	if start < len(s) {
		dst = append(dst, s[start:]...)
	}
	return dst
}

func decodeString(src []byte, noQuotes bool) ([]byte, uint, error) {
	major := src[0] & maskOutAdditionalType
	minor := src[0] & maskOutMajorType
	if major != majorTypeByteString {
		return []byte{}, 0, fmt.Errorf("Major type is: %d in decodeString", major)
	}
	result := []byte{'"'}
	if noQuotes {
		result = []byte{}
	}
	length, bytesRead, err := decodeIntAdditonalType(src[1:], minor)
	if err != nil {
		return []byte{}, 0, err
	}
	bytesRead++
	st := bytesRead
	len := uint(length)
	bytesRead += len

	result = append(result, src[st:st+len]...)
	if noQuotes {
		return result, bytesRead, nil
	}
	return append(result, '"'), bytesRead, nil
}

func decodeUTF8String(src []byte) ([]byte, uint, error) {
	major := src[0] & maskOutAdditionalType
	minor := src[0] & maskOutMajorType
	if major != majorTypeUtf8String {
		return []byte{}, 0, fmt.Errorf("Major type is: %d in decodeUTF8String", major)
	}
	result := []byte{'"'}
	length, bytesRead, err := decodeIntAdditonalType(src[1:], minor)
	if err != nil {
		return []byte{}, 0, err
	}
	bytesRead++
	st := bytesRead
	len := uint(length)
	bytesRead += len

	for i := st; i < bytesRead; i++ {
		// Check if the character needs encoding. Control characters, slashes,
		// and the double quote need json encoding. Bytes above the ascii
		// boundary needs utf8 encoding.
		if src[i] < 0x20 || src[i] > 0x7e || src[i] == '\\' || src[i] == '"' {
			// We encountered a character that needs to be encoded. Switch
			// to complex version of the algorithm.
			dst := []byte{'"'}
			dst = decodeStringComplex(dst, string(src[st:st+len]), i-st)
			return append(dst, '"'), bytesRead, nil
		}
	}
	// The string has no need for encoding an therefore is directly
	// appended to the byte slice.
	result = append(result, src[st:st+len]...)
	return append(result, '"'), bytesRead, nil
}

func array2Json(src []byte, dst io.Writer) (uint, error) {
	dst.Write([]byte{'['})
	major := (src[0] & maskOutAdditionalType)
	minor := src[0] & maskOutMajorType
	if major != majorTypeArray {
		return 0, fmt.Errorf("Major type is: %d in array2Json", major)
	}
	len := 0
	bytesRead := uint(0)
	unSpecifiedCount := false
	if minor == additionalTypeInfiniteCount {
		unSpecifiedCount = true
		bytesRead = 1
	} else {
		var length int64
		var err error
		length, bytesRead, err = decodeIntAdditonalType(src[1:], minor)
		if err != nil {
			fmt.Println("Error!!!")
			return 0, err
		}
		len = int(length)
		bytesRead++
	}
	curPos := bytesRead
	for i := 0; unSpecifiedCount || i < len; i++ {
		bc, err := Cbor2JsonOneObject(src[curPos:], dst)
		if err != nil {
			if src[curPos] == byte(majorTypeSimpleAndFloat|additionalTypeBreak) {
				bytesRead++
				break
			}
			return 0, err
		}
		curPos += bc
		bytesRead += bc
		if unSpecifiedCount {
			if src[curPos] == byte(majorTypeSimpleAndFloat|additionalTypeBreak) {
				bytesRead++
				break
			}
			dst.Write([]byte{','})
		} else if i+1 < len {
			dst.Write([]byte{','})
		}
	}
	dst.Write([]byte{']'})
	return bytesRead, nil
}

func map2Json(src []byte, dst io.Writer) (uint, error) {
	major := (src[0] & maskOutAdditionalType)
	minor := src[0] & maskOutMajorType
	if major != majorTypeMap {
		return 0, fmt.Errorf("Major type is: %d in map2Json", major)
	}
	len := 0
	bytesRead := uint(0)
	unSpecifiedCount := false
	if minor == additionalTypeInfiniteCount {
		unSpecifiedCount = true
		bytesRead = 1
	} else {
		var length int64
		var err error
		length, bytesRead, err = decodeIntAdditonalType(src[1:], minor)
		if err != nil {
			fmt.Println("Error!!!")
			return 0, err
		}
		len = int(length)
		bytesRead++
	}
	if len%2 == 1 {
		return 0, fmt.Errorf("Invalid Length of map %d - has to be even", len)
	}
	dst.Write([]byte{'{'})
	curPos := bytesRead
	for i := 0; unSpecifiedCount || i < len; i++ {
		bc, err := Cbor2JsonOneObject(src[curPos:], dst)
		if err != nil {
			//We hit the BREAK
			if src[curPos] == byte(majorTypeSimpleAndFloat|additionalTypeBreak) {
				bytesRead++
				break
			}
			return 0, err
		}
		curPos += bc
		bytesRead += bc
		if i%2 == 0 {
			//Even position values are keys
			dst.Write([]byte{':'})
		} else {
			if unSpecifiedCount {
				if src[curPos] == byte(majorTypeSimpleAndFloat|additionalTypeBreak) {
					bytesRead++
					break
				}
				dst.Write([]byte{','})
			} else if i+1 < len {
				dst.Write([]byte{','})
			}
		}
	}
	dst.Write([]byte{'}'})
	return bytesRead, nil
}

func decodeTagData(src []byte) ([]byte, uint, error) {
	major := (src[0] & maskOutAdditionalType)
	minor := src[0] & maskOutMajorType
	if major != majorTypeTags {
		return nil, 0, fmt.Errorf("Major type is: %d in decodeTagData", major)
	}
	if minor == additionalTypeTimestamp {
		tsMajor := src[1] & maskOutAdditionalType
		if tsMajor == majorTypeUnsignedInt || tsMajor == majorTypeNegativeInt {
			n, bc, err := decodeInteger(src[1:])
			if err != nil {
				return []byte{}, 0, err
			}
			t := time.Unix(n, 0)
			if decodeTimeZone != nil {
				t = t.In(decodeTimeZone)
			} else {
				t = t.In(time.UTC)
			}
			tsb := []byte{}
			tsb = append(tsb, '"')
			tsb = t.AppendFormat(tsb, IntegerTimeFieldFormat)
			tsb = append(tsb, '"')
			return tsb, 1 + bc, nil
		} else if tsMajor == majorTypeSimpleAndFloat {
			n, bc, err := decodeFloat(src[1:])
			if err != nil {
				return []byte{}, 0, err
			}
			secs := int64(n)
			n -= float64(secs)
			n *= float64(1e9)
			t := time.Unix(secs, int64(n))
			if decodeTimeZone != nil {
				t = t.In(decodeTimeZone)
			} else {
				t = t.In(time.UTC)
			}
			tsb := []byte{}
			tsb = append(tsb, '"')
			tsb = t.AppendFormat(tsb, NanoTimeFieldFormat)
			tsb = append(tsb, '"')
			return tsb, 1 + bc, nil
		} else {
			return nil, 0, fmt.Errorf("TS format is neigther int nor float: %d", tsMajor)
		}
	} else if minor == additionalTypeEmbeddedJSON {
		dataMajor := src[1] & maskOutAdditionalType
		if dataMajor == majorTypeByteString {
			emb, bc, err := decodeString(src[1:], true)
			if err != nil {
				return nil, 0, err
			}
			return emb, 1 + bc, nil
		}
		return nil, 0, fmt.Errorf("Unsupported embedded Type: %d in decodeEmbeddedJSON", dataMajor)
	} else if minor == additionalTypeIntUint16 {
            val,_,_ := decodeIntAdditonalType(src[1:], minor)
            if uint16(val) == additionalTypeTagHexString {
                emb, bc, _ := decodeString(src[3:], true)
                dst := []byte{'"'}
                for _, v := range emb {
                        dst = append(dst, hexTable[v>>4], hexTable[v&0x0f])
                }
                return append(dst, '"'), 3+bc, nil
            }
        }
	return nil, 0, fmt.Errorf("Unsupported Additional Type: %d in decodeTagData", minor)
}

func decodeSimpleFloat(src []byte) ([]byte, uint, error) {
	major := (src[0] & maskOutAdditionalType)
	minor := src[0] & maskOutMajorType
	if major != majorTypeSimpleAndFloat {
		return nil, 0, fmt.Errorf("Major type is: %d in decodeSimpleFloat", major)
	}
	switch minor {
	case additionalTypeBoolTrue:
		return []byte("true"), 1, nil
	case additionalTypeBoolFalse:
		return []byte("false"), 1, nil
	case additionalTypeNull:
		return []byte("null"), 1, nil

	case additionalTypeFloat16:
		fallthrough
	case additionalTypeFloat32:
		fallthrough
	case additionalTypeFloat64:
		v, bc, err := decodeFloat(src)
		if err != nil {
			return nil, 0, err
		}
		ba := []byte{}
		switch {
		case math.IsNaN(v):
			return []byte("\"NaN\""), bc, nil
		case math.IsInf(v, 1):
			return []byte("\"+Inf\""), bc, nil
		case math.IsInf(v, -1):
			return []byte("\"-Inf\""), bc, nil
		}
		if bc == 5 {
			ba = strconv.AppendFloat(ba, v, 'f', -1, 32)
		} else {
			ba = strconv.AppendFloat(ba, v, 'f', -1, 64)
		}
		return ba, bc, nil
	default:
		return nil, 0, fmt.Errorf("Invalid Additional Type: %d in decodeSimpleFloat", minor)
	}
}

// Cbor2JsonOneObject takes in byte array and decodes ONE CBOR Object
// usually a MAP. Use this when only ONE CBOR object needs decoding.
// Decoded string is written to the dst.
// Returns the bytes decoded and if any error was encountered.
func Cbor2JsonOneObject(src []byte, dst io.Writer) (uint, error) {
	var err error
	major := (src[0] & maskOutAdditionalType)
	bc := uint(0)
	var s []byte
	switch major {
	case majorTypeUnsignedInt:
		fallthrough
	case majorTypeNegativeInt:
		var n int64
		n, bc, err = decodeInteger(src)
		dst.Write([]byte(strconv.Itoa(int(n))))

	case majorTypeByteString:
		s, bc, err = decodeString(src, false)
		dst.Write(s)

	case majorTypeUtf8String:
		s, bc, err = decodeUTF8String(src)
		dst.Write(s)

	case majorTypeArray:
		bc, err = array2Json(src, dst)

	case majorTypeMap:
		bc, err = map2Json(src, dst)

	case majorTypeTags:
		s, bc, err = decodeTagData(src)
		dst.Write(s)

	case majorTypeSimpleAndFloat:
		s, bc, err = decodeSimpleFloat(src)
		dst.Write(s)
	}
	return bc, err
}

// Cbor2JsonManyObjects decodes all the CBOR Objects present in the
// source byte array. It keeps on decoding until it runs out of bytes.
// Decoded string is written to the dst. At the end of every CBOR Object
// newline is written to the output stream.
// Returns the number of bytes decoded and if any error was encountered.
func Cbor2JsonManyObjects(src []byte, dst io.Writer) (uint, error) {
	curPos := uint(0)
	totalBytes := uint(len(src))
	for curPos < totalBytes {
		bc, err := Cbor2JsonOneObject(src[curPos:], dst)
		if err != nil {
			return curPos, err
		}
		dst.Write([]byte("\n"))
		curPos += bc
	}
	return curPos, nil
}

// Detect if the bytes to be printed is Binary or not.
func binaryFmt(p []byte) bool {
	if len(p) > 0 && p[0] > 0x7F {
		return true
	}
	return false
}

// DecodeIfBinaryToString converts a binary formatted log msg to a
// JSON formatted String Log message - suitable for printing to Console/Syslog.
func DecodeIfBinaryToString(in []byte) string {
	if binaryFmt(in) {
		var b bytes.Buffer
		Cbor2JsonManyObjects(in, &b)
		return b.String()
	}
	return string(in)
}

// DecodeObjectToStr checks if the input is a binary format, if so,
// it will decode a single Object and return the decoded string.
func DecodeObjectToStr(in []byte) string {
	if binaryFmt(in) {
		var b bytes.Buffer
		Cbor2JsonOneObject(in, &b)
		return b.String()
	}
	return string(in)
}

// DecodeIfBinaryToBytes checks if the input is a binary format, if so,
// it will decode all Objects and return the decoded string as byte array.
func DecodeIfBinaryToBytes(in []byte) []byte {
	if binaryFmt(in) {
		var b bytes.Buffer
		Cbor2JsonManyObjects(in, &b)
		return b.Bytes()
	}
	return in
}
