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
	GET_ASYNC_TASK_FAILED                  = "something wrong happened when get async task in pending"
	ASYNC_TASK_USER_NOT_FOUND              = "async task's launcher isn't exist"
	ASYNC_TASK_LAUNCHER_PERMISSION_DENIED  = "async task's launcher has no right to execute the task"
	ACCESS_TO_OSS_FAILED                   = "get access to oss failed"
	GET_IMPORT_FILE_FAILED                 = "cannot download import file"
	FILE_MUST_BE_XLSX                      = "import file must be a valid excel file"
	READ_FILE_FAILED                       = "failed to read file, please check if file is valid or retry later"
	FILE_FORMAT_ERROR                      = "file format is incorrect"
	FIELD_NAME_ERROR_FORMAT                = "field Name error on cell A%d"
	FIELD_PRICE_ERROR_FORMAT               = "field Price error on cell B%d"
	FIELD_CLASSID_ERROR_FORMAT             = "field ClassID error on cell C%d"
	FIELD_CLASSID_NOT_FOUND_FORMAT         = "value of field ClassID of cell C%d not found among asset classes"
	FIELD_CLASSID_NOT_IN_DEPARTMENT_FORMAT = "value of field ClassID of cell C%d not in your department"
	FIELD_TYPE_ERROR_FORMAT                = "field Type error on cell D%d"
	FIELD_COUNT_ERROR_FORMAT               = "field Count error on cell E%d"
	FIELD_EXPIRE_ERROR_FORMAT              = "field Expire error on cell F%d"
	FIELD_THRESHOLD_ERROR_FORMAT           = "field Threshold error on cell G%d"
	PARSE_ASSET_FAILED_FORMAT              = "parse asset info failed on line %d"
	INSERT_ASSET_FAILED_FORMAT             = "insert asset failed on line %d"
	IMPORT_ASSET_SUCCESS                   = "Successfully import assets!"
	GET_LOG_FAILED                         = "failed to get logs"
)

const (
	endpoint          = "https://oss-cn-beijing.aliyuncs.com"
	myAccessKeyId     = "LTAI5tCpT5SSksUNe355TY8V"
	myAccessKeySecret = "WgkwIjagXCfiu0ykLmrZu1bcXQswV5"
	importBucketName  = "import-bucket"
	exportBucketName  = "export-bucket-1"
)

var (
	fieldList = []string{"Name", "Price", "ClassID", "Type", "Count", "Expire", "Threshold", "Description", "Position"}
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

	thisUser, err := dao.UserDao.GetUserByID(asyncTask.UserID)
	if err != nil {
		log.Println(err.Error())
		return
	} else if thisUser == nil {
		err = dao.AsyncDao.ModifyAsyncTaskInfo(asyncTask.ID, map[string]interface{}{
			"state":   3,
			"message": ASYNC_TASK_USER_NOT_FOUND,
		})
		if err != nil {
			log.Println(err.Error())
		}
		return
	}

	err = dao.AsyncDao.ModifyAsyncTaskInfo(asyncTask.ID, map[string]interface{}{
		"state": 1,
	})
	if err != nil {
		log.Println(err.Error())
		return
	}

	if asyncTask.Type == 0 {
		if !thisUser.DepartmentSuper || thisUser.DepartmentID != asyncTask.DepartmentID {
			err = dao.AsyncDao.ModifyAsyncTaskInfo(asyncTask.ID, map[string]interface{}{
				"state":   3,
				"message": ASYNC_TASK_LAUNCHER_PERMISSION_DENIED,
			})
			if err != nil {
				log.Println(err.Error())
			}
			return
		}

		err = ImportAssets(asyncTask)
		if err != nil {
			err = dao.AsyncDao.ModifyAsyncTaskInfo(asyncTask.ID, map[string]interface{}{
				"state":   3,
				"message": err.Error(),
			})
			log.Println(err.Error())
		} else {
			err = dao.AsyncDao.ModifyAsyncTaskInfo(asyncTask.ID, map[string]interface{}{
				"state":   2,
				"message": IMPORT_ASSET_SUCCESS,
			})
			log.Println(err.Error())
		}
	} else {
		if !thisUser.EntitySuper || thisUser.EntityID != asyncTask.EntityID {
			err = dao.AsyncDao.ModifyAsyncTaskInfo(asyncTask.ID, map[string]interface{}{
				"state":   3,
				"message": ASYNC_TASK_LAUNCHER_PERMISSION_DENIED,
			})
			if err != nil {
				log.Println(err.Error())
			}
			return
		}

		err = ExportLogs(asyncTask)
		if err != nil {
			err = dao.AsyncDao.ModifyAsyncTaskInfo(asyncTask.ID, map[string]interface{}{
				"state":   3,
				"message": err.Error(),
			})
			log.Println(err.Error())
		} else {
			err = dao.AsyncDao.ModifyAsyncTaskInfo(asyncTask.ID, map[string]interface{}{
				"state":   2,
				"message": IMPORT_ASSET_SUCCESS,
			})
			log.Println(err.Error())
		}
	}
}

