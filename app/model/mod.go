package model

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"

	couchdb "github.com/leesper/couchdb-golang"
)

/*
*	Structs für die Datenhaltung
 */

//KarteikastenData für die Datenhaltung
type KarteikastenData struct {
	ID                       string        `json:"_id"`
	Rev                      string        `json:"_rev"`
	Type                     string        `json:"type"`
	Ueberkategorie           string        `json:"ueberkategorie"`
	UnterKategorie           string        `json:"unterkategorie"`
	KarteikastenName         string        `json:"karteikastenname"`
	ErstellerName            string        `json:"erstellername"`
	Sichtbarkeit             string        `json:"sichtbarkeit"`
	AnzKarten                string        `json:"anzkarten"`
	KarteikastenBeschreibung string        `json:"karteikastenbeschreibung"`
	Karteikarten             []Karteikarte `json:"karteikarten"`
	Gelerntvon               string        `json:"gelerntvon"`
	Fortschritt              string        `json:"fortschritt"`
	KeinEigentum             bool          `json:"keineigentum"`
}

//Karteikarte Datenhaltung für die Karten
type Karteikarte struct {
	KartenName string `json:"kartenname"`
	Fach       string `json:"fach"`
	Frage      string `json:"frage"`
	Antwort    string `json:"antwort"`
	KartenNr   string `json:"kartennr"`
	IstGruen   bool
	Fach0      bool
	Fach1      bool
	Fach2      bool
	Fach3      bool
	Fach4      bool
}

/*
*	HelferStrucs
 */

//KarteienData (Struct für das Karteikasten Tmpl.)
type KarteienData struct {
	Karteititel       string
	KastenBelegt      bool
	KarteikastenDaten []KarteikastenData
}

//BaseDat (Struct für das Base Tmpl.)
type BaseDat struct {
	LoggedIn bool
	Link     string
	Name     string
	TagMK    string
	TagK     string
}

//IndexDat (Struct für das Index Tmpl.)
type IndexDat struct {
	BaseDaten  BaseDat
	ErrorMsg   string
	Anzuserges string
	Anzkartges string
	Anzkatges  string
}

//RegisterDat (Struct für das Register Tmpl.)
type RegisterDat struct {
	BaseDaten     BaseDat
	ErrorMsg      string
	ErrorUsername string
	ErrorEmail    string
}

//KarteikastenDat (Struct für das Karteikasten Tmpl.)
type KarteikastenDat struct {
	BaseDaten BaseDat
	Karteien  [5]KarteienData
	ErrorMsg  string
}

//MeineKarteikastenDat (Struct für das MeineKarteikasten Tmpl.)
type MeineKarteikastenDat struct {
	BaseDaten BaseDat
	Karteien  [2]KarteienData
}

//ViewDat (Struct für das View Tmpl.)
type ViewDat struct {
	BaseDaten    BaseDat
	Karteikasten KarteikastenData
	ErrorMsg     string
	KartenName   string
	Fach0        bool
	Fach1        bool
	Fach2        bool
	Fach3        bool
	Fach4        bool
	Frage        string
	Antwort      string
}

//LernDat (Struct für das Lern Tmpl.)
type LernDat struct {
	BaseDaten    BaseDat
	Karteikasten KarteikastenData
	ErrorMsg     string
	KartenName   string
	FachAnz0     string
	FachAnz1     string
	FachAnz2     string
	FachAnz3     string
	FachAnz4     string
	Fach0        bool
	Fach1        bool
	Fach2        bool
	Fach3        bool
	Fach4        bool
	Frage        string
	Antwort      string
	Aufgedeckt   bool
	KartenNr     string
}

//EditDat (Struct für das Profil Tmpl.)
type EditDat struct {
	BaseDaten    BaseDat
	ErrorMsg     string
	Karteikasten KarteikastenData
}

//EditNDat Struct für das EditN Tmpl
type EditNDat struct {
	BaseDaten    BaseDat
	Karteikasten KarteikastenData
	ErrorMsg     string
	KartenName   string
	Frage        string
	Antwort      string
	KartenNr     string
}

//ProfilDat (Struct für das Profil Tmpl.)
type ProfilDat struct {
	BaseDaten   BaseDat
	ID          string
	ErrorPW     string
	ErrorEmail  string
	ErrorMsg    string
	Name        string
	Email       string
	MeineKarten string
	TagMK       string
	Erstellt    string
}

