# Copyright 2024-2025 NetCracker Technology Corporation
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import os

DEFAULT_SECRETS_DIR = '/etc/secrets/zookeeper-integration-tests-pod-secrets'


def _secrets_dir(environ) -> str:
    return environ.get('INTEGRATION_TESTS_SECRETS_DIR', DEFAULT_SECRETS_DIR)


def check_that_parameters_are_presented(environ, *variable_names) -> bool:
    for variable in variable_names:
        if not environ.get(variable):
            return False
    return True


def secret_is_present(environ, name) -> bool:
    path = os.path.join(_secrets_dir(environ), name)
    return os.path.isfile(path) and os.path.getsize(path) > 0


def get_excluded_tags(environ) -> list:
    if not check_that_parameters_are_presented(environ, 'PROMETHEUS_URL'):
        return ['prometheus']
    if not (secret_is_present(environ, 'PROMETHEUS_USER')
            and secret_is_present(environ, 'PROMETHEUS_PASSWORD')):
        return ['prometheus']
    return []
