package api

import (
	"asset-management/myerror"
	"asset-management/utils"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	sts20150401 "github.com/alibabacloud-go/sts-20150401/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

type ossApi struct {
}

var OssApi *ossApi

func newOssApi() *ossApi {
	return &ossApi{}
}

func init() {
	OssApi = newOssApi()
}

const (
	roleArn           = "acs:ram::1237783510428187:role/xiongjc21"
	roleSessionName   = "xiongjc21"
	durationSeconds   = 900
	myAccessKeyId     = "LTAI5tCpT5SSksUNe355TY8V"
	myAccessKeySecret = "WgkwIjagXCfiu0ykLmrZu1bcXQswV5"
)

/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */
func (oss *ossApi) createClient(accessKeyId *string, accessKeySecret *string) (_result *sts20150401.Client, _err error) {
	config := &openapi.Config{
		// 必填，您的 AccessKey ID
		AccessKeyId: accessKeyId,
		// 必填，您的 AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String("sts.cn-qingdao.aliyuncs.com")
	_result = &sts20150401.Client{}
	_result, _err = sts20150401.NewClient(config)
	return _result, _err
}

/*
Handle func for GET /oss/key
*/
func (oss *ossApi) GetTempKey(ctx *utils.Context) {
	client, err := oss.createClient(tea.String("LTAI5tCpT5SSksUNe355TY8V"), tea.String("WgkwIjagXCfiu0ykLmrZu1bcXQswV5"))

	if err != nil {
		ctx.BadRequest(myerror.OSS_REQUEST_FAILED, myerror.OSS_REQUEST_FAILED_INFO)
		return
	}

	assumeRoleRequest := &sts20150401.AssumeRoleRequest{
		DurationSeconds: tea.Int64(durationSeconds),
		RoleArn:         tea.String(roleArn),
		RoleSessionName: tea.String(roleSessionName),
	}
	runtime := &util.RuntimeOptions{}

	resp, err := client.AssumeRoleWithOptions(assumeRoleRequest, runtime)
	if err != nil {
		ctx.BadRequest(myerror.OSS_REQUEST_FAILED, myerror.OSS_REQUEST_FAILED_INFO)
		return
	}

	ctx.Success(resp.Body.Credentials)
}
