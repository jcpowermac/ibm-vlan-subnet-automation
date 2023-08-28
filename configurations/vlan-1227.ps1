
#$vdswitchName = "vcs8e-vcs-ci-workload-private"
# $pgrolename = "openshift_portgroup"
# $principal = "VSPHERE.LOCAL\openshift"

$pgrole = get-virole $pgrolename


$vdswitch = Get-VDSwitch -Name $vdswitchName
$notes = "vlan: 1227 gateway: 10.38.252.1 cidr: 25 mask: 255.255.255.128"
$newpg = New-VDPortgroup -Name "ci-vlan-1227" -Notes $notes -VDSwitch $vdswitch -VLanId 1227

New-VIPermission -Entity $newpg -Principal $principal -Role $pgrole -Propagate $True}
