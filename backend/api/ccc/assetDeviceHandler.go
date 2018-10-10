package handler

import (
	"nyota/backend/model"
	"nyota/backend/model/config"
	card "nyota/backend/model/dashboard"
	"nyota/backend/uicomponent"
	"nyota/backend/utils"
	"encoding/json"
	"goprizm/httputils"
	"net/http"
	"strconv"

	"github.com/gobs/simplejson"
)

const (
	unclassifiedParamFilter = "Unclassified Device"
	genericParamFilter      = "Generic"
	areaChart               = "area"
)

// GetAssetDeviceDashboardPanelLayout - Defines structure of the Dashboard Panel Layout
func GetAssetDeviceDashboardPanelLayout(s *model.SessionContext, w http.ResponseWriter, req *http.Request) {
	httputils.ServeJSON(w, uicomponent.GetAssetDeviceDefaultDahboardlayout())
}

// GetAssetDeviceDashboardPanelDetails - Defines structure of the Dashboard Panel Details
func GetAssetDeviceDashboardPanelDetails(s *model.SessionContext, w http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	id, err := strconv.Atoi(params.Get("id"))
	if err != nil {
		s.Err = &model.AppError{Type: utils.ValidatationError,
			Message: string("ID missing"), Code: http.StatusUnprocessableEntity}
		return
	}
	dashBoardCard := uicomponent.GetAssetDeviceDashboardCard(id, s)
	deviceUpdateDashboardCardData(id, s, req, dashBoardCard.Charts)
	handleNoDataForChartAndCard(dashBoardCard)
	httputils.ServeJSON(w, dashBoardCard)
}

func handleNoDataForChartAndCard(panelCard *card.DashboardCard) {
	hide := true
	for viewIndex := range panelCard.Charts {
		var chartArray []*card.DashboardChart
		for chartIndex := range panelCard.Charts[viewIndex] {
			if panelCard.Charts[viewIndex][chartIndex].Count > 0 {
				chartArray = append(chartArray, panelCard.Charts[viewIndex][chartIndex])
			}
			if panelCard.Charts[viewIndex][chartIndex].Type == areaChart &&
				panelCard.Charts[viewIndex][chartIndex].Count == 1 &&
				panelCard.Charts[viewIndex][chartIndex].Sum == 0 {
				panelCard.Charts[viewIndex][chartIndex].Filter = card.NoFilter()
			}
		}
		if len(chartArray) > 0 {
			hide = false
			panelCard.Charts[viewIndex] = chartArray
		}
	}
	panelCard.Hide = hide
}

// deviceUpdateDashboardCardData - Returns panel card details based on panel id
func deviceUpdateDashboardCardData(cardID int, s *model.SessionContext, req *http.Request,
	charts map[string][]*card.DashboardChart) {
	data := []model.APIDataArray{}
	switch cardID {
	case 1:
		data := getProfileOverview(s, req)
		uicomponent.AppendAPIDataArrayToChart(data, charts[top5Const])
		return
	case 2:
		data = getProfileDeviceCategory(s, req)
		break
	case 3:
		data = getProfileDeviceFamily(s, req)
		break
	case 4:
		data = getProfileDeviceName(s, req)
		break
	case 5:
		data = getProfileDeviceMacVendor(s, req)
		break
	default:
		return
	}
	// all records
	uicomponent.AppendTreeMapAPIFullDataArrayToChart(data[:], charts[allConst])
	// Top 5
	uicomponent.AppendAPIDataArrayToChart(data, charts[top5Const])
}

func getUserClassifiedFilter(view string, s *model.SessionContext, req *http.Request) map[string]string {
	paramMap := utils.PrepareMapWithTenantInfoAndActiveFilter(s, req)
	utils.UpdateDeviceClassificationFilter(view, paramMap)
	return paramMap
}

