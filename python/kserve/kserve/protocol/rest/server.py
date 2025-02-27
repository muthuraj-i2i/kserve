# Copyright 2022 The KServe Authors.
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

import os
import asyncio
import logging
from typing import Dict, Optional, Union

import uvicorn
from fastapi import FastAPI, Request, Response
from fastapi.routing import APIRouter
from prometheus_client import REGISTRY, exposition, Counter
from timing_asgi import TimingClient, TimingMiddleware
from timing_asgi.integrations import StarletteScopeToName

from kserve.errors import (
    InferenceError,
    InvalidInput,
    ModelNotFound,
    ModelNotReady,
    ServerNotLive,
    ServerNotReady,
    UnsupportedProtocol,
    generic_exception_handler,
    inference_error_handler,
    invalid_input_handler,
    model_not_found_handler,
    model_not_ready_handler,
    not_implemented_error_handler,
    server_not_live_handler,
    server_not_ready_handler,
    unsupported_protocol_error_handler,
)
from kserve.logging import trace_logger, logger
from kserve.protocol.dataplane import DataPlane

from .openai.config import maybe_register_openai_endpoints
from .v1_endpoints import register_v1_endpoints
from .v2_endpoints import register_v2_endpoints
from ..model_repository_extension import ModelRepositoryExtension

from opentelemetry import metrics
from opentelemetry.exporter.otlp.proto.http.metric_exporter import OTLPMetricExporter
from opentelemetry.sdk.metrics import MeterProvider
from opentelemetry.sdk.metrics.export import PeriodicExportingMetricReader
from opentelemetry.sdk.resources import Resource
from starlette.middleware.base import BaseHTTPMiddleware

# Configure OpenTelemetry
# OTEL_COLLECTOR_ENDPOINT = "http://localhost:4318/v1/metrics"  # Replace with your collector endpoint

# Read the endpoint from the environment variable
otel_collector_endpoint = os.getenv("OTEL_COLLECTOR_ENDPOINT", "localhost:4317")

logger.info(f"OTEL_COLLECTOR_ENDPOINT --------------------->: {otel_collector_endpoint}")

resource = Resource.create({"service.name": "sklearn-v2-iris-predictor"})
exporter = OTLPMetricExporter(endpoint=otel_collector_endpoint)
reader = PeriodicExportingMetricReader(exporter)
provider = MeterProvider(resource=resource, metric_readers=[reader])
metrics.set_meter_provider(provider)
meter = metrics.get_meter(__name__)

http_request_counter = meter.create_counter(
    "http_server_request_count",
    description="Counts the number of HTTP requests received",
)

# Define the counter for HTTP requests
pro_http_request_counter = Counter(
    "http_server_request_count",
    "Counts the number of HTTP requests received",
    ["method", "path"]
)

class RequestCountMiddleware(BaseHTTPMiddleware):
    async def dispatch(self, request: Request, call_next):
        # Increment the counter for each request
        pro_http_request_counter.labels(method=request.method, path=request.url.path).inc()
        http_request_counter.add(1, {"method": request.method, "path": request.url.path})
        response = await call_next(request)
        return response
    
# Hardcoded counter metric
counter = meter.create_counter("test_requests_total", unit="requests")
counter.add(10, {"environment": "test", "region": "us-east-1"}) #added labels here.

# Hardcoded gauge metric
gauge = meter.create_observable_gauge("test_cpu_usage", unit="percentage", callbacks=[lambda options: [(55.5, {"environment": "test", "region": "us-east-1"})]])

# Hardcoded histogram Metric
histogram = meter.create_histogram("test_request_latency", unit="ms")
histogram.record(100, {"environment": "test", "region": "us-east-1"})
histogram.record(250, {"environment": "test", "region": "us-east-1"})
histogram.record(500, {"environment": "test", "region": "us-east-1"})

# Hardcoded summary metric
summary_counter = meter.create_counter("test_summary_count", unit="requests")
summary_sum = meter.create_counter("test_summary_sum", unit="ms")
summary_counter.add(3, {"environment": "test", "region": "us-east-1"})
summary_sum.add(850, {"environment": "test", "region": "us-east-1"})
logger.info(f" pushing metrics to otel --------------------->: {otel_collector_endpoint}")




async def push_prometheus_to_otel():
    """Pushes Prometheus metrics to OpenTelemetry."""
    try:
        # Collect Prometheus metrics
        metric_families = REGISTRY.collect()

        # Convert Prometheus metrics to OpenTelemetry metrics
        for family in metric_families:
            for sample in family.samples:

                name = sample.name
                labels = sample.labels
                value = sample.value
                unit = family.unit

                if family.type == "counter":
                    counter = meter.create_counter(name, unit=unit)
                    bound_counter = counter.bind(labels)
                    bound_counter.add(int(value))

                elif family.type == "gauge":
                    gauge = meter.create_observable_gauge(name, unit=unit, callbacks=[lambda options: [(value, labels)]])

                elif family.type == "summary":
                    count_name = f"{name}_count"
                    sum_name = f"{name}_sum"

                    for sample2 in family.samples:
                      if sample2.name == f"{name}_count":
                        count_value = sample2.value
                        counter = meter.create_counter(count_name, unit=unit)
                        bound_counter = counter.bind(labels)
                        bound_counter.add(int(count_value))

                      if sample2.name == f"{name}_sum":
                        sum_value = sample2.value;
                        counter = meter.create_counter(sum_name, unit=unit)
                        bound_counter = counter.bind(labels)
                        bound_counter.add(sum_value)

                elif family.type == "histogram":
                    sum_name = f"{name}_sum"
                    count_name = f"{name}_count"

                    for sample2 in family.samples:
                      if sample2.name == f"{name}_sum":
                        sum_value = sample2.value;
                        counter = meter.create_counter(sum_name, unit=unit)
                        bound_counter = counter.bind(labels)
                        bound_counter.add(sum_value)
                      if sample2.name == f"{name}_count":
                        count_value = sample2.value;
                        counter = meter.create_counter(count_name, unit=unit)
                        bound_counter = counter.bind(labels)
                        bound_counter.add(count_value)

    except Exception as e:
        logger.error(f"Error converting Prometheus metrics to OpenTelemetry: {e}")
        print(f"Error pushing metrics to OpenTelemetry: {e}")

