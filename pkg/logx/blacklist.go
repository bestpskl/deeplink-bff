package logx

import (
	"fmt"
	"strings"
)

// blackList contains a list of sensitive field names that should be redacted in logs.
// The list is case-sensitive and should contain only lowercase values.
// To avoid duplicates and ensure consistent redaction, all entries must be:
//   - Unique (no duplicates allowed)
//   - Lowercase (for consistent matching)
//   - No leading/trailing whitespace
//
// Common sensitive fields include:
//   - Authentication credentials (passwords, tokens)
//   - Personal identifiable information (email, username)
//   - Security keys and secrets

var blackList = map[string]struct{}{
	"password":      {},
	"email":         {},
	"username":      {},
	"token":         {},
	"api_key":       {},
	"access_token":  {},
	"refresh_token": {},
	"authorization": {},
	"cookie":        {},
	"session":       {},
	"jwt":           {},
	"bearer":        {},
	"basic":         {},
	"digest":        {},
	"oauth":         {},
	"client_id":     {},
	"client_secret": {},
	"private_key":   {},
	"public_key":    {},
	"full_name":     {},
	"first_name":    {},
	"last_name":     {},
	"phone":         {},
	"address":       {},
	"city":          {},
	"state":         {},
	"country":       {},
	"zip":           {},
	"postcode":      {},
	"ssn":           {},
	"sin":           {},
	"nino":          {},
	"license":       {},
	"passport":      {},
	"driver":        {},
	"ssn_last4":     {},
	"sin_last4":     {},
	"nino_last4":    {},
	"phone_last4":   {},
	"card":          {},
	"cc":            {},
	"card_number":   {},
	"cc_number":     {},
	"ccn":           {},
	"cvv":           {},
	"cvv2":          {},
	"cv2":           {},
	"expiration":    {},
	"exp_date":      {},
	"exp_month":     {},
	"exp_year":      {},
	"exp":           {},
	"birthdate":     {},
	"birth_date":    {},
	"birth_year":    {},
	"birth_month":   {},
	"birth_day":     {},
	"birth":         {},
	"ssn_hash":      {},
	"sin_hash":      {},
	"nino_hash":     {},
	"phone_hash":    {},
	"cvv_hash":      {},
	"cvv2_hash":     {},
	"cv2_hash":      {},
	"cid":           {},
	"customer_id":   {},
}

// validateBlackList ensures the blacklist contains valid entries.
// It checks for:
//   - Duplicate entries
//   - Non-lowercase entries
//   - Entries with leading/trailing whitespace

func validateBlackList(blackList map[string]struct{}) error {
	seen := make(map[string]struct{})

	for field := range blackList {
		// Check for whitespace
		if strings.TrimSpace(field) != field {
			return fmt.Errorf("blacklist entry '%s' contains leading or trailing whitespace", field)
		}

		// Check for lowercase
		if strings.ToLower(field) != field {
			return fmt.Errorf("blacklist entry '%s' is not lowercase", field)
		}

		// Check for duplicates (technically unnecessary for a map, but retained for completeness)
		if _, exists := seen[field]; exists {
			return fmt.Errorf("duplicate entry '%s' in blacklist", field)
		}
		seen[field] = struct{}{}
	}

	return nil
}
