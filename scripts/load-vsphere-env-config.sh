source /var/run/vault/vsphere-config/ibm.sh
source /var/run/vault/vsphere-config/ibm-8.sh


if [ -z ${vsphere_url:-} ]; then
    echo "$(date -u --rfc-3339=seconds) - lease $LEASE_NUMBER does not have a valid definition"
    exit 1
fi
echo "$(date -u --rfc-3339=seconds) - selected vCenter ${vsphere_url}"
