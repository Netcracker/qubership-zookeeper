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

package controllers

import (
	"context"
	"fmt"
	zookeeperservice "github.com/Netcracker/qubership-zookeeper/operator/api/v1"
	"github.com/Netcracker/qubership-zookeeper/operator/controllers/provider"
	"github.com/Netcracker/qubership-zookeeper/operator/util"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"time"
)

const (
	backupDaemonConditionReason = "ZooKeeperBackupDaemonReadinessStatus"
	backupDaemonHashName        = "spec.backupDaemon"
)

type ReconcileBackupDaemon struct {
	reconciler           *ZooKeeperServiceReconciler
	cr                   *zookeeperservice.ZooKeeperService
	logger               logr.Logger
	backupDaemonProvider provider.BackupDaemonResourceProvider
}

func NewReconcileBackupDaemon(r *ZooKeeperServiceReconciler, cr *zookeeperservice.ZooKeeperService, logger logr.Logger) ReconcileBackupDaemon {
	return ReconcileBackupDaemon{
		reconciler:           r,
		cr:                   cr,
		logger:               logger,
		backupDaemonProvider: provider.NewBackupDaemonResourceProvider(cr, logger),
	}
}

func (r ReconcileBackupDaemon) Status() error {
	if err := r.reconciler.updateConditions(r.cr,
		NewCondition(statusFalse,
			typeInProgress,
			backupDaemonConditionReason,
			"ZooKeeper Backup Daemon health check")); err != nil {
		return err
	}
	r.logger.Info("Start checking for ZooKeeper Backup Daemon pod")
	err := wait.PollUntilContextTimeout(context.Background(), 10*time.Second, time.Duration(r.cr.Spec.Global.PodsReadyTimeout)*time.Second, true, func(ctx context.Context) (done bool, err error) {
		if r.reconciler.isDeploymentReady(r.backupDaemonProvider.GetServiceName(), r.cr.Namespace, r.logger) {
			return true, nil
		}
		r.logger.Info(fmt.Sprintf("%s is not ready yet", r.backupDaemonProvider.GetServiceName()))
		return false, nil
	})
	if err != nil {
		return r.reconciler.updateConditions(r.cr, NewCondition(statusFalse,
			typeFailed,
			backupDaemonConditionReason,
			"ZooKeeper Backup Daemon pod is not ready"))
	}
	return r.reconciler.updateConditions(r.cr, NewCondition(statusTrue,
		typeReady,
		backupDaemonConditionReason,
		"ZooKeeper Backup Daemon pod is ready"))
}

