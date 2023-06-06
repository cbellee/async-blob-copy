LOCATION='australiaeast'
PREFIX='async-blob-copy-test'
RESOURCE_GROUP_NAME="$PREFIX-rg"

# deploy resource group
az group create --location $LOCATION --name $RESOURCE_GROUP_NAME

# deploy infrastructure
funcName=`az deployment group create \
  --resource-group $RESOURCE_GROUP_NAME \
  --template-file deploy.bicep \
  --parameters location=$LOCATION \
  --parameters prefix=$PREFIX \
  --query properties.outputs.funcName.value -o tsv`

cd ./handler

# compile binary
GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -trimpath handler.go

# publish function
func azure functionapp publish $funcName --list-ignored-files
