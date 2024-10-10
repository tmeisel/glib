package azure

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

type Client struct {
	client *armresources.ResourceGroupsClient
}

func NewClient(subscriptionID string) (*Client, error) {
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	clientFactory, err := armresources.NewClientFactory(subscriptionID, credential, nil)
	if err != nil {
		return nil, err
	}

	client := clientFactory.NewResourceGroupsClient()

	return &Client{client: client}, nil
}

func (c *Client) Do() {

}
