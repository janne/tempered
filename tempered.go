package tempered

import (
	"fmt"
)

// #cgo LDFLAGS: -ltempered -lhidapi-libusb
// #include <stdlib.h>
// #include <tempered.h>
import "C"
import "unsafe"

type Device struct {
	VendorId        uint16
	ProductId       uint16
	InterfaceNumber int
	Path            string
	TypeName        string
	handle          *C.tempered_device
}

type Tempered struct {
	Devices []Device
}

type Sensing struct {
	TempC float32
	RelHum float32
}

func (d Device) Sense() (s Sensing, err error) {
	if C.tempered_read_sensors(d.handle) {
		count := int(C.tempered_get_sensor_count(d.handle))
		for sensor := 0; sensor < count; sensor++ {
			typ := int(C.tempered_get_sensor_type(d.handle, C.int(sensor)))
			if typ == C.TEMPERED_SENSOR_TYPE_NONE {
				err = fmt.Errorf("No such sensor, or type is not supported.\n")
				return
			}
			if (typ & C.TEMPERED_SENSOR_TYPE_TEMPERATURE) > 0 {
				var TempC C.float
				if C.tempered_get_temperature(d.handle, C.int(sensor), &TempC) {
					s.TempC = float32(TempC)
				} else {
					err = fmt.Errorf("Temperature failed: %s\n", C.GoString(C.tempered_error(d.handle)))
				}
			}

			if (typ & C.TEMPERED_SENSOR_TYPE_HUMIDITY) > 0 {
				var RelHum C.float
				if C.tempered_get_humidity(d.handle, C.int(sensor), &RelHum) {
					s.RelHum = float32(RelHum)
				} else {
					err = fmt.Errorf("Humidity failed: %s\n", C.GoString(C.tempered_error(d.handle)))
				}
			}
		}
	} else {
		err = fmt.Errorf("Failed to read the sensors: %s\n",
			C.GoString(C.tempered_error(d.handle)))
	}
	return
}

func New() (t *Tempered, err error) {
	var error *C.char
	t = &Tempered{}
	if !C.tempered_init(&error) {
		err = fmt.Errorf(C.GoString(error))
		C.free(unsafe.Pointer(error))
	}

	list := C.tempered_enumerate(&error)
	if list == nil {
		err = fmt.Errorf(C.GoString(error))
		C.free(unsafe.Pointer(error))
		return
	}

	t.Devices = make([]Device, 0)
	for dev := list; dev != nil; dev = dev.next {
		handle := C.tempered_open(dev, &error)
		if handle == nil {
			err = fmt.Errorf(C.GoString(error))
			C.free(unsafe.Pointer(error))
			t = nil
			return
		}
		t.Devices = append(t.Devices, Device{
			uint16(dev.vendor_id),
			uint16(dev.product_id),
			int(dev.interface_number),
			C.GoString(dev.path),
			C.GoString(dev.type_name),
			handle})
	}

	C.tempered_free_device_list(list)
	return
}

func (t *Tempered) Close() (err error) {
	var error *C.char
	for _, dev := range t.Devices {
		C.tempered_close(dev.handle)
	}

	if !C.tempered_exit(&error) {
		err = fmt.Errorf(C.GoString(error))
		C.free(unsafe.Pointer(error))
	}
	return
}
