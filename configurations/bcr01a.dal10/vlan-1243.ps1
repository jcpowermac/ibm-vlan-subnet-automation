
#$vdswitchName = "vcs8e-vcs-ci-workload-private"
# $pgrolename = "openshift_portgroup"
# $principal = "VSPHERE.LOCAL\openshift"

$pgrole = get-virole $pgrolename


$vdswitch = Get-VDSwitch -Name $vdswitchName
$notes = "vlan: 1243 gateway: 10.176.121.225 cidr: 27 mask: 255.255.255.224"
$newpg = New-VDPortgroup -Name "ci-vlan-1243" -Notes $notes -VDSwitch $vdswitch -VLanId 1243

New-VIPermission -Entity $newpg -Principal $principal -Role $pgrole -Propagate $True
