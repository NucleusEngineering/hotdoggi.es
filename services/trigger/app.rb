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

require 'json'
require 'functions_framework'
require 'opencensus-stackdriver'
require 'cloud_events'

PREFIX_IDENTIFIER = "es.hotdoggi"
SERVICE_NAME = "trigger"

FunctionsFramework.on_startup do
  # Create clients
  set_global :firestore_client do
    require 'google/cloud/firestore'
    Google::Cloud::Firestore.new project_id: ENV['GOOGLE_CLOUD_PROJECT']
  end
  set_global :pubsub_client do
    require 'google/cloud/pubsub'
    Google::Cloud::PubSub.new project_id: ENV['GOOGLE_CLOUD_PROJECT']
  end
  OpenCensus.configure do |c|
    c.trace.exporter = OpenCensus::Trace::Exporters::Stackdriver.new project_id: ENV['GOOGLE_CLOUD_PROJECT']
  end
end

FunctionsFramework.cloud_event 'function' do |fs_event|

  prefix = "#{PREFIX_IDENTIFIER}.services.#{SERVICE_NAME}/"

  # Pickup W3C trace context
  traceparent = fs_event.data['value']['fields']['traceparent']['stringValue']
  trace_context = OpenCensus::Trace::TraceContextData.new(
    traceparent.split('-')[1],
    traceparent.split('-')[2],
    traceparent.split('-')[3].to_i
  )
  trace = OpenCensus::Trace::SpanContext.create_root(trace_context: trace_context)

  trace.in_span "#{prefix}trigger.handler:event" do |_span|
    event_type = fs_event.subject.split('/')[-2]
    event_id = fs_event.subject.split('/')[-1]
    logger.info "detected change to: #{event_type}:#{event_id}"

    event = nil
    trace.in_span "#{prefix}trigger.data:load" do |_subspan|
      doc = global(:firestore_client).col(event_type).doc(event_id).get
      event = CloudEvents::Event.create(
        id: event_id,
        type: event_type,
        traceparent: traceparent,
        source: doc[:source],
        spec_version: doc[:specversion],
        data_content_type: doc[:datacontenttype],
        subject: doc[:subject],
        time: doc[:time],
        data: doc[:data]
      )
      logger.info "received event: #{event.to_h}"
    end

    trace.in_span "#{prefix}trigger.message:publish" do |_subspan|
      topic = global(:pubsub_client).topic ENV['TOPIC']
      result = topic.publish event.to_h.to_json,
        type: event.type,
        source: event.source
      logger.info "publish event message: #{result.message_id}"
    end
  end
end
