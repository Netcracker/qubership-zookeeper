// Copyright 2024-2025 NetCracker Technology Corporation
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

package provider

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

const (
	AutoRestartAnnotation             = "zookeeperservice.netcracker.com/auto-restart"
	ResourceVersionAnnotationTemplate = "%s/resource-version"
	podSecretsMountPathTemplate       = "/etc/secrets/%s-pod-secrets"
)

func GetPodSecretsMountPath(service string) string {
	return fmt.Sprintf(podSecretsMountPathTemplate, service)
}

func getPodSecretsDefaultMode() *int32 {
	mode := int32(0644)
	return &mode
}

type ProjectedSecretSource struct {
	SecretName string
	Items      []corev1.KeyToPath
}

func NewPodSecretsProjectedVolume(volumeName string, sources []ProjectedSecretSource) corev1.Volume {
	var projectedSources []corev1.VolumeProjection
	for _, source := range sources {
		if source.SecretName == "" || len(source.Items) == 0 {
			continue
		}
		projectedSources = append(projectedSources, corev1.VolumeProjection{
			Secret: &corev1.SecretProjection{
				LocalObjectReference: corev1.LocalObjectReference{Name: source.SecretName},
				Items:                source.Items,
			},
		})
	}
	return corev1.Volume{
		Name: volumeName,
		VolumeSource: corev1.VolumeSource{
			Projected: &corev1.ProjectedVolumeSource{
				DefaultMode: getPodSecretsDefaultMode(),
				Sources:     projectedSources,
			},
		},
	}
}

func NewPodSecretsVolumeMount(volumeName, mountPath string) corev1.VolumeMount {
	return corev1.VolumeMount{
		Name:      volumeName,
		MountPath: mountPath,
		ReadOnly:  true,
	}
}

func SecretKeyToPath(key, path string) corev1.KeyToPath {
	return corev1.KeyToPath{Key: key, Path: path}
}
