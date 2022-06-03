FROM gcr.io/endpoints-release/endpoints-runtime-serverless:2
USER root
ENV ENDPOINTS_SERVICE_PATH /etc/endpoints/service.json
COPY service.json /etc/endpoints/service.json
RUN chown -R envoy:envoy /etc/endpoints/service.json && chmod -R 755 /etc/endpoints/service.json
USER envoy
ENTRYPOINT ["/env_start_proxy.py"]

# more config: https://cloud.google.com/endpoints/docs/openapi/specify-esp-v2-startup-options