func ExportLogs(task *model.AsyncTask) error {
	if task.Type == 1 {
		_, err := dao.LogDao.GetLoginLogByEntityIDAndTime(task.EntityID, task.FromTime)
		if err != nil {
			return errors.New(GET_LOG_FAILED)
		}
	}

	_, err := getExportBucket()
	if err != nil {
		return errors.New(ACCESS_TO_OSS_FAILED)
	}
	return nil
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

	defer func() {
		if err = rows.Close(); err != nil {
			log.Println(err)
		}
	}()

	if rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			return errors.New(READ_FILE_FAILED)
		}
		isOk := checkRowsKeyValid(row)
		if !isOk {
			return errors.New(FILE_FORMAT_ERROR)
		}
	}

	cols, err := importFile.Cols("Asset")
	if err != nil {
		return errors.New(READ_FILE_FAILED)
	}

	err = checkColsValueValid(cols, task.DepartmentID)
	if err != nil {
		return err
	}

	lineIndex := 1
	for rows.Next() {
		lineIndex += 1
		row, err := rows.Columns()
		if err != nil {
			return errors.New(READ_FILE_FAILED)
		}
		thisAsset, err := parseRowGetAssetInfo(row, task.DepartmentID, task.UserID)
		if err != nil {
			return fmt.Errorf(PARSE_ASSET_FAILED_FORMAT, lineIndex)
		}
		err = dao.AssetDao.Create(*thisAsset)
		if err != nil {
			return fmt.Errorf(INSERT_ASSET_FAILED_FORMAT, lineIndex)
		}
	}

	return nil
}

func parseRowGetAssetInfo(row []string, departmentID uint, userID uint) (*model.Asset, error) {
	if len(row) != len(fieldList) {
		return nil, errors.New(FILE_FORMAT_ERROR)
	}
	assetName := row[0]
	assetPrice, _ := decimal.NewFromString(row[1])
	assetClassID, _ := strconv.ParseUint(row[2], 10, 64)
	assetType := 1
	assetCount := 1
	assetExpire, _ := strconv.ParseUint(row[5], 10, 64)
	assetThreshold, _ := strconv.ParseUint(row[6], 10, 64)
	assetDescription := row[7]
	assetPosition := row[8]

	return &model.Asset{
		Name:         assetName,
		Price:        assetPrice,
		NetWorth:     assetPrice,
		ClassID:      uint(assetClassID),
		Type:         assetType,
		Number:       assetCount,
		Expire:       uint(assetExpire),
		Threshold:    uint(assetThreshold),
		Description:  assetDescription,
		Position:     assetPosition,
		DepartmentID: departmentID,
		UserID:       userID,
	}, nil
}

