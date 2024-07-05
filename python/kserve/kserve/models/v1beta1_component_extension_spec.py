# Copyright 2023 The KServe Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# coding: utf-8

"""
    KServe

    Python SDK for KServe  # noqa: E501

    The version of the OpenAPI document: v0.1
    Generated by: https://openapi-generator.tech
"""


import pprint
import re  # noqa: F401

import six

from kserve.configuration import Configuration


class V1beta1ComponentExtensionSpec(object):
    """NOTE: This class is auto generated by OpenAPI Generator.
    Ref: https://openapi-generator.tech

    Do not edit the class manually.
    """

    """
    Attributes:
      openapi_types (dict): The key is attribute name
                            and the value is attribute type.
      attribute_map (dict): The key is attribute name
                            and the value is json key in definition.
    """
    openapi_types = {
        'annotations': 'dict(str, str)',
        'batcher': 'V1beta1Batcher',
        'canary_traffic_percent': 'int',
        'container_concurrency': 'int',
        'deployment_strategy': 'K8sIoApiAppsV1DeploymentStrategy',
        'labels': 'dict(str, str)',
        'logger': 'V1beta1LoggerSpec',
        'max_replicas': 'int',
        'min_replicas': 'int',
        'query': 'str',
        'query_parameters': 'str',
        'scale_metric': 'str',
        'scale_metric_type': 'str',
        'scale_target': 'int',
        'scaler': 'V1beta1ScalerSpec',
        'server_address': 'str',
        'timeout': 'int'
    }

    attribute_map = {
        'annotations': 'annotations',
        'batcher': 'batcher',
        'canary_traffic_percent': 'canaryTrafficPercent',
        'container_concurrency': 'containerConcurrency',
        'deployment_strategy': 'deploymentStrategy',
        'labels': 'labels',
        'logger': 'logger',
        'max_replicas': 'maxReplicas',
        'min_replicas': 'minReplicas',
        'query': 'query',
        'query_parameters': 'queryParameters',
        'scale_metric': 'scaleMetric',
        'scale_metric_type': 'scaleMetricType',
        'scale_target': 'scaleTarget',
        'scaler': 'scaler',
        'server_address': 'serverAddress',
        'timeout': 'timeout'
    }

    def __init__(self, annotations=None, batcher=None, canary_traffic_percent=None, container_concurrency=None, deployment_strategy=None, labels=None, logger=None, max_replicas=None, min_replicas=None, query=None, query_parameters=None, scale_metric=None, scale_metric_type=None, scale_target=None, scaler=None, server_address=None, timeout=None, local_vars_configuration=None):  # noqa: E501
        """V1beta1ComponentExtensionSpec - a model defined in OpenAPI"""  # noqa: E501
        if local_vars_configuration is None:
            local_vars_configuration = Configuration()
        self.local_vars_configuration = local_vars_configuration

        self._annotations = None
        self._batcher = None
        self._canary_traffic_percent = None
        self._container_concurrency = None
        self._deployment_strategy = None
        self._labels = None
        self._logger = None
        self._max_replicas = None
        self._min_replicas = None
        self._query = None
        self._query_parameters = None
        self._scale_metric = None
        self._scale_metric_type = None
        self._scale_target = None
        self._scaler = None
        self._server_address = None
        self._timeout = None
        self.discriminator = None

        if annotations is not None:
            self.annotations = annotations
        if batcher is not None:
            self.batcher = batcher
        if canary_traffic_percent is not None:
            self.canary_traffic_percent = canary_traffic_percent
        if container_concurrency is not None:
            self.container_concurrency = container_concurrency
        if deployment_strategy is not None:
            self.deployment_strategy = deployment_strategy
        if labels is not None:
            self.labels = labels
        if logger is not None:
            self.logger = logger
        if max_replicas is not None:
            self.max_replicas = max_replicas
        if min_replicas is not None:
            self.min_replicas = min_replicas
        if query is not None:
            self.query = query
        if query_parameters is not None:
            self.query_parameters = query_parameters
        if scale_metric is not None:
            self.scale_metric = scale_metric
        if scale_metric_type is not None:
            self.scale_metric_type = scale_metric_type
        if scale_target is not None:
            self.scale_target = scale_target
        if scaler is not None:
            self.scaler = scaler
        if server_address is not None:
            self.server_address = server_address
        if timeout is not None:
            self.timeout = timeout

    @property
    def annotations(self):
        """Gets the annotations of this V1beta1ComponentExtensionSpec.  # noqa: E501

        Annotations that will be added to the component pod. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/  # noqa: E501

        :return: The annotations of this V1beta1ComponentExtensionSpec.  # noqa: E501
        :rtype: dict(str, str)
        """
        return self._annotations

    @annotations.setter
    def annotations(self, annotations):
        """Sets the annotations of this V1beta1ComponentExtensionSpec.

        Annotations that will be added to the component pod. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/  # noqa: E501

        :param annotations: The annotations of this V1beta1ComponentExtensionSpec.  # noqa: E501
        :type: dict(str, str)
        """

        self._annotations = annotations

    @property
    def batcher(self):
        """Gets the batcher of this V1beta1ComponentExtensionSpec.  # noqa: E501


        :return: The batcher of this V1beta1ComponentExtensionSpec.  # noqa: E501
        :rtype: V1beta1Batcher
        """
        return self._batcher

    @batcher.setter
    def batcher(self, batcher):
        """Sets the batcher of this V1beta1ComponentExtensionSpec.


        :param batcher: The batcher of this V1beta1ComponentExtensionSpec.  # noqa: E501
        :type: V1beta1Batcher
        """

        self._batcher = batcher

    @property
    def canary_traffic_percent(self):
        """Gets the canary_traffic_percent of this V1beta1ComponentExtensionSpec.  # noqa: E501

        CanaryTrafficPercent defines the traffic split percentage between the candidate revision and the last ready revision  # noqa: E501

        :return: The canary_traffic_percent of this V1beta1ComponentExtensionSpec.  # noqa: E501
        :rtype: int
        """
        return self._canary_traffic_percent

    @canary_traffic_percent.setter
    def canary_traffic_percent(self, canary_traffic_percent):
        """Sets the canary_traffic_percent of this V1beta1ComponentExtensionSpec.

        CanaryTrafficPercent defines the traffic split percentage between the candidate revision and the last ready revision  # noqa: E501

        :param canary_traffic_percent: The canary_traffic_percent of this V1beta1ComponentExtensionSpec.  # noqa: E501
        :type: int
        """

        self._canary_traffic_percent = canary_traffic_percent

    @property
    def container_concurrency(self):
        """Gets the container_concurrency of this V1beta1ComponentExtensionSpec.  # noqa: E501

        ContainerConcurrency specifies how many requests can be processed concurrently, this sets the hard limit of the container concurrency(https://knative.dev/docs/serving/autoscaling/concurrency).  # noqa: E501

        :return: The container_concurrency of this V1beta1ComponentExtensionSpec.  # noqa: E501
        :rtype: int
        """
        return self._container_concurrency

    @container_concurrency.setter
    def container_concurrency(self, container_concurrency):
        """Sets the container_concurrency of this V1beta1ComponentExtensionSpec.

        ContainerConcurrency specifies how many requests can be processed concurrently, this sets the hard limit of the container concurrency(https://knative.dev/docs/serving/autoscaling/concurrency).  # noqa: E501

        :param container_concurrency: The container_concurrency of this V1beta1ComponentExtensionSpec.  # noqa: E501
        :type: int
        """

        self._container_concurrency = container_concurrency

    @property
    def deployment_strategy(self):
        """Gets the deployment_strategy of this V1beta1ComponentExtensionSpec.  # noqa: E501


        :return: The deployment_strategy of this V1beta1ComponentExtensionSpec.  # noqa: E501
        :rtype: K8sIoApiAppsV1DeploymentStrategy
        """
        return self._deployment_strategy

    @deployment_strategy.setter
    def deployment_strategy(self, deployment_strategy):
        """Sets the deployment_strategy of this V1beta1ComponentExtensionSpec.


        :param deployment_strategy: The deployment_strategy of this V1beta1ComponentExtensionSpec.  # noqa: E501
        :type: K8sIoApiAppsV1DeploymentStrategy
        """

        self._deployment_strategy = deployment_strategy

    @property
    def labels(self):
        """Gets the labels of this V1beta1ComponentExtensionSpec.  # noqa: E501

        Labels that will be added to the component pod. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/  # noqa: E501

        :return: The labels of this V1beta1ComponentExtensionSpec.  # noqa: E501
        :rtype: dict(str, str)
        """
        return self._labels

    @labels.setter
    def labels(self, labels):
        """Sets the labels of this V1beta1ComponentExtensionSpec.

        Labels that will be added to the component pod. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/  # noqa: E501

        :param labels: The labels of this V1beta1ComponentExtensionSpec.  # noqa: E501
        :type: dict(str, str)
        """

        self._labels = labels

    @property
    def logger(self):
        """Gets the logger of this V1beta1ComponentExtensionSpec.  # noqa: E501


        :return: The logger of this V1beta1ComponentExtensionSpec.  # noqa: E501
        :rtype: V1beta1LoggerSpec
        """
        return self._logger

    @logger.setter
    def logger(self, logger):
        """Sets the logger of this V1beta1ComponentExtensionSpec.


        :param logger: The logger of this V1beta1ComponentExtensionSpec.  # noqa: E501
        :type: V1beta1LoggerSpec
        """

        self._logger = logger

    @property
    def max_replicas(self):
        """Gets the max_replicas of this V1beta1ComponentExtensionSpec.  # noqa: E501

        Maximum number of replicas for autoscaling.  # noqa: E501

        :return: The max_replicas of this V1beta1ComponentExtensionSpec.  # noqa: E501
        :rtype: int
        """
        return self._max_replicas

    @max_replicas.setter
    def max_replicas(self, max_replicas):
        """Sets the max_replicas of this V1beta1ComponentExtensionSpec.

        Maximum number of replicas for autoscaling.  # noqa: E501

        :param max_replicas: The max_replicas of this V1beta1ComponentExtensionSpec.  # noqa: E501
        :type: int
        """

        self._max_replicas = max_replicas

    @property
    def min_replicas(self):
        """Gets the min_replicas of this V1beta1ComponentExtensionSpec.  # noqa: E501

        Minimum number of replicas, defaults to 1 but can be set to 0 to enable scale-to-zero.  # noqa: E501

        :return: The min_replicas of this V1beta1ComponentExtensionSpec.  # noqa: E501
        :rtype: int
        """
        return self._min_replicas

    @min_replicas.setter
    def min_replicas(self, min_replicas):
        """Sets the min_replicas of this V1beta1ComponentExtensionSpec.

        Minimum number of replicas, defaults to 1 but can be set to 0 to enable scale-to-zero.  # noqa: E501

        :param min_replicas: The min_replicas of this V1beta1ComponentExtensionSpec.  # noqa: E501
        :type: int
        """

        self._min_replicas = min_replicas

    @property
    def query(self):
        """Gets the query of this V1beta1ComponentExtensionSpec.  # noqa: E501

        Query to run to get metrics from Prometheus  # noqa: E501

        :return: The query of this V1beta1ComponentExtensionSpec.  # noqa: E501
        :rtype: str
        """
        return self._query

    @query.setter
    def query(self, query):
        """Sets the query of this V1beta1ComponentExtensionSpec.

        Query to run to get metrics from Prometheus  # noqa: E501

        :param query: The query of this V1beta1ComponentExtensionSpec.  # noqa: E501
        :type: str
        """

        self._query = query

    @property
    def query_parameters(self):
        """Gets the query_parameters of this V1beta1ComponentExtensionSpec.  # noqa: E501

        A comma-separated list of query Parameters to include while querying the Prometheus endpoint.  # noqa: E501

        :return: The query_parameters of this V1beta1ComponentExtensionSpec.  # noqa: E501
        :rtype: str
        """
        return self._query_parameters

    @query_parameters.setter
    def query_parameters(self, query_parameters):
        """Sets the query_parameters of this V1beta1ComponentExtensionSpec.

        A comma-separated list of query Parameters to include while querying the Prometheus endpoint.  # noqa: E501

        :param query_parameters: The query_parameters of this V1beta1ComponentExtensionSpec.  # noqa: E501
        :type: str
        """

        self._query_parameters = query_parameters

    @property
    def scale_metric(self):
        """Gets the scale_metric of this V1beta1ComponentExtensionSpec.  # noqa: E501

        ScaleMetric defines the scaling metric type watched by autoscaler possible values are concurrency, rps, cpu, memory. concurrency, rps are supported via Knative Pod Autoscaler(https://knative.dev/docs/serving/autoscaling/autoscaling-metrics).  # noqa: E501

        :return: The scale_metric of this V1beta1ComponentExtensionSpec.  # noqa: E501
        :rtype: str
        """
        return self._scale_metric

    @scale_metric.setter
    def scale_metric(self, scale_metric):
        """Sets the scale_metric of this V1beta1ComponentExtensionSpec.

        ScaleMetric defines the scaling metric type watched by autoscaler possible values are concurrency, rps, cpu, memory. concurrency, rps are supported via Knative Pod Autoscaler(https://knative.dev/docs/serving/autoscaling/autoscaling-metrics).  # noqa: E501

        :param scale_metric: The scale_metric of this V1beta1ComponentExtensionSpec.  # noqa: E501
        :type: str
        """

        self._scale_metric = scale_metric

    @property
    def scale_metric_type(self):
        """Gets the scale_metric_type of this V1beta1ComponentExtensionSpec.  # noqa: E501

        Type of metric to use. Options are Utilization, or AverageValue.  # noqa: E501

        :return: The scale_metric_type of this V1beta1ComponentExtensionSpec.  # noqa: E501
        :rtype: str
        """
        return self._scale_metric_type

    @scale_metric_type.setter
    def scale_metric_type(self, scale_metric_type):
        """Sets the scale_metric_type of this V1beta1ComponentExtensionSpec.

        Type of metric to use. Options are Utilization, or AverageValue.  # noqa: E501

        :param scale_metric_type: The scale_metric_type of this V1beta1ComponentExtensionSpec.  # noqa: E501
        :type: str
        """

        self._scale_metric_type = scale_metric_type

    @property
    def scale_target(self):
        """Gets the scale_target of this V1beta1ComponentExtensionSpec.  # noqa: E501

        ScaleTarget specifies the integer target value of the metric type the Autoscaler watches for. concurrency and rps targets are supported by Knative Pod Autoscaler (https://knative.dev/docs/serving/autoscaling/autoscaling-targets/).  # noqa: E501

        :return: The scale_target of this V1beta1ComponentExtensionSpec.  # noqa: E501
        :rtype: int
        """
        return self._scale_target

    @scale_target.setter
    def scale_target(self, scale_target):
        """Sets the scale_target of this V1beta1ComponentExtensionSpec.

        ScaleTarget specifies the integer target value of the metric type the Autoscaler watches for. concurrency and rps targets are supported by Knative Pod Autoscaler (https://knative.dev/docs/serving/autoscaling/autoscaling-targets/).  # noqa: E501

        :param scale_target: The scale_target of this V1beta1ComponentExtensionSpec.  # noqa: E501
        :type: int
        """

        self._scale_target = scale_target

    @property
    def scaler(self):
        """Gets the scaler of this V1beta1ComponentExtensionSpec.  # noqa: E501


        :return: The scaler of this V1beta1ComponentExtensionSpec.  # noqa: E501
        :rtype: V1beta1ScalerSpec
        """
        return self._scaler

    @scaler.setter
    def scaler(self, scaler):
        """Sets the scaler of this V1beta1ComponentExtensionSpec.


        :param scaler: The scaler of this V1beta1ComponentExtensionSpec.  # noqa: E501
        :type: V1beta1ScalerSpec
        """

        self._scaler = scaler

    @property
    def server_address(self):
        """Gets the server_address of this V1beta1ComponentExtensionSpec.  # noqa: E501

        Address of Prometheus server.  # noqa: E501

        :return: The server_address of this V1beta1ComponentExtensionSpec.  # noqa: E501
        :rtype: str
        """
        return self._server_address

    @server_address.setter
    def server_address(self, server_address):
        """Sets the server_address of this V1beta1ComponentExtensionSpec.

        Address of Prometheus server.  # noqa: E501

        :param server_address: The server_address of this V1beta1ComponentExtensionSpec.  # noqa: E501
        :type: str
        """

        self._server_address = server_address

    @property
    def timeout(self):
        """Gets the timeout of this V1beta1ComponentExtensionSpec.  # noqa: E501

        TimeoutSeconds specifies the number of seconds to wait before timing out a request to the component.  # noqa: E501

        :return: The timeout of this V1beta1ComponentExtensionSpec.  # noqa: E501
        :rtype: int
        """
        return self._timeout

    @timeout.setter
    def timeout(self, timeout):
        """Sets the timeout of this V1beta1ComponentExtensionSpec.

        TimeoutSeconds specifies the number of seconds to wait before timing out a request to the component.  # noqa: E501

        :param timeout: The timeout of this V1beta1ComponentExtensionSpec.  # noqa: E501
        :type: int
        """

        self._timeout = timeout

    def to_dict(self):
        """Returns the model properties as a dict"""
        result = {}

        for attr, _ in six.iteritems(self.openapi_types):
            value = getattr(self, attr)
            if isinstance(value, list):
                result[attr] = list(map(
                    lambda x: x.to_dict() if hasattr(x, "to_dict") else x,
                    value
                ))
            elif hasattr(value, "to_dict"):
                result[attr] = value.to_dict()
            elif isinstance(value, dict):
                result[attr] = dict(map(
                    lambda item: (item[0], item[1].to_dict())
                    if hasattr(item[1], "to_dict") else item,
                    value.items()
                ))
            else:
                result[attr] = value

        return result

    def to_str(self):
        """Returns the string representation of the model"""
        return pprint.pformat(self.to_dict())

    def __repr__(self):
        """For `print` and `pprint`"""
        return self.to_str()

    def __eq__(self, other):
        """Returns true if both objects are equal"""
        if not isinstance(other, V1beta1ComponentExtensionSpec):
            return False

        return self.to_dict() == other.to_dict()

    def __ne__(self, other):
        """Returns true if both objects are not equal"""
        if not isinstance(other, V1beta1ComponentExtensionSpec):
            return True

        return self.to_dict() != other.to_dict()
