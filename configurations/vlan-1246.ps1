
#$vdswitchName = "vcs8e-vcs-ci-workload-private"
# $pgrolename = "openshift_portgroup"
# $principal = "VSPHERE.LOCAL\openshift"

$pgrole = get-virole $pgrolename


$vdswitch = Get-VDSwitch -Name $vdswitchName
$notes = "vlan: 1246 gateway: 10.94.72.129 cidr: 25 mask: 255.255.255.128"
$newpg = New-VDPortgroup -Name "ci-vlan-1246" -Notes $notes -VDSwitch $vdswitch -VLanId 1246

New-VIPermission -Entity $newpg -Principal $principal -Role $pgrole -Propagate $True}