func checkColsValueValid(cols *excelize.Cols, departmentID uint) error {
	// Check field Name: cannot be empty
	if !cols.Next() {
		return errors.New(READ_FILE_FAILED)
	}

	col, err := cols.Rows()
	if err != nil {
		return errors.New(READ_FILE_FAILED)
	}

	for k, name := range col[1:] {
		if name == "" {
			return fmt.Errorf(FIELD_NAME_ERROR_FORMAT, k+2)
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
	for k, priceStr := range col[1:] {
		if priceStr == "" {
			return fmt.Errorf(FIELD_PRICE_ERROR_FORMAT, k+2)
		}
		price, err := decimal.NewFromString(priceStr)
		if err != nil {
			return fmt.Errorf(FIELD_PRICE_ERROR_FORMAT, k+2)
		} else if minimalPrice.Cmp(price) == 1 || maxiumPrice.Cmp(price) == -1 {
			return fmt.Errorf(FIELD_PRICE_ERROR_FORMAT, k+2)
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

	for k, classIDStr := range col[1:] {
		if classIDStr == "" {
			return fmt.Errorf(FIELD_CLASSID_ERROR_FORMAT, k+2)
		}
		classID, err := strconv.ParseUint(classIDStr, 10, 64)
		if err != nil {
			return fmt.Errorf(FIELD_CLASSID_ERROR_FORMAT, k+2)
		}

		thisClass, err := service.AssetClassService.GetAssetClassByID(uint(classID))
		if err != nil {
			return fmt.Errorf(FIELD_CLASSID_ERROR_FORMAT, k+2)
		} else if thisClass == nil {
			return fmt.Errorf(FIELD_CLASSID_NOT_FOUND_FORMAT, k+2)
		} else if thisClass.DepartmentID != departmentID {
			return fmt.Errorf(FIELD_CLASSID_NOT_IN_DEPARTMENT_FORMAT, k+2)
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

	for k, assetTypeStr := range col[1:] {
		if assetTypeStr != "1" {
			return fmt.Errorf(FIELD_TYPE_ERROR_FORMAT, k+2)
		}
	}

	if !cols.Next() {
		return errors.New(READ_FILE_FAILED)
	}
	col, err = cols.Rows()
	if err != nil {
		return errors.New(READ_FILE_FAILED)
	}

	for k, assetCountStr := range col[1:] {
		if assetCountStr != "1" {
			return fmt.Errorf(FIELD_COUNT_ERROR_FORMAT, k+2)
		}
	}

	// Check field Expire, Threshold: uint
	if !cols.Next() {
		return errors.New(READ_FILE_FAILED)
	}
	col, err = cols.Rows()
	if err != nil {
		return errors.New(READ_FILE_FAILED)
	}

	for k, expireStr := range col[1:] {
		if expireStr == "" {
			return fmt.Errorf(FIELD_EXPIRE_ERROR_FORMAT, k+2)
		}
		_, err = strconv.ParseUint(expireStr, 10, 64)
		if err != nil {
			return fmt.Errorf(FIELD_EXPIRE_ERROR_FORMAT, k+2)
		}
	}

	if !cols.Next() {
		return errors.New(READ_FILE_FAILED)
	}

	col, err = cols.Rows()
	if err != nil {
		return errors.New(READ_FILE_FAILED)
	}

	for k, thresholdStr := range col[1:] {
		if thresholdStr == "" {
			return fmt.Errorf(FIELD_THRESHOLD_ERROR_FORMAT, k+2)
		}
		_, err = strconv.ParseUint(thresholdStr, 10, 64)
		if err != nil {
			return fmt.Errorf(FIELD_THRESHOLD_ERROR_FORMAT, k+2)
		}
	}

	return nil
}

func checkRowsKeyValid(rowHeaders []string) bool {
	lengthKeys := len(fieldList)
	lengthHeaders := len(rowHeaders)
	if lengthHeaders != lengthKeys {
		return false
	}

	for i := 0; i < lengthKeys; i++ {
		if fieldList[i] != rowHeaders[i] {
			return false
		}
	}

	return true
}
