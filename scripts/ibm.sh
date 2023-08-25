echo Scheduling job on IBM Cloud instance

if [[ ${LEASED_RESOURCE} == *"vlan"* ]]; then
  if ! [ -z "${vsphere_url+x}" ]; then
    echo "vsphere_url is not defined, this should have been assigned in the step, exiting."
    exit 1
  fi

  case $vsphere_url in
  "vcs8e-vc.ocp2.dev.cluster.com")
    VCENTER_AUTH_PATH=/var/run/vault/ibmcloud/secrets.sh
    vsphere_url="vcs8e-vc.ocp2.dev.cluster.com"
    vsphere_datacenter="IBMCloud"
    cloud_where_run="IBM"

    vsphere_resource_pool="/IBMCloud/host/vcs-ci-workload/Resources"
    vsphere_cluster="vcs-ci-workload"
    vsphere_datastore="vsanDatastore"
    ;;

  "v8c-2-vcenter.ocp2.dev.cluster.com")
    VCENTER_AUTH_PATH=/var/run/vault/ibmcloud-2/secrets.sh
    vsphere_url="v8c-2-vcenter.ocp2.dev.cluster.com"
    vsphere_datacenter="IBMCloud"
    cloud_where_run="IBM"
    dns_server="10.38.76.172"
    vsphere_resource_pool="/IBMCloud/host/vcs-ci-workload/Resources"
    vsphere_cluster="vcs-ci-workload"
    vsphere_datastore="vsanDatastore"
    ;;

  "vcenter.ibmc.devcluster.openshift.com")
    VCENTER_AUTH_PATH=/var/run/vault/vsphere8-secrets/secrets.sh
    vsphere_url="vcenter.ibmc.devcluster.openshift.com"
    vsphere_datacenter="IBMCdatacenter"
    cloud_where_run="IBM8"
    dns_server="192.168.${LEASE_NUMBER}.1"
    vsphere_resource_pool="/IBMCdatacenter/host/IBMCcluster/Resources/ipi-ci-clusters"
    vsphere_cluster="IBMCcluster"
    vsphere_datastore="vsanDatastore"
    ;;

  "vcenter.devqe.ibmc.devcluster.openshift.com")
    VCENTER_AUTH_PATH=/var/run/vault/devqe-secrets/secrets.sh
    vsphere_url="vcenter.devqe.ibmc.devcluster.openshift.com"
    vsphere_datacenter="DEVQEdatacenter"
    cloud_where_run="IBMC-DEVQE"
    dns_server="192.168.${LEASE_NUMBER}.1"
    vsphere_resource_pool="/DEVQEdatacenter/host/DEVQEcluster/Resources/ipi-ci-clusters"
    vsphere_cluster="DEVQEcluster"
    vsphere_datastore="vsanDatastore"
    ;;

  *)
    echo "vsphere_url: ${vsphere_url} is not configured, exiting."
    exit 1
    ;;
  esac
else
  if [ "${LEASE_NUMBER}" -ge 88 ] && [ "${LEASE_NUMBER}" -lt 130 ]; then
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
  if [ "${LEASE_NUMBER}" -ge 130 ] && [ "${LEASE_NUMBER}" -lt 140 ]; then
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
  if [ "${LEASE_NUMBER}" -ge 151 ] && [ "${LEASE_NUMBER}" -lt 160 ]; then
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

  # For leases >= 200, run on the IBM Cloud vSphere 8 env
  if [ "${LEASE_NUMBER}" -ge 200 ] && [ "${LEASE_NUMBER}" -lt 221 ]; then
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

  # For leases >= 221, run on the IBM Cloud vSphere env
  if [ "${LEASE_NUMBER}" -ge 221 ]; then
    echo Scheduling job on IBM Cloud instance
    VCENTER_AUTH_PATH=/var/run/vault/devqe-secrets/secrets.sh
    vsphere_url="vcenter.devqe.ibmc.devcluster.openshift.com"
    vsphere_datacenter="DEVQEdatacenter"
    cloud_where_run="IBMC-DEVQE"
    dns_server="192.168.${LEASE_NUMBER}.1"
    vsphere_resource_pool="/DEVQEdatacenter/host/DEVQEcluster/Resources/ipi-ci-clusters"
    vsphere_cluster="DEVQEcluster"
    vsphere_datastore="vsanDatastore"
  fi
fi
