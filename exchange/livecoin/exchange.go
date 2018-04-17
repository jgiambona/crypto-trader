package livecoin

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ffimnsr/trader/exchange"
	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/labstack/echo"
)

// Base API URL.
const (
	LiveCoinAPIURL = "https://api.livecoin.net"
)

// The API error codes that are being returned.
const (
	UnknownError             = 1
	SystemError              = 2
	AuthenticationError      = 10
	AuthenticationIsRequired = 11
	AuthenticationFailed     = 12
	SignatureIncorrect       = 20
	AccessDenied             = 30
	APIDisabled              = 31
	APIRestrictedByIP        = 32
	IncorrectParameters      = 100
	IncorrectAPIKey          = 101
	IncorrectUserID          = 102
	IncorrectCurrency        = 103
	IncorrectAmount          = 104
	UnableToBlockFunds       = 150
)

// Currencies that are allowed on this exchange.
const (
	CurrencyPairAllowed     = "BTC/USD BTC/EUR BTC/RUR LTC/BTC LTC/USD EMC/BTC EMC/USD EMC/RUR EMC/ETH EMC/DASH EMC/XMR DASH/BTC DASH/USD ETH/BTC ETH/USD ETH/RUR DOGE/BTC DOGE/USD CURE/BTC SIB/BTC SIB/RUR TX/BTC RBIES/BTC ADZ/BTC ADZ/USD BSD/BTC SXC/BTC BTA/BTC VOX/BTC MOJO/BTC CRBIT/BTC CRBIT/LEO CRBIT/ETH SHIFT/BTC SHIFT/USD SHIFT/ETH YOC/BTC YOC/USD YOC/RUR YOC/ETH CREVA/BTC LSK/BTC LSK/USD EL/BTC EL/USD EL/RUR HNC/BTC HNC/ETH HNC/USD HNC/EUR CLOAK/BTC CLOAK/USD MOIN/BTC BLU/BTC BLU/USD LEO/BTC LEO/USD LEO/RUR LEO/ETH PPC/BTC PPC/USD NMC/BTC MONA/BTC REE/BTC REE/USD REE/ETH GAME/BTC BLK/BTC SYS/BTC DGB/BTC THS/BTC THS/USD THS/RUR THS/ETH VRC/BTC SLR/BTC DBIX/BTC XMR/BTC XMR/USD BTS/BTC GB/BTC VRM/BTC ATX/BTC ENT/BTC BURST/BTC NXT/BTC POST/BTC POST/ETH EDR/BTC EDR/USD EDR/RUR KRB/BTC KRB/USD KRB/RUR ARC/BTC ARC/USD ARC/RUR ARC/ETH ARC/BCH GYC/BTC DMC/BTC DMC/USD VRS/BTC VRS/USD XRC/BTC XRC/USD BIT/BTC DOLLAR/BTC OD/BTC XAUR/BTC GOLOS/BTC UNC/BTC VLTC/BTC CCRB/BTC CCRB/ETH BPC/BTC BPC/ETH EUR/USD USD/RUR PRES/BTC DIME/BTC DIME/USD DIME/EUR DIME/RUR DIME/ETH ZBC/BTC ZBC/USD ZBC/EUR ZBC/RUR ZBC/ETH ZBC/XMR ZBC/DANC DIBC/BTC DIBC/USD DIBC/ETH VSL/BTC ICN/BTC NVC/BTC NVC/USD XSPEC/BTC LUNA/BTC ACN/BTC LDC/BTC MSCN/BTC MSCN/USD MSCN/EUR MSCN/RUR MSCN/ETH POSW/BTC POSW/USD POSW/EUR POSW/RUR POSW/ETH POSW/DASH POSW/XMR POSW/LTC OBITS/BTC OBITS/USD OBITS/ETH TIME/BTC TIME/ETH TIME/USD WAVES/BTC DANC/BTC DANC/USD INCNT/BTC TAAS/BTC TAAS/USD XMS/BTC SOAR/BTC SOAR/ETH PIVX/BTC PIVX/USD PIVX/EUR PIVX/RUR PIVX/ETH FUNC/BTC FUNC/USD FUNC/ETH ITI/BTC ITI/EUR ITI/RUR ITI/ETH PUT/BTC PUT/USD PUT/RUR PUT/ETH GUP/BTC GUP/ETH MNE/BTC MNE/ETH WINGS/BTC WINGS/ETH UNRC/BTC UNRC/USD RLT/BTC RLT/USD RLT/RUR RLT/ETH FORTYTWO/BTC FORTYTWO/USD FORTYTWO/ETH STRAT/BTC STRAT/USD STRAT/ETH INSN/BTC INSN/USD INSN/ETH QAU/BTC QAU/USD QAU/ETH TRUMP/BTC TRUMP/ETH FNC/BTC FNC/ETH FNC/USD MCO/BTC MCO/USD MCO/ETH VOISE/BTC VOISE/USD VOISE/EUR VOISE/ETH PPY/BTC PPY/USD PPY/ETH ASAFE2/BTC ASAFE2/USD PLBT/BTC PLBT/EMC PLBT/ETH PLBT/USD KPL/BTC KPL/USD KPL/ETH BCH/BTC BCH/USD BCH/RUR BCH/ETH BCH/ZBC MCR/BTC MCR/ETH PIPL/BTC PIPL/ETH HVN/BTC HVN/USD HVN/ETH XRL/BTC XRL/ETH MGO/BTC MGO/ETH FU/BTC FU/ETH WIC/BTC WIC/ETH CTR/BTC CTR/USD CTR/ETH GRS/BTC GRS/USD GRS/ETH PRO/BTC PRO/USD PRO/ETH XEM/BTC XEM/USD XEM/ETH CPC/BTC CPC/USD CPC/ETH wETT/BTC wETT/USD wETT/ETH eETT/BTC eETT/USD eETT/ETH SUMO/BTC SUMO/ETH QTUM/BTC QTUM/USD QTUM/ETH OMG/BTC OMG/USD OMG/ETH PAY/BTC PAY/USD PAY/ETH KNC/BTC KNC/USD KNC/ETH GNT/BTC GNT/USD GNT/ETH EOS/BTC EOS/USD EOS/ETH BAT/BTC BAT/USD BAT/ETH REP/BTC REP/USD REP/ETH MTL/BTC MTL/USD MTL/ETH DGD/BTC DGD/USD DGD/ETH CVC/BTC CVC/USD CVC/ETH SNGLS/BTC SNGLS/USD SNGLS/ETH SNT/BTC SNT/USD SNT/ETH GNO/BTC GNO/USD GNO/ETH ZRX/BTC ZRX/USD ZRX/ETH BNT/BTC BNT/USD BNT/ETH FUN/BTC FUN/USD FUN/ETH EDG/BTC EDG/USD EDG/ETH ANT/BTC ANT/USD ANT/ETH ETHOS/BTC ETHOS/USD ETHOS/ETH STORJ/BTC STORJ/USD STORJ/ETH RLC/BTC RLC/USD RLC/ETH TKN/BTC TKN/USD TKN/ETH MLN/BTC MLN/USD MLN/ETH TRST/BTC TRST/USD TRST/ETH FirstBlood/BTC FirstBlood/USD FirstBlood/ETH VIB/BTC VIB/USD VIB/ETH MTCoin/BTC MTCoin/ETH BIO/BTC BIO/ETH BIO/RUR BIO/USD NEO/BTC NEO/USD NEO/ETH OTN/BTC OTN/USD OTN/ETH MNX/BTC MNX/USD MNX/ETH ENJ/BTC ENJ/ETH DAY/BTC DAY/ETH ETHP/BTC ETHP/ETH ATM/BTC ATM/ETH DMD/BTC DMD/USD DMD/ETH OXY/BTC OXY/USD OXY/ETH CLD/BTC CLD/ETH ARTE/BTC ARTE/ETH CDX/BTC CDX/ETH CLPC/BTC CLPC/USD CLPC/ETH ESP/BTC ESP/ETH BTB/BTC BTB/ETH ESC/BTC ESC/ETH PRG/BTC PRG/USD PRG/ETH AMM/BTC AMM/USD AMM/ETH HST/BTC HST/USD HST/ETH ERO/BTC ERO/ETH KICK/BTC KICK/USD KICK/RUR KICK/ETH UQC/BTC UQC/USD UQC/ETH GRX/BTC GRX/ETH INS/BTC INS/ETH ICOS/BTC ICOS/USD ICOS/ETH TER/BTC TER/ETH FLP/BTC FLP/ETH RBM/BTC RBM/USD RBM/EUR RBM/ETH RBM/LTC FLIXX/BTC FLIXX/ETH DTR/BTC DTR/ETH EVC/BTC EVC/ETH SPF/BTC SPF/ETH B2B/BTC B2B/ETH TFL/BTC TFL/ETH CHSB/BTC CHSB/ETH PIN/BTC PIN/ETH GOAL/BTC GOAL/ETH GOAL/USD COV/BTC COV/ETH IFAN/BTC IFAN/ETH IPL/BTC IPL/ETH NOX/BTC NOX/ETH TWC/BTC TWC/TER TWC/ETH KAPU/BTC KAPU/ETH TRX/BTC TRX/ETH TRX/USD CRC/BTC CRC/ETH CRC/USD CRC/EUR ESR/BTC ESR/ETH ESR/BCH ESR/DASH IPBC/BTC IPBC/ETH IPBC/USD IPBC/XMR SPA/BTC SPA/ETH LTT/BTC LTT/ETH AMB/BTC AMB/ETH AMB/USD AMB/RUR ECHO/BTC ECHO/ETH ECHO/USD HNR/BTC HNR/ETH HNR/USD NAM/BTC NAM/ETH NAM/USD ORE/BTC ORE/ETH ORE/USD PPT/BTC PPT/ETH PPT/USD VEN/BTC VEN/ETH VEN/USD DIG/BTC DIG/ETH DIG/USD DIG/LTC ARK/BTC ARK/ETH ARK/USD ECIO/BTC ECIO/ETH XSN/BTC XSN/ETH XSN/LTC BPTN/BTC BPTN/ETH IDH/BTC IDH/ETH XBT/BTC XBT/ETH USC/BTC USC/ETH FXT/BTC FXT/ETH CBR/BTC CBR/ETH CBR/USD CBR/EUR VIEW/BTC VIEW/ETH"
	CurrencyAllowed         = "BTC LTC EMC DASH DOGE MONA PPC NMC CURE ETH SIB TX RBIES ADZ MOJO BSD SXC BTA VOX CRBIT SHIFT YOC CREVA LSK EL HNC CLOAK MOIN BLU LEO REE GAME BLK SYS DGB THS VRC SLR DBIX SCN XMR BTS GB VRM ATX ENT BURST NXT POST EDR KRB ARC GYC DMC VRS XRC BIT OD XAUR GOLOS DOLLAR UNC VLTC CCRB BPC VSL DIME PRES ICN DIBC NVC XSPEC LUNA ZBC ACN LDC MSCN POSW OBITS TIME WAVES DANC INCNT TAAS XMS SOAR PIVX FUNC WINGS GUP MNE PUT ITI UNRC RLT STRAT FORTYTWO INSN QAU TRUMP FNC MCO VOISE PPY ASAFE2 PLBT KPL BCH MCR FU WIC PIPL HVN XRL MGO CTR GRS PRO XEM CPC wETT eETT SUMO VIB QTUM OMG PAY KNC GNT EOS BAT REP MTL DGD CVC SNGLS SNT GNO ZRX BNT FUN EDG ANT ETHOS STORJ RLC TKN MLN TRST FirstBlood MTCoin NEO BIO OTN MNX ENJ DAY ETHP ATM OXY DMD CLD ARTE CDX CLPC BTB ESP ESC PRG AMM HST ERO KICK UQC GRX INS ICOS TER FLP RBM FLIXX DTR EVC SPF B2B TFL CHSB PIN GOAL COV IFAN IPL NOX TWC KAPU TRX CRC ESR IPBC SPA LTT AMB ECHO HNR NAM ORE PPT VEN DIG ARK ECIO XSN BPTN IDH XBT USC FXT VIEW CBR"
	CurrencyAllowedFiatOnly = "USD EUR RUR"
)

