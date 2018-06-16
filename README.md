# Container Invoker
This is a spike into running and invoking containers via an API.

## TODO
- Query params
- If a user shuts down the server, there needs to either be a 30 second wait, to allow the containers to expire and be terminated correctly. Or we would need to manually kill all the active containers on kill signal. Otherwise containers are left orphaned 😭.
- Add prometheus metrics, or some other kind of metric tracking. Or even a generic metric interface.
- Further abstract logging so that the invoker and the API don't necessarily know anything about the logging implementation.
- Add response caching, to avoid firing up containers, or calling containers for the same data.
- Create CLI invoker to test modules locally without running the api for example. `$ invoker invoker module-a` for example.
- Make timeout time configurable.