var karteienNamen = [5]string{"Naturwissenschaften", "Sprachen", "Gesellschaft", "Wirtschaft", "Geisteswissenschaften"}
var meinekarteienNamen = [2]string{"Selbst erstellte Karteikästen", "Gelernte Karteikästen andere Nutzer"}

// Db handle
var btDBS *couchdb.Database

func init() {
	var err error
	btDBS, err = couchdb.NewDatabase("http://localhost:5984/braintrain")
	if err != nil {
		panic(err)
	}
}

//CreateKarteikasten erstellt einen Karteikasten und speichern Ihn in der Datenbank
func CreateKarteikasten(username string, name string, beschreibung string, kategorie string, sichtbarkeit string) (string, error) {

	var karteikasten KarteikastenData

	karteikasten.AnzKarten = "0"
	karteikasten.ErstellerName = username
	karteikasten.KarteikastenName = name
	karteikasten.KarteikastenBeschreibung = beschreibung
	karteikasten.Sichtbarkeit = sichtbarkeit
	karteikasten.Type = "Karteikasten"

	var slice = strings.Split(kategorie, ",")

	karteikasten.Ueberkategorie = slice[0]
	karteikasten.UnterKategorie = slice[1]

	var karteikastenMap, _ = kartei2Map(karteikasten)

	delete(karteikastenMap, "_id")
	delete(karteikastenMap, "_rev")
	delete(karteikastenMap, "fortschritt")
	delete(karteikastenMap, "keineigentum")

	id, _, err := btDBS.Save(karteikastenMap, nil)

	return id, err
}

//UpdateKarteikasten updated einen gegeben Karteikasten, der über seine ID gesucht wird.
func UpdateKarteikasten(id string, name string, beschreibung string, kategorie string, sichtbarkeit string) error {
	var karteikasten KarteikastenData

	karteikastenMap, _ := btDBS.Get(id, nil)

	karteikasten, _ = map2kartei(karteikastenMap)

	karteikasten.KarteikastenName = name
	karteikasten.KarteikastenBeschreibung = beschreibung
	karteikasten.Sichtbarkeit = sichtbarkeit
	karteikasten.Type = "Karteikasten"

	var slice = strings.Split(kategorie, ",")

	karteikasten.Ueberkategorie = slice[0]
	karteikasten.UnterKategorie = slice[1]

	karteikastenMap, _ = kartei2Map(karteikasten)

	delete(karteikastenMap, "fortschritt")
	delete(karteikastenMap, "keineigentum")

	err := btDBS.Set(id, karteikastenMap)

	return err
}

//UpdateKarteiLernen Um die Karte nach dem Lernen zu aktualisieren
func UpdateKarteiLernen(id string, choise string, kartennr string) {

	karteiMap, _ := btDBS.Get(id, nil)
	var kartei, _ = map2kartei(karteiMap)
	var index, _ = strconv.ParseInt(kartennr, 0, 0)
	index--
	var fachnr, _ = strconv.ParseInt(kartei.Karteikarten[index].Fach, 0, 0)

	if choise == "true" {
		if fachnr < 4 {
			fachnr++
		}
	} else {
		fachnr = 0
	}

	kartei.Karteikarten[index].Fach = strconv.Itoa(int(fachnr))

	karteiMap, _ = kartei2Map(kartei)

	btDBS.Set(id, karteiMap)
}

//ZufallsKarte liefert den Index einer Zufallskarte
func (karteikasten KarteikastenData) ZufallsKarte() (ret int64) {
	var f string
	r := rand.Intn(14)
	if r == 0 {
		f = "4"
	} else if r == 1 || r == 2 {
		f = "3"
	} else if r == 3 || r == 4 || r == 5 {
		f = "2"
	} else if r == 6 || r == 7 || r == 8 || r == 9 {
		f = "1"
	} else if r == 10 || r == 11 || r == 12 || r == 13 || r == 14 {
		f = "0"
	}
	fmt.Println(f)
	var len = len(karteikasten.Karteikarten)

	var kartenNr = make([]string, len)
	var i = 0

	for j := 0; j < len; j++ {
		if karteikasten.Karteikarten[j].Fach == f {
			kartenNr[i] = karteikasten.Karteikarten[j].KartenNr
			i++
		}
	}

	fmt.Println(kartenNr)
	fmt.Println(i)

	var index int
	if i > 1 {
		i--
		index = rand.Intn(i)
	} else if i == 1 {
		index = 0
	} else {
		var tmp = (len - 1)
		index = rand.Intn(tmp)
	}

	fmt.Println(index)

	ret, _ = strconv.ParseInt(kartenNr[index], 0, 0)
	ret--
	return ret
}

