package serving

import (
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/kubeflow/arena/pkg/apis/types"
	"github.com/kubeflow/arena/pkg/util/kubectl"
)

const (
	ResourceGPU       v1.ResourceName = "nvidia.com/gpu"
	ResourceGPUMemory v1.ResourceName = "aliyun.com/gpu-mem"
	ResourceGPUCore   v1.ResourceName = "aliyun.com/gpu-core.percentage"
)

func UpdateTensorflowServing(args *types.UpdateTensorFlowServingArgs) error {
	deploy, err := findAndBuildDeployment(&args.CommonUpdateServingArgs)
	if err != nil {
		return err
	}

	if args.Command == "" {
		containerArgs := deploy.Spec.Template.Spec.Containers[0].Args
		if len(containerArgs) > 0 {
			servingArgs := containerArgs[0]

			if strings.HasSuffix(servingArgs, "\n") {
				servingArgs = strings.TrimSpace(servingArgs[:len(servingArgs)-1])
			}
			arr := strings.Split(servingArgs, "--")
			params := make(map[string]string)
			for index, argItem := range arr {
				if index == 0 {
					continue
				}
				pair := strings.Split(argItem, "=")
				if len(pair) <= 1 {
					continue
				}
				params[fmt.Sprintf("--%s", pair[0])] = argItem[len(pair[0])+1:]
			}
			if args.ModelName != "" {
				params["--model_name"] = args.ModelName
			}
			if args.ModelPath != "" {
				params["--model_base_path"] = args.ModelPath
			}
			if args.ModelConfigFile != "" {
				params["--model_config_file"] = args.ModelConfigFile
			}
			if args.MonitoringConfigFile != "" {
				params["--monitoring_config_file"] = args.MonitoringConfigFile
			}

			var newArgs []string
			newArgs = append(newArgs, "/usr/bin/tensorflow_model_server")
			for k, v := range params {
				newArgs = append(newArgs, fmt.Sprintf("%s=%s", k, v))
			}

			deploy.Spec.Template.Spec.Containers[0].Args = []string{strings.Join(newArgs, " ")}
		}
	}

	if args.Annotations != nil && len(args.Annotations) > 0 {
		for k, v := range args.Annotations {
			deploy.Annotations[k] = v
			deploy.Spec.Template.Annotations[k] = v
		}
	}

	if args.Labels != nil && len(args.Labels) > 0 {
		for k, v := range args.Labels {
			deploy.Labels[k] = v
			deploy.Spec.Template.Labels[k] = v
		}
	}

	if args.NodeSelectors != nil && len(args.NodeSelectors) > 0 {
		if deploy.Spec.Template.Spec.NodeSelector == nil {
			deploy.Spec.Template.Spec.NodeSelector = map[string]string{}
		}
		for k, v := range args.NodeSelectors {
			deploy.Spec.Template.Spec.NodeSelector[k] = v
		}
	}

	if args.Tolerations != nil && len(args.Tolerations) > 0 {
		if deploy.Spec.Template.Spec.Tolerations == nil {
			deploy.Spec.Template.Spec.Tolerations = []v1.Toleration{}
		}
		mapSet := make(map[string]interface{})
		for _, toleration := range deploy.Spec.Template.Spec.Tolerations {
			mapSet[fmt.Sprintf("%s=%s:%s,%s", toleration.Key,
				toleration.Value,
				toleration.Effect,
				toleration.Operator)] = nil
		}
		for _, toleration := range args.Tolerations {
			if _, ok := mapSet[fmt.Sprintf("%s=%s:%s,%s", toleration.Key,
				toleration.Value,
				toleration.Effect,
				toleration.Operator)]; !ok {
				deploy.Spec.Template.Spec.Tolerations = append(deploy.Spec.Template.Spec.Tolerations, v1.Toleration{
					Key:      toleration.Key,
					Value:    toleration.Value,
					Effect:   v1.TaintEffect(toleration.Effect),
					Operator: v1.TolerationOperator(toleration.Operator),
				})
			}

		}
	}

	return updateDeployment(args.Name, args.Version, deploy)
}

