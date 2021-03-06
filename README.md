# Golang Sensu plugins

> Buildable assets for Sensu 5.0

Sensu [asset](https://docs.sensu.io/sensu-go/5.0/reference/assets/) is a self-contained executable that Sensu is able to distribute to the agents who need it and then use it in checks, handlers etc. Old Ruby-based Sensu checks can't be used, because they require Ruby and shitload of gems to work properly.

The preferred way to implement checks now is to use Golang, because it fullfills all the requirements from the previous paragraph.

## Checks

- [x] Load average - `check-load`
- [x] Memory - `check-memory`
- [x] Filesystem usage `check-fs`
- [x] generic SystemD unit check - `check-systemd-unit`
  - note that it requires `root` privileges to run, hence [#2315](https://github.com/sensu/sensu-go/issues/2315) has to be closed first
- [x] TCP check - `check-tcp`
- [x] Kernel process (aka `ps ...`) - `check-ps`
- [ ] S.M.A.R.T. - `check-smart` _started_
- [ ] mdamd consistency check - `check-mdraid`
- [x] HTTP check - `check-http` _just return codes so far_
- [ ] ICMP check - `check-icmp`
- [ ] Prometheus metric - `check-promql`
- [ ] Docker container status - `check-docker-process`
- [x] Consul/Nomad -`check-consul-service`

## Handlers

- [ ] [Slack](https://github.com/sensu/sensu-slack-handler/)

# Building

1. Clone this repo
2. Have `go-1.11` installed and `make` too
3. Run `make`
4. Use built `.tar.gz` with your `sensuctl asset create ...` command
5. Profit :)

# Notes

- Sensu guys have made [some libs](https://github.com/sensu-plugins/sensu-plugins-go/) to help with commonalities, those are included
- There's a [repo](https://github.com/hico-horiuchi/sensu-plugins-go/) with some old Golang plugins written for Sensu 1.x - it might be worth exploring

