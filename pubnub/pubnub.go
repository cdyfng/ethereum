package pubnub1

import (
	"crypto/rand"
	"errors"
	"fmt"
	"strings"
	"net/http"
	"encoding/json"
	"net/url"
)

type UUID struct {
	bytes []byte
}

func UUIDGen() (UUID, error) {
	u := UUID{make([]byte, 16)}
	n, err := rand.Read(u.bytes)
	if n != 16 {
		return UUID{}, errors.New("cant read 16 bytes from random reader")
	} else if err != nil {
		return UUID{}, err
	}

	return u, nil
}

func (u UUID) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		u.bytes[0:4], u.bytes[4:6], u.bytes[6:8], u.bytes[8:10], u.bytes[10:])
}


const pubnubOrigin = "pubsub.pubnub.com"

var pubnubClientHeaders = map[string]string{
	"V":	"3.3",
	"User-Agent":"Go-Google",
	"Accept": "*/*",
}

type PubNub struct {
	publish_key, subscribe_key string
	secret_key, cipher_key     string
	ssl                        bool
	session_uuid               UUID
	origin_url                 string

	time_token string
}

type PubNubInterface interface{
	Time() (string, error)
	Publish(channel string, message interface{}) (string, error)
	Subscribe(channel string, stopChan <-chan struct{}) (<-chan interface{}, error)
}


func (pn *PubNub) request(urlbits []string, origin string, encode bool, urlparams url.Values) ([]interface{}, error) {
	if urlbits == nil {
		return nil, errors.New("empty urlbits")
	}

	if encode {
		for i, bit := range urlbits {
			urlbits[i] = url.QueryEscape(bit)
		}
	}

	url := pn.origin_url + "/" + strings.Join(urlbits, "/")

	if urlparams != nil {
		url = url + "?" + urlparams.Encode()
	}

	fmt.Printf("url %s\n", url)

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	for header, value := range pubnubClientHeaders {
		req.Header.Set(header, value)
	}

	response, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println("nil respnse")
		return nil, err
	}

	fmt.Printf("body %v\n", response.Body)
	decoder := json.NewDecoder(response.Body)
	defer response.Body.Close()

	var out []interface{}
	if err := decoder.Decode(&out); err != nil {
		fmt.Println("nil decoder")
		return nil, err
	}

	fmt.Printf("out %v\n", out)
	return out, nil
}

func (pn *PubNub) Time() (string, error) {

	resp, err := pn.request([]string{"time", "0"}, pn.origin_url, false, nil)

	if err != nil {
		return "", err
	}

	if len(resp) < 1 {
		return "", fmt.Errorf("Unexpected response : %s", resp)
	}

	time, ok := resp[0].(float64)

	if !ok {
		return "", errors.New("PubNub time response is not a float64")
	}

	return fmt.Sprintf("%.0f", time), nil

}


func (pn *PubNub) Publish(channel string, message interface{}) (string, error){

	json, err := json.Marshal(message)

	fmt.Println(json)
	fmt.Println(string(json))

	if err != nil {
		return "", err
	}

	args := []string{"publish", pn.publish_key, pn.subscribe_key, "0", channel, "0", string(json)}

	query := url.Values{}
	query.Add("uuid", pn.session_uuid.String())

	resp, err := pn.request(args, pn.origin_url, false, query)

	if err != nil {
		return "", err
	}

	if len(resp) < 3 {
		return "", fmt.Errorf("Unexcepted response:", resp)
	}

	if resp[0].(float64) != 1 {
		return "", errors.New(resp[1].(string))
	}

	timestamp := resp[2].(string)
	return timestamp, nil

}

func (pn *PubNub) Subscribe(channel string, stopChan <-chan struct{}) (<-chan interface{}, error){

	out := make(chan interface{}, 1)

	go func(){
		for{
			select {
			case <-stopChan:
				close(out)
				return
			default:
			}

			args := []string{"subscribe", pn.subscribe_key, "0", channel, "0", pn.time_token}

			query := url.Values{}
			query.Add("uuid", pn.session_uuid.String())

			resp, err := pn.request(args, pn.origin_url, true, query)

			if err != nil {
				close(out)
				return
			}

			if len(resp) < 2 {
				fmt.Println("continue because resp < 2")
				continue
			}

			messages := resp[0].([]interface{})

			pn.time_token = resp[1].(string)

			if len(messages) ==0 {
				fmt.Println("continue because len(messages) =0 , messages: %v", messages)
				continue
			}

			for _, msg := range messages{
				select{
				case out <- msg:
				case <-stopChan:
					close(out)
					return
				}
			}


		}
	}()
	return out, nil
}

func NewPubNub(publish_key, subscribe_key, secret_key, cipher_key string, ssl bool) PubNubInterface{
	pn := &PubNub{
		publish_key: publish_key,
		subscribe_key: subscribe_key,
		secret_key: secret_key,
		cipher_key: cipher_key,
		ssl: ssl,

		time_token: "0",
	}

	uuid, err := UUIDGen()
	if err != nil{
		panic(err)
	}

	pn.session_uuid = uuid

	if ssl {
		pn.origin_url = "https://" + pubnubOrigin
	} else {
		pn.origin_url = "http://" + pubnubOrigin
	}

	fmt.Printf("url :%s\n", pn.origin_url)
	return pn
}




