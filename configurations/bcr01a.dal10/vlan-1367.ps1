
#$vdswitchName = "vcs8e-vcs-ci-workload-private"
# $pgrolename = "openshift_portgroup"
# $principal = "VSPHERE.LOCAL\openshift"

$pgrole = get-virole $pgrolename


$vdswitch = Get-VDSwitch -Name $vdswitchName
$notes = "vlan: 1367 gateway: 10.5.203.1 cidr: 24 mask: 255.255.255.0"
$newpg = New-VDPortgroup -Name "ci-vlan-1367" -Notes $notes -VDSwitch $vdswitch -VLanId 1367

New-VIPermission -Entity $newpg -Principal $principal -Role $pgrole -Propagate $True