func UpdateTritonServing(args *types.UpdateTritonServingArgs) error {
	deploy, err := findAndBuildDeployment(&args.CommonUpdateServingArgs)
	if err != nil {
		return err
	}

	if args.Command == "" {
		containerArgs := deploy.Spec.Template.Spec.Containers[0].Args

		servingArgs := containerArgs[0]
		if strings.HasSuffix(servingArgs, "\n") {
			servingArgs = strings.TrimSpace(servingArgs[:len(servingArgs)-1])
		}
		arr := strings.Split(servingArgs, "--")

		params := make(map[string]string)
		for index, argItem := range arr {
			if index == 0 {
				continue
			}
			pair := strings.Split(argItem, "=")
			if len(pair) <= 1 {
				continue
			}
			params[fmt.Sprintf("--%s", pair[0])] = argItem[len(pair[0])+1:]
		}

		if args.ModelRepository != "" {
			params["--model-repository"] = args.ModelRepository
		}

		var newArgs []string
		newArgs = append(newArgs, "tritonserver")
		for k, v := range params {
			newArgs = append(newArgs, fmt.Sprintf("%s=%s", k, v))
		}

		deploy.Spec.Template.Spec.Containers[0].Args = []string{strings.Join(newArgs, " ")}
	}

	if args.Annotations != nil && len(args.Annotations) > 0 {
		for k, v := range args.Annotations {
			deploy.Annotations[k] = v
			deploy.Spec.Template.Annotations[k] = v
		}
	}

	if args.Labels != nil && len(args.Labels) > 0 {
		for k, v := range args.Labels {
			deploy.Labels[k] = v
			deploy.Spec.Template.Labels[k] = v
		}
	}

	if args.NodeSelectors != nil && len(args.NodeSelectors) > 0 {
		if deploy.Spec.Template.Spec.NodeSelector == nil {
			deploy.Spec.Template.Spec.NodeSelector = map[string]string{}
		}
		for k, v := range args.NodeSelectors {
			deploy.Spec.Template.Spec.NodeSelector[k] = v
		}
	}

	if args.Tolerations != nil && len(args.Tolerations) > 0 {
		if deploy.Spec.Template.Spec.Tolerations == nil {
			deploy.Spec.Template.Spec.Tolerations = []v1.Toleration{}
		}
		for _, toleration := range args.Tolerations {
			deploy.Spec.Template.Spec.Tolerations = append(deploy.Spec.Template.Spec.Tolerations, v1.Toleration{
				Key:      toleration.Key,
				Value:    toleration.Value,
				Effect:   v1.TaintEffect(toleration.Effect),
				Operator: v1.TolerationOperator(toleration.Operator),
			})
		}
	}

	return updateDeployment(args.Name, args.Version, deploy)
}

func UpdateCustomServing(args *types.UpdateCustomServingArgs) error {
	deploy, err := findAndBuildDeployment(&args.CommonUpdateServingArgs)
	if err != nil {
		return err
	}

	if args.Annotations != nil && len(args.Annotations) > 0 {
		for k, v := range args.Annotations {
			deploy.Annotations[k] = v
			deploy.Spec.Template.Annotations[k] = v
		}
	}

	if args.Labels != nil && len(args.Labels) > 0 {
		for k, v := range args.Labels {
			deploy.Labels[k] = v
			deploy.Spec.Template.Labels[k] = v
		}
	}

	if args.NodeSelectors != nil && len(args.NodeSelectors) > 0 {
		if deploy.Spec.Template.Spec.NodeSelector == nil {
			deploy.Spec.Template.Spec.NodeSelector = map[string]string{}
		}
		for k, v := range args.NodeSelectors {
			deploy.Spec.Template.Spec.NodeSelector[k] = v
		}
	}

	if args.Tolerations != nil && len(args.Tolerations) > 0 {
		if deploy.Spec.Template.Spec.Tolerations == nil {
			deploy.Spec.Template.Spec.Tolerations = []v1.Toleration{}
		}
		exist := map[string]bool{}
		var tolerations []v1.Toleration
		for _, toleration := range args.Tolerations {
			tolerations = append(tolerations, v1.Toleration{
				Key:      toleration.Key,
				Value:    toleration.Value,
				Effect:   v1.TaintEffect(toleration.Effect),
				Operator: v1.TolerationOperator(toleration.Operator),
			})
			exist[toleration.Key+toleration.Value] = true
		}

		for _, preToleration := range deploy.Spec.Template.Spec.Tolerations {
			if !exist[preToleration.Key+preToleration.Value] {
				tolerations = append(tolerations, preToleration)
			}
		}
		deploy.Spec.Template.Spec.Tolerations = tolerations
	}

	return updateDeployment(args.Name, args.Version, deploy)
}

