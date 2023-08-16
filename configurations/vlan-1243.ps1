
$vdswitchName = "vcs8e-vcs-ci-workload-private"
$vdswitch = Get-VDSwitch -Name $vdswitchName 
$notes = "vlan: 1243 gateway: 10.38.204.129 cidr: 25 mask: 255.255.255.128"
New-VDPortgroup -Name "ci-vlan-1243" -Notes $notes -VDSwitch $vdswitch -VLanId 1243
