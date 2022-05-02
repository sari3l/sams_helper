package sams

import (
	"encoding/json"
	"github.com/tidwall/gjson"
)

type PageContent struct {
	PageContentVO    gjson.Result `json:"pageContentVO"`
	PageModuleVOList gjson.Result `json:"pageModuleVOList"`
}

func parsePageContent(result gjson.Result) (error, PageContent) {
	pageContent := PageContent{}
	pageContent.PageContentVO = result.Get("data.pageContentVO")
	pageContent.PageModuleVOList = result.Get("data.pageModuleVOList")
	return nil, pageContent
}

func (session *Session) GetPageMoreData(pageId string, moduleId string) (error, PageContent) {
	data := GetPageDataMoreParam{
		Uid:           session.Uid,
		PageContentId: pageId,
		PageModuleId:  moduleId,
		ModuleDataNum: 2,
		ApiVersion:    1,
		UseNew:        true,
		AddressInfo:   session.Address,
	}
	dataStr, _ := json.Marshal(data)
	err, result := session.Request.POST(GetPageMoreDataAPI, dataStr)
	if err != nil {
		return err, PageContent{}
	}
	return parsePageContent(result)
}

func (session *Session) GetPageData(pageId string) (error, PageContent) {
	data := GetPageDataParam{
		Uid:           session.Uid,
		PageContentId: pageId,
		Authorize:     true,
		Latitude:      session.Address.Latitude,
		Longitude:     session.Address.Longitude,
		AddressInfo:   session.Address,
	}
	dataStr, _ := json.Marshal(data)
	err, result := session.Request.POST(GetPageDataAPI, dataStr)
	if err != nil {
		return err, PageContent{}
	}
	return parsePageContent(result)
}
