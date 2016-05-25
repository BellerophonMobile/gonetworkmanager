package gonetworkmanager

import (
	"github.com/godbus/dbus"
)

const (
	AccessPointInterface = NetworkManagerInterface + ".AccessPoint"

	AccessPointPropertySSID = AccessPointInterface + ".Ssid"
)

type AccessPoint interface {
	// GetSSID returns the Service Set Identifier identifying the access point.
	GetSSID() string
}

func NewAccessPoint(objectPath dbus.ObjectPath) (AccessPoint, error) {
	var a accessPoint
	return &a, a.init(NetworkManagerInterface, objectPath)
}

type accessPoint struct {
	dbusBase
}

func (a *accessPoint) GetSSID() string {
	return string(a.getSliceByteProperty(AccessPointPropertySSID))
}
