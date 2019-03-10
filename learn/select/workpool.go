package main

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

type (
	Worker interface {
		Heath() bool
		Close() error
	}

	Workshop struct {
		addFn           func() (*WorkerInfo, error)
		maxQuota        int
		maxIdleDuiation time.Duration
		infos           map[Worker]*WorkerInfo
		minLoadInfo     *WorkerInfo
		stats           *WorkshopStats
		statsReader     atomic.Value
		lock            sync.Mutex
		wg              sync.WaitGroup
		closeCh         chan struct{}
		closeLock       sync.Mutex
	}

	WorkerInfo struct {
		worker     Worker
		jobNum     int32
		idleExpire time.Time
	}

	WorkshopStats struct {
		Worker  int32
		Idle    int32
		Created uint64
		Doing   int32
		Done    uint64
		Maxload int32
		Minload int32
	}
)

func (w *Workshop) reportStatusLocked() {
	w.statsReader.Store(*w.stats)
}

func (w *Workshop) gc() {
	for {
		select {
		case <-w.closeCh:
			return
		default:
			time.Sleep(w.maxIdleDuiation)
			w.lock.Lock()
			w.refreshLocked(true)
			w.lock.Unlock()
		}
	}
}

func (w *Workshop) refreshLocked(reportStats bool) {
	var max, min, tmp int32
	min = math.MaxInt32
	var minLoadInfo *WorkerInfo
	for _, info := range w.infos {
		if !w.checkInfoLocked(info) {
			continue
		}
		tmp = info.jobNum
		if tmp > max {
			max = tmp
		}
		if tmp < min {
			min = tmp
		}
		if minLoadInfo != nil && tmp >= minLoadInfo.jobNum {
			continue
		}
		minLoadInfo = info
	}
	if min == math.MaxInt32 {
		min = 0

	}
	w.stats.Minload = min
	w.stats.Maxload = max
	if reportStats {
		w.reportStatusLocked()
	}
	w.minLoadInfo = minLoadInfo
}

const (
	defaultWorkerMaxQuota     = 64              // 最大工人数
	defaultWorkerIdleDuration = 3 * time.Minute // 最大默认时长
)

var (
	ErrWorkshopClosed = fmt.Errorf("%s", "workshop is closed")
)

func NewWorkShop(maxQuota int, MaxIdleDuration time.Duration,
	newWorkerFunc func() (Worker, error)) *Workshop {

	if maxQuota <= 0 {
		maxQuota = defaultWorkerMaxQuota
	}
	if MaxIdleDuration <= 0 {
		MaxIdleDuration = defaultWorkerIdleDuration
	}
	w := new(Workshop)
	w.stats = new(WorkshopStats)
	w.reportStatusLocked()
	w.maxQuota = maxQuota
	w.maxIdleDuiation = MaxIdleDuration
	w.infos = make(map[Worker]*WorkerInfo, maxQuota)
	w.closeCh = make(chan struct{})
	w.addFn = func() (info *WorkerInfo, e error) {
		defer func() {
			if p := recover(); p != nil {
				e = fmt.Errorf("%v", p)
			}
		}()
		worker, err := newWorkerFunc()
		if err != nil {
			return nil, err
		}
		info = &WorkerInfo{
			worker: worker,
		}
		w.infos[worker] = info
		w.stats.Created++
		w.stats.Worker++

		return info, nil
	}

	go w.gc()
	return w
}

func (w *Workshop) Callback(fn func(Worker) error) error {

	select {
	case <-w.closeCh:
		return ErrWorkshopClosed
	default:

	}
	w.lock.Lock()
	info, err := w.hireLocked() // 雇佣一个工人
	w.lock.Unlock()
	if err != nil {
		return err
	}
	worker := info.worker
	defer func() { // 收尾工作
		w.lock.Lock()
		_, ok := w.infos[worker]
		if !ok {
			worker.Close()
		} else {
			w.fireLocked(info) // 解聘
		}
		w.lock.Unlock()
	}()
	return fn(worker)
}

