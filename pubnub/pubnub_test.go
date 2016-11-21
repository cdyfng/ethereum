package pubnub1

import (
	"testing"
	"time"
)

func TestUUIDGen(t *testing.T) {
	u, err := UUIDGen()
	if err != nil {
		t.Error(err)
	}
	t.Logf("UUID: %s", u)
}


func TestPubNubTime(t *testing.T) {
	pubnub := NewPubNub("pub-c-bc7b26a8-f287-4ab2-8407-00f6ecec745b", "sub-c-0047d4b4-aba0-11e6-85a3-02ee2ddab7fe", "", "", false)

	done := make(chan bool)

	go func(){
		time, err := pubnub.Time()

		if err != nil {
			t.Errorf("pubnub.Time: %s", err)
		}

		if time == "" || time == "0" {
			t.Errorf("pubnub.Time returned zero value %q", time)
		} else {
			t.Logf("pubnub.Time sucess: %s", time)
		}

		done <- true
	}()

	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Errorf("No response from pubnub.Time call")
	}
}

func TestPubNub(t *testing.T) {
	channel := "pchat"
	message := map[string]string{"hi":"hello"}// "Hello, pubnub"

	nmsgs := 1

	pubnub := NewPubNub("pub-c-bc7b26a8-f287-4ab2-8407-00f6ecec745b", "sub-c-0047d4b4-aba0-11e6-85a3-02ee2ddab7fe", "", "", false)
	//pubnub := NewPubNub("demo", "demo", "", "", false)


	// publish some messages
	for i := 0; i < nmsgs; i++ {
		t.Logf("--> Publishing %q", message)
		resp, err := pubnub.Publish(channel, message)

		if err != nil {
			t.Errorf("Publish error: %s", err)
		}

		t.Logf("--> Publish response: %#v", resp)
	}

	t.Log("---Publish ok")


	// get a go channel of json objects from a pubnub channel
	subchan, err := pubnub.Subscribe(channel, nil)

	if err != nil {
		t.Errorf("Subscribe error: %s", err)
	}

	// wait a moment..
	time.Sleep(100 * time.Millisecond)



	loop:
	for {

		select {
		// a message completed
		case msg, ok := <-subchan:
			if !ok {
				t.Errorf("Subscriber channel closed")
			}

			t.Logf("<-- Subscriber got message: %q", msg)

			nmsgs--

			if nmsgs == 0 {
				// success; all messages accounted for
				break loop
			}

		// failure; we timed out
		case <-time.After(1 * time.Second):
		// timeout, failure
			t.Errorf("Subscriber timed out on channel %s", channel)
			break loop
		}
	}
}


