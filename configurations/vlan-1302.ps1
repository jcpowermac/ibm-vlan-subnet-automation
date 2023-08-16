
$vdswitchName = "vcs8e-vcs-ci-workload-private"
$vdswitch = Get-VDSwitch -Name $vdswitchName 
$notes = "vlan: 1302 gateway: 10.94.196.1 cidr: 25 mask: 255.255.255.128"
New-VDPortgroup -Name "ci-vlan-1302" -Notes $notes -VDSwitch $vdswitch -VLanId 1302
