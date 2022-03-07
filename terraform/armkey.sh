#! /usr/local/bin/bash
ACCOUNT_KEY=$(az storage account keys list --resource-group tfstate --account-name fichistfstate --query '[0].value' -o tsv --subscription "58731601-e920-43f2-af4b-f1a79a8ade9f")
echo $ACCOUNT_KEY
export ARM_ACCESS_KEY=$ACCOUNT_KEY
