package api

import (
	"asset-management/app/define"
	"asset-management/app/service"
	"asset-management/myerror"
	"asset-management/utils"
)

type statApi struct {
}

var StatApi *statApi

func newStatApi() *statApi {
	return &statApi{}
}

func init() {
	StatApi = newStatApi()
}

/*
Handle func for GET /department/:department_id/asset/stat/total
*/
func (stat *statApi) GetDepartmentStatTotal(ctx *utils.Context) {
	hasIdentity, departmentID, err := AssetClassApi.CheckAssetIdentity(ctx)
	if err != nil {
		return
	} else if !hasIdentity {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}

	stats, err := service.StatService.GetDepartmentStat(departmentID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	statRes := define.DepartmentStatTotalResponse{
		Stats: stats,
	}

	ctx.Success(statRes)
}

/*
Handle func for GET /department/:department_id/asset/stat/distribution
*/
func (stat *statApi) GetDepartmentStatDistribution(ctx *utils.Context) {
	hasIdentity, departmentID, err := AssetClassApi.CheckAssetIdentity(ctx)
	if err != nil {
		return
	} else if !hasIdentity {
		ctx.Forbidden(myerror.PERMISSION_DENIED, myerror.PERMISSION_DENIED_INFO)
		return
	}

	distribution, err := service.StatService.GetDepartmentAssetDistribution(departmentID)
	if err != nil {
		ctx.InternalError(err.Error())
		return
	}

	distributionRes := define.AssetDistributionResponse{
		Distribution: distribution,
	}

	ctx.Success(distributionRes)
}
