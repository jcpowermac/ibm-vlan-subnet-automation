# datacenter to create vlan in
datacenter = "dal10"


# Subnet Size
subnet_capacity = 128

# How many vlans to create with the above subnet size
vlan_quantity = 1


# The datacenter pod that the vlan will be created in.

# look for existing vlan...
# ibmcloud sl vlan list
# ibmcloud sl vlan detail 3136542
# remove .softlayer.com
# or ibmcloud sl vlan options
router = "bcr03a.dal10"

vlan_tags = ["openshift","ci"]
