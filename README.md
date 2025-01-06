# Azure Storage Upload Action

[![Open in GitHub Codespaces](https://github.com/codespaces/badge.svg)](https://codespaces.new/ZEISS/template-action?quickstart=1)

## Example

```yaml
name: Upload to Azure Storage
on:
  push:
    branches:
      - main
jobs:
  upload:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: azure/login@v2
        with:
          client-id: ${{ secrets.AZURE_CLIENT_ID }}
          tenant-id: ${{ secrets.AZURE_TENANT_ID }}
          subscription-id: ${{ secrets.AZURE_SUBSCRIPTION_ID }}
      - uses: zeiss/az-storage-upload@main
        with:
          path: ./dist
          container_name: $www
          account_url: ${{ secrets.AccountUrl }
```

## Inputs

 Key                  | Value                                                                      |
|---------------------|----------------------------------------------------------------------------|
| `container_name`    | The name of the storage account container these assets will be uploaded to |
| `path`              | The path of the files that are uploaded                                    |
| `account_url`       | The URL of the storage account                                             |


## License

[MIT](/LICENSE)
