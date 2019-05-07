package db

import "strconv"

// QRow 查询结果
type QRow map[string]string

// GetInt 获取int类型
func (q QRow) GetInt(n string, def ...int) int {
	s, ok := q[n]
	if !ok {
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}
	r, err := strconv.Atoi(s)
	if err != nil {
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}
	return r
}

// GetInt64 获取int类型
func (q QRow) GetInt64(n string, def ...int64) int64 {
	s, ok := q[n]
	if !ok {
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}
	r, err := strconv.Atoi(s)
	if err != nil {
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}
	return int64(r)
}

// GetFloat32 获取float32类型
func (q QRow) GetFloat32(n string, def ...float32) float32 {
	s, ok := q[n]
	if !ok {
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}
	r, err := strconv.ParseFloat(s, 32)
	if err != nil {
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}
	return float32(r)
}

// GetFloat64 获取float64类型
func (q QRow) GetFloat64(n string, def ...float64) float64 {
	s, ok := q[n]
	if !ok {
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}
	r, err := strconv.ParseFloat(s, 64)
	if err != nil {
		if len(def) > 0 {
			return def[0]
		}
		return 0
	}
	return r
}
