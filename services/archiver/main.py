# Copyright 2022 Google

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
import json

from flask import Flask, request

from cloudevents.http import to_json
from cloudevents.http.conversion import from_dict

from google.cloud import storage

from opentelemetry import trace
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.sampling import StaticSampler, TraceIdRatioBased, Decision
from opentelemetry.exporter.cloud_trace import CloudTraceSpanExporter
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.trace.propagation.tracecontext import TraceContextTextMapPropagator

# Global constants for service identification
PREFIX_IDENTIFIER = "es.hotdoggi"
SERVICE_NAME = "archiver"

def create_tracer():
    """ Create open telemetry tracer and configure sample exporter """
    tracer_provider = TracerProvider()
    cloud_trace_exporter = CloudTraceSpanExporter()
    tracer_provider.add_span_processor(
        BatchSpanProcessor(cloud_trace_exporter)
    )
    trace.set_tracer_provider(tracer_provider)
    tracer = trace.get_tracer(f"${PREFIX_IDENTIFIER}.service.${SERVICE_NAME}/")

    return tracer


app = Flask(__name__)
tracer = create_tracer()

@app.route("/v1/events/", methods=["POST"])
def index():
    """ Single POST endpoint for receiving events """
    event = unwrap(request)
    
    # Explicitly override context from original event trace
    ctx = parent_context(event["traceparent"])

    # Begin event span
    with tracer.start_as_current_span("archiver.handler:event", context=ctx):
        identifier = event["id"]
        type_name = event["type"]
        print(f"processing event: {identifier}")

        # Begin data write span
        with tracer.start_as_current_span("archiver.data:write"):
            # Create GCS client
            client = storage.Client()
            bucket = client.bucket(os.environ["ARCHIVAL_BUCKET"])
            # Set object name
            blob = bucket.blob(f"{type_name}/{identifier}")
            # Set blob to deserialize event
            blob.upload_from_string(data=to_json(event), content_type="application/json")
            # Copy event envelope as object metadata
            blob.metadata = {
                "id": identifier,
                "type": type_name,
                "source": event["source"],
                "specversion": event["specversion"],
                "traceparent": event["traceparent"],
                "time": event["time"],
                "subject": event["subject"],
                "datacontenttype": event["datacontenttype"]
            }
            blob.patch()

    return ("", 204)


def parent_context(traceparent):
    """ Rehydrate trace context propagated and embedded in event """
    carrier = {'traceparent': traceparent}
    ctx = TraceContextTextMapPropagator().extract(carrier=carrier)
    
    print(f"picked up trace: {ctx}")

    return ctx


def unwrap(request):
    """ Unwrap and deserialize CloudEvent from Pub/Sub message """
    envelope = request.get_json()
    print(f"received: {json.dumps(envelope)}")

    # Check for Pub/Sub message format
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
        # Check are OK, decode Pub/Sub message payload
        data = base64.b64decode(pubsub_message["data"]).decode("utf-8").strip()

    event_dict = json.loads(data)

    # Strip empty values from dict
    event_dict = omit_empty(event_dict)

    # Deserialize into CloudEvent
    event = from_dict(event_dict)

    return event


def omit_empty(dict_map):
    """ Recursively drop empty values from dict """
    if type(dict_map) is dict:
        return dict((key, omit_empty(value)) for key, value in dict_map.items() if value and omit_empty(value))
    else:
        return dict_map


if __name__ == "__main__":
    debug = False
    if os.getenv("ENVIRONMENT") == "dev":
        debug = True
    app.run(debug=debug, host="0.0.0.0", port=int(os.environ.get("PORT", 8080)))
