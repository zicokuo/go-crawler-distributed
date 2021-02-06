package master

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"go-crawler-distributed/global"
	"go-crawler-distributed/internal/crontab/common"
)

/**
* @Author: super
* @Date: 2021-02-06 19:25
* @Description:
**/

func EtcdSaveJob(ctx context.Context, job *common.Job) (oldJob *common.Job, err error) {
	jobKey := common.JOB_SAVE_DIR + job.Name
	jobValue, err := job.MarshalJSON()
	if err != nil{
		return
	}
	putResp, err := global.EtcdKV.Put(ctx, jobKey, string(jobValue), clientv3.WithPrevKV())
	if err != nil{
		return
	}
	if putResp.PrevKv != nil{
		oldJobObj := &common.Job{}
		_ = oldJobObj.UnmarshalJSON(putResp.PrevKv.Value)
		oldJob = oldJobObj
	}
	return
}

func EtcdDeleteJob(ctx context.Context, name string) (oldJob *common.Job, err error) {
	jobKey := common.JOB_SAVE_DIR + name

	delResp, err := global.EtcdKV.Delete(ctx, jobKey, clientv3.WithPrevKV())
	if err != nil{
		return
	}
	if len(delResp.PrevKvs) != 0{
		oldJobObj := &common.Job{}
		_ = oldJobObj.UnmarshalJSON(delResp.PrevKvs[0].Value)
		oldJob = oldJobObj
	}
	return
}

func EtcdListJobs(ctx context.Context) (jobList []*common.Job, err error) {
	dirKey := common.JOB_SAVE_DIR

	getResp, err := global.EtcdKV.Get(ctx, dirKey, clientv3.WithPrefix())
	if err != nil{
		return
	}
	jobList = make([]*common.Job, len(getResp.Kvs))
	for i := 0; i<len(getResp.Kvs); i++{
		job := &common.Job{}
		_ = job.UnmarshalJSON(getResp.Kvs[i].Value)
		jobList = append(jobList, job)
	}
	return
}

func EtcdKillJob(ctx context.Context, name string) (err error) {
	killerKey := common.JOB_KILLER_DIR + name

	leaseResp, err := global.EtcdLease.Grant(ctx, 1)
	if err !=  nil{
		return
	}
	leaseId := leaseResp.ID
	_, err = global.EtcdKV.Put(ctx, killerKey, "", clientv3.WithLease(leaseId))
	return
}