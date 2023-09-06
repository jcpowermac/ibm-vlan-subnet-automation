
#$vdswitchName = "vcs8e-vcs-ci-workload-private"
# $pgrolename = "openshift_portgroup"
# $principal = "VSPHERE.LOCAL\openshift"

$pgrole = get-virole $pgrolename


$vdswitch = Get-VDSwitch -Name $vdswitchName
$notes = "vlan: 1364 gateway: 10.23.169.1 cidr: 24 mask: 255.255.255.0"
$newpg = New-VDPortgroup -Name "ci-vlan-1364" -Notes $notes -VDSwitch $vdswitch -VLanId 1364

New-VIPermission -Entity $newpg -Principal $principal -Role $pgrole -Propagate $True
