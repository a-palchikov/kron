![kron](https://github.com/a-palchikov/kron)

## About kron

- [Building](https://github.com/a-palchikov/kron/wiki/Building)
- [Usage](https://github.com/a-palchikov/kron/wiki/Usage)
- [Command reference](https://github.com/a-palchikov/kron/wiki/CommandLine)
- [Contributing](https://github.com/a-palchikov/kron/wiki/Contributing)

## Dependencies (in future)

- [cron expression parser](github.com/gorhill/cronexpr)
- [ectd](github.com/coreos/etcd)
- [Serf](github.com/hashicorp/serf)

kron is a simple distributed job scheduler, based on ideas from the paper "Reliable Cron across the Planet" (https://queue.acm.org/detail.cfm?id=2745840).

A distributed cron service's main responsibility is coping with potential failures and being able to reliably schedule and reschedule jobs irregardless of their idempotency.
The availability aspect of the system is achieved with a master/slave paradigm whereas a single server node is doing the scheduling while other nodes maintain certain redundancy and are ready to manage the scheduler queue should the current server node fail.

All the nodes participating in scheduling share the same state at all times. The state is kept in a naive persistent key/value store but can be upgraded to anything else (i.e. etcd) in the future.
The nodes are managed by a cluster membership service that handles failures which also can be offloaded to anything else (i.e. Serf).

