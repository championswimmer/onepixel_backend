package db

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
	city := lo.Must(geoipDB.City(net.ParseIP("2402:3a80:43b8:e640:e8d2:94e2:1c54:8847")))

	assert.Equal(t, "IN", city.Country.IsoCode)
	assert.Equal(t, "DL", city.Subdivisions[0].IsoCode)
	assert.Equal(t, "India", city.Country.Names["en"])
	assert.Equal(t, "New Delhi", city.City.Names["en"])

}
