## About kron

- [Building](https://github.com/a-palchikov/kron/wiki/Building)
- [Usage](https://github.com/a-palchikov/kron/wiki/Usage)
- [Command reference](https://github.com/a-palchikov/kron/wiki/CommandLine)
- [Design](https://github.com/a-palchikov/kron/wiki/Design-overview)

## Dependencies

- [libzmq](https://github.com/zeromq/libzmq)

kron is a simple distributed job scheduler.

kron will eventually:
- [ ] be able to schedule jobs across a cluster of nodes
- [ ] be resilient to node failures by employing a consesus algorithm
- [ ] be able to spawn jobs with resource restrictions and in isolation
- [ ] have an intuitive gRPC-based API
