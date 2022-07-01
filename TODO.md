# Todo List

- create loaders on CR jobs
- resume trace context in go services
- write README

- implement sync GET of `/dogs?stream` && `/dogs/{key}?stream`
    - return current materialized snapshot from Firestore
    - create listener connection https://firebase.google.com/docs/firestore/query-data/listen#listen_to_multiple_documents_in_a_collection
    - serve back concurrently