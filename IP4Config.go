package gonetworkmanager

import (
	"encoding/json"

	"github.com/godbus/dbus"
)

const (
	IP4ConfigInterface = NetworkManagerInterface + ".IP4Config"

	IP4ConfigPropertyAddresses   = IP4ConfigInterface + ".Addresses"
	IP4ConfigPropertyRoutes      = IP4ConfigInterface + ".Routes"
	IP4ConfigPropertyNameservers = IP4ConfigInterface + ".Nameservers"
	IP4ConfigPropertyDomains     = IP4ConfigInterface + ".Domains"
)

type IP4Address struct {
	Address string
	Prefix  uint8
	Gateway string
}

type IP4Route struct {
	Route   string
	Prefix  uint8
	NextHop string
	Metric  uint8
}

type IP4Config interface {
	// GetAddresses gets an array of tuples of IPv4 address/prefix/gateway. All 3
	// elements of each tuple are in network byte order. Essentially: [(addr,
	// prefix, gateway), (addr, prefix, gateway), ...]
	GetAddresses() ([]IP4Address, error)

	// GetRoutes gets tuples of IPv4 route/prefix/next-hop/metric. All 4 elements
	// of each tuple are in network byte order. 'route' and 'next hop' are IPv4
	// addresses, while prefix and metric are simple unsigned integers.
	// Essentially: [(route, prefix, next-hop, metric), (route, prefix, next-hop,
	// metric), ...]
	GetRoutes() ([]IP4Route, error)

	// GetNameservers gets the nameservers in use.
	GetNameservers() ([]string, error)

	// GetDomains gets a list of domains this address belongs to.
	GetDomains() ([]string, error)

	MarshalJSON() ([]byte, error)
}

func NewIP4Config(objectPath dbus.ObjectPath) (IP4Config, error) {
	var c ip4Config
	return &c, c.init(NetworkManagerInterface, objectPath)
}

type ip4Config struct {
	dbusBase
}

func (c *ip4Config) GetAddresses() ([]IP4Address, error) {
	addresses, err := c.getSliceSliceUint32Property(IP4ConfigPropertyAddresses)
	if err != nil {
		return nil, err
	}
	ret := make([]IP4Address, len(addresses))

	for i, parts := range addresses {
		ret[i] = IP4Address{
			Address: ip4ToString(parts[0]),
			Prefix:  uint8(parts[1]),
			Gateway: ip4ToString(parts[2]),
		}
	}

	return ret, nil
}

func (c *ip4Config) GetRoutes() ([]IP4Route, error) {
	routes, err := c.getSliceSliceUint32Property(IP4ConfigPropertyRoutes)
	if err != nil {
		return nil, err
	}
	ret := make([]IP4Route, len(routes))

	for i, parts := range routes {
		ret[i] = IP4Route{
			Route:   ip4ToString(parts[0]),
			Prefix:  uint8(parts[1]),
			NextHop: ip4ToString(parts[2]),
			Metric:  uint8(parts[3]),
		}
	}

	return ret, nil
}

func (c *ip4Config) GetNameservers() ([]string, error) {
	nameservers, err := c.getSliceUint32Property(IP4ConfigPropertyNameservers)
	if err != nil {
		return nil, err
	}
	ret := make([]string, len(nameservers))

	for i, ns := range nameservers {
		ret[i] = ip4ToString(ns)
	}

	return ret, nil
}

func (c *ip4Config) GetDomains() ([]string, error) {
	return c.getSliceStringProperty(IP4ConfigPropertyDomains)
}

func (c *ip4Config) MarshalJSON() ([]byte, error) {
	Addresses, err := c.GetAddresses()
	if err != nil {
		return nil, err
	}
	Routes, err := c.GetRoutes()
	if err != nil {
		return nil, err
	}
	Nameservers, err := c.GetNameservers()
	if err != nil {
		return nil, err
	}
	Domains, err := c.GetDomains()
	if err != nil {
		return nil, err
	}

	return json.Marshal(map[string]interface{}{
		"Addresses":   Addresses,
		"Routes":      Routes,
		"Nameservers": Nameservers,
		"Domains":     Domains,
	})
}
