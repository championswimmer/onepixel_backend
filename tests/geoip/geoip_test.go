package geoip

import (
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"net"
	"onepixel_backend/src/db"
	"testing"
)

func Test_GeoIPResolutionIPV4(t *testing.T) {
	geoipDB := db.GetGeoIPDB()

	// Test IPV4
	city := lo.Must(geoipDB.City(net.ParseIP("42.108.28.82")))

	assert.Equal(t, "IN", city.Country.IsoCode)
	assert.Equal(t, "DL", city.Subdivisions[0].IsoCode)
	assert.Equal(t, "India", city.Country.Names["en"])
	assert.Equal(t, "Delhi", city.City.Names["en"])

}

func Test_GeoIPResolutionIPV6(t *testing.T) {
	geoipDB := db.GetGeoIPDB()

	// Test IPV4
	city := lo.Must(geoipDB.City(net.ParseIP("2406:7400:63:66d8:48f5:9eed:2a3b:f286")))

	assert.Equal(t, "IN", city.Country.IsoCode)
	assert.Equal(t, "KA", city.Subdivisions[0].IsoCode)
	assert.Equal(t, "India", city.Country.Names["en"])
	assert.Equal(t, "Bengaluru", city.City.Names["en"])

}
