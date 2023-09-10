# angeltrax
Tooling around the AngelTrax camera system.

## Development

### Reverse Engineering
Capture all the traffic related to your account:

```
sudo tcpdump -i any host ${your-site}.pro8cms.com -s0 -w/tmp/pro8.pcap
```

Get a feel for the HTTP traffic inside a PCAP file:

```
tshark -r /tmp/pro8.pcap -O http -T fields -e frame.number -e http.request.method -e http.response_in -e http.request.full_uri -e http.request.line -e http.response.code.desc -e http.request_in -e http.response.line -e http.file_data | sed -e 's/\\r\\n,/\n/g' -e 's/\\r\\n/\n/g' -e 's/\\n/\n/g' | less
```
