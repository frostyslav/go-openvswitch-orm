package parser

import "math"

func minValue(data map[string]interface{}) int64 {
	if min, ok := data[minKey]; ok {
		val, ok := min.(float64)
		if !ok {
			return 0
		}
		return int64(val)
	}

	return 0
}

func maxValue(data map[string]interface{}) int64 {
	if max, ok := data[maxKey]; ok {
		if max == "unlimited" {
			return math.MaxInt64
		}
		val, ok := max.(float64)
		if ok {
			return int64(val)
		}
	}

	return 0
}

func minIntegerValue(data map[string]interface{}) int64 {
	if min, ok := data[minIntegerKey]; ok {
		val, ok := min.(float64)
		if !ok {
			return -1
		}
		return int64(val)
	}

	return -1
}

func maxIntegerValue(data map[string]interface{}) int64 {
	if max, ok := data[maxIntegerKey]; ok {
		val, ok := max.(float64)
		if ok {
			return int64(val)
		}
	}

	return 0
}

func maxLengthValue(data map[string]interface{}) int64 {
	if max, ok := data[maxLengthKey]; ok {
		val, ok := max.(float64)
		if ok {
			return int64(val)
		}
	}

	return 0
}

func enumValues(enum interface{}) []string {
	var values []string
	e, ok := enum.([]interface{})
	if ok {
		if len(e) == 2 {
			if eValues, ok := e[1].([]interface{}); ok {
				for _, v := range eValues {
					values = append(values, v.(string))
				}
			}
		}
	}
	return values
}
