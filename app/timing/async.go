package timing

import (
	"asset-management/app/dao"
	"log"
)

const (
	GET_ASYNC_TASK_FAILED = "Something wrong happened when get async task in pending"
	// NO_PENDING_TASK = "No tasks in pending"
)

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

	} else if asyncTask.Type == 1 {

	} else {

	}
}
