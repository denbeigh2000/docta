**NOTICE**: I am probably not going to do much more "meaningful" developemnt this on this due to other ongoing projects, but it can be kinda useful as-is so I am open-sourcing it. If you want to have some simple, sane health checks along with your typical HTTP endpoint, this may be useful to you.

# docta

docta is a simple health checker service. It aims to provide checks that you
probably need (HTTP, file, CPU, memory, etc.) out-of-the-box, and produce output
compatible with existing monitoring solutions with meaningful information.

It is lightweight, and designed to sit on the same host as your application as a
health check endpoint.