func (r ReconcileBackupDaemon) Reconcile() error {
	backupDaemonProvider := r.backupDaemonProvider

	backupDaemonSecret, err := r.reconciler.watchSecret(r.cr.Spec.BackupDaemon.SecretName, r.cr, r.logger)
	if err != nil {
		if provider.IsVaultSecretManagementEnabled(r.cr) && errors.IsNotFound(err) {
			log.Info("Secret was removed. Cannot watch for it to rotate secrets")
			backupDaemonSecret = &corev1.Secret{}
		} else {
			return err
		}
	}

	var s3Secret *corev1.Secret
	if r.cr.Spec.BackupDaemon.S3 != nil && r.cr.Spec.BackupDaemon.S3.Enabled && r.cr.Spec.BackupDaemon.S3.SecretName != "" {
		s3Secret, err = r.reconciler.watchSecret(r.cr.Spec.BackupDaemon.S3.SecretName, r.cr, r.logger)
		if err != nil && !errors.IsNotFound(err) {
			return err
		}
	}

	s3AliasesSecretName := fmt.Sprintf("%s-backup-daemon-s3-aliases", r.cr.Name)
	s3AliasesSecret, aliasesErr := r.reconciler.findSecret(s3AliasesSecretName, r.cr.Namespace, r.logger)
	s3AliasesEnabled := aliasesErr == nil

	backupDaemonSpecHash, err := util.Hash(r.cr.Spec.BackupDaemon)
	if err != nil {
		return err
	}
	if r.reconciler.ResourceHashes[backupDaemonHashName] != backupDaemonSpecHash ||
		r.reconciler.ResourceHashes[globalHashName] != globalSpecHash ||
		secretVersionChanged(r.reconciler.ResourceVersions, backupDaemonSecret) ||
		secretVersionChanged(r.reconciler.ResourceVersions, s3Secret) ||
		secretVersionChanged(r.reconciler.ResourceVersions, s3AliasesSecret) {
		if r.cr.Spec.BackupDaemon.BackupStorage.PersistentVolumeType != "" {
			backupStorage := r.cr.Spec.BackupDaemon.BackupStorage.DeepCopy()
			if backupStorage.PersistentVolumeClaimName == "" {
				backupStorage.PersistentVolumeClaimName = fmt.Sprintf(provider.SnapshotsPersistentVolumeClaimPattern, r.cr.Name)
			}

			// Persistent volume claim for snapshots could be created in ZooKeeper
			_, err := r.reconciler.findPersistentVolumeClaim(backupStorage.PersistentVolumeClaimName, r.cr.Namespace, r.logger)
			if err != nil {
				backupPersistentVolumeClaim, err := r.reconciler.processSnapshotsPersistentVolumeClaim(*backupStorage, r.cr, r.logger)
				if err != nil {
					return err
				}
				if backupPersistentVolumeClaim != nil {
					if err := r.reconciler.createPersistentVolumeClaim(backupPersistentVolumeClaim, r.logger); err != nil {
						return err
					}
				}
			}
		}

		clientService := backupDaemonProvider.NewBackupDaemonClientService()
		if err := controllerutil.SetControllerReference(r.cr, clientService, r.reconciler.Scheme); err != nil {
			return err
		}
		if err := r.reconciler.createOrUpdateService(clientService, r.logger); err != nil {
			return err
		}

		serviceAccount := provider.NewServiceAccount(r.backupDaemonProvider.GetServiceAccountName(), r.cr.Namespace)
		if err := r.reconciler.createServiceAccount(serviceAccount, r.logger); err != nil {
			return err
		}

		if provider.IsVaultSecretManagementEnabled(r.cr) {
			err := r.processVaultSecrets(backupDaemonSecret)
			if err != nil {
				return err
			}
		}

		deployment := backupDaemonProvider.NewBackupDaemonDeployment(s3AliasesEnabled)
		if err := controllerutil.SetControllerReference(r.cr, deployment, r.reconciler.Scheme); err != nil {
			return err
		}
		if err := r.reconciler.createOrUpdateDeployment(deployment, r.logger); err != nil {
			return err
		}

		r.logger.Info("Updating ZooKeeper Backup Daemon status")
		if err := r.updateBackupDaemonStatus(r.cr); err != nil {
			return err
		}
	} else {
		r.logger.Info("Backup Daemon configuration didn't change, skipping reconcile loop")
	}

	if err := updateDeploymentSecretRestartAnnotations(
		r.reconciler.Client, r.cr.Namespace, r.backupDaemonProvider.GetServiceName(), r.logger,
		backupDaemonSecret, s3Secret, s3AliasesSecret); err != nil {
		return err
	}

	r.reconciler.ResourceVersions[backupDaemonSecret.Name] = backupDaemonSecret.ResourceVersion
	if s3Secret != nil {
		r.reconciler.ResourceVersions[s3Secret.Name] = s3Secret.ResourceVersion
	}
	if s3AliasesSecret != nil {
		r.reconciler.ResourceVersions[s3AliasesSecret.Name] = s3AliasesSecret.ResourceVersion
	}
	r.reconciler.ResourceHashes[backupDaemonHashName] = backupDaemonSpecHash
	return nil
}

// updateBackupDaemonStatus updates the status of ZooKeeper Backup Daemon
func (r ReconcileBackupDaemon) updateBackupDaemonStatus(cr *zookeeperservice.ZooKeeperService) error {
	labels := r.backupDaemonProvider.GetBackupDaemonSelectorLabels()
	foundPodList, err := r.reconciler.findPodList(r.cr.Namespace, labels)
	if err != nil {
		return err
	}
	r.cr.Status.BackupDaemonStatus.Nodes = getPodNames(foundPodList.Items)
	return r.reconciler.Client.Status().Update(context.TODO(), cr)
}
