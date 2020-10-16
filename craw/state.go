package craw

import (
	"../net"
	"log"
)

var (
	InitialViewState = ""
	InitialEventvalId = ""

)

func init() {
	doc,err := net.HttpGet("https://www.cdggzy.com/site/LandTrade/LandList.aspx")
	if err != nil {
		log.Fatal(err.Error())
	}

	InitialEventvalId, _ = doc.Find("form#form1 input#__EVENTVALIDATION").Eq(0).Attr("value")
	InitialViewState, _ = doc.Find("form#form1 input#__VIEWSTATE").Eq(0).Attr("value")
}


