package hedera

type FeatureCollection struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

type Feature struct {
	Type       string     `json:"type"`
	Geometry   Geometry   `json:"geometry"`
	Properties Properties `json:"properties"`
}

type Geometry struct {
	Type        string        `json:"type"`
	Coordinates [][][]float64 `json:"coordinates"`
}
type Properties struct {
	Gid            int    `json:"gid"`
	Idfarmer       string `json:"idfarmer"`
	Companyid      string `json:"companyid"`
	Regionid       string `json:"regionid"`
	Countryid      string `json:"countryid"`
	Stateid        string `json:"stateid"`
	Municipalityid string `json:"municipalityid"`
	Technicalid    string `json:"technicalid"`
	Status         string `json:"status"`
	Activity       string `json:"activity"`
	Bsow           int    `json:"bsow"`
	Product        string `json:"product"`
	Eharvest       string `json:"eharvest"`
	Latcenter      string `json:"latcenter"`
	Loncenter      string `json:"loncenter"`
	Hash           string `json:"hash,omitempty"`
}
