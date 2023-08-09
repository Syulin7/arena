package serving

import (
	"fmt"
	"strings"

	"github.com/kubeflow/arena/pkg/apis/types"
	"github.com/kubeflow/arena/pkg/argsbuilder"
)

type UpdateKServeJobBuilder struct {
	args      *types.UpdateKServeArgs
	argValues map[string]interface{}
	argsbuilder.ArgsBuilder
}

func NewUpdateKServeJobBuilder() *UpdateKServeJobBuilder {
	args := &types.UpdateKServeArgs{
		CommonUpdateServingArgs: types.CommonUpdateServingArgs{
			Replicas: 1,
		},
	}
	return &UpdateKServeJobBuilder{
		args:        args,
		argValues:   map[string]interface{}{},
		ArgsBuilder: argsbuilder.NewUpdateKServeArgsBuilder(args),
	}
}

// Name is used to set job name,match option --name
func (b *UpdateKServeJobBuilder) Name(name string) *UpdateKServeJobBuilder {
	if name != "" {
		b.args.Name = name
	}
	return b
}

// Namespace is used to set job namespace,match option --namespace
func (b *UpdateKServeJobBuilder) Namespace(namespace string) *UpdateKServeJobBuilder {
	if namespace != "" {
		b.args.Namespace = namespace
	}
	return b
}

// Command is used to set job command
func (b *UpdateKServeJobBuilder) Command(args []string) *UpdateKServeJobBuilder {
	b.args.Command = strings.Join(args, " ")
	return b
}

// Image is used to set job image,match the option --image
func (b *UpdateKServeJobBuilder) Image(image string) *UpdateKServeJobBuilder {
	if image != "" {
		b.args.Image = image
	}
	return b
}

// Envs is used to set env of job containers,match option --env
func (b *UpdateKServeJobBuilder) Envs(envs map[string]string) *UpdateKServeJobBuilder {
	if envs != nil && len(envs) != 0 {
		envSlice := []string{}
		for key, value := range envs {
			envSlice = append(envSlice, fmt.Sprintf("%v=%v", key, value))
		}
		b.argValues["env"] = &envSlice
	}
	return b
}

// Annotations is used to add annotations for job pods,match option --annotation
func (b *UpdateKServeJobBuilder) Annotations(annotations map[string]string) *UpdateKServeJobBuilder {
	if annotations != nil && len(annotations) != 0 {
		s := []string{}
		for key, value := range annotations {
			s = append(s, fmt.Sprintf("%v=%v", key, value))
		}
		b.argValues["annotation"] = &s
	}
	return b
}

// Labels is used to add labels for job
func (b *UpdateKServeJobBuilder) Labels(labels map[string]string) *UpdateKServeJobBuilder {
	if labels != nil && len(labels) != 0 {
		s := []string{}
		for key, value := range labels {
			s = append(s, fmt.Sprintf("%v=%v", key, value))
		}
		b.argValues["label"] = &s
	}
	return b
}

// Replicas is used to set serving job replicas,match the option --replicas
func (b *UpdateKServeJobBuilder) Replicas(count int) *UpdateKServeJobBuilder {
	if count > 0 {
		b.args.Replicas = count
	}
	return b
}

// Build is used to build the job
func (b *UpdateKServeJobBuilder) Build() (*Job, error) {
	for key, value := range b.argValues {
		b.AddArgValue(key, value)
	}
	if err := b.PreBuild(); err != nil {
		return nil, err
	}
	if err := b.ArgsBuilder.Build(); err != nil {
		return nil, err
	}
	return NewJob(b.args.Name, types.KServeJob, b.args), nil
}
