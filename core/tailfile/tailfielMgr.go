package tailfile

import (
	"github.com/sirupsen/logrus"
	"logAgent/common"
)

/*
@author RandySun
@create 2021-09-12-23:46
*/

// tailTask 的管理者
type tailTaskMgr struct {
	tailTaskMap      map[string]*tailTask       // 收集日志对象 所有的tailTask任务
	collectEntryList []common.CollectEntry      // 启动时候的日志文件配置,所有配置项
	confChan         chan []common.CollectEntry // etcd日志配置发生了变化的监控, 等待新配置项

}

var (
	ttMgr *tailTaskMgr
)

// 初始化tailTask,为每一个日志文件构造一个单独的tailTask
func Init(allEtcdConf []common.CollectEntry) (err error) {
	// allEtcdConf里面的内容: [{"path":"G:/goproject/go/logAgent/demo/web.log","topic":"web_log"},{"path":"G:/goproject/go/logAgent/demo/randy.log","topic":"randy_log"}]
	ttMgr = &tailTaskMgr{
		tailTaskMap:      make(map[string]*tailTask, 20),
		collectEntryList: allEtcdConf,
		confChan:         make(chan []common.CollectEntry), //  // 初始化新配置的管道 做一个阻塞的channel
	}
	// 打开文件读取日志  针对每一个日志收集项创建一个taskObj
	for _, conf := range allEtcdConf {

		tt := newTailTask(conf.Path, conf.Topic) // 创建一个日志收集任务
		err := tt.Init()                         // 去打开日志文件准备去读
		if err != nil {
			logrus.Errorf("create tailObj for path:%s, failed, err: %v", conf.Path, err)
		}
		logrus.Infof("create atail task for path: %s success", conf.Path)
		ttMgr.tailTaskMap[tt.path] = tt // 把创建这个tailTask任务登记在册,方便后续管理
		// 启一个后台的goroutine去收集日志
		go tt.run()
	}

	go ttMgr.watch() // 在后台等待新的配置过来
	return
}

func (t *tailTaskMgr) watch() {
	for {
		// 派一个小弟等待新配置来
		newConf := <-t.confChan // 取到值说明新的配置来了
		// 新配置来了后应该管理一下我之前启动的那些tailTask

		logrus.Infof("get new conf from etcd , conf: %v, start manage tailTask....", newConf)
		for _, conf := range newConf {
			// 1. 原来已经存在的任务就不用动
			if t.isTailTaskExist(conf) {
				continue
			}
			// 2. 原来没有的我要创建一个新的tailTask任务
			tt := newTailTask(conf.Path, conf.Topic) // 创建一个日志收集任务
			err := tt.Init()                         // 去打开日志文件准备去读
			if err != nil {
				logrus.Errorf("new create tailObj for path:%s, failed, err: %v", conf.Path, err)
				continue
			}
			logrus.Infof("new create atail task for path: %s success", conf.Path)
			t.tailTaskMap[tt.path] = tt // 把创建这个tailTask任务登记在册,方便后续管理
			// 启一个后台的goroutine去收集日志
			go tt.run()

		}
		// 3. 原来有的现在没有的要tailTask任务停掉
		// 3.1 找出tailTaskMap中存在,但是newConf不存在的那些tailTask,把他们都停掉
		for key, task := range t.tailTaskMap {
			var flage bool
			for _, conf := range newConf {
				if key == conf.Path {
					flage = true
				}
			}
			if !flage {
				// 这个tailTask要停掉
				logrus.Infof("the task collect path: %s need to stop", task.path)
				delete(t.tailTaskMap, key) // 从管理类中删除掉
				task.cancel()
			}
		}
	}

}

// 判断tailTaskMap是否存在该项收集,新的配置
func (t *tailTaskMgr) isTailTaskExist(conf common.CollectEntry) bool {
	_, ok := t.tailTaskMap[conf.Path]
	return ok

}

// 将监听的配置发生变化保存到chan中
func SendNewConf(newConf []common.CollectEntry) {
	ttMgr.confChan <- newConf
}
