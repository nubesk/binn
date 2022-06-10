package binn

import (
	"fmt"
	"time"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestUseID(t *testing.T) {
	idStorage := DefaultIDStorage()

	_ = idStorage.Add(
		"1c7a8201-cdf7-11ec-a9b3-0242ac110004",
		time.Now().Add(time.Duration(10)*time.Minute),
	)

	err := idStorage.Use("1c7a8201-cdf7-11ec-a9b3-0242ac110004")
	assert.Nil(t, err)
}

func TestUseExpiredID(t *testing.T) {
	idStorage := DefaultIDStorage()
	_ = idStorage.Add(
		"1c7a8201-cdf7-11ec-a9b3-0242ac110004",
		time.Now().Add(time.Duration(1)*time.Millisecond),
	)

	time.Sleep(time.Duration(1)*time.Millisecond)

	err := idStorage.Use("1c7a8201-cdf7-11ec-a9b3-0242ac110004")
	assert.Error(t, err, "this id (\"1c7a8201-cdf7-11ec-a9b3-0242ac110004\") is expired")
}

func TestUseInvalidID(t *testing.T) {
	idStorage := DefaultIDStorage()
	err := idStorage.Use("1c7a8201-cdf7-11ec-a9b3-0242ac110004")
	assert.Error(t, err, "this ID (\"1c7a8201-cdf7-11ec-a9b3-0242ac110004\") is invalid")
}

func TestUseIDUpdatedExpiredAt(t *testing.T) {
	idStorage := DefaultIDStorage()
	_ = idStorage.Add(
		"1c7a8201-cdf7-11ec-a9b3-0242ac110004",
		time.Now().Add(time.Duration(1)*time.Millisecond),
	)

	time.Sleep(time.Duration(1)*time.Millisecond)

	idStorage.Update(
		"1c7a8201-cdf7-11ec-a9b3-0242ac110004",
		time.Now().Add(time.Duration(1)*time.Minute),
	)

	err := idStorage.Use("1c7a8201-cdf7-11ec-a9b3-0242ac110004")
	assert.Nil(t, err)
}

func TestAddBottleToStorage(t *testing.T) {
	idStorage := DefaultIDStorage()
	containerStorage := NewContainerStorage(true, 0, idStorage)
	idStorage.Add(
		"1c7a8201-cdf7-11ec-a9b3-0242ac110004",
		time.Now().Add(time.Duration(10)*time.Minute),
	)

	bottle := NewBottle(
		"1c7a8201-cdf7-11ec-a9b3-0242ac110004",
		"This is a Test Message",
		nil,
	)

	err := containerStorage.Add(bottle)

	assert.Nil(t, err)
}

func TestAddBottleOnNoValidation(t *testing.T) {
	containerStorage := NewContainerStorage(false, 0, nil)
	
	bottle := NewBottle(
		"1c7a8201-cdf7-11ec-a9b3-0242ac110004",
		"This is a Test Message",
		nil,
	)

	err := containerStorage.Add(bottle)

	assert.Nil(t, err)
}

func TestGetBottleFromStorage(t *testing.T) {
	containerStorage := NewContainerStorage(false, 0, nil)
	_ = containerStorage.Add(NewBottle(
		"1c7a8201-cdf7-11ec-a9b3-0242ac110004",
		"This is a Test Message",
		nil,
	))

	bottle, _ := containerStorage.Get()
	assert.NotEqual(t, "1c7a8201-cdf7-11ec-a9b3-0242ac110004", bottle.ID())
}

func TestGetBottleFromEmptyContainerStorage(t *testing.T) {
	containerStorage := NewContainerStorage(false, 0, nil)

	bottle, err := containerStorage.Get()

	assert.EqualError(t, err, "this storage has no containers")
	assert.Nil(t, bottle)
}

func TestAddBottleOverflowStorage(t *testing.T) {
	containerStorage := NewContainerStorage(false, 0, nil)
	for i := 0; i < 1000; i++ {
		_ = containerStorage.Add(
			NewBottle(
				fmt.Sprintf("%d", i),
				fmt.Sprintf("%d", i),
				nil,
			))
	}
	
	containerStorage.Add(NewBottle("1000", "1000", nil))

	b, _ := containerStorage.Get()
	assert.Equal(t, "1", b.Message().Text)
}

func TestAddExpiredBottle(t *testing.T) {
	idStorage := DefaultIDStorage()
	storage := NewContainerStorage(true, time.Duration(10) * time.Millisecond, idStorage)
	idStorage.Add(
		"1c7a8201-cdf7-11ec-a9b3-0242ac110004",
		time.Now().Add(time.Duration(10) * time.Minute),
	)
	_ = storage.Add(NewBottle(
		"1c7a8201-cdf7-11ec-a9b3-0242ac110004",
		"Empty Bottle",
		nil,
	))
	_, _ = storage.Get()

	time.Sleep(time.Duration(20) * time.Millisecond)
	
	err := storage.Add(NewBottle(
		"1c7a8201-cdf7-11ec-a9b3-0242ac110004",
		"This bottles is expired",
		nil,
	))
	assert.Error(t, err, "this container is expired")
}

func TestAddLongMessageMaxMessageTextLength(t *testing.T) {
	idStorage := DefaultIDStorage()
	storage := NewContainerStorage(true, time.Duration(10) * time.Minute, idStorage)
	idStorage.Add(
		"1c7a8201-cdf7-11ec-a9b3-0242ac110004",
		time.Now().Add(time.Duration(10) * time.Minute),
	)
	longText := `Lorem ipsum dolor sit amet,
consectetur adipiscing elit.Pellentesque
sapien purus, rhoncus a consectetur vitae,
consectetur a sem. Vestibulum ante
ipsum primis in faucibus orci luctus
et ultrices posuere cubilia curae;`
	
	_ = storage.Add(NewBottle(
		"1c7a8201-cdf7-11ec-a9b3-0242ac110004",
		longText,
		nil,
	))

	bottle, _ := storage.Get()
	assert.Equal(t, 200, len(bottle.Message().Text))
}

func TestGenerateID(t *testing.T) {
	_ = GenerateID()
}
