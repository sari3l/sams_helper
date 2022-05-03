package sams

func (session *Session) GetGuaranteedSupplyMoreGoods() (error, []ShowGoods) {
	var goodsList = make([]ShowGoods, 0)
	err, result := session.GetPageMoreData("1187641882302384150", "1210005874370846742")
	if err != nil {
		return err, nil
	}
	for _, v := range result.PageModuleVOList.Array() {
		if v.Get("moduleSign").Str == "goodsModule" {
			for _, v2 := range v.Get("renderContent.goodsList").Array() {
				_, goods := parseShowGoods(v2)
				goodsList = append(goodsList, goods)
			}
		}
	}
	return nil, goodsList
}

func (session *Session) GetGuaranteedSupplyGoods() (error, []ShowGoods) {
	var goodsList = make([]ShowGoods, 0)
	err, result := session.GetPageData("1187641882302384150")
	if err != nil {
		return err, nil
	}
	for _, v := range result.PageModuleVOList.Array() {
		if v.Get("moduleSign").Str == "goodsModule" {
			for _, v2 := range v.Get("renderContent.goodsList").Array() {
				_, goods := parseShowGoods(v2)
				goodsList = append(goodsList, goods)
			}
		}
	}
	return nil, goodsList
}

func (session *Session) GetGuaranteedSupplyGoodsAll() (error, []ShowGoods) {
	var goodsList = make([]ShowGoods, 0)
	err, result := session.GetGuaranteedSupplyGoods()
	if err != nil {
		return err, nil
	}
	goodsList = append(goodsList, result...)
	err, resultMore := session.GetGuaranteedSupplyMoreGoods()
	if err != nil {
		return err, nil
	}
	goodsList = append(goodsList, resultMore...)
	return nil, goodsList
}
