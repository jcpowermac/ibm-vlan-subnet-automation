
$vdswitchName = "vcs8e-vcs-ci-workload-private"
$vdswitch = Get-VDSwitch -Name $vdswitchName 
$notes = "vlan: 1237 gateway: 10.38.247.1 cidr: 25 mask: 255.255.255.128"
New-VDPortgroup -Name "ci-vlan-1237" -Notes $notes -VDSwitch $vdswitch -VLanId 1237