func UpdateKServe(args *types.UpdateKServeArgs) error {

	return nil
}

func findAndBuildDeployment(args *types.CommonUpdateServingArgs) (*appsv1.Deployment, error) {
	job, err := SearchServingJob(args.Namespace, args.Name, args.Version, args.Type)
	if err != nil {
		return nil, err
	}
	if args.Version == "" {
		jobInfo := job.Convert2JobInfo()
		args.Version = jobInfo.Version
	}

	var suffix string
	switch args.Type {
	case types.TFServingJob:
		suffix = "tensorflow-serving"
		break
	case types.TritonServingJob:
		suffix = "tritoninferenceserver"
		break
	case types.CustomServingJob:
		suffix = "custom-serving"
		break
	default:
		return nil, fmt.Errorf("invalid serving job type [%s]", args.Type)
	}

	deployName := fmt.Sprintf("%s-%s-%s", args.Name, args.Version, suffix)
	deploy, err := kubectl.GetDeployment(deployName, args.Namespace)
	if err != nil {
		return nil, err
	}

	if args.Image != "" {
		deploy.Spec.Template.Spec.Containers[0].Image = args.Image
	}

	if args.Replicas >= 0 {
		replicas := int32(args.Replicas)
		deploy.Spec.Replicas = &replicas
	}

	resourceLimits := deploy.Spec.Template.Spec.Containers[0].Resources.Limits
	if resourceLimits == nil {
		resourceLimits = make(map[v1.ResourceName]resource.Quantity)
	}

	if args.GPUCount > 0 {
		resourceLimits[ResourceGPU] = resource.MustParse(strconv.Itoa(args.GPUCount))
		if _, ok := resourceLimits[ResourceGPUMemory]; ok {
			delete(resourceLimits, ResourceGPUMemory)
		}
	}

	if args.GPUMemory > 0 {
		resourceLimits[ResourceGPUMemory] = resource.MustParse(strconv.Itoa(args.GPUMemory))
		if _, ok := resourceLimits[ResourceGPU]; ok {
			delete(resourceLimits, ResourceGPU)
		}
	}

	if args.GPUCore > 0 && args.GPUCore%5 == 0 {
		resourceLimits[ResourceGPUCore] = resource.MustParse(strconv.Itoa(args.GPUCore))
		if _, ok := resourceLimits[ResourceGPU]; ok {
			delete(resourceLimits, ResourceGPU)
		}
	}

	if args.Cpu != "" {
		resourceLimits[v1.ResourceCPU] = resource.MustParse(args.Cpu)
	}

	if args.Memory != "" {
		resourceLimits[v1.ResourceMemory] = resource.MustParse(args.Memory)
	}
	deploy.Spec.Template.Spec.Containers[0].Resources.Limits = resourceLimits

	var newEnvs []v1.EnvVar
	exist := map[string]bool{}
	if args.Envs != nil {
		for k, v := range args.Envs {
			envVar := v1.EnvVar{
				Name:  k,
				Value: v,
			}
			newEnvs = append(newEnvs, envVar)
			exist[k] = true
		}
	}
	for _, env := range deploy.Spec.Template.Spec.Containers[0].Env {
		if !exist[env.Name] {
			newEnvs = append(newEnvs, env)
		}
	}
	deploy.Spec.Template.Spec.Containers[0].Env = newEnvs

	if args.Command != "" {
		// commands: sh -c xxx
		commands := deploy.Spec.Template.Spec.Containers[0].Command
		shell := commands[0]
		newCommands := []string{shell, "-c", args.Command}
		deploy.Spec.Template.Spec.Containers[0].Command = newCommands
		deploy.Spec.Template.Spec.Containers[0].Args = []string{}
	}

	return deploy, nil
}

func updateDeployment(name, version string, deploy *appsv1.Deployment) error {
	err := kubectl.UpdateDeployment(deploy)
	if err == nil {
		log.Infof("The serving job %s with version %s has been updated successfully", name, version)
	} else {
		log.Errorf("The serving job %s with version %s update failed", name, version)
	}
	return err
}

func findAndBuildInferenceService(args *types.CommonUpdateServingArgs) error {

	return nil
}