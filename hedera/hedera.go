// Accesses hedera's API
package hedera

import (
	"encoding/json"
	"time"

	hederaSDK "github.com/hashgraph/hedera-sdk-go/v2"
)

type HederaAPI interface {
	SubmitPolygon(Feature) error
	Close() error
}

type hederaAPI struct {
	client  *hederaSDK.Client
	topicID hederaSDK.TopicID
	nodeID  hederaSDK.AccountID
}

// Returns a new instance of the API
func NewHederaAPI(accountID string, privateKey string, ch chan []byte) (HederaAPI, error) {
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

	api := hederaAPI{
		client: client,
	}

	err = api.initClient(ch)

	return api, err
}

func (api *hederaAPI) initClient(ch chan []byte) error {
	// Create a topic
	transactionResponse, err := hederaSDK.NewTopicCreateTransaction().
		SetTransactionMemo("hash polygon").
		SetAdminKey(api.client.GetOperatorPublicKey()).
		Execute(api.client)

	if err != nil {
		return err
	}

	api.nodeID = transactionResponse.NodeID

	// Get topic receipt
	transactionReceipt, err := transactionResponse.GetReceipt(api.client)
	if err != nil {
		return err
	}

	topicID := *transactionReceipt.TopicID
	api.topicID = topicID

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

	return nil
}

func (api hederaAPI) Close() error {
	// Delete topic
	transactionResponse, err := hederaSDK.NewTopicDeleteTransaction().
		SetTopicID(api.topicID).
		SetNodeAccountIDs([]hederaSDK.AccountID{api.nodeID}).
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

// Send the polygon to the HCS and returns it's hash
// TODO: change this to a pub sub architecture
func (api hederaAPI) SubmitPolygon(polygon Feature) error {
	jsonValue, err := json.Marshal(polygon)
	if err != nil {
		return err
	}

	// Submit to topic
	_, err = hederaSDK.NewTopicMessageSubmitTransaction().
		SetMessage(jsonValue).
		SetTopicID(api.topicID).
		Execute(api.client)

	return err
}
