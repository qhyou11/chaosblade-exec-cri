/*
 * Copyright 1999-2019 Alibaba Group Holding Ltd.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package exec

import (
	"context"
	"github.com/chaosblade-io/chaosblade-spec-go/log"

	"github.com/chaosblade-io/chaosblade-spec-go/spec"
)

const (
	ForceFlag = "force"
)

type ContainerCommandModelSpec struct {
	spec.BaseExpModelCommandSpec
}

func NewContainerCommandSpec() spec.ExpModelCommandSpec {
	return &ContainerCommandModelSpec{
		spec.BaseExpModelCommandSpec{
			ExpActions: []spec.ExpActionCommandSpec{
				NewRemoveActionCommand(),
			},
			ExpFlags: []spec.ExpFlagSpec{},
		},
	}
}

func (cms *ContainerCommandModelSpec) Name() string {
	return "container"
}

func (cms *ContainerCommandModelSpec) ShortDesc() string {
	return `Execute a container experiment`
}

func (cms *ContainerCommandModelSpec) LongDesc() string {
	return `Execute a container experiment.`
}

type RemoveActionCommand struct {
	spec.BaseExpActionCommandSpec
}

func NewRemoveActionCommand() spec.ExpActionCommandSpec {
	return &RemoveActionCommand{
		spec.BaseExpActionCommandSpec{
			ActionMatchers: []spec.ExpFlagSpec{},
			ActionFlags: []spec.ExpFlagSpec{
				&spec.ExpFlag{
					Name:   ForceFlag,
					Desc:   "force remove",
					NoArgs: true,
				},
			},
			ActionExecutor: &removeActionExecutor{},
			ActionExample: `# Delete the container id that is a76d53933d3f",
blade create cri container remove --container-id a76d53933d3f. If container-runtime is contained, the container-id shoud be full id`,
			ActionCategories: []string{CategorySystemContainer},
		},
	}
}

func (*RemoveActionCommand) Name() string {
	return "remove"
}

func (*RemoveActionCommand) Aliases() []string {
	return []string{"rm"}
}

func (*RemoveActionCommand) ShortDesc() string {
	return "remove a container"
}

func (r *RemoveActionCommand) LongDesc() string {
	if r.ActionLongDesc != "" {
		return r.ActionLongDesc
	}
	return "remove a container"
}

type removeActionExecutor struct {
}

func (*removeActionExecutor) Name() string {
	return "remove"
}

func (e *removeActionExecutor) SetChannel(channel spec.Channel) {
}

func (e *removeActionExecutor) Exec(uid string, ctx context.Context, model *spec.ExpModel) *spec.Response {
	if _, ok := spec.IsDestroy(ctx); ok {
		return spec.ReturnSuccess(uid)
	}
	flags := model.ActionFlags
	client, err := GetClientByRuntime(model)
	if err != nil {
		log.Errorf(ctx, spec.ContainerExecFailed.Sprintf("GetClient", err))
		return spec.ResponseFailWithFlags(spec.ContainerExecFailed, "GetClient", err)
	}
	containerId := flags[ContainerIdFlag.Name]
	// containerName := flags[ContainerNameFlag.Name]
	// containerLabelSelector := parseContainerLabelSelector(flags[ContainerNameFlag.Name])
	// container, _ := GetContainer(ctx, client, uid, containerId, containerName, containerLabelSelector)
	// if !response.Success {
	// 	return response
	// }
	forceFlag := flags[ForceFlag]

	err = client.RemoveContainer(ctx, containerId, judgeForce(forceFlag))
	if err != nil {
		log.Errorf(ctx, spec.ContainerExecFailed.Sprintf("ContainerRemove", err))
		return spec.ResponseFailWithFlags(spec.ContainerExecFailed, "ContainerRemove", err)
	}
	return spec.ReturnSuccess(uid)
}

func judgeForce(forceflag string) bool {
	if forceflag != "" {
		return true
	}
	return false
}
