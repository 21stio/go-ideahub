package middleware

import (
	"net/http"
	"strings"
	"github.com/oschwald/geoip2-golang"
	//"github.com/tomasen/realip"
	"net"
	"github.com/21stio/go-ideahub/types"
	"crypto/sha256"
	"bytes"
	"time"
	"log"
	"github.com/davecgh/go-spew/spew"
	"github.com/21stio/go-ideahub/routes"
	"github.com/21stio/go-ideahub/queries"
	"encoding/base64"
	"github.com/tomasen/realip"
	"github.com/21stio/go-ideahub/utils"
)

var salt = []byte(utils.GetEnv("SALT"))
var isoToTz = map[string]string{}
var locs = map[string]*time.Location{}
var chLen = 1000
var visitsCh = make(chan types.Visit, chLen)

func init() {
	isoToTz["AD"] = "Europe/Andorra"
	isoToTz["AE"] = "Asia/Dubai"
	isoToTz["AF"] = "Asia/Kabul"
	isoToTz["AG"] = "America/Antigua"
	isoToTz["AI"] = "America/Anguilla"
	isoToTz["AL"] = "Europe/Tirane"
	isoToTz["AM"] = "Asia/Yerevan"
	isoToTz["AO"] = "Africa/Luanda"
	isoToTz["AQ"] = "Antarctica/Vostok"
	isoToTz["AR"] = "America/Argentina/Ushuaia"
	isoToTz["AS"] = "Pacific/Pago_Pago"
	isoToTz["AT"] = "Europe/Vienna"
	isoToTz["AU"] = "Australia/Darwin"
	isoToTz["AW"] = "America/Aruba"
	isoToTz["AX"] = "Europe/Mariehamn"
	isoToTz["AZ"] = "Asia/Baku"
	isoToTz["BA"] = "Europe/Sarajevo"
	isoToTz["BB"] = "America/Barbados"
	isoToTz["BD"] = "Asia/Dhaka"
	isoToTz["BE"] = "Europe/Brussels"
	isoToTz["BF"] = "Africa/Ouagadougou"
	isoToTz["BG"] = "Europe/Sofia"
	isoToTz["BH"] = "Asia/Bahrain"
	isoToTz["BI"] = "Africa/Bujumbura"
	isoToTz["BJ"] = "Africa/Porto-Novo"
	isoToTz["BL"] = "America/St_Barthelemy"
	isoToTz["BM"] = "Atlantic/Bermuda"
	isoToTz["BN"] = "Asia/Brunei"
	isoToTz["BO"] = "America/La_Paz"
	isoToTz["BQ"] = "America/Kralendijk"
	isoToTz["BR"] = "America/Cuiaba"
	isoToTz["BS"] = "America/Nassau"
	isoToTz["BT"] = "Asia/Thimphu"
	isoToTz["BW"] = "Africa/Gaborone"
	isoToTz["BY"] = "Europe/Minsk"
	isoToTz["BZ"] = "America/Belize"
	isoToTz["CA"] = "America/Winnipeg"
	isoToTz["CC"] = "Indian/Cocos"
	isoToTz["CD"] = "Africa/Lubumbashi"
	isoToTz["CF"] = "Africa/Bangui"
	isoToTz["CG"] = "Africa/Brazzaville"
	isoToTz["CH"] = "Europe/Zurich"
	isoToTz["CI"] = "Africa/Abidjan"
	isoToTz["CK"] = "Pacific/Rarotonga"
	isoToTz["CL"] = "America/Santiago"
	isoToTz["CM"] = "Africa/Douala"
	isoToTz["CN"] = "Asia/Shanghai"
	isoToTz["CO"] = "America/Bogota"
	isoToTz["CR"] = "America/Costa_Rica"
	isoToTz["CU"] = "America/Havana"
	isoToTz["CV"] = "Atlantic/Cape_Verde"
	isoToTz["CW"] = "America/Curacao"
	isoToTz["CX"] = "Indian/Christmas"
	isoToTz["CY"] = "Asia/Famagusta"
	isoToTz["CZ"] = "Europe/Prague"
	isoToTz["DE"] = "Europe/Berlin"
	isoToTz["DJ"] = "Africa/Djibouti"
	isoToTz["DK"] = "Europe/Copenhagen"
	isoToTz["DM"] = "America/Dominica"
	isoToTz["DO"] = "America/Santo_Domingo"
	isoToTz["DZ"] = "Africa/Algiers"
	isoToTz["EC"] = "Pacific/Galapagos"
	isoToTz["EE"] = "Europe/Tallinn"
	isoToTz["EG"] = "Africa/Cairo"
	isoToTz["EH"] = "Africa/El_Aaiun"
	isoToTz["ER"] = "Africa/Asmara"
	isoToTz["ES"] = "Europe/Madrid"
	isoToTz["ET"] = "Africa/Addis_Ababa"
	isoToTz["FI"] = "Europe/Helsinki"
	isoToTz["FJ"] = "Pacific/Fiji"
	isoToTz["FK"] = "Atlantic/Stanley"
	isoToTz["FM"] = "Pacific/Pohnpei"
	isoToTz["FO"] = "Atlantic/Faroe"
	isoToTz["FR"] = "Europe/Paris"
	isoToTz["GA"] = "Africa/Libreville"
	isoToTz["GB"] = "Europe/London"
	isoToTz["GD"] = "America/Grenada"
	isoToTz["GE"] = "Asia/Tbilisi"
	isoToTz["GF"] = "America/Cayenne"
	isoToTz["GG"] = "Europe/Guernsey"
	isoToTz["GH"] = "Africa/Accra"
	isoToTz["GI"] = "Europe/Gibraltar"
	isoToTz["GL"] = "America/Thule"
	isoToTz["GM"] = "Africa/Banjul"
	isoToTz["GN"] = "Africa/Conakry"
	isoToTz["GP"] = "America/Guadeloupe"
	isoToTz["GQ"] = "Africa/Malabo"
	isoToTz["GR"] = "Europe/Athens"
	isoToTz["GS"] = "Atlantic/South_Georgia"
	isoToTz["GT"] = "America/Guatemala"
	isoToTz["GU"] = "Pacific/Guam"
	isoToTz["GW"] = "Africa/Bissau"
	isoToTz["GY"] = "America/Guyana"
	isoToTz["HK"] = "Asia/Hong_Kong"
	isoToTz["HN"] = "America/Tegucigalpa"
	isoToTz["HR"] = "Europe/Zagreb"
	isoToTz["HT"] = "America/Port-au-Prince"
	isoToTz["HU"] = "Europe/Budapest"
	isoToTz["ID"] = "Asia/Jayapura"
	isoToTz["IE"] = "Europe/Dublin"
	isoToTz["IL"] = "Asia/Jerusalem"
	isoToTz["IM"] = "Europe/Isle_of_Man"
	isoToTz["IN"] = "Asia/Kolkata"
	isoToTz["IO"] = "Indian/Chagos"
	isoToTz["IQ"] = "Asia/Baghdad"
	isoToTz["IR"] = "Asia/Tehran"
	isoToTz["IS"] = "Atlantic/Reykjavik"
	isoToTz["IT"] = "Europe/Rome"
	isoToTz["JE"] = "Europe/Jersey"
	isoToTz["JM"] = "America/Jamaica"
	isoToTz["JO"] = "Asia/Amman"
	isoToTz["JP"] = "Asia/Tokyo"
	isoToTz["KE"] = "Africa/Nairobi"
	isoToTz["KG"] = "Asia/Bishkek"
	isoToTz["KH"] = "Asia/Phnom_Penh"
	isoToTz["KI"] = "Pacific/Kiritimati"
	isoToTz["KM"] = "Indian/Comoro"
	isoToTz["KN"] = "America/St_Kitts"
	isoToTz["KP"] = "Asia/Pyongyang"
	isoToTz["KR"] = "Asia/Seoul"
	isoToTz["KW"] = "Asia/Kuwait"
	isoToTz["KY"] = "America/Cayman"
	isoToTz["KZ"] = "Asia/Qyzylorda"
	isoToTz["LA"] = "Asia/Vientiane"
	isoToTz["LB"] = "Asia/Beirut"
	isoToTz["LC"] = "America/St_Lucia"
	isoToTz["LI"] = "Europe/Vaduz"
	isoToTz["LK"] = "Asia/Colombo"
	isoToTz["LR"] = "Africa/Monrovia"
	isoToTz["LS"] = "Africa/Maseru"
	isoToTz["LT"] = "Europe/Vilnius"
	isoToTz["LU"] = "Europe/Luxembourg"
	isoToTz["LV"] = "Europe/Riga"
	isoToTz["LY"] = "Africa/Tripoli"
	isoToTz["MA"] = "Africa/Casablanca"
	isoToTz["MC"] = "Europe/Monaco"
	isoToTz["MD"] = "Europe/Chisinau"
	isoToTz["ME"] = "Europe/Podgorica"
	isoToTz["MF"] = "America/Marigot"
	isoToTz["MG"] = "Indian/Antananarivo"
	isoToTz["MH"] = "Pacific/Majuro"
	isoToTz["MK"] = "Europe/Skopje"
	isoToTz["ML"] = "Africa/Bamako"
	isoToTz["MM"] = "Asia/Yangon"
	isoToTz["MN"] = "Asia/Ulaanbaatar"
	isoToTz["MO"] = "Asia/Macau"
	isoToTz["MP"] = "Pacific/Saipan"
	isoToTz["MQ"] = "America/Martinique"
	isoToTz["MR"] = "Africa/Nouakchott"
	isoToTz["MS"] = "America/Montserrat"
	isoToTz["MT"] = "Europe/Malta"
	isoToTz["MU"] = "Indian/Mauritius"
	isoToTz["MV"] = "Indian/Maldives"
	isoToTz["MW"] = "Africa/Blantyre"
	isoToTz["MX"] = "America/Tijuana"
	isoToTz["MY"] = "Asia/Kuching"
	isoToTz["AD"] = "Europe/Andorra"
	isoToTz["AE"] = "Asia/Dubai"
	isoToTz["AF"] = "Asia/Kabul"
	isoToTz["AG"] = "America/Antigua"
	isoToTz["AI"] = "America/Anguilla"
	isoToTz["AL"] = "Europe/Tirane"
	isoToTz["AM"] = "Asia/Yerevan"
	isoToTz["AO"] = "Africa/Luanda"
	isoToTz["AQ"] = "Antarctica/Vostok"
	isoToTz["AR"] = "America/Argentina/Ushuaia"
	isoToTz["AS"] = "Pacific/Pago_Pago"
	isoToTz["AT"] = "Europe/Vienna"
	isoToTz["AU"] = "Australia/Darwin"
	isoToTz["AW"] = "America/Aruba"
	isoToTz["AX"] = "Europe/Mariehamn"
	isoToTz["PE"] = "America/Lima"
	isoToTz["PF"] = "Pacific/Tahiti"
	isoToTz["PG"] = "Pacific/Bougainville"
	isoToTz["PG"] = "Pacific/Port_Moresby"
	isoToTz["PH"] = "Asia/Manila"
	isoToTz["PK"] = "Asia/Karachi"
	isoToTz["PL"] = "Europe/Warsaw"
	isoToTz["PM"] = "America/Miquelon"
	isoToTz["PN"] = "Pacific/Pitcairn"
	isoToTz["PR"] = "America/Puerto_Rico"
	isoToTz["PS"] = "Asia/Hebron"
	isoToTz["PT"] = "Europe/Lisbon"
	isoToTz["PW"] = "Pacific/Palau"
	isoToTz["PY"] = "America/Asuncion"
	isoToTz["QA"] = "Asia/Qatar"
	isoToTz["RE"] = "Indian/Reunion"
	isoToTz["RO"] = "Europe/Bucharest"
	isoToTz["RS"] = "Europe/Belgrade"
	isoToTz["RU"] = "Asia/Yakutsk"
	isoToTz["RW"] = "Africa/Kigali"
	isoToTz["SA"] = "Asia/Riyadh"
	isoToTz["SB"] = "Pacific/Guadalcanal"
	isoToTz["SC"] = "Indian/Mahe"
	isoToTz["SD"] = "Africa/Khartoum"
	isoToTz["SE"] = "Europe/Stockholm"
	isoToTz["SG"] = "Asia/Singapore"
	isoToTz["SH"] = "Atlantic/St_Helena"
	isoToTz["SI"] = "Europe/Ljubljana"
	isoToTz["SJ"] = "Arctic/Longyearbyen"
	isoToTz["SK"] = "Europe/Bratislava"
	isoToTz["SL"] = "Africa/Freetown"
	isoToTz["SM"] = "Europe/San_Marino"
	isoToTz["SN"] = "Africa/Dakar"
	isoToTz["SO"] = "Africa/Mogadishu"
	isoToTz["SR"] = "America/Paramaribo"
	isoToTz["SS"] = "Africa/Juba"
	isoToTz["ST"] = "Africa/Sao_Tome"
	isoToTz["SV"] = "America/El_Salvador"
	isoToTz["SX"] = "America/Lower_Princes"
	isoToTz["SY"] = "Asia/Damascus"
	isoToTz["SZ"] = "Africa/Mbabane"
	isoToTz["TC"] = "America/Grand_Turk"
	isoToTz["TD"] = "Africa/Ndjamena"
	isoToTz["TF"] = "Indian/Kerguelen"
	isoToTz["TG"] = "Africa/Lome"
	isoToTz["TH"] = "Asia/Bangkok"
	isoToTz["TJ"] = "Asia/Dushanbe"
	isoToTz["TK"] = "Pacific/Fakaofo"
	isoToTz["TL"] = "Asia/Dili"
	isoToTz["TM"] = "Asia/Ashgabat"
	isoToTz["TN"] = "Africa/Tunis"
	isoToTz["TO"] = "Pacific/Tongatapu"
	isoToTz["TR"] = "Europe/Istanbul"
	isoToTz["TT"] = "America/Port_of_Spain"
	isoToTz["TV"] = "Pacific/Funafuti"
	isoToTz["TW"] = "Asia/Taipei"
	isoToTz["TZ"] = "Africa/Dar_es_Salaam"
	isoToTz["UA"] = "Europe/Simferopol"
	isoToTz["UG"] = "Africa/Kampala"
	isoToTz["UM"] = "Pacific/Wake"
	isoToTz["US"] = "Pacific/Honolulu"
	isoToTz["UY"] = "America/Montevideo"
	isoToTz["UZ"] = "Asia/Tashkent"
	isoToTz["VA"] = "Europe/Vatican"
	isoToTz["VC"] = "America/St_Vincent"
	isoToTz["VE"] = "America/Caracas"
	isoToTz["VI"] = "America/St_Thomas"
	isoToTz["VN"] = "Asia/Ho_Chi_Minh"
	isoToTz["VU"] = "Pacific/Efate"
	isoToTz["WF"] = "Pacific/Wallis"
	isoToTz["WS"] = "Pacific/Apia"
	isoToTz["YE"] = "Asia/Aden"
	isoToTz["YT"] = "Indian/Mayotte"
	isoToTz["ZA"] = "Africa/Johannesburg"
	isoToTz["ZM"] = "Africa/Lusaka"
	isoToTz["ZW"] = "Africa/Harare"

	for k, v := range isoToTz {
		loc, err := time.LoadLocation(v)
		if err != nil {
			log.Fatal(err)
		}
		locs[k] = loc
	}

	insertVisits()
}

