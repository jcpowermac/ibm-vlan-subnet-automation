

```bash
vault kv get -address https://vault.ci.openshift.org -field ibm.sh kv/selfservice/vsphere-vmc/config > ibm.sh
vault kv get -address https://vault.ci.openshift.org -field ibm-8.sh kv/selfservice/vsphere-vmc/config > ibm-8.sh
vault kv get -address https://vault.ci.openshift.org -field load-vsphere-env-config.sh kv/selfservice/vsphere-vmc/config > load-vsphere-env-config.sh
```

