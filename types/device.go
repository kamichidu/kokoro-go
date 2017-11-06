package types

import (
	"time"
)

type DeviceKind string

const (
	DeviceUnknown     DeviceKind = "unknown"
	DeviceIOS                    = "ios"
	DeviceAndroid                = "android"
	DeviceUWP                    = "uwp"
	DeviceChrome                 = "chrome"
	DeviceFirefox                = "firefox"
	DeviceOfficialWeb            = "official_web"
)

type Device struct {
	// Device name
	Name string `json:"name,omitempty"`

	// Kind of device
	Kind DeviceKind `json:"kind,omitempty"`

	// Unique identifier for classification a device
	DeviceIdentifier string `json:"device_identifier,omitempty"`

	NotificationIdentifier string `json:"notification_identifier,omitempty"`

	SubscribeNotification bool `json:"subscribe_notification,omitempty"`

	LastActivityAt time.Time `json:"last_activity_at,omitempty"`

	PushRegistered bool `json:"push_registered,omitempty"`

	// Access token for this device
	AccessToken *AccessToken `json:"access_token,omitempty"`
}
