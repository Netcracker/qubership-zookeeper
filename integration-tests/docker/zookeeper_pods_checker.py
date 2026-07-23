import os
import time

from PlatformLibrary import PlatformLibrary

environ = os.environ
namespace = environ.get("ZOOKEEPER_OS_PROJECT")
service = environ.get("ZOOKEEPER_HOST")
backup_daemon = environ.get("ZOOKEEPER_BACKUP_DAEMON_HOST")
managed_by_operator = environ.get("ZOOKEEPER_IS_MANAGED_BY_OPERATOR")
timeout = 500

CR_API_VERSION = "netcracker.com/v1"
CR_KIND = "ZooKeeperService"


def get_desired_zk_replicas(k8s_lib):
    cr = k8s_lib.get_custom_resource(CR_API_VERSION, CR_KIND, namespace, service)
    return cr["spec"]["zooKeeper"]["replicas"]


def check_cr_reconciled(k8s_lib):
    cr = k8s_lib.get_custom_resource(CR_API_VERSION, CR_KIND, namespace, service)
    for cond in cr.get("status", {}).get("conditions", []):
        if cond.get("type") == "Ready" and cond.get("reason") == "ZooKeeperReadinessStatus" and cond.get("status") == "True":
            return True
    conditions = [(c.get("type"), c.get("status"), c.get("reason"))
                  for c in cr.get("status", {}).get("conditions", [])]
    print(f"[CR] not reconciled, conditions: {conditions}")
    return False


if __name__ == '__main__':
    time.sleep(20)
    try:
        k8s_lib = PlatformLibrary(managed_by_operator)
        desired_replicas = get_desired_zk_replicas(k8s_lib)
        print(f"[CR] desired ZooKeeper replicas: {desired_replicas}")
    except Exception as e:
        print(e)
        exit(1)

    timeout_start = time.time()
    while time.time() < timeout_start + timeout:
        try:
            deployments = k8s_lib.get_deployment_entities_count_for_service(namespace, service)
            ready_deployments = k8s_lib.get_active_deployment_entities_count_for_service(namespace, service)
            if backup_daemon is not None and len(backup_daemon) != 0:
                deployments += k8s_lib.get_deployment_entities_count_for_service(namespace, backup_daemon, 'component')
                ready_deployments += k8s_lib.get_active_deployment_entities_count_for_service(namespace, backup_daemon, 'component')

            operator_name = f"{service}-service-operator"
            print(f"[Operator] name: {operator_name}")
            operator_total = k8s_lib.get_deployment_entities_count_for_service(namespace, operator_name, 'component')
            operator_ready = k8s_lib.get_active_deployment_entities_count_for_service(namespace, operator_name,'component')

            print(f"[Check status] deployments: {deployments}, ready deployments: {ready_deployments}")
            print(f"[Operator] active: {operator_ready}/{operator_total}")
        except Exception as e:
            print(e)
            time.sleep(10)
            continue

        zk_replicas_ok = ready_deployments >= desired_replicas
        operator_ok = operator_ready == operator_total and operator_total > 0
        cr_ok = check_cr_reconciled(k8s_lib)

        if deployments == ready_deployments and deployments != 0 and zk_replicas_ok and operator_ok and cr_ok:
            print("ZooKeeper deployments are ready")
            exit(0)
        time.sleep(10)

    print(f'ZooKeeper deployments are not ready at least {timeout} seconds')
    exit(1)
