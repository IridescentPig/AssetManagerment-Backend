package timing

import (
	"asset-management/app/dao"
	"asset-management/app/model"
	"asset-management/app/service"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/shopspring/decimal"
	"github.com/xuri/excelize/v2"
	"gorm.io/datatypes"
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
	CREATE_LOG_EXCEL_FILE_ERROR            = "fail to create log excel file"
	UPLOAD_TO_OSS_FAILED                   = "fail to upload log file to oss"
	EXPORT_ASSET_SUCCESS                   = "Successfully export logs!"
	ASYNC_TASK_SUCCESS                     = "Successfully finish async task!"
)

const (
	endpoint             = "https://oss-cn-beijing.aliyuncs.com"
	downloadLinkPrefix   = "https://export-bucket-1.oss-cn-beijing.aliyuncs.com/"
	myAccessKeyId        = "LTAI5tCpT5SSksUNe355TY8V"
	myAccessKeySecret    = "WgkwIjagXCfiu0ykLmrZu1bcXQswV5"
	importBucketName     = "import-bucket"
	exportBucketName     = "export-bucket-1"
	exportSheetName      = "log"
	TEMP_LOG_FILE_FORMAT = "/var/tmp/export_log_%d.xlsx"
	OSS_LOG_FILE_FORMAT  = "logs/log_%d_%s.xlsx"
	TIME_FORMAT          = "2006-01-02_15-04-05"
)

var (
	fieldList           = []string{"Name", "Price", "ClassID", "Type", "Count", "Expire", "Threshold", "Description", "Position"}
	exportFileFiledList = []string{"ID", "Method", "URL", "Status", "ErrorCode", "ErrorMessage", "UserID", "Username", "Time"}
	cellIndexFormatList = []string{"A%d", "B%d", "C%d", "D%d", "E%d", "F%d", "G%d", "H%d", "I%d"}
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
			if err != nil {
				log.Println(err.Error())
				return
			}
		} else {
			err = dao.AsyncDao.ModifyAsyncTaskInfo(asyncTask.ID, map[string]interface{}{
				"state":   2,
				"message": IMPORT_ASSET_SUCCESS,
			})
			if err != nil {
				log.Println(err.Error())
				return
			}
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
			if err != nil {
				log.Println(err.Error())
				return
			}
		} else {
			err = dao.AsyncDao.ModifyAsyncTaskInfo(asyncTask.ID, map[string]interface{}{
				"state":         2,
				"message":       EXPORT_ASSET_SUCCESS,
				"download_link": asyncTask.DownloadLink,
			})
			if err != nil {
				log.Println(err.Error())
				return
			}
		}
	}
	log.Println(ASYNC_TASK_SUCCESS)
}

