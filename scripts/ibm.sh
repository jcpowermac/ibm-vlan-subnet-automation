# For leases >= than 88, run on the IBM Cloud


if [[ ${LEASED_RESOURCE} == *"vlan"* ]]; then
  # for testing we know there is only a single port group with vlan in the name
  echo Scheduling job on IBM Cloud instance
  VCENTER_AUTH_PATH=/var/run/vault/ibmcloud/secrets.sh
  vsphere_url="ibmvcenter.vmc-ci.devcluster.openshift.com"
  vsphere_datacenter="IBMCloud"
  cloud_where_run="IBM"
  dns_server="10.38.76.172"
  vsphere_resource_pool="/IBMCloud/host/vcs-ci-workload/Resources"
  vsphere_cluster="vcs-ci-workload"
  vsphere_datastore="vsanDatastore"
fi


if [ ${LEASE_NUMBER} -ge 88 ] && [ ${LEASE_NUMBER} -lt 130 ]; then
  echo Scheduling job on IBM Cloud instance
  VCENTER_AUTH_PATH=/var/run/vault/ibmcloud/secrets.sh
  vsphere_url="ibmvcenter.vmc-ci.devcluster.openshift.com"
  vsphere_datacenter="IBMCloud"
  cloud_where_run="IBM"
  dns_server="10.38.76.172"
  vsphere_resource_pool="/IBMCloud/host/vcs-ci-workload/Resources"
  vsphere_cluster="vcs-ci-workload"
  vsphere_datastore="vsanDatastore"
fi
# For leases >= than 130 and < 140, jobs are directed towards an overflow vSphere 7 vCenter
if [ ${LEASE_NUMBER} -ge 130 ] && [ ${LEASE_NUMBER} -lt 140 ]; then
  VCENTER_AUTH_PATH=/var/run/vault/ibmcloud-2/secrets.sh
  vsphere_url="v8c-2-vcenter.ocp2.dev.cluster.com"
  vsphere_datacenter="IBMCloud"
  cloud_where_run="IBM"
  dns_server="10.38.76.172"
  vsphere_resource_pool="/IBMCloud/host/vcs-ci-workload/Resources"
  vsphere_cluster="vcs-ci-workload"
  vsphere_datastore="vsanDatastore"
fi
# For leases >= than 151 and < 160, pull from alternate DNS which has context
# relevant to NSX-T segments.
if [ ${LEASE_NUMBER} -ge 151 ] && [ ${LEASE_NUMBER} -lt 160 ]; then
  echo Scheduling job on IBM Cloud instance
  VCENTER_AUTH_PATH=/var/run/vault/ibmcloud/secrets.sh
  vsphere_url="ibmvcenter.vmc-ci.devcluster.openshift.com"
  vsphere_datacenter="IBMCloud"
  cloud_where_run="IBM"
  vsphere_resource_pool="/IBMCloud/host/vcs-ci-workload/Resources"
  vsphere_cluster="vcs-ci-workload"
  vsphere_datastore="vsanDatastore"
  dns_server="192.168.133.73"
fi
