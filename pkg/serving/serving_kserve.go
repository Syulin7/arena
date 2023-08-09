package serving

import (
	"fmt"

	"github.com/kubeflow/arena/pkg/apis/config"
	"github.com/kubeflow/arena/pkg/apis/types"
	"github.com/kubeflow/arena/pkg/util"
	"github.com/kubeflow/arena/pkg/workflow"
	log "github.com/sirupsen/logrus"
)

// KServeProcesser use the default processer
type KServeProcesser struct {
	*processer
}

func NewKServeProcesser() Processer {
	p := &processer{
		processerType:   types.KServeJob,
		client:          config.GetArenaConfiger().GetClientSet(),
		enable:          true,
		useIstioGateway: true,
	}
	return &KServeProcesser{
		processer: p,
	}
}

func SubmitKServeJob(namespace string, args *types.KServeArgs) (err error) {
	nameWithVersion := fmt.Sprintf("%v-%v", args.Name, args.Version)
	args.Namespace = namespace
	processers := GetAllProcesser()
	processer, ok := processers[args.Type]
	if !ok {
		return fmt.Errorf("not found processer whose type is %v", args.Type)
	}
	jobs, err := processer.GetServingJobs(args.Namespace, args.Name, args.Version)
	if err != nil {
		return err
	}
	if err := ValidateJobsBeforeSubmiting(jobs, args.Name); err != nil {
		return err
	}
	// the master is also considered as a worker
	chart := util.GetChartsFolder() + "/kserve"
	err = workflow.SubmitJob(nameWithVersion, string(types.KServeJob), namespace, args, chart, args.HelmOptions...)
	if err != nil {
		return err
	}
	log.Infof("The Job %s has been submitted successfully", args.Name)
	log.Infof("You can run `arena serve get %s --type %s -n %s` to check the job status", args.Name, args.Type, args.Namespace)
	return nil
}
