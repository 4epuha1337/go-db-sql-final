package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
)

var (
	// randSource источник псевдо случайных чисел.
	// Для повышения уникальности в качестве seed
	// используется текущее время в unix формате (в виде числа)
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")// настройте подключение к БД
	if err != nil {
		require.NoError(t, err)
	}
    defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
    id, err := store.Add(parcel)
	if err != nil {
		require.NoError(t, err)
	}

	assert.NotEmpty(t, id)

	// get
	// получите только что добавленную посылку, убедитесь в отсутствии ошибки
	// проверьте, что значения всех полей в полученном объекте совпадают со значениями полей в переменной parcel
	prc, err := store.Get(id)
	if err != nil {
		require.NoError(t, err)
	}

	assert.Equal(t, parcel.Client, prc.Client)
	assert.Equal(t, parcel.Address, prc.Address)
	assert.Equal(t, parcel.Status, prc.Status)
	assert.Equal(t, parcel.CreatedAt, prc.CreatedAt)

	// delete
	// удалите добавленную посылку, убедитесь в отсутствии ошибки
	// проверьте, что посылку больше нельзя получить из БД
	err = store.Delete(id)
	if err != nil {
		require.NoError(t, err)
	}
	
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db") // настройте подключение к БД
	if err != nil {
		require.NoError(t, err)
	}
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
    id, err := store.Add(parcel)
	if err != nil {
		require.NoError(t, err)
	}

	assert.NotEmpty(t, id)

	// set address
	// обновите адрес, убедитесь в отсутствии ошибки
	newAddress := "new test address"

	err = store.SetAddress(id, newAddress)

	if err != nil {
		require.NoError(t, err)
	}

	// check
	// получите добавленную посылку и убедитесь, что адрес обновился
	prc, err := store.Get(id)
	if err != nil {
		require.NoError(t, err)
	}

	assert.Equal(t, parcel.Client, prc.Client)
	assert.Equal(t, prc.Address, newAddress)
	assert.Equal(t, parcel.Status, prc.Status)
	assert.Equal(t, parcel.CreatedAt, prc.CreatedAt)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")// настройте подключение к БД
    if err != nil {
		require.NoError(t, err)
	}
    defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	id, err := store.Add(parcel)
	if err != nil {
		require.NoError(t, err)
	}

	assert.NotEmpty(t, id)

	// set status
	// обновите статус, убедитесь в отсутствии ошибки
	newStatus := ParcelStatusSent
	err = store.SetStatus(id, newStatus)
	if err != nil {
		require.NoError(t, err)
	}

	// check
	// получите добавленную посылку и убедитесь, что статус обновился
	prc, err := store.Get(id)
	if err != nil {
		require.NoError(t, err)
	}

	assert.Equal(t, parcel.Client, prc.Client)
	assert.Equal(t, parcel.Address, prc.Address)
	assert.Equal(t, prc.Status, newStatus)
	assert.Equal(t, parcel.CreatedAt, prc.CreatedAt)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")// настройте подключение к БД
    if err != nil {
		require.NoError(t, err)
	}
    defer db.Close()

	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	// задаём всем посылкам один и тот же идентификатор клиента
	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	parcelsNum := 3

	// add
	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
		if err != nil {
			require.NoError(t, err)
		}
	
		assert.NotEmpty(t, id)

		// обновляем идентификатор добавленной у посылки
		parcels[i].Number = id

		// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[id] = parcels[i]
	}

	// get by client
	storedParcels, err := store.GetByClient(client) // получите список посылок по идентификатору клиента, сохранённого в переменной client
	// убедитесь в отсутствии ошибки
	if err != nil {
		require.NoError(t, err)
	}
	// убедитесь, что количество полученных посылок совпадает с количеством добавленных

	assert.Equal(t, parcelsNum, len(storedParcels))


	// check
	for i, parcel := range storedParcels {
		// в parcelMap лежат добавленные посылки, ключ - идентификатор посылки, значение - сама посылка
		// убедитесь, что все посылки из storedParcels есть в parcelMap
		// убедитесь, что значения полей полученных посылок заполнены верно
		assert.Equal(t, parcels[i].Address, parcel.Address)
		assert.Equal(t, parcels[i].Client, parcel.Client)
		assert.Equal(t, parcels[i].CreatedAt, parcel.CreatedAt)
		assert.Equal(t, parcels[i].Status, parcel.Status)
	}
}
