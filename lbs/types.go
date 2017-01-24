package lbs

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"strings"
)

const E float64 = 6378.137 // 地球半径，单位为 km

type LngLat struct {
	Longitude float64 `json:"longitude"` //经度
	Latitude  float64 `json:"latitude"`  //纬度
}

func (L LngLat) String() string {
	return fmt.Sprintf("%f,%f", L.Longitude, L.Latitude)
}

func rad(d float64) float64 {
	return d * math.Pi / 180.0
}

// 单位 km
func (L LngLat) Distance(to LngLat) float64 {
	lat1 := rad(L.Latitude)
	lat2 := rad(to.Latitude)
	a := lat1 - lat2
	b := rad(L.Longitude) - rad(to.Longitude)
	s := 2 * math.Asin(math.Sqrt(math.Pow(math.Sin(a/2), 2)+
		math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(b/2), 2)))
	s = s * E
	return s
}

func LngLatFromString(value string) LngLat {
	v := LngLat{}
	vs := strings.Split(value, ",")
	count := len(vs)
	if count > 0 {
		v.Longitude, _ = strconv.ParseFloat(vs[0], 64)
	}
	if count > 1 {
		v.Latitude, _ = strconv.ParseFloat(vs[1], 64)
	}
	return v
}

type Polyline []LngLat

func (L Polyline) String() string {
	w := bytes.NewBuffer(nil)
	for i, v := range L {
		if i != 0 {
			w.WriteString("|")
		}
		w.WriteString(v.String())
	}
	return w.String()
}

func PolylineFromString(value string) Polyline {
	v := Polyline{}
	vs := strings.Split(value, "|")
	for _, vv := range vs {
		v = append(v, LngLatFromString(vv))
	}
	return v
}

type Polygon Polyline

func (P Polygon) In(loc LngLat) bool {

	count := len(P)

	rs := 0

	i, j := 0, count-1

	for i = 0; i < count; i++ {
		v1 := P[i]
		v2 := P[j]
		if ((v1.Longitude < loc.Longitude && v2.Longitude >= loc.Longitude) ||
			(v2.Longitude < loc.Longitude && v1.Longitude >= loc.Longitude)) && (v1.Latitude < loc.Latitude && v2.Latitude >= loc.Latitude) {
			if v1.Latitude+(loc.Longitude-v1.Longitude)/(v2.Longitude-v1.Longitude)*(v2.Latitude-v1.Latitude) < loc.Latitude {
				rs = rs ^ 1
			} else {
				rs = rs ^ 0
			}

			j = i
		}
	}

	return rs != 0
}

type Box struct {
	Min LngLat
	Max LngLat
}

const D = 111.044736

func BoxFromCenter(loc LngLat, distance float64) Box {

	v := Box{}

	lng := distance / math.Abs(math.Cos(rad(loc.Latitude))*D)
	lat := distance / D

	v.Min.Longitude = loc.Longitude - lng
	v.Max.Longitude = loc.Longitude + lng
	v.Min.Latitude = loc.Latitude - lat
	v.Max.Longitude = loc.Longitude + lat

	return v
}
