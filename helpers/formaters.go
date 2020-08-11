package helpers

import (
	"unsafe"
	"reflect"
	"fmt"
	"time"
)

//FileSizeFormating gives file size format to the closest unit (B, kB, MB, etc.)
func FileSizeFormating(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

//FormatDate formats a date in Unix (int64) to a string with the format YYYY-MM-DDTHH:MM:SS
func FormatDate(t int64) string {
	expTime := time.Unix(t, 0)
	expTimeString := expTime.Format(time.RFC3339)
	// Remove the last 6 characters (-03:00)
	expTimeString = expTimeString[:len(expTimeString)-6]

	return expTimeString
}

//BytesToString converts a slice of bytes to the corresponding string
func BytesToString(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{bh.Data, bh.Len}
	return *(*string)(unsafe.Pointer(&sh))
}