import base64
import os
import io
import json

from flask import Flask, request

from cloudevents.http import CloudEvent, from_json, to_json

from google.cloud import bigquery
from google.cloud.exceptions import NotFound

from opencensus.ext.stackdriver import trace_exporter as stackdriver_exporter
from opencensus.trace.propagation import trace_context_http_header_format
import opencensus.trace.tracer

app = Flask(__name__)

@app.route("/", methods=["POST"])
def index():
    event = unwrap(request)
    tracer = pickup_trace(event["traceparent"])
    with tracer.span(name="analytics.handler.event"):
        identifier = event["id"]
        type_name = event["type"]
        print(f"processing event: {identifier}")

        client = bigquery.Client()

        project_id = os.environ["GOOGLE_CLOUD_PROJECT"]
        dataset_name = os.environ["BQ_DATASET"]
        table_name = type_name.replace(".", "_")

        try:
            with tracer.span(name="analytics.check"):
                client.get_table(f"{project_id}.{dataset_name}.{table_name}")
        except NotFound:
            # Job insertion with schema auto detection
            with tracer.span(name="analytics.load"):
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
        with tracer.span(name="analytics.insert"):
            print(f"streaming insert: {identifier}:{table_name}")
            rows = [
                json.load(io.BytesIO(to_json(event)))
            ]
            errors = client.insert_rows_json(f"{project_id}.{dataset_name}.{table_name}", rows)
            if errors != []:
                print("insert errors: {}".format(errors))

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

    print(f"picked up trace: {tracer.span_context.trace_id}")
    print(f"picked up span: {tracer.span_context.span_id}")

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