/*
Copyright 2024.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
	http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

// InstanceState describes the state of an ECS instance.
type InstanceState string

var (
	// InstanceStatePending is the string representing an instance in a pending state.
	InstanceStatePending = InstanceState("pending")

	// InstanceStateRunning is the string representing an instance in a running state.
	InstanceStateRunning = InstanceState("running")

	// InstanceStateShuttingDown is the string representing an instance shutting down.
	InstanceStateShuttingDown = InstanceState("shutting-down")

	// InstanceStateTerminated is the string representing an instance that has been terminated.
	InstanceStateTerminated = InstanceState("terminated")

	// InstanceStateStopping is the string representing an instance
	// that is in the process of being stopped and can be restarted.
	InstanceStateStopping = InstanceState("stopping")

	// InstanceStateStopped is the string representing an instance
	// that has been stopped and can be restarted.
	InstanceStateStopped = InstanceState("stopped")
)
