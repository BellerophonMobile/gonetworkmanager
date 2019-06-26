package gonetworkmanager

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/godbus/dbus"
)

const (
	dbusMethodAddMatch = "org.freedesktop.DBus.AddMatch"
)

type dbusBase struct {
	conn *dbus.Conn
	obj  dbus.BusObject
}

func (d *dbusBase) init(iface string, objectPath dbus.ObjectPath) error {
	var err error

	d.conn, err = dbus.SystemBus()
	if err != nil {
		return err
	}

	d.obj = d.conn.Object(iface, objectPath)

	return nil
}

func (d *dbusBase) call(value interface{}, method string, args ...interface{}) error {
	return d.obj.Call(method, 0, args...).Store(value)
}

func (d *dbusBase) call2(value1 interface{}, value2 interface{}, method string, args ...interface{}) error {
	return d.obj.Call(method, 0, args...).Store(value1, value2)
}

func (d *dbusBase) subscribe(iface, member string) {
	rule := fmt.Sprintf("type='signal',interface='%s',path='%s',member='%s'",
		iface, d.obj.Path(), NetworkManagerInterface)
	d.conn.BusObject().Call(dbusMethodAddMatch, 0, rule)
}

func (d *dbusBase) subscribeNamespace(namespace string) {
	rule := fmt.Sprintf("type='signal',path_namespace='%s'", namespace)
	d.conn.BusObject().Call(dbusMethodAddMatch, 0, rule)
}

func (d *dbusBase) getProperty(iface string) (interface{}, error) {
	variant, err := d.obj.GetProperty(iface)
	if err != nil {
		return nil, err
	}
	return variant.Value(), nil
}

func (d *dbusBase) getObjectProperty(iface string) (dbus.ObjectPath, error) {
	value, err := d.getProperty(iface)
	if err != nil {
		return "", makeErrVariantType(iface)
	}
	return value.(dbus.ObjectPath), nil
}

func (d *dbusBase) getSliceObjectProperty(iface string) ([]dbus.ObjectPath, error) {
	value, err := d.getProperty(iface)
	if err != nil {
		return nil, makeErrVariantType(iface)
	}
	return value.([]dbus.ObjectPath), nil
}

func (d *dbusBase) getStringProperty(iface string) (string, error) {
	value, err := d.getProperty(iface)
	if err != nil {
		return "", makeErrVariantType(iface)
	}
	return value.(string), nil
}

func (d *dbusBase) getSliceStringProperty(iface string) ([]string, error) {
	value, err := d.getProperty(iface)
	if err != nil {
		return nil, makeErrVariantType(iface)
	}
	return value.([]string), nil
}

func (d *dbusBase) getMapStringVariantProperty(iface string) (map[string]dbus.Variant, error) {
	value, err := d.getProperty(iface)
	if err != nil {
		return nil, makeErrVariantType(iface)
	}
	return value.(map[string]dbus.Variant), nil
}

func (d *dbusBase) getUint8Property(iface string) (uint8, error) {
	value, err := d.getProperty(iface)
	if err != nil {
		return 0, makeErrVariantType(iface)
	}
	return value.(uint8), nil
}

func (d *dbusBase) getUint32Property(iface string) (uint32, error) {
	value, err := d.getProperty(iface)
	if err != nil {
		return 0, makeErrVariantType(iface)
	}
	return value.(uint32), nil
}

func (d *dbusBase) getSliceUint32Property(iface string) ([]uint32, error) {
	value, err := d.getProperty(iface)
	if err != nil {
		return nil, makeErrVariantType(iface)
	}
	return value.([]uint32), nil
}

func (d *dbusBase) getSliceSliceUint32Property(iface string) ([][]uint32, error) {
	value, err := d.getProperty(iface)
	if err != nil {
		return nil, makeErrVariantType(iface)
	}
	return value.([][]uint32), nil
}

func (d *dbusBase) getSliceByteProperty(iface string) ([]byte, error) {
	value, err := d.getProperty(iface)
	if err != nil {
		return nil, makeErrVariantType(iface)
	}
	return value.([]byte), nil
}

func makeErrVariantType(iface string) error {
	return fmt.Errorf("unexpected variant type for '%s'", iface)
}

func ip4ToString(ip uint32) string {
	bs := []byte{0, 0, 0, 0}
	binary.LittleEndian.PutUint32(bs, ip)
	return net.IP(bs).String()
}
