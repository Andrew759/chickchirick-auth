package config

type contextKey string

const (
	UserUserKey     contextKey = "validateUser"
	UserPropertyKey contextKey = "validateUserProperty"
	UserPhotoKey    contextKey = "validateUserPhoto"
	UserMetaKey     contextKey = "validateUserMeta"
	UserBanKey      contextKey = "validateUserBan"
)
