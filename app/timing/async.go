package timing

import (
	"asset-management/app/dao"
	"asset-management/app/model"
	"asset-management/app/service"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/shopspring/decimal"
	"github.com/xuri/excelize/v2"
)

const (
	GET_ASYNC_TASK_FAILED          = "something wrong happened when get async task in pending"
	ACCESS_TO_OSS_FAILED           = "get access to oss failed"
	GET_IMPORT_FILE_FAILED         = "cannot download import file"
	FILE_MUST_BE_XLSX              = "import file must be a valid excel file"
	READ_FILE_FAILED               = "failed to read file, please check if file is valid or retry later"
	FILE_FORMAT_ERROR              = "file format is incorrect"
	FIELD_NAME_ERROR_FORMAT        = "field Name error on cell A%d"
	FIELD_PRICE_ERROR_FORMAT       = "field Price error on cell B%d"
	FIELD_CLASSID_ERROR_FORMAT     = "field ClassID error on cell C%d"
	FIELD_CLASSID_NOT_FOUND_FORMAT = "field ClassID not found on cell C%d"
)

const (
	endpoint          = "https://oss-cn-beijing.aliyuncs.com"
	myAccessKeyId     = "LTAI5tCpT5SSksUNe355TY8V"
	myAccessKeySecret = "WgkwIjagXCfiu0ykLmrZu1bcXQswV5"
	importBucketName  = "import-bucket"
	exportBucketName  = "export-bucket-1"
)

var (
	keyList = []string{"Name", "Price", "Expire", "ClassID", "Type", "Count", "Threshold", "Description", "Position"}
)

func getClient() (*oss.Client, error) {
	return oss.New(endpoint, myAccessKeyId, myAccessKeySecret)
}

func getImportBucket() (*oss.Bucket, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	return client.Bucket(importBucketName)
}

func getExportBucket() (*oss.Bucket, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}

	return client.Bucket(exportBucketName)
}

type GetPendingAsyncTask struct {
}

func (task *GetPendingAsyncTask) Run() {
	asyncTask, err := dao.AsyncDao.GetPendingTask()
	if err != nil {
		log.Println(GET_ASYNC_TASK_FAILED)
		return
	} else if asyncTask == nil {
		return
	}

	if asyncTask.Type == 0 {
		dao.AsyncDao.ModifyAsyncTaskInfo(asyncTask.ID, map[string]interface{}{
			"state": 1,
		})
		ImportAssets(asyncTask)
		if err != nil {

		}
	} else if asyncTask.Type == 1 {

	} else {

	}
}

func ImportAssets(task *model.AsyncTask) error {
	importBucket, err := getImportBucket()
	if err != nil {
		return errors.New(ACCESS_TO_OSS_FAILED)
	}

	importFileReader, err := importBucket.GetObject(task.ObjectKey)
	if err != nil {
		return errors.New(GET_IMPORT_FILE_FAILED)
	}
	defer func() {
		if err := importFileReader.Close(); err != nil {
			log.Println(err.Error())
		}
	}()

	importFile, err := excelize.OpenReader(importFileReader)
	if err != nil {
		return errors.New(FILE_MUST_BE_XLSX)
	}
	defer func() {
		if err := importFile.Close(); err != nil {
			log.Println(err)
		}
	}()

	rows, err := importFile.Rows("Asset")
	if err != nil {
		return errors.New(READ_FILE_FAILED)
	}
	if rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			return errors.New(READ_FILE_FAILED)
		}
		isOk := checkRowsKeyValid(row)
		if !isOk {
			return errors.New(FILE_FORMAT_ERROR)
		}
		if err = rows.Close(); err != nil {
			log.Println(err)
		}
	} else {
		if err = rows.Close(); err != nil {
			log.Println(err)
		}
		return errors.New(FILE_FORMAT_ERROR)
	}

	_, err = importFile.Cols("Asset")
	if err != nil {
		return errors.New(READ_FILE_FAILED)
	}

	return nil
}

func checkColsValueValid(cols *excelize.Cols) error {
	// Check field Name: cannot be empty
	if !cols.Next() {
		return errors.New(READ_FILE_FAILED)
	}

	col, err := cols.Rows()
	if err != nil {
		return errors.New(READ_FILE_FAILED)
	}

	for k, name := range col {
		if name == "" {
			return fmt.Errorf(FIELD_NAME_ERROR_FORMAT, k+1)
		}
	}

	// Check field Price: must be decimal(10, 2)
	if !cols.Next() {
		return errors.New(READ_FILE_FAILED)
	}

	col, err = cols.Rows()
	if err != nil {
		return errors.New(READ_FILE_FAILED)
	}

	minimalPrice := decimal.NewFromFloat(0)
	maxiumPrice, _ := decimal.NewFromString("99999999.99")
	for k, priceStr := range col {
		if priceStr == "" {
			return fmt.Errorf(FIELD_PRICE_ERROR_FORMAT, k+1)
		}
		price, err := decimal.NewFromString(priceStr)
		if err != nil {
			return fmt.Errorf(FIELD_PRICE_ERROR_FORMAT, k+1)
		} else if minimalPrice.Cmp(price) == 1 || maxiumPrice.Cmp(price) == -1 {
			return fmt.Errorf(FIELD_PRICE_ERROR_FORMAT, k+1)
		}
	}

	// Check field ClassID: uint, exists
	if !cols.Next() {
		return errors.New(READ_FILE_FAILED)
	}

	col, err = cols.Rows()
	if err != nil {
		return errors.New(READ_FILE_FAILED)
	}

	for k, classIDStr := range col {
		if classIDStr == "" {
			return fmt.Errorf(FIELD_CLASSID_ERROR_FORMAT, k+1)
		}
		classID, err := strconv.ParseUint(classIDStr, 10, 64)
		if err != nil {
			return fmt.Errorf(FIELD_CLASSID_ERROR_FORMAT, k+1)
		}

		isExist, err := service.AssetClassService.ExistsAssetClass(uint(classID))
		if err != nil {
			return fmt.Errorf(FIELD_CLASSID_ERROR_FORMAT, k+1)
		} else if !isExist {
			return fmt.Errorf(FIELD_CLASSID_NOT_FOUND_FORMAT, k+1)
		}
	}

	// Check field Type, Count: int, can only value 1 so far
	if !cols.Next() {
		return errors.New(READ_FILE_FAILED)
	}

	col, err = cols.Rows()
	if err != nil {
		return errors.New(READ_FILE_FAILED)
	}

	return nil
}

func checkRowsKeyValid(rowHeaders []string) bool {
	lengthKeys := len(keyList)
	lengthHeaders := len(rowHeaders)
	if lengthHeaders != lengthKeys {
		return false
	}

	for i := 0; i < lengthKeys; i++ {
		if keyList[i] != rowHeaders[i] {
			return false
		}
	}

	return true
}
