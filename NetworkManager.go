package gonetworkmanager

import (
	"encoding/json"

	"github.com/godbus/dbus"
)

const (
	NetworkManagerInterface  = "org.freedesktop.NetworkManager"
	NetworkManagerObjectPath = "/org/freedesktop/NetworkManager"

	NetworkManagerGetDevices               = NetworkManagerInterface + ".GetDevices"
	NetworkManagerPropertyState            = NetworkManagerInterface + ".State"
	NetworkManagerPropertyActiveConnection = NetworkManagerInterface + ".ActiveConnections"
)

type NetworkManager interface {

	// GetDevices gets the list of network devices.
	GetDevices() []Device

	// GetState returns the overall networking state as determined by the
	// NetworkManager daemon, based on the state of network devices under it's
	// management.
	GetState() NmState

	GetActiveConnections() []ActiveConnection

	Subscribe() <-chan *dbus.Signal
	Unsubscribe()

	MarshalJSON() ([]byte, error)
}

func NewNetworkManager() (NetworkManager, error) {
	var nm networkManager
	return &nm, nm.init(NetworkManagerInterface, NetworkManagerObjectPath)
}

type networkManager struct {
	dbusBase

	sigChan chan *dbus.Signal
}

func (n *networkManager) GetDevices() []Device {
	var devicePaths []dbus.ObjectPath

	n.call(&devicePaths, NetworkManagerGetDevices)
	devices := make([]Device, len(devicePaths))

	var err error
	for i, path := range devicePaths {
		devices[i], err = DeviceFactory(path)
		if err != nil {
			panic(err)
		}
	}

	return devices
}

func (n *networkManager) GetState() NmState {
	return NmState(n.getUint32Property(NetworkManagerPropertyState))
}

func (n *networkManager) GetActiveConnections() []ActiveConnection {
	acPaths := n.getSliceObjectProperty(NetworkManagerPropertyActiveConnection)
	ac := make([]ActiveConnection, len(acPaths))

	var err error
	for i, path := range acPaths {
		ac[i], err = NewActiveConnection(path)
		if err != nil {
			panic(err)
		}
	}

	return ac
}

func (n *networkManager) Subscribe() <-chan *dbus.Signal {
	if n.sigChan != nil {
		return n.sigChan
	}

	n.subscribeNamespace(NetworkManagerObjectPath)
	n.sigChan = make(chan *dbus.Signal, 10)
	n.conn.Signal(n.sigChan)

	return n.sigChan
}

func (n *networkManager) Unsubscribe() {
	n.conn.RemoveSignal(n.sigChan)
	n.sigChan = nil
}

func (n *networkManager) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"NetworkState": n.GetState().String(),
		"Devices":      n.GetDevices(),
	})
}
