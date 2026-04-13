package internal

import (
	"encoding/base64"
	"fmt"
	"io"
	"math"
	"os"
	"reflect"
	"regexp"
	"text/template"

	"github.com/cespare/xxhash/v2"
)

// UtilFuncs returns a map of utility functions for templates.
func UtilFuncs() template.FuncMap {
	return template.FuncMap{
		// Numeric / hex
		"hexString": hexString,
		"add":       add,
		"sub":       sub,
		"mul":       mul,
		"div":       div,
		"mod":       mod,

		// Hashing
		"hash":     hashString,
		"fileHash": fileHash,

		// File utilities
		"fileSize": fileSize,
		"readFile": readFileFunc,

		// General utilities
		"default": defaultVal,

		// Sequences
		"seq": seq,

		// Encoding
		"base64Encode": base64Encode,
		"base64Decode": base64Decode,

		// Regex
		"regexMatch":   regexMatch,
		"regexReplace": regexReplace,

		// Environment
		"env": os.Getenv,
	}
}

// ---------------------------------------------------------------------------
// Numeric helpers
// ---------------------------------------------------------------------------

// toFloat64 normalises any numeric type to float64.
// JSON numbers arrive as float64 in Go templates, but this also handles every
// other built-in numeric type for safety.
func toFloat64(v any) (float64, error) {
	switch n := v.(type) {
	case int:
		return float64(n), nil
	case int8:
		return float64(n), nil
	case int16:
		return float64(n), nil
	case int32:
		return float64(n), nil
	case int64:
		return float64(n), nil
	case uint:
		return float64(n), nil
	case uint8:
		return float64(n), nil
	case uint16:
		return float64(n), nil
	case uint32:
		return float64(n), nil
	case uint64:
		return float64(n), nil
	case float32:
		return float64(n), nil
	case float64:
		return n, nil
	default:
		return 0, fmt.Errorf("unsupported numeric type %T", v)
	}
}

// hexString converts a numeric value to its hex string representation without
// any "0x" prefix.  Floats are truncated to int64 before formatting.
func hexString(v any) (string, error) {
	f, err := toFloat64(v)
	if err != nil {
		return "", fmt.Errorf("hexString: %w", err)
	}
	return fmt.Sprintf("%x", int64(f)), nil
}

// add returns the sum of two numeric values.
func add(a, b any) (float64, error) {
	fa, err := toFloat64(a)
	if err != nil {
		return 0, fmt.Errorf("add: %w", err)
	}
	fb, err := toFloat64(b)
	if err != nil {
		return 0, fmt.Errorf("add: %w", err)
	}
	return fa + fb, nil
}

// sub returns a − b.
func sub(a, b any) (float64, error) {
	fa, err := toFloat64(a)
	if err != nil {
		return 0, fmt.Errorf("sub: %w", err)
	}
	fb, err := toFloat64(b)
	if err != nil {
		return 0, fmt.Errorf("sub: %w", err)
	}
	return fa - fb, nil
}

// mul returns a × b.
func mul(a, b any) (float64, error) {
	fa, err := toFloat64(a)
	if err != nil {
		return 0, fmt.Errorf("mul: %w", err)
	}
	fb, err := toFloat64(b)
	if err != nil {
		return 0, fmt.Errorf("mul: %w", err)
	}
	return fa * fb, nil
}

// div returns a ÷ b, returning an error on division by zero.
func div(a, b any) (float64, error) {
	fa, err := toFloat64(a)
	if err != nil {
		return 0, fmt.Errorf("div: %w", err)
	}
	fb, err := toFloat64(b)
	if err != nil {
		return 0, fmt.Errorf("div: %w", err)
	}
	if fb == 0 {
		return 0, fmt.Errorf("div: division by zero")
	}
	return fa / fb, nil
}

// mod returns the floating-point remainder of a ÷ b (uses math.Mod).
func mod(a, b any) (float64, error) {
	fa, err := toFloat64(a)
	if err != nil {
		return 0, fmt.Errorf("mod: %w", err)
	}
	fb, err := toFloat64(b)
	if err != nil {
		return 0, fmt.Errorf("mod: %w", err)
	}
	if fb == 0 {
		return 0, fmt.Errorf("mod: division by zero")
	}
	return math.Mod(fa, fb), nil
}

