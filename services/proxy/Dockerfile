# Copyright 2022 Daniel Stamer

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#     http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM gcr.io/endpoints-release/endpoints-runtime-serverless:2.37.0
USER root
# Disable tracing
ENV ESPv2_ARGS "^++^--disable_tracing"
# Bump HTTP timeouts to an hour to account for sockets
ENV ESPv2_ARGS "${ESPv2_ARGS}++--http_request_timeout_s=3600"
# Set config
ENV ENDPOINTS_SERVICE_PATH /etc/endpoints/service.json
COPY service.json /etc/endpoints/service.json
RUN chown -R envoy:envoy /etc/endpoints/service.json && chmod -R 755 /etc/endpoints/service.json

USER envoy
ENTRYPOINT ["/env_start_proxy.py"]

# more config: https://cloud.google.com/endpoints/docs/openapi/specify-esp-v2-startup-options