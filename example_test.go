// +build more_test

package kokoro

import (
	"fmt"
	"log"
	"os"

	"github.com/kamichidu/kokoro-go/types"
)

func ExampleRegisterDevice() {
	// retrieve user's credential
	email := os.Getenv("KOKOROGO_EMAIL")
	password := os.Getenv("KOKOROGO_PASSWORD")

	cli := NewClient("https://kokoro.io")
	// registering device (aka login with device) requires user's credential
	device, err := cli.RegisterDevice(email, password, &types.Device{
		// display name
		Name: "kokoro-go-example",
		// device kind
		Kind: types.DeviceUnknown,
		// your app name
		DeviceIdentifier: "kokoro-go-example",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(device.Name)
	// Output: kokoro-go-example
	fmt.Println(device.Kind)
	// Output: unknown
	fmt.Println(device.DeviceIdentifier)
	// Output: kokoro-go-example
}