// ---------------------------------------------------------------------------
// Hashing
// ---------------------------------------------------------------------------

// hashString returns the xxHash 64-bit digest of the given string.
func hashString(s string) uint64 {
	return xxhash.Sum64String(s)
}

// fileHash returns the xxHash 64-bit digest of the contents of the named file.
func fileHash(filename string) (uint64, error) {
	f, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	h := xxhash.New()
	if _, err := io.Copy(h, f); err != nil {
		return 0, err
	}
	return h.Sum64(), nil
}

// ---------------------------------------------------------------------------
// File utilities
// ---------------------------------------------------------------------------

// fileSize returns the size in bytes of the named file.
func fileSize(filename string) (int64, error) {
	info, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// readFileFunc reads the named file and returns its contents as a string.
func readFileFunc(filename string) (string, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ---------------------------------------------------------------------------
// Default value
// ---------------------------------------------------------------------------

// defaultVal returns val unless it is "empty", in which case it returns def.
// Designed for piping: {{ .Name | default "unnamed" }}
func defaultVal(def, val any) any {
	if isEmpty(val) {
		return def
	}
	return val
}

// isEmpty reports whether v is considered empty: nil, zero-value numerics,
// false booleans, empty strings, and empty slices/maps.
func isEmpty(v any) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.String:
		return rv.Len() == 0
	case reflect.Bool:
		return !rv.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return rv.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return rv.Float() == 0
	case reflect.Slice, reflect.Map, reflect.Array:
		return rv.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return rv.IsNil()
	}
	return false
}

// ---------------------------------------------------------------------------
// Sequences
// ---------------------------------------------------------------------------

// seq generates a slice of integers.
//
//	seq end          → 1..end  (or -1..end when end < 0)
//	seq start end    → start..end (step auto-detected)
//	seq start end step
func seq(args ...int) ([]int, error) {
	var start, end, step int
	switch len(args) {
	case 1:
		end = args[0]
		if end > 0 {
			start = 1
			step = 1
		} else if end < 0 {
			start = -1
			step = -1
		} else {
			return nil, nil
		}
	case 2:
		start = args[0]
		end = args[1]
		if end >= start {
			step = 1
		} else {
			step = -1
		}
	case 3:
		start = args[0]
		end = args[1]
		step = args[2]
		if step == 0 {
			return nil, fmt.Errorf("seq: step cannot be zero")
		}
		if (end > start && step < 0) || (end < start && step > 0) {
			return nil, fmt.Errorf("seq: step %d goes in the wrong direction for range %d..%d", step, start, end)
		}
	default:
		return nil, fmt.Errorf("seq: expected 1-3 arguments, got %d", len(args))
	}

	var result []int
	if step > 0 {
		for i := start; i <= end; i += step {
			result = append(result, i)
		}
	} else {
		for i := start; i >= end; i += step {
			result = append(result, i)
		}
	}
	return result, nil
}

// ---------------------------------------------------------------------------
// Encoding helpers
// ---------------------------------------------------------------------------

// base64Encode returns the standard base64 encoding of s.
func base64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

// base64Decode decodes a standard base64 encoded string.
func base64Decode(s string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", fmt.Errorf("base64Decode: %w", err)
	}
	return string(b), nil
}

// ---------------------------------------------------------------------------
// Regex helpers
// ---------------------------------------------------------------------------

// regexMatch reports whether the string s matches the regular expression pattern.
func regexMatch(pattern, s string) (bool, error) {
	return regexp.MatchString(pattern, s)
}

// regexReplace replaces all matches of pattern in s with replacement.
// The replacement string may use $1, $2, etc. for sub-match references.
func regexReplace(pattern, replacement, s string) (string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", fmt.Errorf("regexReplace: %w", err)
	}
	return re.ReplaceAllString(s, replacement), nil
}