func insertVisits() {
	for i := 0; i < 10; i++ {
		go func() {
			for {
				visit := <-visitsCh

				err := queries.InsertVisit(visit)
				if err != nil {
					spew.Dump(err)
					spew.Dump("dropped visit")
				}
			}
		}()
	}
}

func Tracking(reader *geoip2.Reader) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		err := func() (err error) {
			if strings.Contains(r.URL.Path, routes.PUBLIC) ||
				strings.Contains(r.URL.Path, routes.AUTHENTICATED+routes.ME) ||
				strings.Contains(r.URL.Path, routes.AUTHENTICATED+routes.SUBMIT) ||
				strings.Contains(r.URL.Path, routes.AUTHENTICATED+routes.COMMENT) {
				return
			}

			ipString := realip.FromRequest(r)
			ip := net.ParseIP(ipString)

			city, err := reader.City(ip)
			if err != nil {
				return
			}

			loc, ok := locs[city.Country.IsoCode]
			if !ok {
				spew.Dump(city.Country.IsoCode)
				spew.Dump("did not find location")
				return
			}

			visit := types.Visit{}
			visit.City = city.City.Names["en"]
			visit.CountryCode = city.Country.IsoCode
			visit.ContinentCode = city.Continent.Code
			visit.Path = r.URL.Path
			visit.CreatedAt = time.Now().UTC()

			var buffer bytes.Buffer
			buffer.Write(salt)
			buffer.WriteString(ipString)
			buffer.WriteString(time.Now().In(loc).Format("Jan 02 2006"))

			h := sha256.New()
			h.Write(buffer.Bytes())

			visit.VisitorId = base64.URLEncoding.EncodeToString(h.Sum(nil))

			if len(visitsCh) < chLen-10 {
				visitsCh <- visit
			}

			return
		}()
		if err != nil {
			return
		}

		next(w, r)
	}
}