//GetFortschritt liefert den Fortschritt des Karteikastens und die Anzahl der Karten pro Fach zurück
func (karteikasten KarteikastenData) GetFortschritt() (fortschritt int, fach [5]int) {

	if karteikasten.AnzKarten == "0" {
		return 0, fach
	}
	var faecher [5]int
	for i := 0; i < len(karteikasten.Karteikarten); i++ {
		value, _ := strconv.ParseInt(karteikasten.Karteikarten[i].Fach, 0, 0)
		switch value {
		case 0:
			faecher[0]++
		case 1:
			faecher[1]++
		case 2:
			faecher[2]++
		case 3:
			faecher[3]++
		case 4:
			faecher[4]++
		}
	}
	fach = faecher

	for j := 1; j < len(faecher); j++ {
		fortschritt += j * faecher[j]
	}

	wert, _ := strconv.ParseInt(karteikasten.AnzKarten, 0, 0)
	var intwert = int(wert)

	fortschritt *= 100
	fortschritt /= (4 * intwert)

	return fortschritt, faecher
}

//GetBaseData lieft Struct für BaseDat
func GetBaseData(link string, name string) BaseDat {
	var base []BaseDat
	var ret BaseDat
	ret.Link = link

	if name == "" {
		ret.LoggedIn = false

	} else {
		var karteikaesten []KarteikastenData
		userMap, _ := btDBS.QueryJSON(`
						{
							"selector": {
								"type": {
										"$eq": "Benutzer"
								},
								"name": {
									"$eq": "` + name + `"
								}
							}
						}`)

		mapstructure.Decode(userMap, &base)

		for i := 0; i < len(base); i++ {
			if base[i].Name == name {
				ret = base[i]
				break
			}
		}

		karteikastenMap, _ := btDBS.QueryJSON(`
		{
			"selector": {
				"type": {
						"$eq": "Karteikasten"
				},
				"erstellername":{
					"$eq":"` + name + `"
				}
			}
		}`)

		mapstructure.Decode(karteikastenMap, &karteikaesten)

		ret.TagMK = strconv.Itoa(len(karteikaesten))

	}
	var kaesten []KarteikastenData

	kastenMap, _ := btDBS.QueryJSON(`
	{
		"selector": {
			 "type": {
					"$eq": "Karteikasten"
			 },
			 "sichtbarkeit":{
				"$eq": "Öffentlich"
			 }
		}
	 }`)

	mapstructure.Decode(kastenMap, &kaesten)

	ret.TagK = strconv.Itoa(len(kaesten))

	return ret

}

//GetIndexData liefter Struct für das Index Tmpl.
func GetIndexData(name string) (IndexDat, error) {
	var indexdata IndexDat

	userMap, _ := btDBS.QueryJSON(`
	{
		"selector": {
			 "type": {
					"$eq": "Benutzer"
			 } 
		}
	 }`)

	var base []BaseDat

	mapstructure.Decode(userMap, &base)

	indexdata.BaseDaten = GetBaseData("", name)

	indexdata.Anzuserges = strconv.Itoa(len(base))

	karteikaestenMap, _ := btDBS.QueryJSON(`
	{
		"selector": {
			 "type": {
					"$eq": "Karteikasten"
			 }
		},
		"fields": ["anzkarten"]
	 }`)

	var karteikaestenStruct []KarteikastenData
	mapstructure.Decode(karteikaestenMap, &karteikaestenStruct)

	indexdata.Anzkatges = strconv.Itoa(len(karteikaestenStruct))

	var anzkarttmp int

	for i := 0; i < len(karteikaestenStruct); i++ {
		tmp, _ := strconv.ParseInt(karteikaestenStruct[i].AnzKarten, 0, 0)
		anzkarttmp += int(tmp)
	}

	indexdata.Anzkartges = strconv.Itoa(anzkarttmp)

	return indexdata, nil
}

