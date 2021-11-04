# Specify the provider and access details
#provider "ksyun" {
#  region = "cn-beijing-6"
#}

resource "ksyun_ssh_key" "default1" {
  key_name="ssh_key_tf"
  public_key="ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEA3gI6clIzRTFwl6ASM+Phr4PpedpYhrGBIUJUw7Qg3cLy8ZxGosgnx/QzIcromLRvEUc1TTKJgBRxkGvbM4bqO6X5ju0QizosRLADHU5Ive+cmDRDRwnAn3d2jqurAwjCfhVsutQRXVU5qgCPYfM7wEvMqkkMtYDZ87eM60amwSsmnSFfR97lWKPMZRx7QKZOwKK5cimNEAOvZo9SWGyRyCFcY2glvpT5F6ZWfyyIb0w0z0DRDgMGk0g3iOmEjVHTqxFN/aVXam74W4ExrbM88rq9ycaAd1nCCOGXQiHjRd54xHZODOJB3a3WNVbM9HRLOzC4R4S2CBs5Qp0vaWB7jw== ksyun@bjzjm01-op-bs-monitor007213.bjzjm01.ksyun.com"
}
