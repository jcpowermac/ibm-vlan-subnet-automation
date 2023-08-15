# For leases >= 200, run on the IBM Cloud vSphere 8 env
if [ ${LEASE_NUMBER} -ge 200 ] && [ ${LEASE_NUMBER} -lt 221 ]; then
  echo Scheduling job on IBM Cloud instance
  VCENTER_AUTH_PATH=/var/run/vault/vsphere8-secrets/secrets.sh
  vsphere_url="vcenter.ibmc.devcluster.openshift.com"
  vsphere_datacenter="IBMCdatacenter"
  cloud_where_run="IBM8"
  dns_server="192.168.${LEASE_NUMBER}.1"
  vsphere_resource_pool="/IBMCdatacenter/host/IBMCcluster/Resources/ipi-ci-clusters"
  vsphere_cluster="IBMCcluster"
  vsphere_datastore="vsanDatastore"
fi