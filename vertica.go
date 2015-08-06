package govertica

import "fmt"
import "time"
import "strings"
import "strconv"
import "reflect"

const ProtocolVersion int = 3 << 16

func Quote(val interface{}) string {
	if val == nil {
		return "NULL"
	}

	switch v := val.(type) {
	case bool:
		if v {
			return "TRUE"
		}

		return "FALSE"
	case time.Time:
		return v.Format("'2006-01-02 15:04:05'::timestamp")
	case string:
		return fmt.Sprintf("'%s'", strings.Replace(v, "'", "''", -1))
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case uint8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case uint16:
		return strconv.FormatInt(int64(v), 10)
	case int:
		return strconv.FormatInt(int64(v), 10)
	case uint:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case uint32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(int64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 64)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case fmt.Stringer:
		return v.String()
	default:
		var rval = reflect.ValueOf(val)
		if rval.Kind() == reflect.Array {
			quotedVals := []string{}

			valCount := rval.Len()
			for i := 0; i < valCount; i++ {
				quotedVals = append(quotedVals, Quote(rval.Index(i)))
			}

			return strings.Join(quotedVals, ", ")
		}

		return ""
	}
}
