
#$vdswitchName = "vcs8e-vcs-ci-workload-private"
# $pgrolename = "openshift_portgroup"
# $principal = "VSPHERE.LOCAL\openshift"

$pgrole = get-virole $pgrolename


$vdswitch = Get-VDSwitch -Name $vdswitchName
$notes = "vlan: 1254 gateway: 10.5.183.1 cidr: 25 mask: 255.255.255.128"
$newpg = New-VDPortgroup -Name "ci-vlan-1254" -Notes $notes -VDSwitch $vdswitch -VLanId 1254

New-VIPermission -Entity $newpg -Principal $principal -Role $pgrole -Propagate $True
