
```
terraform init
terraform plan -var-file=./secrets/terraform.tfvars -var-file=./terraform.tfvars
terraform apply -var-file=./secrets/terraform.tfvars -var-file=./terraform.tfvars
```


```
terraform state list | grep vlan | xargs -I {} terraform state show {} | grep vlan_number | awk -F\= '{print $2}' | tr -d '[:blank:]'
```
