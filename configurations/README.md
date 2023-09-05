
Wireguard subnets
```
$ grep subnet vyatta-*-01.conf | awk '{print $3}' | paste -sd "," -
```

```
# set these three vars
  15        0.000 $pgrolename = "Administrator"
  16        0.000 $principal = "VSPHERE.LOCAL\Administrators"
  17        0.000 $vdswitchName = "VManagement"

# run all the scripts
  28       25.201 Get-ChildItem *.ps1  | %{Invoke-Expression $_.FullName}

```