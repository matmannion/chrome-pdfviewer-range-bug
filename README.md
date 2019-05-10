# Chrome bug reproduction - PDF viewer doesn't send SameSite cookies in Range requests
MRE for Chrome bug [961617](https://crbug.com/961617)

To run, use `go run app.go` and visit http://localhost:8000/

The image and video display correctly (with both cookies being read) but for the PDF, range requests
don't include the SameSite=Lax cookie, which causes the proof-of-concept to return a 303. The PDF appears
to be loading forever and every page past the first is just blank.

Logs display the request URI path, the Range header and the Cookie header; in the logs below we can see that
the Range request doesn't include the Auth1 header which has been set with SameSite=Lax:

```
➜  crbug go run app.go
http: 2019/05/10 10:45:22 Server is starting on http://localhost:8000/...
http: 2019/05/10 10:45:25 GET /  
http: 2019/05/10 10:45:29 GET /pdf.pdf  Auth1=auth1; Auth2=auth2
http: 2019/05/10 10:45:29 GET /pdf.pdf bytes=1900544-2542335 Auth2=auth2
```

Version 74.0.3729.131 (Official Build) (64-bit)
