
$vdswitchName = "vcs8e-vcs-ci-workload-private"
$vdswitch = Get-VDSwitch -Name $vdswitchName 
$notes = "vlan: 956 gateway: 10.93.134.1 cidr: 25 mask: 255.255.255.128"
New-VDPortgroup -Name "ci-vlan-956" -Notes $notes -VDSwitch $vdswitch -VLanId 956
