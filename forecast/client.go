package forecast

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mitchellh/mapstructure"
)

type js map[string]interface{}

type Response struct {
	OperationalMode   string `mapstructure:"operationalMode"`
	SrsName           string `mapstructure:"srsName"`
	CreationDate      string `mapstructure:"creationDate"`
	CreationDateLocal string `mapstructure:"creationDateLocal"`
	ProductionCenter  string `mapstructure:"productionCenter"`
	Credit            string `mapstructure:"credit"`
	MoreInformation   string `mapstructure:"moreInformation"`
	Location          struct {
		Latitude        float64 `mapstructure:"latitude"`
		Longitude       float64 `mapstructure:"longitude"`
		Elevation       int     `mapstructure:"elevation"`
		WFO             string  `mapstructure:"wfo"`
		Timezone        string  `mapstructure:"timezone"`
		AreaDescription string  `mapstructure:"areaDescription"`
		Radar           string  `mapstructure:"radar"`
		Zone            string  `mapstructure:"zone"`
		County          string  `mapstructure:"county"`
		FireZone        string  `mapstructure:"firezone"`
		Metar           string  `mapstructure:"metar"`
	}
	Time struct {
		LayoutKey       string   `mapstructure:"layoutKey"`
		StartPeriodName []string `mapstructure:"startPeriodName"`
		StartValidTime  []string `mapstructure:"startValidTime"`
		TempLabel       []string `mapstructure:"tempLabel"`
	}
	Data struct {
		Temperature []int     `mapstructure:"temperature"`
		POP         []float64 `mapstructure:"pop"`
		Weather     []string  `mapstructure:"weather"`
		IconLink    []string  `mapstructure:"iconLink"`
		Hazard      []string  `mapstructure:"hazard"`
		HazardURL   []string  `mapstructure:"hazardUrl"`
		Text        []string  `mapstructure:"text"`
	}
	Currentobservation struct {
		ID           string  `mapstructure:"id"`
		Name         string  `mapstructure:"name"`
		Elev         int     `mapstructure:"elev"`
		Latitude     float64 `mapstructure:"latitude"`
		Longitude    float64 `mapstructure:"longitude"`
		Date         string  `mapstructure:"Date"`
		Temp         int     `mapstructure:"Temp"`
		Dewp         int     `mapstructure:"Dewp"`
		Relh         int     `mapstructure:"Relh"`
		Winds        int     `mapstructure:"Winds"`
		WindD        int     `mapstructure:"Windd"`
		Gust         string  `mapstructure:"Gust"`
		Weather      string  `mapstructure:"Weather"`
		Weatherimage string  `mapstructure:"Weatherimage"`
		Visibility   float64 `mapstructure:"Visibility"`
		Altimeter    float64 `mapstructure:"Altimeter"`
		SLP          float64 `mapstructure:"SLP"`
		Timezone     string  `mapstructure:"timezone"`
		State        string  `mapstructure:"state"`
		WindChill    int     `mapstructure:"WindChill"`
	}
}

type Client struct {
	base string
	lang string
	unit int
}

func (c *Client) Metric() *Client {
	c.unit = 1
	return c
}

func (c *Client) Imperial() *Client {
	c.unit = 0
	return c
}

func (c *Client) Lang(lang string) *Client {
	c.lang = lang
	return c
}

func (c *Client) Get(lat, lon float64) (*Response, error) {
	var j js
	var rval Response
	resp, err := http.Get(fmt.Sprintf(c.base, lat, lon, c.unit, c.lang))
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if string(data[:1]) != "{" {
		// NOAA does this thing where they return a 200 OK
		// with just a <script> with a window.location
		// redirect in it when there's a problem (like no
		// data for that lat/lon pair...
		return nil, errors.New("no data available")
	}
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &j)
	if err != nil {
		return nil, err
	}
	decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           &rval,
	})
	err = decoder.Decode(j)
	if err != nil {
		return nil, err
	}
	return &rval, nil
}

func New() *Client {
	return &Client{
		base: "http://forecast.weather.gov/MapClick.php?lat=%f&lon=%f&unit=%d&lg=%s&FcstType=json",
		lang: "english",
		unit: 0,
	}
}
