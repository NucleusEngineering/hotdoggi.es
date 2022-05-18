import base64
import os
import json

from flask import Flask, request

from cloudevents.http import CloudEvent, from_json, to_json

from google.cloud import storage

from opencensus.ext.stackdriver import trace_exporter as stackdriver_exporter
from opencensus.trace.propagation import trace_context_http_header_format
import opencensus.trace.tracer

app = Flask(__name__)

@app.route("/", methods=["POST"])
def index():
    event = unwrap(request)
    tracer = pickup_trace(event["traceparent"])
    with tracer.span(name="hotdoggies-archiver"):
        identifier = event["id"]
        type_name = event["type"]
        print(f"HOTDOGGIES processing event: {identifier}")

        with tracer.span(name="hotdoggies-archiver.write"):
            client = storage.Client()
            bucket = client.bucket(os.environ["HOTDOGGIES_ARCHIVAL_BUCKET"])
            blob = bucket.blob(f"{type_name}/{identifier}")
            blob.upload_from_string(data=to_json(event), content_type="application/json")
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


def pickup_trace(traceparent):
    exporter = stackdriver_exporter.StackdriverExporter(
        project_id=os.environ["GOOGLE_CLOUD_PROJECT"]
    )
    propagator = trace_context_http_header_format.TraceContextPropagator()
    headers = {
        "traceparent": traceparent
    }
    span_context = propagator.from_headers(headers)
    tracer = opencensus.trace.tracer.Tracer(
        span_context=span_context,
        exporter=exporter
    )

    print(f"HOTDOGGIES picked up trace: {tracer.span_context.trace_id}")
    print(f"HOTDOGGIES picked up span: {tracer.span_context.span_id}")

    return tracer


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
    app.run(debug=True, host="0.0.0.0", port=int(os.environ.get("PORT", 8080)))