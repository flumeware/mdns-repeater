mDNS Repeater
---

 - Repeat mDNS information between VLANs or otherwise across broadcast domains, whilst controlling what is shared.
 - Listens on multiple interfaces
 - Control which devices are available via the repeater

Name Matching
===

When a query is received every name in it is compared to the configured names, where the configured name is a substring a unicast query is sent to that target and the response it sent to the client. `*` can be used a name to always send responses from a target.

Example Config
====
This config will proxy the dns response from sending any query which includes `example_name` to `192.168.1.1:5353`

```
{
	"ifaces": ["eth0"],
	"single_targets":[
		{
			"names": ["example_name"],
			"ip": "192.168.1.1",
			"port": 5353
		}
	]
}
```

Multiple targets can cover a name, which will cause the query to be sent to all matching targets. Generally a target will cover both the Service Discovery names and its own specific name.