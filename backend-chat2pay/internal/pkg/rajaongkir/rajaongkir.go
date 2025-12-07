package rajaongkir

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const baseURL = "https://api.rajaongkir.com/starter"

type RajaOngkir struct {
	apiKey string
}

func NewRajaOngkir(apiKey string) *RajaOngkir {
	return &RajaOngkir{apiKey: apiKey}
}

type Province struct {
	ProvinceID string `json:"province_id"`
	Province   string `json:"province"`
}

type City struct {
	CityID     string `json:"city_id"`
	ProvinceID string `json:"province_id"`
	Province   string `json:"province"`
	Type       string `json:"type"`
	CityName   string `json:"city_name"`
	PostalCode string `json:"postal_code"`
}

type CostResult struct {
	Code  string        `json:"code"`
	Name  string        `json:"name"`
	Costs []ServiceCost `json:"costs"`
}

type ServiceCost struct {
	Service     string `json:"service"`
	Description string `json:"description"`
	Cost        []Cost `json:"cost"`
}

type Cost struct {
	Value int    `json:"value"`
	Etd   string `json:"etd"`
	Note  string `json:"note"`
}

type ProvinceResponse struct {
	Rajaongkir struct {
		Results []Province `json:"results"`
	} `json:"rajaongkir"`
}

type CityResponse struct {
	Rajaongkir struct {
		Results []City `json:"results"`
	} `json:"rajaongkir"`
}

type CostResponse struct {
	Rajaongkir struct {
		Results []CostResult `json:"results"`
	} `json:"rajaongkir"`
}

func (r *RajaOngkir) GetProvinces() ([]Province, error) {
	req, err := http.NewRequest("GET", baseURL+"/province", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("key", r.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result ProvinceResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.Rajaongkir.Results, nil
}

func (r *RajaOngkir) GetCities(provinceID string) ([]City, error) {
	url := baseURL + "/city"
	if provinceID != "" {
		url += "?province=" + provinceID
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("key", r.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result CityResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.Rajaongkir.Results, nil
}

func (r *RajaOngkir) GetCost(origin, destination string, weight int, courier string) ([]CostResult, error) {
	data := url.Values{}
	data.Set("origin", origin)
	data.Set("destination", destination)
	data.Set("weight", fmt.Sprintf("%d", weight))
	data.Set("courier", courier)

	req, err := http.NewRequest("POST", baseURL+"/cost", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("key", r.apiKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result CostResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.Rajaongkir.Results, nil
}

// GetAllCouriers gets shipping cost from multiple couriers
func (r *RajaOngkir) GetAllCouriers(origin, destination string, weight int) ([]CostResult, error) {
	couriers := []string{"jne", "pos", "tiki"}
	var allResults []CostResult

	for _, courier := range couriers {
		results, err := r.GetCost(origin, destination, weight, courier)
		if err != nil {
			continue
		}
		allResults = append(allResults, results...)
	}

	return allResults, nil
}

// Waybill tracking types
type WaybillDetail struct {
	Waybill   string `json:"waybill_number"`
	Courier   string `json:"courier"`
	Service   string `json:"service_code"`
	Status    string `json:"status"`
	Origin    string `json:"origin"`
	Dest      string `json:"destination"`
	ShipDate  string `json:"waybill_date"`
	Receiver  string `json:"receiver_name"`
	Manifests []Manifest `json:"manifest"`
}

type Manifest struct {
	Description string `json:"manifest_description"`
	Date        string `json:"manifest_date"`
	Time        string `json:"manifest_time"`
	City        string `json:"city_name"`
}

type WaybillResponse struct {
	Rajaongkir struct {
		Result WaybillResult `json:"result"`
		Status struct {
			Code        int    `json:"code"`
			Description string `json:"description"`
		} `json:"status"`
	} `json:"rajaongkir"`
}

type WaybillResult struct {
	Delivered     bool   `json:"delivered"`
	Summary       Summary `json:"summary"`
	Details       Details `json:"details"`
	DeliveryStatus DeliveryStatus `json:"delivery_status"`
	Manifest      []Manifest `json:"manifest"`
}

type Summary struct {
	CourierCode   string `json:"courier_code"`
	CourierName   string `json:"courier_name"`
	WaybillNumber string `json:"waybill_number"`
	ServiceCode   string `json:"service_code"`
	WaybillDate   string `json:"waybill_date"`
	ShipperName   string `json:"shipper_name"`
	ReceiverName  string `json:"receiver_name"`
	Origin        string `json:"origin"`
	Destination   string `json:"destination"`
	Status        string `json:"status"`
}

type Details struct {
	WaybillNumber string `json:"waybill_number"`
	WaybillDate   string `json:"waybill_date"`
	Weight        string `json:"weight"`
	Origin        string `json:"origin"`
	Destination   string `json:"destination"`
	ShipperName   string `json:"shipper_name"`
	ShipperAddr1  string `json:"shipper_address1"`
	ReceiverName  string `json:"receiver_name"`
	ReceiverAddr1 string `json:"receiver_address1"`
}

type DeliveryStatus struct {
	Status      string `json:"status"`
	PodReceiver string `json:"pod_receiver"`
	PodDate     string `json:"pod_date"`
	PodTime     string `json:"pod_time"`
}

// TrackWaybill tracks a shipment by waybill number
func (r *RajaOngkir) TrackWaybill(waybill, courier string) (*WaybillResult, error) {
	data := url.Values{}
	data.Set("waybill", waybill)
	data.Set("courier", strings.ToLower(courier))

	req, err := http.NewRequest("POST", baseURL+"/waybill", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("key", r.apiKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result WaybillResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.Rajaongkir.Status.Code != 200 {
		return nil, fmt.Errorf("tracking failed: %s", result.Rajaongkir.Status.Description)
	}

	return &result.Rajaongkir.Result, nil
}
