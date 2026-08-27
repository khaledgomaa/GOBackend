package main

import (
	"errors"
	"sync"
)

var ErrTruckNotFound = errors.New("truck not found")

type FleetManager interface {
	AddTruck(id string, cargo int) error
	GetTruck(id string) (Truck, error) // having *Truck as return type will return a pointer to the truck, which is more efficient than returning a copy of the truck struct but it will break the mutex as the mutex will be locked when the truck is returned and it will not be unlocked until the truck is no longer used, which can lead to deadlocks. So we will return a copy of the truck struct instead of a pointer to it.
	RemoveTruck(id string) error
	UpdateTruckCargo(id string, cargo int) error
}

type Truck struct {
	ID    string
	Cargo int
}

type truckManager struct {
	trucks       map[string]*Truck
	sync.RWMutex // without variable defined the truckManager struct will have a mutex for thread-safe operations
}

func (t *truckManager) UpdateTruckCargo(id string, cargo int) error {
	_, ok := t.trucks[id]
	if !ok {
		return ErrTruckNotFound
	}
	t.trucks[id].Cargo = cargo
	return nil
}

func (t *truckManager) RemoveTruck(id string) error {
	t.Lock()
	defer t.Unlock()
	_, ok := t.trucks[id]
	if !ok {
		return ErrTruckNotFound
	}
	delete(t.trucks, id)
	return nil
}

func (t *truckManager) GetTruck(id string) (Truck, error) {
	t.RLock()
	defer t.RUnlock()
	truck, ok := t.trucks[id]
	if !ok {
		return Truck{}, ErrTruckNotFound
	}
	return *truck, nil
}

func (t *truckManager) AddTruck(id string, cargo int) error {
	t.Lock()
	defer t.Unlock()
	_, ok := t.trucks[id]
	if ok {
		return errors.New("truck already exists")
	} else {
		t.trucks[id] = &Truck{ID: id, Cargo: cargo}
	}
	return nil
}

func NewTruckManager() truckManager {
	return truckManager{
		trucks: make(map[string]*Truck),
	}
}
