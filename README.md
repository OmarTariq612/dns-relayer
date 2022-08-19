
# dns relayer

`dns-relayer` is a program that relays dns packets.

## Build

```
go build .
```

## Usage
```
usage: dns-relayer [bind_addr:bind_port] server_addr:server_port
```
`bind_addr:bind_port` is `:53` by default.


* It can be used to analyze dns packets by relaying dns to a machine for analysis:
  ```
  dns client <--> dns-relayer <--> dns server
  ```
  ```
  ./dns-relayer 192.168.0.10:53 8.8.8.8:53
  ```
  ```
  2022/08/19 17:23:52 Serving on 192.168.0.10:53
  2022/08/19 17:23:52 the provided server address is 8.8.8.8:53



  ```

* You can also chain multiple relays:
  ```
  dns client <--> dns-relayer(1) <--> dns-relayer(2) <--> ... <--> dns server
  ```
  and this technique can be used to make use of a different port (other than 53) to prevent censorship firewalls from analyzing the dns packet and dropping it.

  | local | remote |
  |:------------------------------------------:|:-------------------------------------:|
  |`./dns-relayer 192.168.0.10:53 x.y.z.n:9090`|`./dns-relayer x.y.z.n:9090 8.8.8.8:53`|

  now change dns settings to forward dns packets to `192.168.0.10:53`.

## Notes
* Serving on all interfaces can conflict with the local resolver that runs on `localhost` on 53:
  ```
  2022/08/19 18:00:49 listen udp :53: bind: address already in use
  ```
  to prevent this issue use a specific interface addresss (ex `192.168.0.10`).
* Binding to ports below 1024 requires root priviliges.

* There are other techniques that can be used to bypass censorship firewalls from dropping dns packets:
  - DNS over HTTPS (DoH): https://en.wikipedia.org/wiki/DNS_over_HTTPS
  - DNSCrypt: https://en.wikipedia.org/wiki/DNSCrypt