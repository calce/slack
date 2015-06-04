package slack

import "net/url"

type Engine interface {
	GetName() string
	Do(action string, params url.Values) (interface{}, error)
}