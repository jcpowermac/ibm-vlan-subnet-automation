
#$vdswitchName = "vcs8e-vcs-ci-workload-private"

$vdswitch = Get-VDSwitch -Name $vdswitchName 
$notes = "vlan: 1255 gateway: 10.38.220.1 cidr: 25 mask: 255.255.255.128"
New-VDPortgroup -Name "ci-vlan-1255" -Notes $notes -VDSwitch $vdswitch -VLanId 1255
