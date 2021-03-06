// Copyright 2017 The Kubernetes Dashboard Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package job

import (
	"github.com/kubernetes/dashboard/src/app/backend/api"
	"github.com/kubernetes/dashboard/src/app/backend/integration/metric/heapster"
	"github.com/kubernetes/dashboard/src/app/backend/resource/common"
	"github.com/kubernetes/dashboard/src/app/backend/resource/dataselect"
	"github.com/kubernetes/dashboard/src/app/backend/resource/pod"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sClient "k8s.io/client-go/kubernetes"
	batch "k8s.io/client-go/pkg/apis/batch/v1"
)

// JobDetail is a presentation layer view of Kubernetes Job resource. This means
// it is Job plus additional augmented data we can get from other sources
// (like services that target the same pods).
type JobDetail struct {
	ObjectMeta api.ObjectMeta `json:"objectMeta"`
	TypeMeta   api.TypeMeta   `json:"typeMeta"`

	// Aggregate information about pods belonging to this Job.
	PodInfo common.PodInfo `json:"podInfo"`

	// Detailed information about Pods belonging to this Job.
	PodList pod.PodList `json:"podList"`

	// Container images of the Job.
	ContainerImages []string `json:"containerImages"`

	// List of events related to this Job.
	EventList common.EventList `json:"eventList"`

	// Parallelism specifies the maximum desired number of pods the job should run at any given
	// time.
	Parallelism *int32 `json:"parallelism"`

	// Completions specifies the desired number of successfully finished pods the job should be
	// run with.
	Completions *int32 `json:"completions"`
}

// GetJobDetail gets job details.
func GetJobDetail(client k8sClient.Interface, heapsterClient heapster.HeapsterClient,
	namespace, name string) (*JobDetail, error) {

	// TODO(floreks): Use channels.
	jobData, err := client.BatchV1().Jobs(namespace).Get(name, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}

	podList, err := GetJobPods(client, heapsterClient, dataselect.DefaultDataSelectWithMetrics, namespace, name)
	if err != nil {
		return nil, err
	}

	podInfo, err := getJobPodInfo(client, jobData)
	if err != nil {
		return nil, err
	}

	eventList, err := GetJobEvents(client, dataselect.DefaultDataSelect, jobData.Namespace, jobData.Name)
	if err != nil {
		return nil, err
	}

	job := getJobDetail(jobData, heapsterClient, *eventList, *podList, *podInfo)
	return &job, nil
}

func getJobDetail(job *batch.Job, heapsterClient heapster.HeapsterClient,
	eventList common.EventList, podList pod.PodList, podInfo common.PodInfo) JobDetail {
	return JobDetail{
		ObjectMeta:      api.NewObjectMeta(job.ObjectMeta),
		TypeMeta:        api.NewTypeMeta(api.ResourceKindJob),
		ContainerImages: common.GetContainerImages(&job.Spec.Template.Spec),
		PodInfo:         podInfo,
		PodList:         podList,
		EventList:       eventList,
		Parallelism:     job.Spec.Parallelism,
		Completions:     job.Spec.Completions,
	}
}
