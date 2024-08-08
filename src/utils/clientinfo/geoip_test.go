package clientinfo

import (
	"onepixel_backend/src/db"
	"testing"

	"github.com/oschwald/geoip2-golang"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

var geoipDB *geoip2.Reader

func init() {
	geoipDB = db.GetGeoIPDB()
}

func Test_GeoIPResolutionIPV4(t *testing.T) {

	// Test IPV4
	geoIpData := lo.Must(GetGeoIpDataFromIP(geoipDB, "202.83.26.113"))

	assert.Equal(t, "India (IN)", geoIpData.LocationCountry)
	assert.Equal(t, "Karnataka (KA)", geoIpData.LocationRegion)
	assert.Equal(t, "Bengaluru", geoIpData.LocationCity)

}

func Test_GeoIPResolutionIPV6(t *testing.T) {

	// Test IPV6
	geoIpData := lo.Must(GetGeoIpDataFromIP(geoipDB, "2406:7400:63:66d8:48f5:9eed:2a3b:f286"))

	assert.Equal(t, "India (IN)", geoIpData.LocationCountry)
	assert.Equal(t, "Karnataka (KA)", geoIpData.LocationRegion)
	assert.Equal(t, "Bengaluru", geoIpData.LocationCity)

}
