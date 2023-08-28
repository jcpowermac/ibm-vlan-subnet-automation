
#$vdswitchName = "vcs8e-vcs-ci-workload-private"
# $pgrolename = "openshift_portgroup"
# $principal = "VSPHERE.LOCAL\openshift"

$pgrole = get-virole $pgrolename


$vdswitch = Get-VDSwitch -Name $vdswitchName
$notes = "vlan: 1232 gateway: 10.38.192.129 cidr: 25 mask: 255.255.255.128"
$newpg = New-VDPortgroup -Name "ci-vlan-1232" -Notes $notes -VDSwitch $vdswitch -VLanId 1232

New-VIPermission -Entity $newpg -Principal $principal -Role $pgrole -Propagate $True}