//GetKarteikastenData liefert das Struct für das Karteikasten Tmpl.
func GetKarteikastenData(name string) (KarteikastenDat, error) {
	var karteikaesten KarteikastenDat
	var karteien [5]KarteienData

	karteikaesten.BaseDaten = GetBaseData("/karteikasten", name)

	for i := 0; i < 5; i++ {
		karteikaestenMap, _ := btDBS.QueryJSON(`
		{
			"selector": {
				 "type": {
						"$eq": "Karteikasten"
				 },
				 "ueberkategorie":{
					"$eq": "` + karteienNamen[i] + `"
				 },
				 "sichtbarkeit": {
					 "$eq": "Öffentlich"
				}
			}
		 }`)
		var karteikaestenStruct []KarteikastenData
		mapstructure.Decode(karteikaestenMap, &karteikaestenStruct)

		index := 0
		for _, v := range karteikaestenMap {
			karteikaestenStruct[index].ID = v["_id"].(string)
			index++

		}

		var karteienStruc KarteienData
		karteienStruc.KarteikastenDaten = karteikaestenStruct
		karteienStruc.Karteititel = karteienNamen[i]

		if len(karteikaestenStruct) == 0 {
			karteienStruc.KastenBelegt = false
		} else {
			karteienStruc.KastenBelegt = true
		}

		karteien[i] = karteienStruc

	}

	karteikaesten.Karteien = karteien
	return karteikaesten, nil
}

//GetMeineKarteikastenData lieft das Struct für das MeineKarteikasten Tmpl.
func GetMeineKarteikastenData(name string) (MeineKarteikastenDat, error) {
	var karteikaesten MeineKarteikastenDat
	var karteien [2]KarteienData

	karteikaesten.BaseDaten = GetBaseData("/meinekarteien", name)

	for i := 0; i < len(karteien); i++ {
		var karteienStruc KarteienData
		karteienStruc.Karteititel = meinekarteienNamen[i]

		var karteikaestenStruct []KarteikastenData
		switch i {
		case 0:
			{
				karteikaestenMap, _ := btDBS.QueryJSON(`
		{
			"selector": {
				 "type": {
						"$eq": "Karteikasten"
				 },
				 "erstellername":{
					"$eq": "` + name + `"
				 }
			}
		 }`)
				mapstructure.Decode(karteikaestenMap, &karteikaestenStruct)

				index := 0
				for _, v := range karteikaestenMap {

					tmp, _ := karteikaestenStruct[index].GetFortschritt()

					karteikaestenStruct[index].Fortschritt = strconv.Itoa(tmp)

					if karteikaestenStruct[index].ErstellerName == name {
						karteikaestenStruct[index].KeinEigentum = false
					} else {
						karteikaestenStruct[index].KeinEigentum = true
					}

					karteikaestenStruct[index].ID = v["_id"].(string)
					index++
				}

				karteienStruc.KarteikastenDaten = karteikaestenStruct

				if len(karteikaestenStruct) == 0 {
					karteienStruc.KastenBelegt = false
				} else {
					karteienStruc.KastenBelegt = true
				}
				karteien[i] = karteienStruc
			}

		case 1:
			{
				//Für öffentliche Karteien
				karteikaestenMap, _ := btDBS.QueryJSON(`
			{
				"selector": {
					 "type": { 
							"$eq": "Kartei" 
					 },
					 "gelerntvon":{
						"$eq": "` + name + `"
					 }
				}
			 }`)
				mapstructure.Decode(karteikaestenMap, &karteikaestenStruct)

				index := 0
				for _, v := range karteikaestenMap {

					if karteikaestenStruct[index].ErstellerName == name {
						karteikaestenStruct[index].KeinEigentum = false
					} else {
						karteikaestenStruct[index].KeinEigentum = true
					}

					karteikaestenStruct[index].ID = v["_id"].(string)
					index++
				}
				karteienStruc.KarteikastenDaten = karteikaestenStruct

				if len(karteikaestenStruct) == 0 {
					karteienStruc.KastenBelegt = false
				} else {
					karteienStruc.KastenBelegt = true
				}
				karteien[i] = karteienStruc
			}
		}

	}

	karteikaesten.Karteien = karteien
	return karteikaesten, nil
}

//GetViewData lieft das Struct für das View Tmpl.
func GetViewData(name string, id string, nr string) (ViewDat, error) {
	var viewdata ViewDat

	viewdata.BaseDaten = GetBaseData("/karteikasten", name)

	karteikasten, _ := btDBS.Get(id, nil)

	var karteikastendata KarteikastenData

	mapstructure.Decode(karteikasten, &karteikastendata)
	karteikastendata.ID = karteikasten["_id"].(string)

	tmp, _ := karteikastendata.GetFortschritt()

	karteikastendata.Fortschritt = strconv.Itoa(tmp)

	var index int64
	if nr != "" {
		index, _ = strconv.ParseInt(nr, 0, 0)
		index--
	} else {
		index = 0
	}
	viewdata.Karteikasten = karteikastendata

	switch viewdata.Karteikasten.Karteikarten[index].Fach {

	case "0":
		viewdata.Fach0 = true
	case "1":
		viewdata.Fach1 = true
	case "2":
		viewdata.Fach2 = true
	case "3":
		viewdata.Fach3 = true
	case "4":
		viewdata.Fach4 = true
	}

	viewdata.Karteikasten.Karteikarten[index].IstGruen = true
	viewdata.KartenName = viewdata.Karteikasten.Karteikarten[index].KartenName
	viewdata.Frage = viewdata.Karteikasten.Karteikarten[index].Frage
	viewdata.Antwort = viewdata.Karteikasten.Karteikarten[index].Antwort

	return viewdata, nil
}

