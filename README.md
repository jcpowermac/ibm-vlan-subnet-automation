
```
cat subnets.json | vault kv patch -address https://vault.ci.openshift.org kv/selfservice/vsphere-vmc/config subnets.json=-
```
