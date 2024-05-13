package providers

// Provider defines current memory store provider
type MemoryStoreProvider interface {
	// SetUserSession sets the user session for given user identifier in form recipe:user_id
	SetUserSession(userId, key, token string) error
	// GetUserSession returns the session token for given token
	GetUserSession(userId, key string) (string, error)
	// DeleteUserSession deletes the user session
	DeleteUserSession(userId, key string) error
	// DeleteAllSessions deletes all the sessions from the session store
	DeleteAllUserSessions(userId string) error
	// DeleteSessionForNamespace deletes the session for a given namespace
	DeleteSessionForNamespace(namespace string) error
	// SetMfaSession sets the mfa session with key and value of userId
	SetMfaSession(userId, key, otp string, expiration int64) error
	// GetMfaSession returns value of given mfa session
	GetMfaSession(userId, key string) (string, error)
	// DeleteMfaSession deletes given mfa session from in-memory store.
	DeleteMfaSession(userId, key string) error
}
