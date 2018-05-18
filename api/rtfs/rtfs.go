package rtfs

import (
	"errors"
	"fmt"

	ipfsapi "github.com/ipfs/go-ipfs-api"
)

type IpfsManager struct {
	Shell    *ipfsapi.Shell
	PubSub   *ipfsapi.PubSubSubscription
	PinTopic string
}

func Initialize() *IpfsManager {
	manager := IpfsManager{}
	manager.Shell = establishShellWithNode("")
	manager.PinTopic = "pin"
	return &manager
}

// Pin is a wrapper method to pin a hash to the local node,
// but also alert the rest of the local nodes to pin
// after which the pin will be sent to the cluster
func (im *IpfsManager) Pin(hash string) error {
	err := im.Shell.Pin(hash)
	if err != nil {
		// TODO: add error reporting
		return err
	}
	im.PublishPubSubMessage(im.PinTopic, hash)
	return nil
}

func establishShellWithNode(url string) *ipfsapi.Shell {
	if url == "" {
		shell := ipfsapi.NewLocalShell()
		return shell
	}
	shell := ipfsapi.NewShell(url)
	return shell
}

// SubscribeToPubSubTopic is used to subscribe to a pubsub topic
func (im *IpfsManager) SubscribeToPubSubTopic(topic string) error {
	if topic == "" {
		return errors.New("invalid topic name")
	}
	// create a pubsub subscription according to topic name
	subscription, err := im.Shell.PubSubSubscribe(topic)
	if err != nil {
		return err
	}
	// store the pubsub scription
	im.PubSub = subscription
	return nil
}

// ConsumeSubscription is used to consume a pubsub subscription
// note that it will automatically exit after receiving and processing all the messages
func (im *IpfsManager) ConsumeSubscription(sub *ipfsapi.PubSubSubscription) error {
	count := 0
	for {
		if count == 1000 {
			break
		}
		subRecord, err := sub.Next()
		if err != nil {
			return err
		}
		if subRecord == nil {
			continue
		}
		count++
		fmt.Println(subRecord)
	}
	return nil
}

// PublishPubSubMessage is used to publish a message to the given topic
func (im *IpfsManager) PublishPubSubMessage(topic string, data string) error {
	if topic == "" && data == "" {
		return errors.New("invalid topic and data")
	}
	err := im.Shell.PubSubPublish(topic, data)
	if err != nil {
		return err
	}
	return nil
}