func getProfileOverviewCounts(s *model.SessionContext, req *http.Request) (int, int, int) {
	//get endpoint count for user classfied device
	paramMap1 := getUserClassifiedFilter("userclassifiedcluster", s, req)
	clusterLabelCountChannel := make(chan *model.HTTPResponse)
	go getEndpointCount(s, paramMap1, clusterLabelCountChannel)

	newMap := getUserClassifiedFilter("userclassifiedrule", s, req)
	ruleCountChannel := make(chan *model.HTTPResponse)
	go getEndpointCount(s, newMap, ruleCountChannel)

	clusterLabelCount := getCount(s, <-clusterLabelCountChannel)
	ruleCount := getCount(s, <-ruleCountChannel)

	paramMap := utils.PrepareMapWithTenantInfoAndActiveFilter(s, req)
	resp, err := getEndpointOverviewDetails(s, paramMap)
	contents, success := utils.ReadHTTPResponse(s, resp, err)
	if success == false {
		return 0, 0, 0
	}

	obj, _ := simplejson.LoadBytes(contents)
	jsonObj := obj.Get("count")

	var total, generic, userClassified int
	total = getCountValueFromJson("count", jsonObj)
	generic = getCountValueFromJson("unknown", jsonObj)
	userClassified = clusterLabelCount + ruleCount

	return total, generic, userClassified
}

func getProfileOverview(s *model.SessionContext, req *http.Request) []model.APIDataArray {

	total, generic, userClassified := getProfileOverviewCounts(s, req)

	var data []model.APIDataArray
	data = append(data, prepareAPIDataArr(total-generic))
	data = append(data, prepareAPIDataArr(total-(generic+userClassified)))
	data = append(data, prepareAPIDataArr(userClassified))
	data = append(data, prepareAPIDataArr(generic))
	data = append(data, prepareAPIDataArr(total))
	return data
}

func getCountValueFromJson(key string, jsonObj *simplejson.Json) int {
	dist, _ := jsonObj.Get(key).Array()
	distLen := len(dist)
	if distLen > 0 {
		dateCount := dist[distLen-1]
		dateCountValue := simplejson.AsJson(dateCount)
		return dateCountValue.Get("count").MustInt()
	}
	return 0
}

func prepareAPIDataArr(count int) model.APIDataArray {
	apiDataArray := make(model.APIDataArray, 0)
	apidata := model.KeySortData(
		0, count)
	apiDataArray = append(apiDataArray, apidata)
	return apiDataArray
}

func getProfileDeviceCategory(s *model.SessionContext, req *http.Request) []model.APIDataArray {
	return []model.APIDataArray{getAssetGroupData(s, req, "device_category")}
}

func getProfileDeviceFamily(s *model.SessionContext, req *http.Request) []model.APIDataArray {
	return []model.APIDataArray{getAssetGroupData(s, req, "device_family")}
}

func getProfileDeviceName(s *model.SessionContext, req *http.Request) []model.APIDataArray {
	return []model.APIDataArray{getAssetGroupData(s, req, "device_name")}
}

func getProfileDeviceMacVendor(s *model.SessionContext, req *http.Request) []model.APIDataArray {
	return []model.APIDataArray{getAssetGroupData(s, req, "mac_vendor")}
}

func getAssetGroupData(s *model.SessionContext, req *http.Request, groupBy string) model.APIDataArray {
	return getAssetGroupDataWithParams(s, req, groupBy, nil)
}

func getAssetGroupDataWithParams(s *model.SessionContext, req *http.Request, groupBy string, addnlParamMap map[string]string) model.APIDataArray {
	paramMap := utils.PrepareMapWithTenantInfoAndActiveFilterForTotalClassified(s, req)
	if nil != addnlParamMap {
		for k, v := range addnlParamMap {
			paramMap[k] = v
		}
	}
	paramMap["group_by"] = groupBy
	return getEndpointGroupData(s, paramMap)
}

func getAssetGroupDataWithParamsAsync(s *model.SessionContext, req *http.Request, groupBy string, addnlParamMap map[string]string, result chan model.APIDataArray) {
	result <- getAssetGroupDataWithParams(s, req, groupBy, addnlParamMap)
}

func GetDevices(s *model.SessionContext, w http.ResponseWriter, req *http.Request) {
	paramMap := utils.PrepareMapWithTenantInfo(s)
	paramMap["global"] = "true"
	resp, err := getAllEntities(s, paramMap, deviceListURL)
	contents, success := utils.ReadAndHandleHTTPResponse(s, resp, err)
	var devices = make([]config.TDeviceWithoutCount, 0)
	if success {
		json.Unmarshal(contents, &devices)
	}
	httputils.ServeJSON(w, uicomponent.GetDevices(devices))
}
