// Accesses hedera's API
package hedera

import (
	"encoding/json"
	"time"

	hederaSDK "github.com/hashgraph/hedera-sdk-go/v2"
)

type HederaAPI interface {
	HashPolygon(Feature, chan []byte) error
}

type hederaAPI struct {
	client *hederaSDK.Client
}

// Returns a new instance of the API
func NewHederaAPI(accountID string, privateKey string) (HederaAPI, error) {
	operatorAccountID, err := hederaSDK.AccountIDFromString(accountID)
	if err != nil {
		return nil, err
	}

	operatorKey, err := hederaSDK.PrivateKeyFromString(privateKey)
	if err != nil {
		return nil, err
	}

	client := hederaSDK.ClientForTestnet()
	client.SetOperator(operatorAccountID, operatorKey)
	return hederaAPI{
		client: client,
	}, nil
}

// Send the polygon to the HCS and returns it's hash
// TODO: change this to a pub sub architecture
func (api hederaAPI) HashPolygon(polygon Feature, ch chan []byte) error {
	jsonValue, err := json.Marshal(polygon)
	if err != nil {
		return err
	}

	// Create a topic
	transactionResponse, err := hederaSDK.NewTopicCreateTransaction().
		SetTransactionMemo("hash polygon").
		SetAdminKey(api.client.GetOperatorPublicKey()).
		Execute(api.client)

	if err != nil {
		return err
	}

	// Get topic receipt
	transactionReceipt, err := transactionResponse.GetReceipt(api.client)
	if err != nil {
		return err
	}

	topicID := *transactionReceipt.TopicID

	// Subscribe to topic
	_, err = hederaSDK.NewTopicMessageQuery().
		SetTopicID(topicID).
		SetStartTime(time.Unix(0, 0)).
		Subscribe(api.client, func(message hederaSDK.TopicMessage) {
			ch <- message.RunningHash
		})

	if err != nil {
		return err
	}

	// Submit to topic
	_, err = hederaSDK.NewTopicMessageSubmitTransaction().
		SetMessage(jsonValue).
		SetTopicID(topicID).
		Execute(api.client)

	if err != nil {
		return err
	}

	// Delete topic
	transactionResponse, err = hederaSDK.NewTopicDeleteTransaction().
		SetTopicID(topicID).
		SetNodeAccountIDs([]hederaSDK.AccountID{transactionResponse.NodeID}).
		SetMaxTransactionFee(hederaSDK.NewHbar(5)).
		Execute(api.client)
	if err != nil {
		return err
	}

	// Get receipt from topic deletion
	_, err = transactionResponse.GetReceipt(api.client)
	if err != nil {
		return err
	}

	return nil
}
