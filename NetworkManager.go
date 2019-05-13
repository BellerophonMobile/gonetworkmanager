package gonetworkmanager

import (
	"encoding/json"

	"github.com/godbus/dbus"
)

const (
	NetworkManagerInterface  = "org.freedesktop.NetworkManager"
	NetworkManagerObjectPath = "/org/freedesktop/NetworkManager"

	NetworkManagerGetDevices               = NetworkManagerInterface + ".GetDevices"
	NetworkManagerActivateConnection       = NetworkManagerInterface + ".ActivateConnection"
	NetworkManagerAddAndActivateConnection = NetworkManagerInterface + ".AddAndActivateConnection"
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

	// GetActiveConnections returns the active connection of network devices.
	GetActiveConnections() []ActiveConnection

	// ActivateWirelessConnection requests activating access point to network device
	ActivateWirelessConnection(connection Connection, device Device, accessPoint AccessPoint) ActiveConnection

	// AddAndActivateWirelessConnection adds a new connection profile to the network device it has been
	// passed. It then activates the connection to the passed access point. The first paramter contains
	// additional information for the connection (most propably the credentials).
	// Example contents for connection are:
	// connection := make(map[string]map[string]interface{})
	// connection["802-11-wireless"] = make(map[string]interface{})
	// connection["802-11-wireless"]["security"] = "802-11-wireless-security"
	// connection["802-11-wireless-security"] = make(map[string]interface{})
	// connection["802-11-wireless-security"]["key-mgmt"] = "wpa-psk"
	// connection["802-11-wireless-security"]["psk"] = password
	AddAndActivateWirelessConnection(connection map[string]map[string]interface{}, device Device, accessPoint AccessPoint) (ac ActiveConnection, err error)

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

func (n *networkManager) ActivateWirelessConnection(c Connection, d Device, ap AccessPoint) ActiveConnection {
	var opath dbus.ObjectPath
	n.call(&opath, NetworkManagerActivateConnection, c.GetPath(), d.GetPath(), ap.GetPath())
	return nil
}

func (n *networkManager) AddAndActivateWirelessConnection(connection map[string]map[string]interface{}, d Device, ap AccessPoint) (ac ActiveConnection, err error) {
	var opath1 dbus.ObjectPath
	var opath2 dbus.ObjectPath

	err = n.callError2(&opath1, &opath2, NetworkManagerAddAndActivateConnection, connection, d.GetPath(), ap.GetPath())
	if err != nil {
		return
	}

	ac, err = NewActiveConnection(opath2)
	if err != nil {
		return
	}
	return
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
