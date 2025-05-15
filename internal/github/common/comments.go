package common

import "fmt"

func FormatComment(file, review string) string {
	return fmt.Sprintf("**File:** %s\n\n%s\n\n", file, review)
}
