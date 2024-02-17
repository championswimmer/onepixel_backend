package clientinfo

import (
	"errors"
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"net"
	"onepixel_backend/src/db/models"
)

func GetGeoIpDataFromIP(geoipDB *geoip2.Reader, ipAddr string) (*models.GeoIpData, error) {
	ip := net.ParseIP(ipAddr)
	geoIpData := &models.GeoIpData{}
	if ip == nil {
		return nil, errors.New("invalid ip address")
	}
	// if IP exists, populate city, country info
	city, err := geoipDB.City(ip)
	if err != nil {
		return nil, err
	}

	if city.Country.Names["en"] != "" {
		geoIpData.LocationCountry = fmt.Sprintf("%s (%s)", city.Country.Names["en"], city.Country.IsoCode)
	}

	if city.Subdivisions[0].Names["en"] != "" {
		geoIpData.LocationRegion = fmt.Sprintf("%s (%s)", city.Subdivisions[0].Names["en"], city.Subdivisions[0].IsoCode)
	}

	if city.City.Names["en"] != "" {
		geoIpData.LocationCity = city.City.Names["en"]
	}
	return geoIpData, nil
}
