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

import base64
import os
import io
import json

from flask import Flask, request

from cloudevents.http import from_json, to_json

from google.cloud import bigquery
from google.cloud.exceptions import NotFound

from opentelemetry import trace
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.sampling import StaticSampler, TraceIdRatioBased, Decision
from opentelemetry.exporter.cloud_trace import CloudTraceSpanExporter
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.trace.propagation.tracecontext import TraceContextTextMapPropagator

PREFIX_IDENTIFIER = "es.hotdoggi"
SERVICE_NAME = "analytics"

app = Flask(__name__)

@app.route("/events/", methods=["POST"])
def index():
    event = unwrap(request)
    tracer = create_tracer()
    ctx = parent_context(event["traceparent"])

    with tracer.start_as_current_span("analytics.handler:event", context=ctx):
        identifier = event["id"]
        type_name = event["type"]
        print(f"processing event: {identifier}")

        client = bigquery.Client()

        project_id = os.environ["GOOGLE_CLOUD_PROJECT"]
        dataset_name = os.environ["BQ_DATASET"]
        table_name = type_name.replace(".", "_")

        try:
            with tracer.start_as_current_span("analytics.data:check", context=ctx):
                client.get_table(f"{project_id}.{dataset_name}.{table_name}")
        except NotFound:
            # Job insertion with schema auto detection
            with tracer.start_as_current_span("analytics.data:load", context=ctx):
                print(f"HOTDOGGIES loading job: {identifier}:{table_name}")

                dataset_ref = client.dataset(dataset_name)
                table_ref = dataset_ref.table(table_name)

                job_config = bigquery.LoadJobConfig()
                job_config.source_format = bigquery.SourceFormat.NEWLINE_DELIMITED_JSON
                job_config.autodetect = True

                source_file = io.BytesIO(to_json(event))

                job = client.load_table_from_file(
                    source_file,
                    table_ref,
                    location="EU",
                    job_config=job_config,
                )

                print(f"pushed job: {job.job_id}")
                return ("", 204)

        # Stream insertion
        with tracer.start_as_current_span("analytics.data:insert", context=ctx):
            print(f"streaming insert: {identifier}:{table_name}")
            rows = [
                json.load(io.BytesIO(to_json(event)))
            ]
            errors = client.insert_rows_json(
                f"{project_id}.{dataset_name}.{table_name}", rows)
            if errors != []:
                print("insert errors: {}".format(errors))

    return ("", 204)


def create_tracer():
    # Default to PROD
    sampler = TraceIdRatioBased(0.1)
    if os.getenv("ENVIRONMENT") == "dev":
        # Always sample in DEV
        sampler = StaticSampler(Decision(True))

    tracer_provider = TracerProvider(sampler=sampler)
    cloud_trace_exporter = CloudTraceSpanExporter()
    tracer_provider.add_span_processor(
        BatchSpanProcessor(cloud_trace_exporter)
    )
    trace.set_tracer_provider(tracer_provider)
    tracer = trace.get_tracer(f"${PREFIX_IDENTIFIER}.service.${SERVICE_NAME}/")

    return tracer


def parent_context(traceparent):
    carrier = {'traceparent': traceparent}
    ctx = TraceContextTextMapPropagator().extract(carrier=carrier)
    
    print(f"picked up trace: {ctx}")

    return ctx


def unwrap(request):
    envelope = request.get_json()
    print(f"received: {json.dumps(envelope)}")
    if not envelope:
        msg = "no Pub/Sub message received"
        print(f"error: {msg}")
        return f"Bad Request: {msg}", 400

    if not isinstance(envelope, dict) or "message" not in envelope:
        msg = "invalid Pub/Sub message format"
        print(f"error: {msg}")
        return f"Bad Request: {msg}", 400

    pubsub_message = envelope["message"]

    data = None
    if isinstance(pubsub_message, dict) and "data" in pubsub_message:
        data = base64.b64decode(pubsub_message["data"]).decode("utf-8").strip()
    event = from_json(data)

    return event


if __name__ == "__main__":
    debug = False
    if os.getenv("ENVIRONMENT") == "dev":
        debug = True
    app.run(debug=debug, host="0.0.0.0", port=int(os.environ.get("PORT", 8080)))
