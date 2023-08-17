
#$vdswitchName = "vcs8e-vcs-ci-workload-private"

$vdswitch = Get-VDSwitch -Name $vdswitchName 
$notes = "vlan: 1235 gateway: 10.38.221.129 cidr: 25 mask: 255.255.255.128"
New-VDPortgroup -Name "ci-vlan-1235" -Notes $notes -VDSwitch $vdswitch -VLanId 1235
