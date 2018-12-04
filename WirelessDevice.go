package gonetworkmanager

import (
	"encoding/json"

	"github.com/godbus/dbus"
)

const (
	WirelessDeviceInterface = DeviceInterface + ".Wireless"

	WirelessDeviceGetAccessPoints = WirelessDeviceInterface + ".GetAccessPoints"
	WirelessDeviceRequestScan     = WirelessDeviceInterface + ".RequestScan"
)

type WirelessDevice interface {
	Device

	// GetAccessPoints gets the list of access points visible to this device.
	// Note that this list does not include access points which hide their SSID.
	// To retrieve a list of all access points (including hidden ones) use the
	// GetAllAccessPoints() method.
	GetAccessPoints() []AccessPoint

	RequestScan()
}

func NewWirelessDevice(objectPath dbus.ObjectPath) (WirelessDevice, error) {
	var d wirelessDevice
	return &d, d.init(NetworkManagerInterface, objectPath)
}

type wirelessDevice struct {
	device
}

func (d *wirelessDevice) GetAccessPoints() []AccessPoint {
	var apPaths []dbus.ObjectPath

	d.call(&apPaths, WirelessDeviceGetAccessPoints)
	aps := make([]AccessPoint, len(apPaths))

	var err error
	for i, path := range apPaths {
		aps[i], err = NewAccessPoint(path)
		if err != nil {
			panic(err)
		}
	}

	return aps
}

func (d *wirelessDevice) RequestScan() {
	var options map[string]interface{}
	d.obj.Call(WirelessDeviceRequestScan, 0, options).Store()
}

func (d *wirelessDevice) MarshalJSON() ([]byte, error) {
	m := d.device.marshalMap()
	m["AccessPoints"] = d.GetAccessPoints()
	return json.Marshal(m)
}
