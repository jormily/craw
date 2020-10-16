package main

import (
	"./erl"
	"./svr"
)

func main() {
	noticeSvr := svr.NewNoticeServer()
	erl.Register("notice", noticeSvr)

	pageSvr := svr.NewPageServer("https://www.cdggzy.com/site/LandTrade/LandList.aspx")
	pageSvr.SetConsumer(&noticeSvr.BaseServer)
	pageSvr.SetPager(svr.NewNoticePager(pageSvr))
	erl.Register("notice_page", pageSvr)


	priceSvr := svr.NewPriceServer()
	erl.Register("price", priceSvr)

	pageSvr2 := svr.NewPageServer("https://www.cdggzy.com/site/LandTrade/LandList.aspx")
	pageSvr2.SetConsumer(&priceSvr.BaseServer)
	pageSvr2.SetPager(svr.NewPricePager(pageSvr2))
	erl.Register("price_page", pageSvr2)

	erl.Wait()

	priceSvr.Check(noticeSvr.CNoticeCraw)
	priceSvr.ExportExcel(".\\土地拍卖价格.xlsx")
}