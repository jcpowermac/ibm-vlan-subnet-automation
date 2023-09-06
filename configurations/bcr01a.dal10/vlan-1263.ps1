
#$vdswitchName = "vcs8e-vcs-ci-workload-private"
# $pgrolename = "openshift_portgroup"
# $principal = "VSPHERE.LOCAL\openshift"

$pgrole = get-virole $pgrolename


$vdswitch = Get-VDSwitch -Name $vdswitchName
$notes = "vlan: 1263 gateway: 10.177.145.129 cidr: 25 mask: 255.255.255.128"
$newpg = New-VDPortgroup -Name "ci-vlan-1263" -Notes $notes -VDSwitch $vdswitch -VLanId 1263

New-VIPermission -Entity $newpg -Principal $principal -Role $pgrole -Propagate $True