type (
	// LiveCoin interfaces the LiveCoin Rest API.
	LiveCoin struct {
		Store influx.Client
		exchange.Base
	}

	// TickerResponse stores the pricing information.
	TickerResponse struct {
		Currency     string  `json:"cur"`
		CurrencyPair string  `json:"symbol"`
		Last         float64 `json:"last"`
		High         float64 `json:"high"`
		Low          float64 `json:"low"`
		Volume       float64 `json:"volume"`
		Vwap         float64 `json:"vwap"`
		MaxBid       float64 `json:"max_bid"`
		MinAsk       float64 `json:"min_ask"`
		BestBid      float64 `json:"best_bid"`
		BestAsk      float64 `json:"best_ask"`
	}
)

// NewInstance returns the created exchange instance.
func NewInstance() *LiveCoin {
	x := new(LiveCoin)
	x.Name = "LiveCoin"
	x.Enabled = true
	if len(os.Getenv("T_PROD")) > 0 {
		x.Store, _ = influx.NewHTTPClient(influx.HTTPConfig{
			Addr: "http://ec2-54-169-102-171.ap-southeast-1.compute.amazonaws.com:8086",
		})
	} else {
		x.Store, _ = influx.NewHTTPClient(influx.HTTPConfig{
			Addr: "http://localhost:8086",
		})
	}
	x.AvailableCurrencyPairs = strings.Fields(CurrencyPairAllowed)
	x.BaseCurrencies = strings.Fields(CurrencyAllowed)

	return x
}

