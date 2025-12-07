package handlers

import (
	"chat2pay/internal/api/presenter"
	"chat2pay/internal/pkg/rajaongkir"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type ShippingHandler struct {
	rajaOngkir *rajaongkir.RajaOngkir
	useAPI     bool
}

func NewShippingHandler(apiKey string) *ShippingHandler {
	handler := &ShippingHandler{
		useAPI: apiKey != "",
	}
	if apiKey != "" {
		handler.rajaOngkir = rajaongkir.NewRajaOngkir(apiKey)
	}
	return handler
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

type ShippingCost struct {
	Service     string `json:"service"`
	Description string `json:"description"`
	Cost        []Cost `json:"cost"`
}

type Cost struct {
	Value int    `json:"value"`
	Etd   string `json:"etd"`
	Note  string `json:"note"`
}

type CostResult struct {
	Code  string         `json:"code"`
	Name  string         `json:"name"`
	Costs []ShippingCost `json:"costs"`
}

// Mock data for provinces
var mockProvinces = []Province{
	{ProvinceID: "1", Province: "Bali"},
	{ProvinceID: "2", Province: "Bangka Belitung"},
	{ProvinceID: "3", Province: "Banten"},
	{ProvinceID: "4", Province: "Bengkulu"},
	{ProvinceID: "5", Province: "DI Yogyakarta"},
	{ProvinceID: "6", Province: "DKI Jakarta"},
	{ProvinceID: "7", Province: "Gorontalo"},
	{ProvinceID: "8", Province: "Jambi"},
	{ProvinceID: "9", Province: "Jawa Barat"},
	{ProvinceID: "10", Province: "Jawa Tengah"},
	{ProvinceID: "11", Province: "Jawa Timur"},
	{ProvinceID: "12", Province: "Kalimantan Barat"},
	{ProvinceID: "13", Province: "Kalimantan Selatan"},
	{ProvinceID: "14", Province: "Kalimantan Tengah"},
	{ProvinceID: "15", Province: "Kalimantan Timur"},
	{ProvinceID: "16", Province: "Kalimantan Utara"},
	{ProvinceID: "17", Province: "Kepulauan Riau"},
	{ProvinceID: "18", Province: "Lampung"},
	{ProvinceID: "19", Province: "Maluku"},
	{ProvinceID: "20", Province: "Maluku Utara"},
	{ProvinceID: "21", Province: "Nanggroe Aceh Darussalam"},
	{ProvinceID: "22", Province: "Nusa Tenggara Barat"},
	{ProvinceID: "23", Province: "Nusa Tenggara Timur"},
	{ProvinceID: "24", Province: "Papua"},
	{ProvinceID: "25", Province: "Papua Barat"},
	{ProvinceID: "26", Province: "Riau"},
	{ProvinceID: "27", Province: "Sulawesi Barat"},
	{ProvinceID: "28", Province: "Sulawesi Selatan"},
	{ProvinceID: "29", Province: "Sulawesi Tengah"},
	{ProvinceID: "30", Province: "Sulawesi Tenggara"},
	{ProvinceID: "31", Province: "Sulawesi Utara"},
	{ProvinceID: "32", Province: "Sumatera Barat"},
	{ProvinceID: "33", Province: "Sumatera Selatan"},
	{ProvinceID: "34", Province: "Sumatera Utara"},
}

// Mock data for cities (sample)
var mockCities = map[string][]City{
	"1": { // Bali
		{CityID: "17", ProvinceID: "1", Province: "Bali", Type: "Kabupaten", CityName: "Badung", PostalCode: "80351"},
		{CityID: "32", ProvinceID: "1", Province: "Bali", Type: "Kabupaten", CityName: "Bangli", PostalCode: "80619"},
		{CityID: "94", ProvinceID: "1", Province: "Bali", Type: "Kabupaten", CityName: "Buleleng", PostalCode: "81111"},
		{CityID: "114", ProvinceID: "1", Province: "Bali", Type: "Kota", CityName: "Denpasar", PostalCode: "80227"},
		{CityID: "128", ProvinceID: "1", Province: "Bali", Type: "Kabupaten", CityName: "Gianyar", PostalCode: "80519"},
	},
	"2": { // Bangka Belitung
		{CityID: "27", ProvinceID: "2", Province: "Bangka Belitung", Type: "Kabupaten", CityName: "Bangka", PostalCode: "33212"},
		{CityID: "28", ProvinceID: "2", Province: "Bangka Belitung", Type: "Kabupaten", CityName: "Bangka Barat", PostalCode: "33315"},
		{CityID: "56", ProvinceID: "2", Province: "Bangka Belitung", Type: "Kabupaten", CityName: "Belitung", PostalCode: "33419"},
		{CityID: "334", ProvinceID: "2", Province: "Bangka Belitung", Type: "Kota", CityName: "Pangkal Pinang", PostalCode: "33115"},
	},
	"3": { // Banten
		{CityID: "106", ProvinceID: "3", Province: "Banten", Type: "Kota", CityName: "Cilegon", PostalCode: "42417"},
		{CityID: "232", ProvinceID: "3", Province: "Banten", Type: "Kabupaten", CityName: "Lebak", PostalCode: "42319"},
		{CityID: "331", ProvinceID: "3", Province: "Banten", Type: "Kabupaten", CityName: "Pandeglang", PostalCode: "42212"},
		{CityID: "402", ProvinceID: "3", Province: "Banten", Type: "Kota", CityName: "Serang", PostalCode: "42111"},
		{CityID: "403", ProvinceID: "3", Province: "Banten", Type: "Kabupaten", CityName: "Serang", PostalCode: "42182"},
		{CityID: "455", ProvinceID: "3", Province: "Banten", Type: "Kabupaten", CityName: "Tangerang", PostalCode: "15914"},
		{CityID: "456", ProvinceID: "3", Province: "Banten", Type: "Kota", CityName: "Tangerang", PostalCode: "15111"},
		{CityID: "457", ProvinceID: "3", Province: "Banten", Type: "Kota", CityName: "Tangerang Selatan", PostalCode: "15332"},
	},
	"4": { // Bengkulu
		{CityID: "62", ProvinceID: "4", Province: "Bengkulu", Type: "Kota", CityName: "Bengkulu", PostalCode: "38229"},
		{CityID: "175", ProvinceID: "4", Province: "Bengkulu", Type: "Kabupaten", CityName: "Kaur", PostalCode: "38911"},
		{CityID: "183", ProvinceID: "4", Province: "Bengkulu", Type: "Kabupaten", CityName: "Kepahiang", PostalCode: "39172"},
	},
	"5": { // DI Yogyakarta
		{CityID: "39", ProvinceID: "5", Province: "DI Yogyakarta", Type: "Kabupaten", CityName: "Bantul", PostalCode: "55715"},
		{CityID: "135", ProvinceID: "5", Province: "DI Yogyakarta", Type: "Kabupaten", CityName: "Gunung Kidul", PostalCode: "55812"},
		{CityID: "210", ProvinceID: "5", Province: "DI Yogyakarta", Type: "Kabupaten", CityName: "Kulon Progo", PostalCode: "55611"},
		{CityID: "419", ProvinceID: "5", Province: "DI Yogyakarta", Type: "Kabupaten", CityName: "Sleman", PostalCode: "55513"},
		{CityID: "501", ProvinceID: "5", Province: "DI Yogyakarta", Type: "Kota", CityName: "Yogyakarta", PostalCode: "55222"},
	},
	"6": { // DKI Jakarta
		{CityID: "151", ProvinceID: "6", Province: "DKI Jakarta", Type: "Kota", CityName: "Jakarta Barat", PostalCode: "11220"},
		{CityID: "152", ProvinceID: "6", Province: "DKI Jakarta", Type: "Kota", CityName: "Jakarta Pusat", PostalCode: "10540"},
		{CityID: "153", ProvinceID: "6", Province: "DKI Jakarta", Type: "Kota", CityName: "Jakarta Selatan", PostalCode: "12230"},
		{CityID: "154", ProvinceID: "6", Province: "DKI Jakarta", Type: "Kota", CityName: "Jakarta Timur", PostalCode: "13330"},
		{CityID: "155", ProvinceID: "6", Province: "DKI Jakarta", Type: "Kota", CityName: "Jakarta Utara", PostalCode: "14140"},
		{CityID: "189", ProvinceID: "6", Province: "DKI Jakarta", Type: "Kabupaten", CityName: "Kepulauan Seribu", PostalCode: "14550"},
	},
	"7": { // Gorontalo
		{CityID: "77", ProvinceID: "7", Province: "Gorontalo", Type: "Kabupaten", CityName: "Bone Bolango", PostalCode: "96511"},
		{CityID: "129", ProvinceID: "7", Province: "Gorontalo", Type: "Kabupaten", CityName: "Gorontalo", PostalCode: "96218"},
		{CityID: "130", ProvinceID: "7", Province: "Gorontalo", Type: "Kota", CityName: "Gorontalo", PostalCode: "96115"},
	},
	"8": { // Jambi
		{CityID: "40", ProvinceID: "8", Province: "Jambi", Type: "Kabupaten", CityName: "Batang Hari", PostalCode: "36613"},
		{CityID: "156", ProvinceID: "8", Province: "Jambi", Type: "Kota", CityName: "Jambi", PostalCode: "36137"},
		{CityID: "193", ProvinceID: "8", Province: "Jambi", Type: "Kabupaten", CityName: "Kerinci", PostalCode: "37167"},
	},
	"9": { // Jawa Barat
		{CityID: "22", ProvinceID: "9", Province: "Jawa Barat", Type: "Kabupaten", CityName: "Bandung", PostalCode: "40311"},
		{CityID: "23", ProvinceID: "9", Province: "Jawa Barat", Type: "Kota", CityName: "Bandung", PostalCode: "40111"},
		{CityID: "24", ProvinceID: "9", Province: "Jawa Barat", Type: "Kabupaten", CityName: "Bandung Barat", PostalCode: "40721"},
		{CityID: "34", ProvinceID: "9", Province: "Jawa Barat", Type: "Kota", CityName: "Banjar", PostalCode: "46311"},
		{CityID: "54", ProvinceID: "9", Province: "Jawa Barat", Type: "Kabupaten", CityName: "Bekasi", PostalCode: "17837"},
		{CityID: "55", ProvinceID: "9", Province: "Jawa Barat", Type: "Kota", CityName: "Bekasi", PostalCode: "17121"},
		{CityID: "79", ProvinceID: "9", Province: "Jawa Barat", Type: "Kabupaten", CityName: "Bogor", PostalCode: "16911"},
		{CityID: "80", ProvinceID: "9", Province: "Jawa Barat", Type: "Kota", CityName: "Bogor", PostalCode: "16119"},
		{CityID: "103", ProvinceID: "9", Province: "Jawa Barat", Type: "Kabupaten", CityName: "Ciamis", PostalCode: "46211"},
		{CityID: "104", ProvinceID: "9", Province: "Jawa Barat", Type: "Kabupaten", CityName: "Cianjur", PostalCode: "43217"},
		{CityID: "115", ProvinceID: "9", Province: "Jawa Barat", Type: "Kota", CityName: "Cimahi", PostalCode: "40512"},
		{CityID: "107", ProvinceID: "9", Province: "Jawa Barat", Type: "Kabupaten", CityName: "Cirebon", PostalCode: "45611"},
		{CityID: "108", ProvinceID: "9", Province: "Jawa Barat", Type: "Kota", CityName: "Cirebon", PostalCode: "45116"},
		{CityID: "126", ProvinceID: "9", Province: "Jawa Barat", Type: "Kota", CityName: "Depok", PostalCode: "16416"},
		{CityID: "149", ProvinceID: "9", Province: "Jawa Barat", Type: "Kabupaten", CityName: "Indramayu", PostalCode: "45214"},
		{CityID: "171", ProvinceID: "9", Province: "Jawa Barat", Type: "Kabupaten", CityName: "Karawang", PostalCode: "41311"},
		{CityID: "211", ProvinceID: "9", Province: "Jawa Barat", Type: "Kabupaten", CityName: "Kuningan", PostalCode: "45511"},
		{CityID: "252", ProvinceID: "9", Province: "Jawa Barat", Type: "Kabupaten", CityName: "Majalengka", PostalCode: "45412"},
		{CityID: "332", ProvinceID: "9", Province: "Jawa Barat", Type: "Kabupaten", CityName: "Pangandaran", PostalCode: "46511"},
		{CityID: "376", ProvinceID: "9", Province: "Jawa Barat", Type: "Kabupaten", CityName: "Purwakarta", PostalCode: "41119"},
		{CityID: "428", ProvinceID: "9", Province: "Jawa Barat", Type: "Kabupaten", CityName: "Subang", PostalCode: "41215"},
		{CityID: "430", ProvinceID: "9", Province: "Jawa Barat", Type: "Kabupaten", CityName: "Sukabumi", PostalCode: "43311"},
		{CityID: "431", ProvinceID: "9", Province: "Jawa Barat", Type: "Kota", CityName: "Sukabumi", PostalCode: "43114"},
		{CityID: "440", ProvinceID: "9", Province: "Jawa Barat", Type: "Kabupaten", CityName: "Sumedang", PostalCode: "45326"},
		{CityID: "468", ProvinceID: "9", Province: "Jawa Barat", Type: "Kota", CityName: "Tasikmalaya", PostalCode: "46116"},
		{CityID: "469", ProvinceID: "9", Province: "Jawa Barat", Type: "Kabupaten", CityName: "Tasikmalaya", PostalCode: "46411"},
	},
	"10": { // Jawa Tengah
		{CityID: "37", ProvinceID: "10", Province: "Jawa Tengah", Type: "Kabupaten", CityName: "Banjarnegara", PostalCode: "53419"},
		{CityID: "38", ProvinceID: "10", Province: "Jawa Tengah", Type: "Kabupaten", CityName: "Banyumas", PostalCode: "53114"},
		{CityID: "41", ProvinceID: "10", Province: "Jawa Tengah", Type: "Kabupaten", CityName: "Batang", PostalCode: "51211"},
		{CityID: "76", ProvinceID: "10", Province: "Jawa Tengah", Type: "Kabupaten", CityName: "Blora", PostalCode: "58219"},
		{CityID: "91", ProvinceID: "10", Province: "Jawa Tengah", Type: "Kabupaten", CityName: "Boyolali", PostalCode: "57312"},
		{CityID: "92", ProvinceID: "10", Province: "Jawa Tengah", Type: "Kabupaten", CityName: "Brebes", PostalCode: "52212"},
		{CityID: "105", ProvinceID: "10", Province: "Jawa Tengah", Type: "Kabupaten", CityName: "Cilacap", PostalCode: "53211"},
		{CityID: "113", ProvinceID: "10", Province: "Jawa Tengah", Type: "Kabupaten", CityName: "Demak", PostalCode: "59519"},
		{CityID: "134", ProvinceID: "10", Province: "Jawa Tengah", Type: "Kabupaten", CityName: "Grobogan", PostalCode: "58111"},
		{CityID: "163", ProvinceID: "10", Province: "Jawa Tengah", Type: "Kabupaten", CityName: "Jepara", PostalCode: "59419"},
		{CityID: "169", ProvinceID: "10", Province: "Jawa Tengah", Type: "Kabupaten", CityName: "Karanganyar", PostalCode: "57718"},
		{CityID: "177", ProvinceID: "10", Province: "Jawa Tengah", Type: "Kabupaten", CityName: "Kebumen", PostalCode: "54319"},
		{CityID: "181", ProvinceID: "10", Province: "Jawa Tengah", Type: "Kabupaten", CityName: "Kendal", PostalCode: "51314"},
		{CityID: "196", ProvinceID: "10", Province: "Jawa Tengah", Type: "Kabupaten", CityName: "Klaten", PostalCode: "57411"},
		{CityID: "209", ProvinceID: "10", Province: "Jawa Tengah", Type: "Kabupaten", CityName: "Kudus", PostalCode: "59311"},
		{CityID: "249", ProvinceID: "10", Province: "Jawa Tengah", Type: "Kabupaten", CityName: "Magelang", PostalCode: "56519"},
		{CityID: "250", ProvinceID: "10", Province: "Jawa Tengah", Type: "Kota", CityName: "Magelang", PostalCode: "56133"},
		{CityID: "344", ProvinceID: "10", Province: "Jawa Tengah", Type: "Kabupaten", CityName: "Pati", PostalCode: "59114"},
		{CityID: "348", ProvinceID: "10", Province: "Jawa Tengah", Type: "Kabupaten", CityName: "Pekalongan", PostalCode: "51161"},
		{CityID: "349", ProvinceID: "10", Province: "Jawa Tengah", Type: "Kota", CityName: "Pekalongan", PostalCode: "51122"},
		{CityID: "352", ProvinceID: "10", Province: "Jawa Tengah", Type: "Kabupaten", CityName: "Pemalang", PostalCode: "52319"},
		{CityID: "375", ProvinceID: "10", Province: "Jawa Tengah", Type: "Kabupaten", CityName: "Purbalingga", PostalCode: "53312"},
		{CityID: "377", ProvinceID: "10", Province: "Jawa Tengah", Type: "Kabupaten", CityName: "Purworejo", PostalCode: "54111"},
		{CityID: "380", ProvinceID: "10", Province: "Jawa Tengah", Type: "Kabupaten", CityName: "Rembang", PostalCode: "59219"},
		{CityID: "386", ProvinceID: "10", Province: "Jawa Tengah", Type: "Kota", CityName: "Salatiga", PostalCode: "50711"},
		{CityID: "398", ProvinceID: "10", Province: "Jawa Tengah", Type: "Kota", CityName: "Semarang", PostalCode: "50135"},
		{CityID: "399", ProvinceID: "10", Province: "Jawa Tengah", Type: "Kabupaten", CityName: "Semarang", PostalCode: "50511"},
		{CityID: "427", ProvinceID: "10", Province: "Jawa Tengah", Type: "Kota", CityName: "Surakarta (Solo)", PostalCode: "57113"},
		{CityID: "433", ProvinceID: "10", Province: "Jawa Tengah", Type: "Kabupaten", CityName: "Sukoharjo", PostalCode: "57511"},
		{CityID: "445", ProvinceID: "10", Province: "Jawa Tengah", Type: "Kabupaten", CityName: "Tegal", PostalCode: "52419"},
		{CityID: "446", ProvinceID: "10", Province: "Jawa Tengah", Type: "Kota", CityName: "Tegal", PostalCode: "52114"},
		{CityID: "458", ProvinceID: "10", Province: "Jawa Tengah", Type: "Kabupaten", CityName: "Temanggung", PostalCode: "56212"},
		{CityID: "497", ProvinceID: "10", Province: "Jawa Tengah", Type: "Kabupaten", CityName: "Wonogiri", PostalCode: "57619"},
		{CityID: "498", ProvinceID: "10", Province: "Jawa Tengah", Type: "Kabupaten", CityName: "Wonosobo", PostalCode: "56311"},
	},
	"11": { // Jawa Timur
		{CityID: "31", ProvinceID: "11", Province: "Jawa Timur", Type: "Kabupaten", CityName: "Bangkalan", PostalCode: "69118"},
		{CityID: "36", ProvinceID: "11", Province: "Jawa Timur", Type: "Kabupaten", CityName: "Banyuwangi", PostalCode: "68416"},
		{CityID: "42", ProvinceID: "11", Province: "Jawa Timur", Type: "Kota", CityName: "Batu", PostalCode: "65311"},
		{CityID: "74", ProvinceID: "11", Province: "Jawa Timur", Type: "Kabupaten", CityName: "Blitar", PostalCode: "66171"},
		{CityID: "75", ProvinceID: "11", Province: "Jawa Timur", Type: "Kota", CityName: "Blitar", PostalCode: "66124"},
		{CityID: "78", ProvinceID: "11", Province: "Jawa Timur", Type: "Kabupaten", CityName: "Bojonegoro", PostalCode: "62119"},
		{CityID: "84", ProvinceID: "11", Province: "Jawa Timur", Type: "Kabupaten", CityName: "Bondowoso", PostalCode: "68219"},
		{CityID: "133", ProvinceID: "11", Province: "Jawa Timur", Type: "Kabupaten", CityName: "Gresik", PostalCode: "61115"},
		{CityID: "160", ProvinceID: "11", Province: "Jawa Timur", Type: "Kabupaten", CityName: "Jember", PostalCode: "68113"},
		{CityID: "164", ProvinceID: "11", Province: "Jawa Timur", Type: "Kabupaten", CityName: "Jombang", PostalCode: "61415"},
		{CityID: "178", ProvinceID: "11", Province: "Jawa Timur", Type: "Kabupaten", CityName: "Kediri", PostalCode: "64184"},
		{CityID: "179", ProvinceID: "11", Province: "Jawa Timur", Type: "Kota", CityName: "Kediri", PostalCode: "64125"},
		{CityID: "222", ProvinceID: "11", Province: "Jawa Timur", Type: "Kabupaten", CityName: "Lamongan", PostalCode: "62218"},
		{CityID: "243", ProvinceID: "11", Province: "Jawa Timur", Type: "Kabupaten", CityName: "Lumajang", PostalCode: "67319"},
		{CityID: "247", ProvinceID: "11", Province: "Jawa Timur", Type: "Kabupaten", CityName: "Madiun", PostalCode: "63153"},
		{CityID: "248", ProvinceID: "11", Province: "Jawa Timur", Type: "Kota", CityName: "Madiun", PostalCode: "63122"},
		{CityID: "251", ProvinceID: "11", Province: "Jawa Timur", Type: "Kabupaten", CityName: "Magetan", PostalCode: "63319"},
		{CityID: "255", ProvinceID: "11", Province: "Jawa Timur", Type: "Kabupaten", CityName: "Malang", PostalCode: "65163"},
		{CityID: "256", ProvinceID: "11", Province: "Jawa Timur", Type: "Kota", CityName: "Malang", PostalCode: "65119"},
		{CityID: "289", ProvinceID: "11", Province: "Jawa Timur", Type: "Kabupaten", CityName: "Mojokerto", PostalCode: "61382"},
		{CityID: "290", ProvinceID: "11", Province: "Jawa Timur", Type: "Kota", CityName: "Mojokerto", PostalCode: "61316"},
		{CityID: "305", ProvinceID: "11", Province: "Jawa Timur", Type: "Kabupaten", CityName: "Nganjuk", PostalCode: "64414"},
		{CityID: "306", ProvinceID: "11", Province: "Jawa Timur", Type: "Kabupaten", CityName: "Ngawi", PostalCode: "63219"},
		{CityID: "317", ProvinceID: "11", Province: "Jawa Timur", Type: "Kabupaten", CityName: "Pacitan", PostalCode: "63512"},
		{CityID: "330", ProvinceID: "11", Province: "Jawa Timur", Type: "Kabupaten", CityName: "Pamekasan", PostalCode: "69319"},
		{CityID: "342", ProvinceID: "11", Province: "Jawa Timur", Type: "Kabupaten", CityName: "Pasuruan", PostalCode: "67153"},
		{CityID: "343", ProvinceID: "11", Province: "Jawa Timur", Type: "Kota", CityName: "Pasuruan", PostalCode: "67118"},
		{CityID: "363", ProvinceID: "11", Province: "Jawa Timur", Type: "Kabupaten", CityName: "Ponorogo", PostalCode: "63411"},
		{CityID: "369", ProvinceID: "11", Province: "Jawa Timur", Type: "Kabupaten", CityName: "Probolinggo", PostalCode: "67282"},
		{CityID: "370", ProvinceID: "11", Province: "Jawa Timur", Type: "Kota", CityName: "Probolinggo", PostalCode: "67215"},
		{CityID: "390", ProvinceID: "11", Province: "Jawa Timur", Type: "Kabupaten", CityName: "Sampang", PostalCode: "69219"},
		{CityID: "409", ProvinceID: "11", Province: "Jawa Timur", Type: "Kabupaten", CityName: "Sidoarjo", PostalCode: "61219"},
		{CityID: "418", ProvinceID: "11", Province: "Jawa Timur", Type: "Kabupaten", CityName: "Situbondo", PostalCode: "68316"},
		{CityID: "441", ProvinceID: "11", Province: "Jawa Timur", Type: "Kabupaten", CityName: "Sumenep", PostalCode: "69413"},
		{CityID: "444", ProvinceID: "11", Province: "Jawa Timur", Type: "Kota", CityName: "Surabaya", PostalCode: "60119"},
		{CityID: "487", ProvinceID: "11", Province: "Jawa Timur", Type: "Kabupaten", CityName: "Trenggalek", PostalCode: "66312"},
		{CityID: "489", ProvinceID: "11", Province: "Jawa Timur", Type: "Kabupaten", CityName: "Tuban", PostalCode: "62319"},
		{CityID: "492", ProvinceID: "11", Province: "Jawa Timur", Type: "Kabupaten", CityName: "Tulungagung", PostalCode: "66212"},
	},
	"12": { // Kalimantan Barat
		{CityID: "61", ProvinceID: "12", Province: "Kalimantan Barat", Type: "Kabupaten", CityName: "Bengkayang", PostalCode: "79213"},
		{CityID: "170", ProvinceID: "12", Province: "Kalimantan Barat", Type: "Kabupaten", CityName: "Kapuas Hulu", PostalCode: "78719"},
		{CityID: "195", ProvinceID: "12", Province: "Kalimantan Barat", Type: "Kabupaten", CityName: "Ketapang", PostalCode: "78819"},
		{CityID: "365", ProvinceID: "12", Province: "Kalimantan Barat", Type: "Kota", CityName: "Pontianak", PostalCode: "78112"},
		{CityID: "388", ProvinceID: "12", Province: "Kalimantan Barat", Type: "Kabupaten", CityName: "Sambas", PostalCode: "79411"},
		{CityID: "414", ProvinceID: "12", Province: "Kalimantan Barat", Type: "Kota", CityName: "Singkawang", PostalCode: "79117"},
	},
	"13": { // Kalimantan Selatan
		{CityID: "19", ProvinceID: "13", Province: "Kalimantan Selatan", Type: "Kabupaten", CityName: "Balangan", PostalCode: "71611"},
		{CityID: "33", ProvinceID: "13", Province: "Kalimantan Selatan", Type: "Kota", CityName: "Banjar Baru", PostalCode: "70714"},
		{CityID: "35", ProvinceID: "13", Province: "Kalimantan Selatan", Type: "Kabupaten", CityName: "Banjar", PostalCode: "70619"},
		{CityID: "43", ProvinceID: "13", Province: "Kalimantan Selatan", Type: "Kabupaten", CityName: "Barito Kuala", PostalCode: "70511"},
		{CityID: "143", ProvinceID: "13", Province: "Kalimantan Selatan", Type: "Kabupaten", CityName: "Hulu Sungai Selatan", PostalCode: "71212"},
		{CityID: "144", ProvinceID: "13", Province: "Kalimantan Selatan", Type: "Kabupaten", CityName: "Hulu Sungai Tengah", PostalCode: "71313"},
		{CityID: "145", ProvinceID: "13", Province: "Kalimantan Selatan", Type: "Kabupaten", CityName: "Hulu Sungai Utara", PostalCode: "71419"},
		{CityID: "203", ProvinceID: "13", Province: "Kalimantan Selatan", Type: "Kabupaten", CityName: "Kotabaru", PostalCode: "72119"},
		{CityID: "446", ProvinceID: "13", Province: "Kalimantan Selatan", Type: "Kabupaten", CityName: "Tabalong", PostalCode: "71513"},
		{CityID: "452", ProvinceID: "13", Province: "Kalimantan Selatan", Type: "Kabupaten", CityName: "Tanah Bumbu", PostalCode: "72211"},
		{CityID: "454", ProvinceID: "13", Province: "Kalimantan Selatan", Type: "Kabupaten", CityName: "Tanah Laut", PostalCode: "70811"},
		{CityID: "466", ProvinceID: "13", Province: "Kalimantan Selatan", Type: "Kabupaten", CityName: "Tapin", PostalCode: "71119"},
	},
	"14": { // Kalimantan Tengah
		{CityID: "44", ProvinceID: "14", Province: "Kalimantan Tengah", Type: "Kabupaten", CityName: "Barito Selatan", PostalCode: "73711"},
		{CityID: "45", ProvinceID: "14", Province: "Kalimantan Tengah", Type: "Kabupaten", CityName: "Barito Timur", PostalCode: "73671"},
		{CityID: "46", ProvinceID: "14", Province: "Kalimantan Tengah", Type: "Kabupaten", CityName: "Barito Utara", PostalCode: "73881"},
		{CityID: "136", ProvinceID: "14", Province: "Kalimantan Tengah", Type: "Kabupaten", CityName: "Gunung Mas", PostalCode: "74511"},
		{CityID: "167", ProvinceID: "14", Province: "Kalimantan Tengah", Type: "Kabupaten", CityName: "Kapuas", PostalCode: "73511"},
		{CityID: "174", ProvinceID: "14", Province: "Kalimantan Tengah", Type: "Kabupaten", CityName: "Katingan", PostalCode: "74411"},
		{CityID: "205", ProvinceID: "14", Province: "Kalimantan Tengah", Type: "Kabupaten", CityName: "Kotawaringin Barat", PostalCode: "74119"},
		{CityID: "206", ProvinceID: "14", Province: "Kalimantan Tengah", Type: "Kabupaten", CityName: "Kotawaringin Timur", PostalCode: "74364"},
		{CityID: "326", ProvinceID: "14", Province: "Kalimantan Tengah", Type: "Kota", CityName: "Palangka Raya", PostalCode: "73112"},
	},
	"15": { // Kalimantan Timur
		{CityID: "19", ProvinceID: "15", Province: "Kalimantan Timur", Type: "Kota", CityName: "Balikpapan", PostalCode: "76111"},
		{CityID: "66", ProvinceID: "15", Province: "Kalimantan Timur", Type: "Kabupaten", CityName: "Berau", PostalCode: "77311"},
		{CityID: "89", ProvinceID: "15", Province: "Kalimantan Timur", Type: "Kota", CityName: "Bontang", PostalCode: "75311"},
		{CityID: "214", ProvinceID: "15", Province: "Kalimantan Timur", Type: "Kabupaten", CityName: "Kutai Barat", PostalCode: "75711"},
		{CityID: "215", ProvinceID: "15", Province: "Kalimantan Timur", Type: "Kabupaten", CityName: "Kutai Kartanegara", PostalCode: "75511"},
		{CityID: "216", ProvinceID: "15", Province: "Kalimantan Timur", Type: "Kabupaten", CityName: "Kutai Timur", PostalCode: "75611"},
		{CityID: "341", ProvinceID: "15", Province: "Kalimantan Timur", Type: "Kabupaten", CityName: "Paser", PostalCode: "76211"},
		{CityID: "354", ProvinceID: "15", Province: "Kalimantan Timur", Type: "Kabupaten", CityName: "Penajam Paser Utara", PostalCode: "76311"},
		{CityID: "387", ProvinceID: "15", Province: "Kalimantan Timur", Type: "Kota", CityName: "Samarinda", PostalCode: "75133"},
	},
	"16": { // Kalimantan Utara
		{CityID: "96", ProvinceID: "16", Province: "Kalimantan Utara", Type: "Kabupaten", CityName: "Bulungan (Bulongan)", PostalCode: "77211"},
		{CityID: "257", ProvinceID: "16", Province: "Kalimantan Utara", Type: "Kabupaten", CityName: "Malinau", PostalCode: "77511"},
		{CityID: "311", ProvinceID: "16", Province: "Kalimantan Utara", Type: "Kabupaten", CityName: "Nunukan", PostalCode: "77421"},
		{CityID: "450", ProvinceID: "16", Province: "Kalimantan Utara", Type: "Kabupaten", CityName: "Tana Tidung", PostalCode: "77611"},
		{CityID: "467", ProvinceID: "16", Province: "Kalimantan Utara", Type: "Kota", CityName: "Tarakan", PostalCode: "77114"},
	},
	"17": { // Kepulauan Riau
		{CityID: "48", ProvinceID: "17", Province: "Kepulauan Riau", Type: "Kota", CityName: "Batam", PostalCode: "29464"},
		{CityID: "71", ProvinceID: "17", Province: "Kepulauan Riau", Type: "Kabupaten", CityName: "Bintan", PostalCode: "29135"},
		{CityID: "172", ProvinceID: "17", Province: "Kepulauan Riau", Type: "Kabupaten", CityName: "Karimun", PostalCode: "29611"},
		{CityID: "184", ProvinceID: "17", Province: "Kepulauan Riau", Type: "Kabupaten", CityName: "Kepulauan Anambas", PostalCode: "29991"},
		{CityID: "237", ProvinceID: "17", Province: "Kepulauan Riau", Type: "Kabupaten", CityName: "Lingga", PostalCode: "29811"},
		{CityID: "302", ProvinceID: "17", Province: "Kepulauan Riau", Type: "Kabupaten", CityName: "Natuna", PostalCode: "29711"},
		{CityID: "462", ProvinceID: "17", Province: "Kepulauan Riau", Type: "Kota", CityName: "Tanjung Pinang", PostalCode: "29111"},
	},
	"18": { // Lampung
		{CityID: "21", ProvinceID: "18", Province: "Lampung", Type: "Kota", CityName: "Bandar Lampung", PostalCode: "35139"},
		{CityID: "223", ProvinceID: "18", Province: "Lampung", Type: "Kabupaten", CityName: "Lampung Barat", PostalCode: "34814"},
		{CityID: "224", ProvinceID: "18", Province: "Lampung", Type: "Kabupaten", CityName: "Lampung Selatan", PostalCode: "35511"},
		{CityID: "225", ProvinceID: "18", Province: "Lampung", Type: "Kabupaten", CityName: "Lampung Tengah", PostalCode: "34212"},
		{CityID: "226", ProvinceID: "18", Province: "Lampung", Type: "Kabupaten", CityName: "Lampung Timur", PostalCode: "34319"},
		{CityID: "227", ProvinceID: "18", Province: "Lampung", Type: "Kabupaten", CityName: "Lampung Utara", PostalCode: "34516"},
		{CityID: "282", ProvinceID: "18", Province: "Lampung", Type: "Kota", CityName: "Metro", PostalCode: "34111"},
	},
	"19": { // Maluku
		{CityID: "14", ProvinceID: "19", Province: "Maluku", Type: "Kota", CityName: "Ambon", PostalCode: "97222"},
		{CityID: "99", ProvinceID: "19", Province: "Maluku", Type: "Kabupaten", CityName: "Buru", PostalCode: "97371"},
		{CityID: "258", ProvinceID: "19", Province: "Maluku", Type: "Kabupaten", CityName: "Maluku Barat Daya", PostalCode: "97451"},
		{CityID: "259", ProvinceID: "19", Province: "Maluku", Type: "Kabupaten", CityName: "Maluku Tengah", PostalCode: "97513"},
		{CityID: "260", ProvinceID: "19", Province: "Maluku", Type: "Kabupaten", CityName: "Maluku Tenggara", PostalCode: "97651"},
		{CityID: "261", ProvinceID: "19", Province: "Maluku", Type: "Kabupaten", CityName: "Maluku Tenggara Barat", PostalCode: "97465"},
	},
	"20": { // Maluku Utara
		{CityID: "138", ProvinceID: "20", Province: "Maluku Utara", Type: "Kabupaten", CityName: "Halmahera Barat", PostalCode: "97757"},
		{CityID: "139", ProvinceID: "20", Province: "Maluku Utara", Type: "Kabupaten", CityName: "Halmahera Selatan", PostalCode: "97911"},
		{CityID: "140", ProvinceID: "20", Province: "Maluku Utara", Type: "Kabupaten", CityName: "Halmahera Tengah", PostalCode: "97853"},
		{CityID: "141", ProvinceID: "20", Province: "Maluku Utara", Type: "Kabupaten", CityName: "Halmahera Timur", PostalCode: "97862"},
		{CityID: "142", ProvinceID: "20", Province: "Maluku Utara", Type: "Kabupaten", CityName: "Halmahera Utara", PostalCode: "97762"},
		{CityID: "476", ProvinceID: "20", Province: "Maluku Utara", Type: "Kota", CityName: "Ternate", PostalCode: "97714"},
		{CityID: "478", ProvinceID: "20", Province: "Maluku Utara", Type: "Kota", CityName: "Tidore Kepulauan", PostalCode: "97815"},
	},
	"21": { // Nanggroe Aceh Darussalam (NAD)
		{CityID: "1", ProvinceID: "21", Province: "Nanggroe Aceh Darussalam (NAD)", Type: "Kabupaten", CityName: "Aceh Barat", PostalCode: "23611"},
		{CityID: "2", ProvinceID: "21", Province: "Nanggroe Aceh Darussalam (NAD)", Type: "Kabupaten", CityName: "Aceh Barat Daya", PostalCode: "23764"},
		{CityID: "3", ProvinceID: "21", Province: "Nanggroe Aceh Darussalam (NAD)", Type: "Kabupaten", CityName: "Aceh Besar", PostalCode: "23951"},
		{CityID: "4", ProvinceID: "21", Province: "Nanggroe Aceh Darussalam (NAD)", Type: "Kabupaten", CityName: "Aceh Jaya", PostalCode: "23654"},
		{CityID: "5", ProvinceID: "21", Province: "Nanggroe Aceh Darussalam (NAD)", Type: "Kabupaten", CityName: "Aceh Selatan", PostalCode: "23719"},
		{CityID: "20", ProvinceID: "21", Province: "Nanggroe Aceh Darussalam (NAD)", Type: "Kota", CityName: "Banda Aceh", PostalCode: "23242"},
		{CityID: "59", ProvinceID: "21", Province: "Nanggroe Aceh Darussalam (NAD)", Type: "Kabupaten", CityName: "Bener Meriah", PostalCode: "24581"},
		{CityID: "242", ProvinceID: "21", Province: "Nanggroe Aceh Darussalam (NAD)", Type: "Kota", CityName: "Lhokseumawe", PostalCode: "24352"},
	},
	"22": { // Nusa Tenggara Barat (NTB)
		{CityID: "68", ProvinceID: "22", Province: "Nusa Tenggara Barat (NTB)", Type: "Kabupaten", CityName: "Bima", PostalCode: "84171"},
		{CityID: "69", ProvinceID: "22", Province: "Nusa Tenggara Barat (NTB)", Type: "Kota", CityName: "Bima", PostalCode: "84139"},
		{CityID: "118", ProvinceID: "22", Province: "Nusa Tenggara Barat (NTB)", Type: "Kabupaten", CityName: "Dompu", PostalCode: "84217"},
		{CityID: "238", ProvinceID: "22", Province: "Nusa Tenggara Barat (NTB)", Type: "Kabupaten", CityName: "Lombok Barat", PostalCode: "83311"},
		{CityID: "239", ProvinceID: "22", Province: "Nusa Tenggara Barat (NTB)", Type: "Kabupaten", CityName: "Lombok Tengah", PostalCode: "83511"},
		{CityID: "240", ProvinceID: "22", Province: "Nusa Tenggara Barat (NTB)", Type: "Kabupaten", CityName: "Lombok Timur", PostalCode: "83612"},
		{CityID: "241", ProvinceID: "22", Province: "Nusa Tenggara Barat (NTB)", Type: "Kabupaten", CityName: "Lombok Utara", PostalCode: "83711"},
		{CityID: "276", ProvinceID: "22", Province: "Nusa Tenggara Barat (NTB)", Type: "Kota", CityName: "Mataram", PostalCode: "83131"},
		{CityID: "438", ProvinceID: "22", Province: "Nusa Tenggara Barat (NTB)", Type: "Kabupaten", CityName: "Sumbawa", PostalCode: "84315"},
		{CityID: "439", ProvinceID: "22", Province: "Nusa Tenggara Barat (NTB)", Type: "Kabupaten", CityName: "Sumbawa Barat", PostalCode: "84419"},
	},
	"23": { // Nusa Tenggara Timur (NTT)
		{CityID: "13", ProvinceID: "23", Province: "Nusa Tenggara Timur (NTT)", Type: "Kabupaten", CityName: "Alor", PostalCode: "85811"},
		{CityID: "58", ProvinceID: "23", Province: "Nusa Tenggara Timur (NTT)", Type: "Kabupaten", CityName: "Belu", PostalCode: "85711"},
		{CityID: "122", ProvinceID: "23", Province: "Nusa Tenggara Timur (NTT)", Type: "Kabupaten", CityName: "Ende", PostalCode: "86351"},
		{CityID: "125", ProvinceID: "23", Province: "Nusa Tenggara Timur (NTT)", Type: "Kabupaten", CityName: "Flores Timur", PostalCode: "86213"},
		{CityID: "212", ProvinceID: "23", Province: "Nusa Tenggara Timur (NTT)", Type: "Kota", CityName: "Kupang", PostalCode: "85119"},
		{CityID: "213", ProvinceID: "23", Province: "Nusa Tenggara Timur (NTT)", Type: "Kabupaten", CityName: "Kupang", PostalCode: "85362"},
	},
	"24": { // Papua
		{CityID: "67", ProvinceID: "24", Province: "Papua", Type: "Kabupaten", CityName: "Biak Numfor", PostalCode: "98119"},
		{CityID: "90", ProvinceID: "24", Province: "Papua", Type: "Kabupaten", CityName: "Boven Digoel", PostalCode: "99662"},
		{CityID: "111", ProvinceID: "24", Province: "Papua", Type: "Kabupaten", CityName: "Deiyai (Deliyai)", PostalCode: "98784"},
		{CityID: "117", ProvinceID: "24", Province: "Papua", Type: "Kabupaten", CityName: "Dogiyai", PostalCode: "98866"},
		{CityID: "157", ProvinceID: "24", Province: "Papua", Type: "Kota", CityName: "Jayapura", PostalCode: "99114"},
		{CityID: "158", ProvinceID: "24", Province: "Papua", Type: "Kabupaten", CityName: "Jayapura", PostalCode: "99352"},
		{CityID: "159", ProvinceID: "24", Province: "Papua", Type: "Kabupaten", CityName: "Jayawijaya", PostalCode: "99511"},
	},
	"25": { // Papua Barat
		{CityID: "124", ProvinceID: "25", Province: "Papua Barat", Type: "Kabupaten", CityName: "Fakfak", PostalCode: "98651"},
		{CityID: "165", ProvinceID: "25", Province: "Papua Barat", Type: "Kabupaten", CityName: "Kaimana", PostalCode: "98671"},
		{CityID: "263", ProvinceID: "25", Province: "Papua Barat", Type: "Kabupaten", CityName: "Manokwari", PostalCode: "98311"},
		{CityID: "264", ProvinceID: "25", Province: "Papua Barat", Type: "Kabupaten", CityName: "Manokwari Selatan", PostalCode: "98355"},
		{CityID: "382", ProvinceID: "25", Province: "Papua Barat", Type: "Kabupaten", CityName: "Raja Ampat", PostalCode: "98489"},
		{CityID: "424", ProvinceID: "25", Province: "Papua Barat", Type: "Kota", CityName: "Sorong", PostalCode: "98411"},
		{CityID: "425", ProvinceID: "25", Province: "Papua Barat", Type: "Kabupaten", CityName: "Sorong", PostalCode: "98431"},
		{CityID: "426", ProvinceID: "25", Province: "Papua Barat", Type: "Kabupaten", CityName: "Sorong Selatan", PostalCode: "98454"},
	},
	"26": { // Riau
		{CityID: "60", ProvinceID: "26", Province: "Riau", Type: "Kabupaten", CityName: "Bengkalis", PostalCode: "28719"},
		{CityID: "120", ProvinceID: "26", Province: "Riau", Type: "Kota", CityName: "Dumai", PostalCode: "28811"},
		{CityID: "147", ProvinceID: "26", Province: "Riau", Type: "Kabupaten", CityName: "Indragiri Hilir", PostalCode: "29212"},
		{CityID: "148", ProvinceID: "26", Province: "Riau", Type: "Kabupaten", CityName: "Indragiri Hulu", PostalCode: "29319"},
		{CityID: "166", ProvinceID: "26", Province: "Riau", Type: "Kabupaten", CityName: "Kampar", PostalCode: "28411"},
		{CityID: "187", ProvinceID: "26", Province: "Riau", Type: "Kabupaten", CityName: "Kepulauan Meranti", PostalCode: "28791"},
		{CityID: "207", ProvinceID: "26", Province: "Riau", Type: "Kabupaten", CityName: "Kuantan Singingi", PostalCode: "29519"},
		{CityID: "350", ProvinceID: "26", Province: "Riau", Type: "Kota", CityName: "Pekanbaru", PostalCode: "28112"},
		{CityID: "351", ProvinceID: "26", Province: "Riau", Type: "Kabupaten", CityName: "Pelalawan", PostalCode: "28311"},
		{CityID: "381", ProvinceID: "26", Province: "Riau", Type: "Kabupaten", CityName: "Rokan Hilir", PostalCode: "28992"},
		{CityID: "385", ProvinceID: "26", Province: "Riau", Type: "Kabupaten", CityName: "Rokan Hulu", PostalCode: "28511"},
		{CityID: "406", ProvinceID: "26", Province: "Riau", Type: "Kabupaten", CityName: "Siak", PostalCode: "28623"},
	},
	"27": { // Sulawesi Barat
		{CityID: "253", ProvinceID: "27", Province: "Sulawesi Barat", Type: "Kabupaten", CityName: "Majene", PostalCode: "91411"},
		{CityID: "262", ProvinceID: "27", Province: "Sulawesi Barat", Type: "Kabupaten", CityName: "Mamasa", PostalCode: "91362"},
		{CityID: "265", ProvinceID: "27", Province: "Sulawesi Barat", Type: "Kabupaten", CityName: "Mamuju", PostalCode: "91519"},
		{CityID: "266", ProvinceID: "27", Province: "Sulawesi Barat", Type: "Kabupaten", CityName: "Mamuju Utara", PostalCode: "91571"},
		{CityID: "362", ProvinceID: "27", Province: "Sulawesi Barat", Type: "Kabupaten", CityName: "Polewali Mandar", PostalCode: "91311"},
	},
	"28": { // Sulawesi Selatan
		{CityID: "38", ProvinceID: "28", Province: "Sulawesi Selatan", Type: "Kabupaten", CityName: "Bantaeng", PostalCode: "92411"},
		{CityID: "47", ProvinceID: "28", Province: "Sulawesi Selatan", Type: "Kabupaten", CityName: "Barru", PostalCode: "90719"},
		{CityID: "72", ProvinceID: "28", Province: "Sulawesi Selatan", Type: "Kabupaten", CityName: "Bone", PostalCode: "92713"},
		{CityID: "95", ProvinceID: "28", Province: "Sulawesi Selatan", Type: "Kabupaten", CityName: "Bulukumba", PostalCode: "92511"},
		{CityID: "123", ProvinceID: "28", Province: "Sulawesi Selatan", Type: "Kabupaten", CityName: "Enrekang", PostalCode: "91719"},
		{CityID: "132", ProvinceID: "28", Province: "Sulawesi Selatan", Type: "Kabupaten", CityName: "Gowa", PostalCode: "92111"},
		{CityID: "162", ProvinceID: "28", Province: "Sulawesi Selatan", Type: "Kabupaten", CityName: "Jeneponto", PostalCode: "92319"},
		{CityID: "244", ProvinceID: "28", Province: "Sulawesi Selatan", Type: "Kabupaten", CityName: "Luwu", PostalCode: "91994"},
		{CityID: "245", ProvinceID: "28", Province: "Sulawesi Selatan", Type: "Kabupaten", CityName: "Luwu Timur", PostalCode: "92981"},
		{CityID: "246", ProvinceID: "28", Province: "Sulawesi Selatan", Type: "Kabupaten", CityName: "Luwu Utara", PostalCode: "92911"},
		{CityID: "254", ProvinceID: "28", Province: "Sulawesi Selatan", Type: "Kota", CityName: "Makassar", PostalCode: "90111"},
		{CityID: "270", ProvinceID: "28", Province: "Sulawesi Selatan", Type: "Kabupaten", CityName: "Maros", PostalCode: "90511"},
		{CityID: "325", ProvinceID: "28", Province: "Sulawesi Selatan", Type: "Kota", CityName: "Palopo", PostalCode: "91911"},
		{CityID: "335", ProvinceID: "28", Province: "Sulawesi Selatan", Type: "Kabupaten", CityName: "Pangkajene Kepulauan", PostalCode: "90611"},
		{CityID: "337", ProvinceID: "28", Province: "Sulawesi Selatan", Type: "Kota", CityName: "Parepare", PostalCode: "91123"},
		{CityID: "359", ProvinceID: "28", Province: "Sulawesi Selatan", Type: "Kabupaten", CityName: "Pinrang", PostalCode: "91251"},
	},
	"29": { // Sulawesi Tengah
		{CityID: "30", ProvinceID: "29", Province: "Sulawesi Tengah", Type: "Kabupaten", CityName: "Banggai", PostalCode: "94711"},
		{CityID: "50", ProvinceID: "29", Province: "Sulawesi Tengah", Type: "Kabupaten", CityName: "Buol", PostalCode: "94564"},
		{CityID: "119", ProvinceID: "29", Province: "Sulawesi Tengah", Type: "Kabupaten", CityName: "Donggala", PostalCode: "94341"},
		{CityID: "291", ProvinceID: "29", Province: "Sulawesi Tengah", Type: "Kabupaten", CityName: "Morowali", PostalCode: "94911"},
		{CityID: "324", ProvinceID: "29", Province: "Sulawesi Tengah", Type: "Kota", CityName: "Palu", PostalCode: "94111"},
		{CityID: "340", ProvinceID: "29", Province: "Sulawesi Tengah", Type: "Kabupaten", CityName: "Parigi Moutong", PostalCode: "94411"},
		{CityID: "366", ProvinceID: "29", Province: "Sulawesi Tengah", Type: "Kabupaten", CityName: "Poso", PostalCode: "94615"},
		{CityID: "410", ProvinceID: "29", Province: "Sulawesi Tengah", Type: "Kabupaten", CityName: "Sigi", PostalCode: "94364"},
		{CityID: "482", ProvinceID: "29", Province: "Sulawesi Tengah", Type: "Kabupaten", CityName: "Tojo Una-Una", PostalCode: "94683"},
		{CityID: "483", ProvinceID: "29", Province: "Sulawesi Tengah", Type: "Kabupaten", CityName: "Toli-Toli", PostalCode: "94542"},
	},
	"30": { // Sulawesi Tenggara
		{CityID: "53", ProvinceID: "30", Province: "Sulawesi Tenggara", Type: "Kota", CityName: "Bau-Bau", PostalCode: "93724"},
		{CityID: "85", ProvinceID: "30", Province: "Sulawesi Tenggara", Type: "Kabupaten", CityName: "Bombana", PostalCode: "93771"},
		{CityID: "101", ProvinceID: "30", Province: "Sulawesi Tenggara", Type: "Kabupaten", CityName: "Buton", PostalCode: "93754"},
		{CityID: "102", ProvinceID: "30", Province: "Sulawesi Tenggara", Type: "Kabupaten", CityName: "Buton Utara", PostalCode: "93745"},
		{CityID: "182", ProvinceID: "30", Province: "Sulawesi Tenggara", Type: "Kota", CityName: "Kendari", PostalCode: "93126"},
		{CityID: "198", ProvinceID: "30", Province: "Sulawesi Tenggara", Type: "Kabupaten", CityName: "Kolaka", PostalCode: "93511"},
		{CityID: "199", ProvinceID: "30", Province: "Sulawesi Tenggara", Type: "Kabupaten", CityName: "Kolaka Utara", PostalCode: "93911"},
		{CityID: "200", ProvinceID: "30", Province: "Sulawesi Tenggara", Type: "Kabupaten", CityName: "Konawe", PostalCode: "93411"},
		{CityID: "201", ProvinceID: "30", Province: "Sulawesi Tenggara", Type: "Kabupaten", CityName: "Konawe Selatan", PostalCode: "93811"},
		{CityID: "202", ProvinceID: "30", Province: "Sulawesi Tenggara", Type: "Kabupaten", CityName: "Konawe Utara", PostalCode: "93311"},
		{CityID: "295", ProvinceID: "30", Province: "Sulawesi Tenggara", Type: "Kabupaten", CityName: "Muna", PostalCode: "93611"},
		{CityID: "494", ProvinceID: "30", Province: "Sulawesi Tenggara", Type: "Kabupaten", CityName: "Wakatobi", PostalCode: "93791"},
	},
	"31": { // Sulawesi Utara
		{CityID: "73", ProvinceID: "31", Province: "Sulawesi Utara", Type: "Kota", CityName: "Bitung", PostalCode: "95512"},
		{CityID: "81", ProvinceID: "31", Province: "Sulawesi Utara", Type: "Kabupaten", CityName: "Bolaang Mongondow (Bolmong)", PostalCode: "95755"},
		{CityID: "82", ProvinceID: "31", Province: "Sulawesi Utara", Type: "Kabupaten", CityName: "Bolaang Mongondow Selatan", PostalCode: "95774"},
		{CityID: "83", ProvinceID: "31", Province: "Sulawesi Utara", Type: "Kabupaten", CityName: "Bolaang Mongondow Timur", PostalCode: "95783"},
		{CityID: "188", ProvinceID: "31", Province: "Sulawesi Utara", Type: "Kabupaten", CityName: "Kepulauan Sangihe", PostalCode: "95819"},
		{CityID: "190", ProvinceID: "31", Province: "Sulawesi Utara", Type: "Kabupaten", CityName: "Kepulauan Siau Tagulandang Biaro (Sitaro)", PostalCode: "95862"},
		{CityID: "191", ProvinceID: "31", Province: "Sulawesi Utara", Type: "Kabupaten", CityName: "Kepulauan Talaud", PostalCode: "95885"},
		{CityID: "204", ProvinceID: "31", Province: "Sulawesi Utara", Type: "Kota", CityName: "Kotamobagu", PostalCode: "95711"},
		{CityID: "267", ProvinceID: "31", Province: "Sulawesi Utara", Type: "Kota", CityName: "Manado", PostalCode: "95247"},
		{CityID: "285", ProvinceID: "31", Province: "Sulawesi Utara", Type: "Kabupaten", CityName: "Minahasa", PostalCode: "95614"},
		{CityID: "286", ProvinceID: "31", Province: "Sulawesi Utara", Type: "Kabupaten", CityName: "Minahasa Selatan", PostalCode: "95914"},
		{CityID: "287", ProvinceID: "31", Province: "Sulawesi Utara", Type: "Kabupaten", CityName: "Minahasa Tenggara", PostalCode: "95995"},
		{CityID: "288", ProvinceID: "31", Province: "Sulawesi Utara", Type: "Kabupaten", CityName: "Minahasa Utara", PostalCode: "95316"},
		{CityID: "485", ProvinceID: "31", Province: "Sulawesi Utara", Type: "Kota", CityName: "Tomohon", PostalCode: "95416"},
	},
	"32": { // Sumatera Barat
		{CityID: "8", ProvinceID: "32", Province: "Sumatera Barat", Type: "Kabupaten", CityName: "Agam", PostalCode: "26411"},
		{CityID: "93", ProvinceID: "32", Province: "Sumatera Barat", Type: "Kota", CityName: "Bukittinggi", PostalCode: "26115"},
		{CityID: "116", ProvinceID: "32", Province: "Sumatera Barat", Type: "Kabupaten", CityName: "Dharmasraya", PostalCode: "27612"},
		{CityID: "186", ProvinceID: "32", Province: "Sumatera Barat", Type: "Kabupaten", CityName: "Kepulauan Mentawai", PostalCode: "25771"},
		{CityID: "236", ProvinceID: "32", Province: "Sumatera Barat", Type: "Kabupaten", CityName: "Lima Puluh Kota", PostalCode: "26212"},
		{CityID: "318", ProvinceID: "32", Province: "Sumatera Barat", Type: "Kota", CityName: "Padang", PostalCode: "25112"},
		{CityID: "319", ProvinceID: "32", Province: "Sumatera Barat", Type: "Kabupaten", CityName: "Padang Pariaman", PostalCode: "25583"},
		{CityID: "320", ProvinceID: "32", Province: "Sumatera Barat", Type: "Kota", CityName: "Padang Panjang", PostalCode: "27122"},
		{CityID: "338", ProvinceID: "32", Province: "Sumatera Barat", Type: "Kota", CityName: "Pariaman", PostalCode: "25511"},
		{CityID: "339", ProvinceID: "32", Province: "Sumatera Barat", Type: "Kabupaten", CityName: "Pasaman", PostalCode: "26318"},
		{CityID: "345", ProvinceID: "32", Province: "Sumatera Barat", Type: "Kabupaten", CityName: "Pasaman Barat", PostalCode: "26511"},
		{CityID: "346", ProvinceID: "32", Province: "Sumatera Barat", Type: "Kota", CityName: "Payakumbuh", PostalCode: "26213"},
		{CityID: "357", ProvinceID: "32", Province: "Sumatera Barat", Type: "Kabupaten", CityName: "Pesisir Selatan", PostalCode: "25611"},
		{CityID: "394", ProvinceID: "32", Province: "Sumatera Barat", Type: "Kota", CityName: "Sawah Lunto", PostalCode: "27416"},
		{CityID: "411", ProvinceID: "32", Province: "Sumatera Barat", Type: "Kabupaten", CityName: "Sijunjung (Sawah Lunto Sijunjung)", PostalCode: "27511"},
		{CityID: "420", ProvinceID: "32", Province: "Sumatera Barat", Type: "Kota", CityName: "Solok", PostalCode: "27315"},
		{CityID: "421", ProvinceID: "32", Province: "Sumatera Barat", Type: "Kabupaten", CityName: "Solok", PostalCode: "27365"},
		{CityID: "422", ProvinceID: "32", Province: "Sumatera Barat", Type: "Kabupaten", CityName: "Solok Selatan", PostalCode: "27779"},
		{CityID: "453", ProvinceID: "32", Province: "Sumatera Barat", Type: "Kabupaten", CityName: "Tanah Datar", PostalCode: "27211"},
	},
	"33": { // Sumatera Selatan
		{CityID: "51", ProvinceID: "33", Province: "Sumatera Selatan", Type: "Kabupaten", CityName: "Banyuasin", PostalCode: "30911"},
		{CityID: "121", ProvinceID: "33", Province: "Sumatera Selatan", Type: "Kabupaten", CityName: "Empat Lawang", PostalCode: "31811"},
		{CityID: "220", ProvinceID: "33", Province: "Sumatera Selatan", Type: "Kabupaten", CityName: "Lahat", PostalCode: "31419"},
		{CityID: "242", ProvinceID: "33", Province: "Sumatera Selatan", Type: "Kota", CityName: "Lubuk Linggau", PostalCode: "31614"},
		{CityID: "292", ProvinceID: "33", Province: "Sumatera Selatan", Type: "Kabupaten", CityName: "Muara Enim", PostalCode: "31315"},
		{CityID: "293", ProvinceID: "33", Province: "Sumatera Selatan", Type: "Kabupaten", CityName: "Musi Banyuasin", PostalCode: "30719"},
		{CityID: "294", ProvinceID: "33", Province: "Sumatera Selatan", Type: "Kabupaten", CityName: "Musi Rawas", PostalCode: "31661"},
		{CityID: "313", ProvinceID: "33", Province: "Sumatera Selatan", Type: "Kabupaten", CityName: "Ogan Ilir", PostalCode: "30811"},
		{CityID: "314", ProvinceID: "33", Province: "Sumatera Selatan", Type: "Kabupaten", CityName: "Ogan Komering Ilir", PostalCode: "30618"},
		{CityID: "315", ProvinceID: "33", Province: "Sumatera Selatan", Type: "Kabupaten", CityName: "Ogan Komering Ulu", PostalCode: "32112"},
		{CityID: "316", ProvinceID: "33", Province: "Sumatera Selatan", Type: "Kabupaten", CityName: "Ogan Komering Ulu Selatan", PostalCode: "32211"},
		{CityID: "317", ProvinceID: "33", Province: "Sumatera Selatan", Type: "Kabupaten", CityName: "Ogan Komering Ulu Timur", PostalCode: "32312"},
		{CityID: "323", ProvinceID: "33", Province: "Sumatera Selatan", Type: "Kota", CityName: "Pagar Alam", PostalCode: "31512"},
		{CityID: "327", ProvinceID: "33", Province: "Sumatera Selatan", Type: "Kota", CityName: "Palembang", PostalCode: "30111"},
		{CityID: "367", ProvinceID: "33", Province: "Sumatera Selatan", Type: "Kota", CityName: "Prabumulih", PostalCode: "31119"},
	},
	"34": { // Sumatera Utara
		{CityID: "15", ProvinceID: "34", Province: "Sumatera Utara", Type: "Kabupaten", CityName: "Asahan", PostalCode: "21214"},
		{CityID: "52", ProvinceID: "34", Province: "Sumatera Utara", Type: "Kabupaten", CityName: "Batu Bara", PostalCode: "21655"},
		{CityID: "70", ProvinceID: "34", Province: "Sumatera Utara", Type: "Kota", CityName: "Binjai", PostalCode: "20712"},
		{CityID: "109", ProvinceID: "34", Province: "Sumatera Utara", Type: "Kabupaten", CityName: "Dairi", PostalCode: "22211"},
		{CityID: "112", ProvinceID: "34", Province: "Sumatera Utara", Type: "Kabupaten", CityName: "Deli Serdang", PostalCode: "20511"},
		{CityID: "137", ProvinceID: "34", Province: "Sumatera Utara", Type: "Kota", CityName: "Gunungsitoli", PostalCode: "22813"},
		{CityID: "146", ProvinceID: "34", Province: "Sumatera Utara", Type: "Kabupaten", CityName: "Humbang Hasundutan", PostalCode: "22457"},
		{CityID: "173", ProvinceID: "34", Province: "Sumatera Utara", Type: "Kabupaten", CityName: "Karo", PostalCode: "22119"},
		{CityID: "217", ProvinceID: "34", Province: "Sumatera Utara", Type: "Kabupaten", CityName: "Labuhan Batu", PostalCode: "21412"},
		{CityID: "218", ProvinceID: "34", Province: "Sumatera Utara", Type: "Kabupaten", CityName: "Labuhan Batu Selatan", PostalCode: "21511"},
		{CityID: "219", ProvinceID: "34", Province: "Sumatera Utara", Type: "Kabupaten", CityName: "Labuhan Batu Utara", PostalCode: "21711"},
		{CityID: "229", ProvinceID: "34", Province: "Sumatera Utara", Type: "Kabupaten", CityName: "Langkat", PostalCode: "20811"},
		{CityID: "268", ProvinceID: "34", Province: "Sumatera Utara", Type: "Kabupaten", CityName: "Mandailing Natal", PostalCode: "22916"},
		{CityID: "278", ProvinceID: "34", Province: "Sumatera Utara", Type: "Kota", CityName: "Medan", PostalCode: "20228"},
		{CityID: "307", ProvinceID: "34", Province: "Sumatera Utara", Type: "Kabupaten", CityName: "Nias", PostalCode: "22876"},
		{CityID: "308", ProvinceID: "34", Province: "Sumatera Utara", Type: "Kabupaten", CityName: "Nias Barat", PostalCode: "22895"},
		{CityID: "309", ProvinceID: "34", Province: "Sumatera Utara", Type: "Kabupaten", CityName: "Nias Selatan", PostalCode: "22865"},
		{CityID: "310", ProvinceID: "34", Province: "Sumatera Utara", Type: "Kabupaten", CityName: "Nias Utara", PostalCode: "22856"},
		{CityID: "321", ProvinceID: "34", Province: "Sumatera Utara", Type: "Kota", CityName: "Padang Sidempuan", PostalCode: "22727"},
		{CityID: "322", ProvinceID: "34", Province: "Sumatera Utara", Type: "Kabupaten", CityName: "Pakpak Bharat", PostalCode: "22272"},
		{CityID: "347", ProvinceID: "34", Province: "Sumatera Utara", Type: "Kota", CityName: "Pematang Siantar", PostalCode: "21126"},
		{CityID: "389", ProvinceID: "34", Province: "Sumatera Utara", Type: "Kabupaten", CityName: "Samosir", PostalCode: "22392"},
		{CityID: "404", ProvinceID: "34", Province: "Sumatera Utara", Type: "Kabupaten", CityName: "Serdang Bedagai", PostalCode: "20915"},
		{CityID: "407", ProvinceID: "34", Province: "Sumatera Utara", Type: "Kota", CityName: "Sibolga", PostalCode: "22522"},
		{CityID: "413", ProvinceID: "34", Province: "Sumatera Utara", Type: "Kabupaten", CityName: "Simalungun", PostalCode: "21162"},
		{CityID: "459", ProvinceID: "34", Province: "Sumatera Utara", Type: "Kota", CityName: "Tanjung Balai", PostalCode: "21321"},
		{CityID: "463", ProvinceID: "34", Province: "Sumatera Utara", Type: "Kabupaten", CityName: "Tapanuli Selatan", PostalCode: "22742"},
		{CityID: "464", ProvinceID: "34", Province: "Sumatera Utara", Type: "Kabupaten", CityName: "Tapanuli Tengah", PostalCode: "22611"},
		{CityID: "465", ProvinceID: "34", Province: "Sumatera Utara", Type: "Kabupaten", CityName: "Tapanuli Utara", PostalCode: "22414"},
		{CityID: "470", ProvinceID: "34", Province: "Sumatera Utara", Type: "Kota", CityName: "Tebing Tinggi", PostalCode: "20632"},
		{CityID: "481", ProvinceID: "34", Province: "Sumatera Utara", Type: "Kabupaten", CityName: "Toba Samosir", PostalCode: "22316"},
	},
}

// GetProvinces godoc
// @Summary Get Provinces
// @Description Get all provinces
// @Tags Shipping
// @Produce json
// @Success 200 {object} presenter.SuccessResponseSwagger
// @Router /shipping/provinces [get]
func (h *ShippingHandler) GetProvinces(c *fiber.Ctx) error {
	if h.useAPI && h.rajaOngkir != nil {
		provinces, err := h.rajaOngkir.GetProvinces()
		if err == nil && len(provinces) > 0 {
			// Convert to local type
			result := make([]Province, len(provinces))
			for i, p := range provinces {
				result[i] = Province{
					ProvinceID: p.ProvinceID,
					Province:   p.Province,
				}
			}
			return c.JSON(presenter.SuccessResponse(result))
		}
	}
	// Fallback to mock data
	return c.JSON(presenter.SuccessResponse(mockProvinces))
}

// GetCities godoc
// @Summary Get Cities
// @Description Get cities by province ID
// @Tags Shipping
// @Produce json
// @Param province_id query string false "Province ID"
// @Success 200 {object} presenter.SuccessResponseSwagger
// @Router /shipping/cities [get]
func (h *ShippingHandler) GetCities(c *fiber.Ctx) error {
	provinceID := c.Query("province_id", "")

	if h.useAPI && h.rajaOngkir != nil {
		cities, err := h.rajaOngkir.GetCities(provinceID)
		if err == nil && len(cities) > 0 {
			// Convert to local type
			result := make([]City, len(cities))
			for i, city := range cities {
				result[i] = City{
					CityID:     city.CityID,
					ProvinceID: city.ProvinceID,
					Province:   city.Province,
					Type:       city.Type,
					CityName:   city.CityName,
					PostalCode: city.PostalCode,
				}
			}
			return c.JSON(presenter.SuccessResponse(result))
		}
	}

	// Fallback to mock data
	if provinceID == "" {
		allCities := []City{}
		for _, cities := range mockCities {
			allCities = append(allCities, cities...)
		}
		return c.JSON(presenter.SuccessResponse(allCities))
	}

	cities, ok := mockCities[provinceID]
	if !ok {
		return c.JSON(presenter.SuccessResponse([]City{}))
	}
	return c.JSON(presenter.SuccessResponse(cities))
}

// GetCost godoc
// @Summary Get Shipping Cost
// @Description Calculate shipping cost
// @Tags Shipping
// @Produce json
// @Param origin query string true "Origin City ID"
// @Param destination query string true "Destination City ID"
// @Param weight query int true "Weight in grams"
// @Param courier query string false "Courier code (jne, pos, tiki)"
// @Success 200 {object} presenter.SuccessResponseSwagger
// @Router /shipping/cost [get]
func (h *ShippingHandler) GetCost(c *fiber.Ctx) error {
	origin := c.Query("origin", "")
	destination := c.Query("destination", "")
	weightStr := c.Query("weight", "1000")
	courier := c.Query("courier", "")
	
	weight, _ := strconv.Atoi(weightStr)
	if weight <= 0 {
		weight = 1000
	}

	// Try API first
	if h.useAPI && h.rajaOngkir != nil && origin != "" && destination != "" {
		var apiResults []rajaongkir.CostResult
		var err error
		
		if courier != "" {
			apiResults, err = h.rajaOngkir.GetCost(origin, destination, weight, courier)
		} else {
			apiResults, err = h.rajaOngkir.GetAllCouriers(origin, destination, weight)
		}
		
		if err == nil && len(apiResults) > 0 {
			// Convert to local type
			results := make([]CostResult, len(apiResults))
			for i, r := range apiResults {
				costs := make([]ShippingCost, len(r.Costs))
				for j, sc := range r.Costs {
					costArr := make([]Cost, len(sc.Cost))
					for k, c := range sc.Cost {
						costArr[k] = Cost{Value: c.Value, Etd: c.Etd, Note: c.Note}
					}
					costs[j] = ShippingCost{
						Service:     sc.Service,
						Description: sc.Description,
						Cost:        costArr,
					}
				}
				results[i] = CostResult{
					Code:  r.Code,
					Name:  r.Name,
					Costs: costs,
				}
			}
			return c.JSON(presenter.SuccessResponse(results))
		}
	}

	// Fallback to mock data
	baseCost := weight * 10
	if baseCost < 9000 {
		baseCost = 9000
	}

	results := []CostResult{
		{
			Code: "jne",
			Name: "Jalur Nugraha Ekakurir (JNE)",
			Costs: []ShippingCost{
				{
					Service:     "REG",
					Description: "Layanan Reguler",
					Cost:        []Cost{{Value: baseCost, Etd: "2-3", Note: ""}},
				},
				{
					Service:     "YES",
					Description: "Yakin Esok Sampai",
					Cost:        []Cost{{Value: baseCost + 10000, Etd: "1", Note: ""}},
				},
			},
		},
		{
			Code: "tiki",
			Name: "Citra Van Titipan Kilat (TIKI)",
			Costs: []ShippingCost{
				{
					Service:     "REG",
					Description: "Regular Service",
					Cost:        []Cost{{Value: baseCost - 1000, Etd: "3-4", Note: ""}},
				},
				{
					Service:     "ONS",
					Description: "Over Night Service",
					Cost:        []Cost{{Value: baseCost + 8000, Etd: "1", Note: ""}},
				},
			},
		},
		{
			Code: "sicepat",
			Name: "SiCepat Express",
			Costs: []ShippingCost{
				{
					Service:     "REG",
					Description: "Reguler",
					Cost:        []Cost{{Value: baseCost - 2000, Etd: "2-3", Note: ""}},
				},
				{
					Service:     "BEST",
					Description: "Besok Sampai Tujuan",
					Cost:        []Cost{{Value: baseCost + 5000, Etd: "1", Note: ""}},
				},
			},
		},
	}

	return c.JSON(presenter.SuccessResponse(results))
}

// TrackWaybill godoc
// @Summary Track Waybill
// @Description Track shipment by waybill number
// @Tags Shipping
// @Produce json
// @Param waybill query string true "Waybill/Tracking number"
// @Param courier query string true "Courier code (jne, pos, tiki)"
// @Success 200 {object} presenter.SuccessResponseSwagger
// @Router /shipping/track [get]
func (h *ShippingHandler) TrackWaybill(c *fiber.Ctx) error {
	waybill := c.Query("waybill", "")
	courier := c.Query("courier", "")

	if waybill == "" || courier == "" {
		return c.Status(400).JSON(presenter.ErrorResponse(fiber.NewError(400, "waybill and courier are required")))
	}

	if !h.useAPI || h.rajaOngkir == nil {
		return c.Status(503).JSON(presenter.ErrorResponse(fiber.NewError(503, "Tracking service unavailable")))
	}

	result, err := h.rajaOngkir.TrackWaybill(waybill, courier)
	if err != nil {
		return c.Status(400).JSON(presenter.ErrorResponse(err))
	}

	return c.JSON(presenter.SuccessResponse(result))
}
