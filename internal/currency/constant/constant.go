package constant

import "errors"

type Event string

type AuthorizationType string

const (
	CURRENCY_ADDED_EVENT Event = "currency_added_event"

	EVENT_IDEMPOTENT_ID = "idempotent_id"

	AccessToken AuthorizationType = "access-token"

	CURRENCY_CACHE_KEY = "CURRENCIES"
)

var (
	ErrNotFound error = errors.New("not found")
	ErrConflict error = errors.New("conflict")
)
