package fastighetsbyran

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/byrnedo/homefinder/internal/pkg/agents"
)

type Crawler struct {
	body string
}

func (f Crawler) Name() string {
	//TODO implement me
	return "Fastighetsbryån"
}

func (f *Crawler) fetch() error {
	req, _ := http.NewRequest("POST", "https://www.fastighetsbyran.com/HemsidanAPI/api/v1/soek/objekt/1/false/", strings.NewReader(`{"valdaMaeklarObjektTyper":[0,14,1,3,9999,4],"valdaNyckelord":[],"valdaLaen":[],"valdaKontor":[],"valdaKommuner":["373"],"valdaNaeromraaden":[4188],"valdaPostorter":[],"inkluderaNyproduktion":true,"inkluderaPaaGaang":true}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	b, _ := ioutil.ReadAll(res.Body)
	f.body = string(b)
	return nil
}

type response struct {
	Results []struct {
		MaeklarObjektId                                                             int      `json:"maeklarObjektId"`
		BildUrl                                                                     string   `json:"bildUrl"`
		LitenRubrik                                                                 string   `json:"litenRubrik"`
		StorRubrik                                                                  string   `json:"storRubrik"`
		MetaData                                                                    []string `json:"metaData"`
		XKoordinat                                                                  float64  `json:"xKoordinat"`
		YKoordinat                                                                  float64  `json:"yKoordinat"`
		PaaGang                                                                     bool     `json:"paaGang"`
		BudgivningPagaar                                                            bool     `json:"budgivningPagaar"`
		AerNyproduktion                                                             bool     `json:"aerNyproduktion"`
		AerProjekt                                                                  bool     `json:"aerProjekt"`
		AerReferensobjekt                                                           bool     `json:"aerReferensobjekt"`
		HarDigitalLiveVisning                                                       bool     `json:"harDigitalLiveVisning"`
		MaeklarId                                                                   int      `json:"maeklarId"`
		Avtalsdag                                                                   *string  `json:"avtalsdag"`
		SenasteTidObjektetBlevTillSalu                                              *string  `json:"senasteTidObjektetBlevTillSalu"`
		SenasteTidpunktSomObjektetBlevIntagetOchSkallAnnoserasFastStatusBaraIntaget *string  `json:"senasteTidpunktSomObjektetBlevIntagetOchSkallAnnoserasFastStatusBaraIntaget"`
	} `json:"results"`
	CurrentPage int `json:"currentPage"`
	PageCount   int `json:"pageCount"`
	PageSize    int `json:"pageSize"`
	RowCount    int `json:"rowCount"`
}

func (f *Crawler) GetForSale(target agents.Target) (listings []agents.Listing, err error) {
	if f.body == "" {
		if err := f.fetch(); err != nil {
			return nil, err
		}
	}
	resp := response{}
	if err := json.Unmarshal([]byte(f.body), &resp); err != nil {
		return nil, err
	}

	for _, item := range resp.Results {

		sqm := 0
		for _, meta := range item.MetaData {
			if strings.HasSuffix(meta, " kvm") {
				sqmF, _ := strconv.ParseFloat(strings.ReplaceAll(strings.TrimSuffix(meta, " kvm"), ",", "."), 64)
				sqm = int(sqmF)
			}
		}
		listings = append(listings, agents.Listing{
			Name:         strings.Join([]string{item.StorRubrik, item.LitenRubrik}, " "),
			Link:         fmt.Sprintf("https://www.fastighetsbyran.com/sv/sverige/objekt/?objektid=%d", item.MaeklarObjektId),
			Image:        item.BildUrl,
			Upcoming:     item.PaaGang,
			Facts:        item.MetaData,
			SquareMetres: sqm,
		})
	}
	return

	/*
		curl '' -X POST -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:95.0) Gecko/20100101 Firefox/95.0' -H 'Accept: application/json' -H 'Accept-Language: en-GB,en;q=0.5' -H 'Accept-Encoding: gzip, deflate, br' -H 'Referer: https://www.fastighetsbyran.com/sv/sverige/till-salu' -H 'Content-Type: application/json' -H 'spraak: sv' -H 'webbmarknad: 204' -H 'Origin: https://www.fastighetsbyran.com' -H 'Connection: keep-alive' -H 'Cookie: FSDTDATA=Hasdata=0&Expires=2022-10-27 13:26:44; FSDANONYMOUS=D6bis88m2AEkAAAAODQ2NGNjNzMtYTJiYy00YWRiLTlmNjItZTlkNWRjZmU0NWFilvQPsmxasXgCFiXQrxavk8_k6PQ1; _bamls_usid=2fbf659b-ffc5-4b10-8af0-fd60a3f002f6; _ga=GA1.2.1746526103.1635334015; wsa903281_Language=1; _bamls_cuid=r1XV7TP9FnwnIJqiimMH; _fbc=fb.1.1635334016222.IwAR0ORUNoXf6c650wzKFxPqQULJf7qsBz672g9plN8QJO05k2b5N5Rq__X4U; _fbp=fb.1.1635334016223.1204844456; _gid=GA1.2.1696002276.1639412664; _gat_UA-2769054-53=1' -H 'Sec-Fetch-Dest: empty' -H 'Sec-Fetch-Mode: cors' -H 'Sec-Fetch-Site: same-origin' --data-raw '{"valdaMaeklarObjektTyper":[],"valdaNyckelord":[],"valdaLaen":[],"valdaKontor":[],"valdaKommuner":["373"],"valdaNaeromraaden":[4188],"valdaPostorter":[],"inkluderaNyproduktion":true,"inkluderaPaaGaang":true}'
	*/
}
