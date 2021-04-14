package services

import (
	"runtime"
	"sync"
)

type ParallelResult = struct {
	r      interface{}
	taskId int
}
type ParallelTask = func(taskId int) (*ParallelResult, SkipResult)
type SkipResult = bool

func DoParallel(taskIds []int, task ParallelTask) []*ParallelResult {
	maxRoutines := runtime.NumCPU()
	availableRoutines := make([]int, 0)
	availableRoutinesMut := sync.Mutex{}
	res := make([]*ParallelResult, 0)
	resMut := sync.Mutex{}
	freeRoutine := make(chan *int)
	for routineId := 0; routineId < maxRoutines; routineId++ {
		availableRoutines = append(availableRoutines, routineId)
	}
	for i := 0; i < len(taskIds); i++ {
		var routineIdx *int
		availableRoutinesMut.Lock()
		if len(availableRoutines) > 0 {
			routineIdx = &availableRoutines[0]
			availableRoutines = availableRoutines[1:]
		}
		availableRoutinesMut.Unlock()
		if routineIdx == nil {
			routineIdx = <-freeRoutine
		}
		go func(taskId int, routIdx int) {
			if r, skip := task(taskIds[taskId]); !skip {
				resMut.Lock()
				res = append(res, r)
				resMut.Unlock()
			}
			freeRoutine <- &routIdx
		}(i, *routineIdx)
	}
	for {
		if len(availableRoutines) == maxRoutines {
			break
		} else {
			// being here, some must run and will wait to read from freeRoutine
			freeRoutineIdx := <-freeRoutine
			availableRoutinesMut.Lock()
			availableRoutines = append(availableRoutines, *freeRoutineIdx)
			availableRoutinesMut.Unlock()
		}
	}
	return res
}
