# Todo List

- create source repo on Github, public
- purchase domain
- add external, global HTTPS load balancer with CDN support and serverless NEG pointing into proxy
- move app service into build pipeline and render artifacts into bucket as CDN serving origin
- point POST:/events at ingest service, remove group:allUsers from run.Invokers for ingest (proxy/gateway SA exclusive)
- complete dogs-service implementation for GET:/dogs, GET:/dogs/{x} and event handler
- NOTE: event handler is not to be added to the OpenAPI spec
