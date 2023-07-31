package ip2loc

import (
	"github.com/daqiancode/env"
	"github.com/ip2location/ip2location-go/v9"
	"github.com/scilive/scibase/logs"
)

var db *ip2location.DB

func init() {
	var err error
	dbFile := env.Get("IP2LOCATION_DB")
	db, err = ip2location.OpenDB(dbFile)
	if err != nil {
		logs.Log.Error().Err(err).Str(dbFile, env.Get("IP2LOCATION_DB")).Msg("failed to load ip2location db")
	}
}

type Location struct {
	Country string `json:"country"`
	Region  string `json:"region"`
	City    string `json:"city"`
}

func Ip2Loc(ip string) (Location, error) {
	r, err := db.Get_all(ip)
	var loc Location
	if err != nil {
		return loc, err
	}
	loc.Country = r.Country_short
	loc.Region = r.Region
	loc.City = r.City
	return loc, nil

}
