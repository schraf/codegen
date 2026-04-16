package internal

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// toFloat64
// ---------------------------------------------------------------------------

func TestToFloat64(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected float64
		wantErr  bool
	}{
		{"int", int(42), 42.0, false},
		{"int8", int8(8), 8.0, false},
		{"int16", int16(16), 16.0, false},
		{"int32", int32(32), 32.0, false},
		{"int64", int64(64), 64.0, false},
		{"uint", uint(10), 10.0, false},
		{"uint8", uint8(8), 8.0, false},
		{"uint16", uint16(16), 16.0, false},
		{"uint32", uint32(32), 32.0, false},
		{"uint64", uint64(64), 64.0, false},
		{"float32", float32(3.14), float64(float32(3.14)), false},
		{"float64", float64(2.718), 2.718, false},
		{"string unsupported", "hello", 0, true},
		{"bool unsupported", true, 0, true},
		{"nil unsupported", nil, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := toFloat64(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.InDelta(t, tt.expected, result, 0.0001)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// hexString
// ---------------------------------------------------------------------------

func TestHexString(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
		wantErr  bool
	}{
		{"int 255", int(255), "ff", false},
		{"int 0", int(0), "0", false},
		{"int 16", int(16), "10", false},
		{"int64 4096", int64(4096), "1000", false},
		{"uint32 255", uint32(255), "ff", false},
		{"float64 as int", float64(255), "ff", false},
		{"float64 truncated", float64(255.9), "ff", false},
		{"negative int", int(-1), "-1", false},
		{"large uint64", uint64(16971367927435232718), "eb8669363b3c05ce", false},
		{"unsupported string", "hello", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := hexString(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Arithmetic: add, sub, mul, div, mod
// ---------------------------------------------------------------------------

func TestAdd(t *testing.T) {
	result, err := add(float64(3), float64(4))
	assert.NoError(t, err)
	assert.InDelta(t, 7.0, result, 0.0001)

	result, err = add(int(10), float64(2.5))
	assert.NoError(t, err)
	assert.InDelta(t, 12.5, result, 0.0001)

	_, err = add("a", 1)
	assert.Error(t, err)

	_, err = add(1, "b")
	assert.Error(t, err)
}

func TestSub(t *testing.T) {
	result, err := sub(float64(10), float64(3))
	assert.NoError(t, err)
	assert.InDelta(t, 7.0, result, 0.0001)

	result, err = sub(int(5), float64(2.5))
	assert.NoError(t, err)
	assert.InDelta(t, 2.5, result, 0.0001)

	_, err = sub("a", 1)
	assert.Error(t, err)
}

func TestMul(t *testing.T) {
	result, err := mul(float64(3), float64(4))
	assert.NoError(t, err)
	assert.InDelta(t, 12.0, result, 0.0001)

	result, err = mul(int(5), float64(2.5))
	assert.NoError(t, err)
	assert.InDelta(t, 12.5, result, 0.0001)

	_, err = mul("a", 1)
	assert.Error(t, err)
}

func TestDiv(t *testing.T) {
	result, err := div(float64(10), float64(4))
	assert.NoError(t, err)
	assert.InDelta(t, 2.5, result, 0.0001)

	// Division by zero
	_, err = div(float64(10), float64(0))
	assert.Error(t, err)

	_, err = div(float64(10), int(0))
	assert.Error(t, err)

	_, err = div("a", 1)
	assert.Error(t, err)
}

func TestMod(t *testing.T) {
	result, err := mod(float64(10), float64(3))
	assert.NoError(t, err)
	assert.InDelta(t, 1.0, result, 0.0001)

	result, err = mod(int(7), int(2))
	assert.NoError(t, err)
	assert.InDelta(t, 1.0, result, 0.0001)

	// Mod by zero
	_, err = mod(float64(10), float64(0))
	assert.Error(t, err)

	_, err = mod("a", 1)
	assert.Error(t, err)
}

// ---------------------------------------------------------------------------
// defaultVal
// ---------------------------------------------------------------------------

func TestDefaultVal(t *testing.T) {
	// nil falls back to default
	assert.Equal(t, "fallback", defaultVal("fallback", nil))

	// empty string falls back to default
	assert.Equal(t, "fallback", defaultVal("fallback", ""))

	// non-empty string keeps value
	assert.Equal(t, "actual", defaultVal("fallback", "actual"))

	// zero int falls back to default
	assert.Equal(t, 42, defaultVal(42, 0))

	// non-zero int keeps value
	assert.Equal(t, 7, defaultVal(42, 7))

	// zero float falls back to default
	assert.Equal(t, 3.14, defaultVal(3.14, 0.0))

	// false bool falls back to default
	assert.Equal(t, true, defaultVal(true, false))

	// true bool keeps value
	assert.Equal(t, true, defaultVal(false, true))

	// empty slice falls back to default
	assert.Equal(t, "fallback", defaultVal("fallback", []int{}))

	// non-empty slice keeps value
	val := []int{1, 2, 3}
	assert.Equal(t, val, defaultVal("fallback", val))

	// empty map falls back to default
	assert.Equal(t, "fallback", defaultVal("fallback", map[string]int{}))

	// non-empty map keeps value
	m := map[string]int{"a": 1}
	assert.Equal(t, m, defaultVal("fallback", m))
}

// ---------------------------------------------------------------------------
// seq
// ---------------------------------------------------------------------------

func TestSeq(t *testing.T) {
	// Single arg: 1..n
	result, err := seq(5)
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, result)

	// Single arg negative: -1..n
	result, err = seq(-3)
	assert.NoError(t, err)
	assert.Equal(t, []int{-1, -2, -3}, result)

	// Single arg zero
	result, err = seq(0)
	assert.NoError(t, err)
	assert.Empty(t, result)

	// Two args: start..end
	result, err = seq(3, 7)
	assert.NoError(t, err)
	assert.Equal(t, []int{3, 4, 5, 6, 7}, result)

	// Two args descending
	result, err = seq(5, 2)
	assert.NoError(t, err)
	assert.Equal(t, []int{5, 4, 3, 2}, result)

	// Three args: start, end, step
	result, err = seq(0, 10, 3)
	assert.NoError(t, err)
	assert.Equal(t, []int{0, 3, 6, 9}, result)

	// Three args descending
	result, err = seq(10, 0, -2)
	assert.NoError(t, err)
	assert.Equal(t, []int{10, 8, 6, 4, 2, 0}, result)

	// Step of zero is an error
	_, err = seq(1, 5, 0)
	assert.Error(t, err)

	// No arguments is an error
	_, err = seq()
	assert.Error(t, err)

	// Too many arguments is an error
	_, err = seq(1, 2, 3, 4)
	assert.Error(t, err)
}

// ---------------------------------------------------------------------------
// base64Encode / base64Decode
// ---------------------------------------------------------------------------

func TestBase64Encode(t *testing.T) {
	assert.Equal(t, "aGVsbG8gd29ybGQ=", base64Encode("hello world"))
	assert.Equal(t, "", base64Encode(""))
}

func TestBase64Decode(t *testing.T) {
	result, err := base64Decode("aGVsbG8gd29ybGQ=")
	assert.NoError(t, err)
	assert.Equal(t, "hello world", result)

	result, err = base64Decode("")
	assert.NoError(t, err)
	assert.Equal(t, "", result)

	// Invalid base64
	_, err = base64Decode("!!!not-valid-base64!!!")
	assert.Error(t, err)
}

func TestBase64RoundTrip(t *testing.T) {
	original := "The quick brown fox jumps over the lazy dog 🦊"
	encoded := base64Encode(original)
	decoded, err := base64Decode(encoded)
	assert.NoError(t, err)
	assert.Equal(t, original, decoded)
}

// ---------------------------------------------------------------------------
// regexMatch / regexReplace
// ---------------------------------------------------------------------------

func TestRegexMatch(t *testing.T) {
	matched, err := regexMatch(`^\d+$`, "12345")
	assert.NoError(t, err)
	assert.True(t, matched)

	matched, err = regexMatch(`^\d+$`, "abc")
	assert.NoError(t, err)
	assert.False(t, matched)

	matched, err = regexMatch(`[a-z]+`, "Hello World")
	assert.NoError(t, err)
	assert.True(t, matched)

	// Invalid regex
	_, err = regexMatch(`[invalid`, "test")
	assert.Error(t, err)
}

func TestRegexReplace(t *testing.T) {
	result, err := regexReplace(`\d+`, "NUM", "abc123def456")
	assert.NoError(t, err)
	assert.Equal(t, "abcNUMdefNUM", result)

	result, err = regexReplace(`\s+`, "_", "hello world  foo")
	assert.NoError(t, err)
	assert.Equal(t, "hello_world_foo", result)

	// No match — string unchanged
	result, err = regexReplace(`xyz`, "replaced", "hello")
	assert.NoError(t, err)
	assert.Equal(t, "hello", result)

	// Invalid regex
	_, err = regexReplace(`[invalid`, "rep", "test")
	assert.Error(t, err)
}

// ---------------------------------------------------------------------------
// env
// ---------------------------------------------------------------------------

func TestEnv(t *testing.T) {
	t.Setenv("CODEGEN_TEST_VAR", "test_value")
	assert.Equal(t, "test_value", os.Getenv("CODEGEN_TEST_VAR"))

	// Non-existent env var returns empty string
	assert.Equal(t, "", os.Getenv("CODEGEN_NONEXISTENT_VAR_XYZ"))
}

// ---------------------------------------------------------------------------
// hashString
// ---------------------------------------------------------------------------

func TestHashString(t *testing.T) {
	// Deterministic — same input always produces same hash
	h1 := hashString("hello")
	h2 := hashString("hello")
	assert.Equal(t, h1, h2)

	// Different inputs produce different hashes
	h3 := hashString("world")
	assert.NotEqual(t, h1, h3)

	// Non-zero for non-empty input
	assert.NotZero(t, h1)

	// Empty string still produces a hash
	h4 := hashString("")
	assert.NotZero(t, h4)
}

// ---------------------------------------------------------------------------
// crc32String
// ---------------------------------------------------------------------------

func TestCrc32String(t *testing.T) {
	// Deterministic
	c1 := crc32String("hello")
	c2 := crc32String("hello")
	assert.Equal(t, c1, c2)

	// Known value for "hello"
	// echo -n "hello" | python3 -c "import binascii; print(hex(binascii.crc32(open(0, 'rb').read())))"
	// should be 0x3610a686
	assert.Equal(t, uint32(0x3610a686), c1)

	// Different inputs
	c3 := crc32String("world")
	assert.NotEqual(t, c1, c3)
}

// ---------------------------------------------------------------------------
// fileSize
// ---------------------------------------------------------------------------

func TestFileSize(t *testing.T) {
	// Create a temp file with known content
	f, err := os.CreateTemp("", "filesize_test")
	require.NoError(t, err)
	defer os.Remove(f.Name())

	_, err = f.WriteString("hello world")
	require.NoError(t, err)
	f.Close()

	size, err := fileSize(f.Name())
	assert.NoError(t, err)
	assert.Equal(t, int64(11), size)

	// Non-existent file
	_, err = fileSize("/tmp/nonexistent_file_that_should_not_exist_codegen_test")
	assert.Error(t, err)
}

// ---------------------------------------------------------------------------
// fileHash
// ---------------------------------------------------------------------------

func TestFileHash(t *testing.T) {
	// Create a temp file with known content
	f, err := os.CreateTemp("", "filehash_test")
	require.NoError(t, err)
	defer os.Remove(f.Name())

	_, err = f.WriteString("hello world")
	require.NoError(t, err)
	f.Close()

	hash, err := fileHash(f.Name())
	assert.NoError(t, err)
	assert.NotZero(t, hash)

	// File hash should equal string hash of the same content
	assert.Equal(t, hashString("hello world"), hash)

	// Non-existent file
	_, err = fileHash("/tmp/nonexistent_file_that_should_not_exist_codegen_test")
	assert.Error(t, err)
}

// ---------------------------------------------------------------------------
// readFile
// ---------------------------------------------------------------------------

func TestReadFileFunc(t *testing.T) {
	// Create a temp file with known content
	f, err := os.CreateTemp("", "readfile_test")
	require.NoError(t, err)
	defer os.Remove(f.Name())

	expected := "line1\nline2\nline3"
	_, err = f.WriteString(expected)
	require.NoError(t, err)
	f.Close()

	content, err := readFileFunc(f.Name())
	assert.NoError(t, err)
	assert.Equal(t, expected, content)

	// Non-existent file
	_, err = readFileFunc("/tmp/nonexistent_file_that_should_not_exist_codegen_test")
	assert.Error(t, err)
}

// ---------------------------------------------------------------------------
// UtilFuncs registration check
// ---------------------------------------------------------------------------

func TestUtilFuncsRegistered(t *testing.T) {
	funcs := UtilFuncs()

	expectedKeys := []string{
		"hexString",
		"add", "sub", "mul", "div", "mod",
		"hash", "crc32", "fileHash",
		"fileSize", "readFile",
		"default",
		"seq",
		"base64Encode", "base64Decode",
		"regexMatch", "regexReplace",
		"env",
	}

	for _, key := range expectedKeys {
		_, ok := funcs[key]
		assert.True(t, ok, "UtilFuncs should contain %q", key)
	}
}