async def background_metrics_push():
    """Background task to periodically push metrics."""
    while True:
        await push_hardcoded_metrics()
        await asyncio.sleep(1)  # Push metrics every 15 seconds (adjust as needed)

async def metrics_handler(request: Request) -> Response:
    encoder, content_type = exposition.choose_encoder(request.headers.get("accept"))
    return Response(content=encoder(REGISTRY), headers={"content-type": content_type})


class PrintTimings(TimingClient):
    def timing(self, metric_name, timing, tags):
        trace_logger.info(f"{metric_name}: {timing} {tags}")


class RESTServer:
    def __init__(
        self,
        app: FastAPI,
        data_plane: DataPlane,
        model_repository_extension: ModelRepositoryExtension,
        http_port: int,
        log_config: Optional[Union[str, Dict]] = None,
        access_log_format: Optional[str] = None,
        workers: int = 1,
    ):
        self.app = app
        self.dataplane = data_plane
        self.model_repository_extension = model_repository_extension
        self.access_log_format = access_log_format
        self._server = uvicorn.Server(
            config=uvicorn.Config(
                app="kserve.model_server:app",
                host="0.0.0.0",
                log_config=log_config,
                port=http_port,
                workers=workers,
            )
        )

        # Add lifespan middleware to handle startup and shutdown events
        # self.app.add_middleware(LifespanMiddleware, on_startup=self.on_startup, on_shutdown=self.on_shutdown)
        self.app.add_middleware(RequestCountMiddleware)

    async def on_startup(self):
        # Start the background task
        self.app.state.background_task = asyncio.create_task(background_metrics_push())

    async def on_shutdown(self):
        # Cancel the background task
        self.app.state.background_task.cancel()
        await self.app.state.background_task

    def _register_endpoints(self):
        root_router = APIRouter()
        root_router.add_api_route(r"/", self.dataplane.live)
        root_router.add_api_route(r"/metrics", metrics_handler, methods=["GET"])
        self.app.include_router(root_router)
        register_v1_endpoints(self.app, self.dataplane, self.model_repository_extension)
        register_v2_endpoints(self.app, self.dataplane, self.model_repository_extension)
        # Register OpenAI endpoints if any of the models in the registry implement the OpenAI interface
        # This adds /openai/v1/completions and /openai/v1/chat/completions routes to the
        # REST server.
        maybe_register_openai_endpoints(self.app, self.dataplane.model_registry)

    def _add_exception_handlers(self):
        self.app.add_exception_handler(InvalidInput, invalid_input_handler)
        self.app.add_exception_handler(InferenceError, inference_error_handler)
        self.app.add_exception_handler(ModelNotFound, model_not_found_handler)
        self.app.add_exception_handler(ModelNotReady, model_not_ready_handler)
        self.app.add_exception_handler(
            NotImplementedError, not_implemented_error_handler
        )
        self.app.add_exception_handler(
            UnsupportedProtocol, unsupported_protocol_error_handler
        )
        self.app.add_exception_handler(ServerNotLive, server_not_live_handler)
        self.app.add_exception_handler(ServerNotReady, server_not_ready_handler)
        self.app.add_exception_handler(Exception, generic_exception_handler)

    def _add_middlewares(self):
        self.app.add_middleware(
            TimingMiddleware,
            client=PrintTimings(),
            metric_namer=StarletteScopeToName(
                prefix="kserve.io", starlette_app=self.app
            ),
        )

        # More context in https://github.com/encode/uvicorn/pull/947
        # At the time of writing the ASGI specs are not clear when it comes
        # to change the access log format, and hence the Uvicorn upstream devs
        # chose to create a custom middleware for this.
        # The allowed log format is specified in https://github.com/Kludex/asgi-logger#usage
        if self.access_log_format:
            from asgi_logger import AccessLoggerMiddleware

            # As indicated by the asgi-logger docs, we need to clear/unset
            # any setting for uvicorn.access to avoid log duplicates.
            logging.getLogger("uvicorn.access").handlers = []
            self.app.add_middleware(
                AccessLoggerMiddleware, format=self.access_log_format
            )
            # The asgi-logger settings don't set propagate to False,
            # so we get duplicates if we don't set it explicitly.
            logging.getLogger("access").propagate = False

    def create_application(self):
        self._add_middlewares()
        self._register_endpoints()
        self._add_exception_handlers()

    async def start(self):
        self.create_application()
        logger.info(f"Starting otel scaler ----------------------->: {otel_collector_endpoint}")
        logger.info(f"Starting uvicorn with {self._server.config.workers} workers")
        await self._server.serve()