func ExportLogs(task *model.AsyncTask) error {
	var logList []*model.Log
	var err error
	if task.Type == 1 {
		logList, err = dao.LogDao.GetLoginLogsForExport(task.EntityID, task.FromTime, task.LogType)
		if err != nil {
			return errors.New(GET_LOG_FAILED)
		}
	} else {
		logList, err = dao.LogDao.GetDataLogsForExport(task.EntityID, task.FromTime, task.LogType)
		if err != nil {
			return errors.New(GET_LOG_FAILED)
		}
	}

	exportFile := excelize.NewFile()
	defer func() {
		if err := exportFile.Close(); err != nil {
			log.Println(err.Error())
		}
	}()

	sheetID, err := exportFile.NewSheet(exportSheetName)
	if err != nil {
		return errors.New(CREATE_LOG_EXCEL_FILE_ERROR)
	}

	exportFile.SetActiveSheet(sheetID)
	// for k, v := range exportFileFiledList {
	// 	_ = exportFile.SetCellStr(exportSheetName, fmt.Sprintf(cellIndexFormatList[k], 1), v)
	// }
	exportFile.SetSheetRow(exportSheetName, fmt.Sprintf(cellIndexFormatList[0], 1), &exportFileFiledList)

	for k, thisLog := range logList {
		// // Set ID
		// _ = exportFile.SetCellInt(exportSheetName, fmt.Sprintf(cellIndexFormatList[0], k+2), int(thisLog.ID))
		// // Set Method
		// _ = exportFile.SetCellStr(exportSheetName, fmt.Sprintf(cellIndexFormatList[1], k+2), thisLog.Method)
		// // Set URL
		// _ = exportFile.SetCellStr(exportSheetName, fmt.Sprintf(cellIndexFormatList[2], k+2), thisLog.URL)
		// // Set Status
		// _ = exportFile.SetCellInt(exportSheetName, fmt.Sprintf(cellIndexFormatList[3], k+2), thisLog.Status)
		// // Set ErrorCodr
		// _ = exportFile.SetCellInt(exportSheetName, fmt.Sprintf(cellIndexFormatList[4], k+2), thisLog.ErrorCode)
		// // Set ErrorMessage
		// _ = exportFile.SetCellStr(exportSheetName, fmt.Sprintf(cellIndexFormatList[5], k+2), thisLog.ErrorMessage)
		// // Set UserID
		// _ = exportFile.SetCellInt(exportSheetName, fmt.Sprintf(cellIndexFormatList[6], k+2), int(thisLog.UserID))
		// // Set Username
		// _ = exportFile.SetCellStr(exportSheetName, fmt.Sprintf(cellIndexFormatList[7], k+2), thisLog.Username)
		// // Set Time
		// _ = exportFile.SetCellStr(exportSheetName, fmt.Sprintf(cellIndexFormatList[8], k+2), thisLog.Time.String())

		_ = exportFile.SetSheetRow(exportSheetName,
			fmt.Sprintf(cellIndexFormatList[0], k+2),
			&[]interface{}{int(thisLog.ID), thisLog.Method, thisLog.URL, thisLog.Status, thisLog.ErrorCode, thisLog.ErrorMessage,
				int(thisLog.UserID), thisLog.Username, thisLog.Time.String()})
	}

	tempLogFilePath := fmt.Sprintf(TEMP_LOG_FILE_FORMAT, task.ID)
	err = exportFile.SaveAs(tempLogFilePath)
	defer func() {
		isExist := false
		_, err := os.Stat(tempLogFilePath)
		if err == nil {
			isExist = true
		} else {
			isExist = !os.IsNotExist(err)
		}

		if isExist {
			err := os.Remove(tempLogFilePath)
			if err != nil {
				log.Println(err.Error())
				return
			}
		}
	}()
	if err != nil {
		return errors.New(CREATE_LOG_EXCEL_FILE_ERROR)
	}

	exportBucket, err := getExportBucket()
	if err != nil {
		return errors.New(ACCESS_TO_OSS_FAILED)
	}

	objectKey := fmt.Sprintf(OSS_LOG_FILE_FORMAT, task.ID, time.Now().Format(TIME_FORMAT))
	err = exportBucket.PutObjectFromFile(objectKey, tempLogFilePath)
	if err != nil {
		return errors.New(UPLOAD_TO_OSS_FAILED)
	}

	task.DownloadLink = downloadLinkPrefix + objectKey
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
	// if len(row) != len(fieldList) {
	// 	log.Println(row)
	// 	log.Println(len(row))
	// 	log.Println(len(fieldList))
	// 	return nil, errors.New(FILE_FORMAT_ERROR)
	// }
	assetName := row[0]
	assetPrice, _ := decimal.NewFromString(row[1])
	assetClassID, _ := strconv.ParseUint(row[2], 10, 64)
	assetType := 1
	assetCount := 1
	assetExpire, _ := strconv.ParseUint(row[5], 10, 64)
	assetThreshold, _ := strconv.ParseUint(row[6], 10, 64)
	assetDescription := ""
	if len(row) >= 8 {
		assetDescription = row[7]
	}
	assetPosition := ""
	if len(row) >= 9 {
		assetPosition = row[8]
	}
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
		Property:     datatypes.JSON([]byte(`{}`)),
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
