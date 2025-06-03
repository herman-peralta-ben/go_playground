# Simple Golang example for a REST API with cert pinning

1. Generate certificates
```bash
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout key.pem -out cert.pem -config cert.conf
```
* key.pem -> Private, only server
* cert.pem -> Public


2. Encode the generated certificate into Base64
```bash
openssl x509 -in cert.pem -outform der | openssl dgst -sha256 -binary | openssl base64
```

3. Configure call
* Use above base64
* Dart: https://localhost:8443
* Android Emu: https://10.0.2.2:8443

4. Run
```bash
go run mock_pinning.go
```

## Optionally, instead generate base64 for self signed cert:
```bash
openssl req -x509 -newkey rsa:2048 -sha256 -nodes -keyout key.pem -out cert.pem -days 365 -subj "/CN=localhost"
openssl x509 -in cert.pem -noout -pubkey | openssl pkey -pubin -outform DER | openssl dgst -sha256 -binary | openssl enc -base64
```