// GetFee returns the current fee for the exchange.
func (e *LiveCoin) GetFee(maker bool) float64 {
	if maker {
		return 0.18 / 100
	}
	return 0.18 / 100
}

// UpdateTicker updates and returns ticker for a currency pair.
func (e *LiveCoin) UpdateTicker() echo.Map {
	p, err := e.GetTicker("BTC/USD")
	if err != nil {
		log.Println(err.Error())
	}

	bp, err := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database:  "trader",
		Precision: "s",
	})
	if err != nil {
		log.Println(err.Error())
	}

	tags := map[string]string{
		"type":     "ticker",
		"pair":     "btc_usd",
		"exchange": "livecoin",
	}
	fields := echo.Map{
		"symbol":        p.Currency,
		"high":          p.High,
		"low":           p.Low,
		"volume":        p.Volume,
		"ask":           p.BestAsk,
		"askVolume":     -1,
		"bid":           p.BestBid,
		"bidVolume":     -1,
		"vwap":          p.Vwap,
		"open":          -1,
		"close":         p.Last,
		"previousClose": -1,
		"change":        -1,
		"percentage":    -1,
		"average":       -1,
		"baseVolume":    p.Volume,
		"quoteVolume":   p.Volume * p.Vwap,
	}

	pt, err := influx.NewPoint("stream", tags, fields, time.Now())
	bp.AddPoint(pt)
	err = e.Store.Write(bp)
	if err != nil {
		log.Println(err.Error())
	}

	return fields
}

// According to the API Examples of LiveCoin:
// Signature is a HMAC-SHA256 encoded message. The HMAC-SHA256
// code is generated using a secret key that was generated
// with your API key. Generated signatures must be converted
// into hexadecimal format and uppercase characters
func createSignature(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	d := hex.EncodeToString(h.Sum(nil))
	return strings.ToUpper(d)
}