//GetLernData lieft das Struct für das View Tmpl.
func GetLernData(name string, id string, aufgedeckt string, kartennr string) (LernDat, error) {

	var lerndata LernDat

	lerndata.BaseDaten = GetBaseData("/karteikasten/lern", name)

	karteikasten, _ := btDBS.Get(id, nil)

	var karteikastendata KarteikastenData

	mapstructure.Decode(karteikasten, &karteikastendata)
	karteikastendata.ID = karteikasten["_id"].(string)

	var index int64

	if kartennr == "" {
		index = karteikastendata.ZufallsKarte()
	} else {
		index, _ = strconv.ParseInt(kartennr, 0, 0)
		index--
	}

	lerndata.Karteikasten = karteikastendata

	aufge, _ := strconv.ParseBool(aufgedeckt)

	if aufge == true {
		lerndata.Aufgedeckt = true
	}

	fortschritt, Fachanz := lerndata.Karteikasten.GetFortschritt()

	lerndata.Karteikasten.Fortschritt = strconv.Itoa(fortschritt)

	lerndata.FachAnz0 = strconv.Itoa(Fachanz[0])
	lerndata.FachAnz1 = strconv.Itoa(Fachanz[1])
	lerndata.FachAnz2 = strconv.Itoa(Fachanz[2])
	lerndata.FachAnz3 = strconv.Itoa(Fachanz[3])
	lerndata.FachAnz4 = strconv.Itoa(Fachanz[4])

	switch lerndata.Karteikasten.Karteikarten[index].Fach {

	case "0":
		lerndata.Fach0 = true
	case "1":
		lerndata.Fach1 = true
	case "2":
		lerndata.Fach2 = true
	case "3":
		lerndata.Fach3 = true
	case "4":
		lerndata.Fach4 = true
	}

	lerndata.KartenName = lerndata.Karteikasten.Karteikarten[index].KartenName
	lerndata.Frage = lerndata.Karteikasten.Karteikarten[index].Frage
	lerndata.Antwort = lerndata.Karteikasten.Karteikarten[index].Antwort
	lerndata.KartenNr = lerndata.Karteikasten.Karteikarten[index].KartenNr

	return lerndata, nil
}

//GetEditData Liefert das Struct für das Edit Tmpl.
func GetEditData(name string, id string) (EditDat, error) {
	var editdata EditDat

	editdata.BaseDaten = GetBaseData("", name)

	if id == "" {
		return editdata, nil
	}
	karteikastenMap, _ := btDBS.Get(id, nil)

	var karteikastenStruct, _ = map2kartei(karteikastenMap)

	editdata.Karteikasten = karteikastenStruct

	return editdata, nil
}

//GetEditNextData Liefert das Struct für das Edit Tmpl.
func GetEditNextData(name string, id string, nr string) (EditNDat, error) {
	var editNextdata EditNDat
	editNextdata.BaseDaten = GetBaseData("", name)

	editNextdataMap, _ := btDBS.Get(id, nil)

	var karteikastenStruct, _ = map2kartei(editNextdataMap)

	fortschritt, _ := karteikastenStruct.GetFortschritt()

	karteikastenStruct.Fortschritt = strconv.Itoa(fortschritt)

	editNextdata.Karteikasten = karteikastenStruct

	var index int64

	if nr == "" {

		if editNextdata.Karteikasten.AnzKarten == "0" {
			editNextdata.KartenNr = ""
			editNextdata.KartenName = ""
			editNextdata.Frage = ""
			editNextdata.Antwort = ""

			return editNextdata, nil
		} else {
			index = 0
		}

	} else {
		index, _ = strconv.ParseInt(nr, 0, 0)
		if index != 0 {
			index--
		}

	}

	editNextdata.KartenNr = karteikastenStruct.Karteikarten[index].KartenNr
	editNextdata.KartenName = karteikastenStruct.Karteikarten[index].KartenName
	editNextdata.Frage = karteikastenStruct.Karteikarten[index].Frage
	editNextdata.Antwort = karteikastenStruct.Karteikarten[index].Antwort

	return editNextdata, nil

}

