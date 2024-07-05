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


class V1beta1ScalerSpec(object):
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
        'max_replicas': 'int',
        'min_replicas': 'int',
        'query': 'str',
        'query_parameters': 'str',
        'scale_metric': 'str',
        'scale_metric_type': 'str',
        'scale_target': 'int',
        'server_address': 'str'
    }

    attribute_map = {
        'max_replicas': 'maxReplicas',
        'min_replicas': 'minReplicas',
        'query': 'query',
        'query_parameters': 'queryParameters',
        'scale_metric': 'scaleMetric',
        'scale_metric_type': 'scaleMetricType',
        'scale_target': 'scaleTarget',
        'server_address': 'serverAddress'
    }

    def __init__(self, max_replicas=None, min_replicas=None, query=None, query_parameters=None, scale_metric=None, scale_metric_type=None, scale_target=None, server_address=None, local_vars_configuration=None):  # noqa: E501
        """V1beta1ScalerSpec - a model defined in OpenAPI"""  # noqa: E501
        if local_vars_configuration is None:
            local_vars_configuration = Configuration()
        self.local_vars_configuration = local_vars_configuration

        self._max_replicas = None
        self._min_replicas = None
        self._query = None
        self._query_parameters = None
        self._scale_metric = None
        self._scale_metric_type = None
        self._scale_target = None
        self._server_address = None
        self.discriminator = None

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
        if server_address is not None:
            self.server_address = server_address

    @property
    def max_replicas(self):
        """Gets the max_replicas of this V1beta1ScalerSpec.  # noqa: E501

        Maximum number of replicas for autoscaling.  # noqa: E501

        :return: The max_replicas of this V1beta1ScalerSpec.  # noqa: E501
        :rtype: int
        """
        return self._max_replicas

    @max_replicas.setter
    def max_replicas(self, max_replicas):
        """Sets the max_replicas of this V1beta1ScalerSpec.

        Maximum number of replicas for autoscaling.  # noqa: E501

        :param max_replicas: The max_replicas of this V1beta1ScalerSpec.  # noqa: E501
        :type: int
        """

        self._max_replicas = max_replicas

    @property
    def min_replicas(self):
        """Gets the min_replicas of this V1beta1ScalerSpec.  # noqa: E501

        Minimum number of replicas, defaults to 1 but can be set to 0 to enable scale-to-zero.  # noqa: E501

        :return: The min_replicas of this V1beta1ScalerSpec.  # noqa: E501
        :rtype: int
        """
        return self._min_replicas

    @min_replicas.setter
    def min_replicas(self, min_replicas):
        """Sets the min_replicas of this V1beta1ScalerSpec.

        Minimum number of replicas, defaults to 1 but can be set to 0 to enable scale-to-zero.  # noqa: E501

        :param min_replicas: The min_replicas of this V1beta1ScalerSpec.  # noqa: E501
        :type: int
        """

        self._min_replicas = min_replicas

    @property
    def query(self):
        """Gets the query of this V1beta1ScalerSpec.  # noqa: E501

        Query to run to get metrics from Prometheus  # noqa: E501

        :return: The query of this V1beta1ScalerSpec.  # noqa: E501
        :rtype: str
        """
        return self._query

    @query.setter
    def query(self, query):
        """Sets the query of this V1beta1ScalerSpec.

        Query to run to get metrics from Prometheus  # noqa: E501

        :param query: The query of this V1beta1ScalerSpec.  # noqa: E501
        :type: str
        """

        self._query = query

    @property
    def query_parameters(self):
        """Gets the query_parameters of this V1beta1ScalerSpec.  # noqa: E501

        A comma-separated list of query Parameters to include while querying the Prometheus endpoint.  # noqa: E501

        :return: The query_parameters of this V1beta1ScalerSpec.  # noqa: E501
        :rtype: str
        """
        return self._query_parameters

    @query_parameters.setter
    def query_parameters(self, query_parameters):
        """Sets the query_parameters of this V1beta1ScalerSpec.

        A comma-separated list of query Parameters to include while querying the Prometheus endpoint.  # noqa: E501

        :param query_parameters: The query_parameters of this V1beta1ScalerSpec.  # noqa: E501
        :type: str
        """

        self._query_parameters = query_parameters

    @property
    def scale_metric(self):
        """Gets the scale_metric of this V1beta1ScalerSpec.  # noqa: E501

        ScaleMetric defines the scaling metric type watched by autoscaler possible values are concurrency, rps, cpu, memory. concurrency, rps are supported via Knative Pod Autoscaler(https://knative.dev/docs/serving/autoscaling/autoscaling-metrics).  # noqa: E501

        :return: The scale_metric of this V1beta1ScalerSpec.  # noqa: E501
        :rtype: str
        """
        return self._scale_metric

    @scale_metric.setter
    def scale_metric(self, scale_metric):
        """Sets the scale_metric of this V1beta1ScalerSpec.

        ScaleMetric defines the scaling metric type watched by autoscaler possible values are concurrency, rps, cpu, memory. concurrency, rps are supported via Knative Pod Autoscaler(https://knative.dev/docs/serving/autoscaling/autoscaling-metrics).  # noqa: E501

        :param scale_metric: The scale_metric of this V1beta1ScalerSpec.  # noqa: E501
        :type: str
        """

        self._scale_metric = scale_metric

    @property
    def scale_metric_type(self):
        """Gets the scale_metric_type of this V1beta1ScalerSpec.  # noqa: E501

        Type of metric to use. Options are Utilization, or AverageValue.  # noqa: E501

        :return: The scale_metric_type of this V1beta1ScalerSpec.  # noqa: E501
        :rtype: str
        """
        return self._scale_metric_type

    @scale_metric_type.setter
    def scale_metric_type(self, scale_metric_type):
        """Sets the scale_metric_type of this V1beta1ScalerSpec.

        Type of metric to use. Options are Utilization, or AverageValue.  # noqa: E501

        :param scale_metric_type: The scale_metric_type of this V1beta1ScalerSpec.  # noqa: E501
        :type: str
        """

        self._scale_metric_type = scale_metric_type

    @property
    def scale_target(self):
        """Gets the scale_target of this V1beta1ScalerSpec.  # noqa: E501

        ScaleTarget specifies the integer target value of the metric type the Autoscaler watches for. concurrency and rps targets are supported by Knative Pod Autoscaler (https://knative.dev/docs/serving/autoscaling/autoscaling-targets/).  # noqa: E501

        :return: The scale_target of this V1beta1ScalerSpec.  # noqa: E501
        :rtype: int
        """
        return self._scale_target

    @scale_target.setter
    def scale_target(self, scale_target):
        """Sets the scale_target of this V1beta1ScalerSpec.

        ScaleTarget specifies the integer target value of the metric type the Autoscaler watches for. concurrency and rps targets are supported by Knative Pod Autoscaler (https://knative.dev/docs/serving/autoscaling/autoscaling-targets/).  # noqa: E501

        :param scale_target: The scale_target of this V1beta1ScalerSpec.  # noqa: E501
        :type: int
        """

        self._scale_target = scale_target

    @property
    def server_address(self):
        """Gets the server_address of this V1beta1ScalerSpec.  # noqa: E501

        Address of Prometheus server.  # noqa: E501

        :return: The server_address of this V1beta1ScalerSpec.  # noqa: E501
        :rtype: str
        """
        return self._server_address

    @server_address.setter
    def server_address(self, server_address):
        """Sets the server_address of this V1beta1ScalerSpec.

        Address of Prometheus server.  # noqa: E501

        :param server_address: The server_address of this V1beta1ScalerSpec.  # noqa: E501
        :type: str
        """

        self._server_address = server_address

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
        if not isinstance(other, V1beta1ScalerSpec):
            return False

        return self.to_dict() == other.to_dict()

    def __ne__(self, other):
        """Returns true if both objects are not equal"""
        if not isinstance(other, V1beta1ScalerSpec):
            return True

        return self.to_dict() != other.to_dict()
