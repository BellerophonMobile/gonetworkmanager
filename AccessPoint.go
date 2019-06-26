package gonetworkmanager

import (
	"encoding/json"

	"github.com/godbus/dbus"
)

const (
	AccessPointInterface = NetworkManagerInterface + ".AccessPoint"

	AccessPointPropertyFlags      = AccessPointInterface + ".Flags"
	AccessPointPropertyWPAFlags   = AccessPointInterface + ".WpaFlags"
	AccessPointPropertyRSNFlags   = AccessPointInterface + ".RsnFlags"
	AccessPointPropertySSID       = AccessPointInterface + ".Ssid"
	AccessPointPropertyFrequency  = AccessPointInterface + ".Frequency"
	AccessPointPropertyHWAddress  = AccessPointInterface + ".HwAddress"
	AccessPointPropertyMode       = AccessPointInterface + ".Mode"
	AccessPointPropertyMaxBitrate = AccessPointInterface + ".MaxBitrate"
	AccessPointPropertyStrength   = AccessPointInterface + ".Strength"
)

type AccessPoint interface {
	GetPath() dbus.ObjectPath

	// GetFlags gets flags describing the capabilities of the access point.
	GetFlags() (uint32, error)

	// GetWPAFlags gets flags describing the access point's capabilities
	// according to WPA (Wifi Protected Access).
	GetWPAFlags() (uint32, error)

	// GetRSNFlags gets flags describing the access point's capabilities
	// according to the RSN (Robust Secure Network) protocol.
	GetRSNFlags() (uint32, error)

	// GetSSID returns the Service Set Identifier identifying the access point.
	GetSSID() (string, error)

	// GetFrequency gets the radio channel frequency in use by the access point,
	// in MHz.
	GetFrequency() (uint32, error)

	// GetHWAddress gets the hardware address (BSSID) of the access point.
	GetHWAddress() (string, error)

	// GetMode describes the operating mode of the access point.
	GetMode() (Nm80211Mode, error)

	// GetMaxBitrate gets the maximum bitrate this access point is capable of, in
	// kilobits/second (Kb/s).
	GetMaxBitrate() (uint32, error)

	// GetStrength gets the current signal quality of the access point, in
	// percent.
	GetStrength() (uint8, error)

	MarshalJSON() ([]byte, error)
}

func NewAccessPoint(objectPath dbus.ObjectPath) (AccessPoint, error) {
	var a accessPoint
	return &a, a.init(NetworkManagerInterface, objectPath)
}

type accessPoint struct {
	dbusBase
}

func (a *accessPoint) GetPath() dbus.ObjectPath {
	return a.obj.Path()
}

func (a *accessPoint) GetFlags() (uint32, error) {
	return a.getUint32Property(AccessPointPropertyFlags)
}

func (a *accessPoint) GetWPAFlags() (uint32, error) {
	return a.getUint32Property(AccessPointPropertyWPAFlags)
}

func (a *accessPoint) GetRSNFlags() (uint32, error) {
	return a.getUint32Property(AccessPointPropertyRSNFlags)
}

func (a *accessPoint) GetSSID() (string, error) {
	r, err := a.getSliceByteProperty(AccessPointPropertySSID)
	if err != nil {
		return "", err
	}
	return string(r), nil
}

func (a *accessPoint) GetFrequency() (uint32, error) {
	return a.getUint32Property(AccessPointPropertyFrequency)
}

func (a *accessPoint) GetHWAddress() (string, error) {
	return a.getStringProperty(AccessPointPropertyHWAddress)
}

func (a *accessPoint) GetMode() (Nm80211Mode, error) {
	r, err := a.getUint32Property(AccessPointPropertyMode)
	if err != nil {
		return Nm80211ModeUnknown, err
	}
	return Nm80211Mode(r), nil
}

func (a *accessPoint) GetMaxBitrate() (uint32, error) {
	return a.getUint32Property(AccessPointPropertyMaxBitrate)
}

func (a *accessPoint) GetStrength() (uint8, error) {
	return a.getUint8Property(AccessPointPropertyStrength)
}

func (a *accessPoint) MarshalJSON() ([]byte, error) {
	Flags, err := a.GetFlags()
	if err != nil {
		return nil, err
	}
	WPAFlags, err := a.GetWPAFlags()
	if err != nil {
		return nil, err
	}
	RSNFlags, err := a.GetRSNFlags()
	if err != nil {
		return nil, err
	}
	SSID, err := a.GetSSID()
	if err != nil {
		return nil, err
	}
	Frequency, err := a.GetFrequency()
	if err != nil {
		return nil, err
	}
	HWAddress, err := a.GetHWAddress()
	if err != nil {
		return nil, err
	}
	Mode, err := a.GetMode()
	if err != nil {
		return nil, err
	}
	MaxBitrate, err := a.GetMaxBitrate()
	if err != nil {
		return nil, err
	}
	Strength, err := a.GetStrength()
	if err != nil {
		return nil, err
	}

	return json.Marshal(map[string]interface{}{
		"Flags":      Flags,
		"WPAFlags":   WPAFlags,
		"RSNFlags":   RSNFlags,
		"SSID":       SSID,
		"Frequency":  Frequency,
		"HWAddress":  HWAddress,
		"Mode":       Mode.String(),
		"MaxBitrate": MaxBitrate,
		"Strength":   Strength,
	})
}