//GetProfilData Liefert das Struct für das Profil Tmpl.
func GetProfilData(name string) (ProfilDat, error) {
	var profildata ProfilDat

	profildata.BaseDaten = GetBaseData("", name)
	profildata.Name = name

	var userdaten []UserData
	var userdat UserData

	var userMap, _ = btDBS.QueryJSON(`
	{
		"selector": {
			"type": {
					"$eq": "Benutzer"
			},
			"name":{
				"$eq":"` + name + `"
			}
		}
	}`)

	mapstructure.Decode(userMap, &userdaten)

	userdat = userdaten[0]

	profildata.Email = userdat.Email
	profildata.Erstellt = userdat.Erstellt
	profildata.TagMK = profildata.BaseDaten.TagMK

	profildata.ID = userMap[0]["_id"].(string)

	var karteikaesten []KarteikastenData
	var karteikastenMap, _ = btDBS.QueryJSON(`
		{
			"selector": {
				"type": {
						"$eq": "Karteikasten"
				},
				"erstellername":{
					"$eq":"` + name + `"
				}
			}
		}`)

	mapstructure.Decode(karteikastenMap, &karteikaesten)

	for i := 0; i < len(karteikaesten); i++ {
		profildata.MeineKarten += karteikaesten[i].AnzKarten
	}

	return profildata, nil
}

//UpdateKarteikastenKarten erneuert die Karten im Karteikasten
func UpdateKarteikastenKarten(id string, KartenName string, Frage string, Antwort string, KartenNr string) error {

	kastenMap, _ := btDBS.Get(id, nil)

	var karteikastenStruct, _ = map2kartei(kastenMap)

	kartennr, _ := strconv.ParseInt(KartenNr, 0, 0)

	if kartennr == 0 {

		kartenNRneu, _ := strconv.ParseInt(karteikastenStruct.AnzKarten, 0, 0)
		kartenNRneuInt := int(kartenNRneu + 1)
		karteikastenStruct.AnzKarten = strconv.Itoa(kartenNRneuInt)

		var Karten = make([]Karteikarte, kartenNRneuInt)

		for i := 0; i < int(kartenNRneu); i++ {
			Karten[i] = karteikastenStruct.Karteikarten[i]
		}

		if Antwort == "" {
			Karten[kartenNRneu].Antwort = ""
		} else {
			Karten[kartenNRneu].Antwort = Antwort
		}
		if Frage == "" {
			Karten[kartenNRneu].Frage = ""
		} else {
			Karten[kartenNRneu].Frage = Frage
		}
		if KartenName == "" {
			Karten[kartenNRneu].KartenName = ""
		} else {
			Karten[kartenNRneu].KartenName = KartenName
		}

		Karten[kartenNRneu].KartenNr = strconv.Itoa(kartenNRneuInt)
		Karten[kartenNRneu].Fach = "0"

		karteikastenStruct.Karteikarten = Karten

		KasteMap, _ := kartei2Map(karteikastenStruct)

		fmt.Println(KasteMap)

		ret := btDBS.Set(id, KasteMap)

		return ret

	} else {

		karteikastenStruct.Karteikarten[kartennr-1].KartenName = KartenName
		karteikastenStruct.Karteikarten[kartennr-1].Frage = Frage
		karteikastenStruct.Karteikarten[kartennr-1].Antwort = Antwort

		KasteMap, _ := kartei2Map(karteikastenStruct)

		ret := btDBS.Set(id, KasteMap)
		return ret
	}

}

//Für alle DatenhaltungsStructs

// Convert from User struct to map[string]interface{} as required by golang-couchdb methods
func kartei2Map(karteiStruc KarteikastenData) (karteiMap map[string]interface{}, err error) {
	uJSON, err := json.Marshal(karteiStruc)
	json.Unmarshal(uJSON, &karteiMap)

	return karteiMap, err
}

// Convert from map[string]interface{} to User struct as required by golang-couchdb methods
func map2kartei(karteiMap map[string]interface{}) (karteiStruc KarteikastenData, err error) {
	uJSON, err := json.Marshal(karteiMap)
	json.Unmarshal(uJSON, &karteiStruc)

	return karteiStruc, err
}
