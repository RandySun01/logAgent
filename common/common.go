package common

/*
@author RandySun
@create 2021-09-12-20:30
*/

// 要收集的日志配置结构体

type CollectEntry struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}