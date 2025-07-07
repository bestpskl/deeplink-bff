package session

import (
	"github.com/google/uuid"
)

// Info represents session information for a request.
// It encapsulates request-specific data like request ID and language preferences.
//
// Example usage:
//
//	info := &Info{}
//	info.SetRequestID("")      // Generates new UUID
//	info.SetLanguage("en-US")
type Info struct {
	session session
}

type session struct {
	requestID string
	language  string
}

// MustGetRequestID returns the request ID from the session.
// Panics if request ID is not set.
//
// Example:
//
//	rid := info.MustGetRequestID()
func (i *Info) MustGetRequestID() string {
	if i.session.requestID == "" {
		panic("request id does not exist")
	}
	return i.session.requestID
}

// MustGetLanguage returns the language setting from the session.
// Panics if language is not set.
//
// Example:
//
//	lang := info.MustGetLanguage()
func (i *Info) MustGetLanguage() string {
	if i.session.language == "" {
		panic("language does not exist")
	}
	return i.session.language
}

// SetRequestID sets the request ID in the session.
// If an empty string is provided, generates a new UUID.
//
// Example:
//
//	info.SetRequestID("")      // Generates new UUID
//	info.SetRequestID("123")   // Uses provided ID
func (i *Info) SetRequestID(rid string) {
	if rid == "" {
		i.session.requestID = uuid.New().String()
	} else {
		i.session.requestID = rid
	}
}

// SetLanguage sets the language preference in the session.
//
// Example:
//
//	info.SetLanguage("en-US")
func (i *Info) SetLanguage(language string) {
	i.session.language = language
}