func (w *Workshop) Close() {
	w.closeLock.Lock()
	defer w.closeLock.Unlock()
	select {
	case <-w.closeCh: //是否已经关闭了
		return
	default:
		close(w.closeCh) // 主动关闭
	}

	w.wg.Wait()
	w.lock.Lock()
	defer w.lock.Unlock()
	for _, info := range w.infos {
		info.worker.Close() // 资源池所以的work 关闭
	}
	w.infos = nil
	w.stats.Idle = 0
	w.stats.Worker = 0
	w.refreshLocked(true)

}

func (w *Workshop) Fire(worker Worker) {
	w.lock.Lock()
	info, ok := w.infos[worker]
	if !ok {
		if worker != nil {
			worker.Close()
		}
		w.lock.Unlock()
		return
	}
	w.fireLocked(info)
	w.lock.Unlock()
}

func (w *Workshop) Hire() (Worker, error) {
	select {
	case <-w.closeCh:
		return nil, ErrWorkshopClosed
	default:

	}
	w.lock.Lock()
	info, err := w.hireLocked()
	if err != nil {
		w.lock.Unlock()
		return nil, err
	}
	w.lock.Unlock()
	return info.worker, nil
}

func (w *Workshop) Stats() WorkshopStats {
	return w.statsReader.Load().(WorkshopStats)
}

// 解聘
func (w *Workshop) fireLocked(info *WorkerInfo) {
	{
		w.stats.Doing--
		w.stats.Done++
		info.jobNum--
		w.wg.Add(-1)
	}
	jobNum := info.jobNum
	if jobNum == 0 {
		//info.idleExpire = coarsetime.CeilingTimeNow().Add(w.maxIdleDuiation)
		w.stats.Idle++
	}

	if jobNum+1 >= w.stats.Maxload {
		w.refreshLocked(true)
		return

	}
	if !w.checkInfoLocked(info) {
		if info == w.minLoadInfo {
			w.refreshLocked(true)
		}
		return
	}

	if jobNum < w.stats.Minload {
		w.stats.Minload = jobNum
		w.minLoadInfo = info

	}
	w.reportStatusLocked()
}

// 雇佣
func (w *Workshop) hireLocked() (*WorkerInfo, error) {

	var info *WorkerInfo
GET:
	info = w.minLoadInfo
	if len(w.infos) >= w.maxQuota || (info != nil && info.jobNum == 0) {
		if !w.checkInfoLocked(info) {
			w.refreshLocked(false)
			goto GET
		}
		if info.jobNum != 0 {
			w.stats.Idle--
		}
		info.jobNum++
		w.setMinLoadInofoLocked()
		w.stats.Minload = w.minLoadInfo.jobNum
		if w.stats.Maxload < info.jobNum {
			w.stats.Maxload = info.jobNum
		}
	} else {
		var err error
		info, err = w.addFn()
		if err != nil {
			return nil, err
		}
		info.jobNum = 1
		w.stats.Minload = 1
		if w.stats.Maxload == 0 {
			w.stats.Maxload = 1
		}
		w.minLoadInfo = info
	}

	w.wg.Add(1)
	w.stats.Doing++
	w.reportStatusLocked()
	return info, nil

}

func (w *Workshop) checkInfoLocked(info *WorkerInfo) bool {

	if !info.worker.Heath() { //||
		//(info.jobNum == 0 && coarsetime.FloorTimeNow().After(info.idleExpire)

		delete(w.infos, info.worker)
		info.worker.Close()
		w.stats.Worker--
		if info.jobNum == 0 {
			w.stats.Idle--
		} else {
			w.wg.Add(-int(info.jobNum))
			w.stats.Doing -= info.jobNum
			w.stats.Done += uint64((info.jobNum))
		}
		return false
	}
	return true
}

func (w *Workshop) setMinLoadInofoLocked() {

	if len(w.infos) == 0 {
		w.minLoadInfo = nil
		return
	}
	var minLoadInfo *WorkerInfo
	for _, info := range w.infos {
		if minLoadInfo != nil && info.jobNum >= minLoadInfo.jobNum {
			continue
		}
		minLoadInfo = info
	}
	w.minLoadInfo = minLoadInfo

}
