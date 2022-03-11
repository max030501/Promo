package promo

import (
	"awesomeProject1/participant"
	"awesomeProject1/prize"
	"errors"
	"fmt"
	"math/rand"
	"sync"
)

type Promo struct {
	Id int `json:"id"`
	Name string `json:"name" binding:"required"`
	Description string `json:"description"`
	Part map[int]*participant.Participant `json:"-"`
	Prizes map[int]*prize.Prize `json:"-"`
	NextIdPart int `json:"-"`
	NextIdPrize int `json:"-"`
}

type RaffleResponse struct {
	participant.Participant `json:"winner"`
	prize.Prize             `json:"prize"`
}
type PromoStore struct {
	mutex sync.Mutex

	promos map[int]*Promo
	nextId int
}

type ResponsePromo struct {

Id int `json:"id"`
Name string `json:"name" binding:"required"`
Description string `json:"description"`
Part map[int]*participant.Participant `json:"participants"`
Prizes map[int]*prize.Prize `json:"prizes"`
	NextIdPart int `json:"-"`
	NextIdPrize int `json:"-"`
}

type RequestPromo struct {
	Name string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func New() *PromoStore {
	store := &PromoStore{}
	store.promos = make(map[int]*Promo)
	store.nextId = 1
	return store
}

func (st *PromoStore) CreatePromo(name string, description string) int {
	st.mutex.Lock()
	defer st.mutex.Unlock()


	promo := Promo{
		Id:           st.nextId,
		Name:         name,
		Description:  description,
		Part: make(map[int]*participant.Participant),
		Prizes: make(map[int]*prize.Prize),
		NextIdPart: 1,
		NextIdPrize: 1,
	}
	st.promos[st.nextId] = &promo
	st.nextId++
	return promo.Id
}

func (st *PromoStore) GetPromos() []Promo {
	st.mutex.Lock()
	defer st.mutex.Unlock()

	promos := make([]Promo, 0, len(st.promos))
	for _, promo := range st.promos {
		promos = append(promos, *promo)
	}
	return promos

}

func (st *PromoStore) GetPromo(id int) (ResponsePromo, error) {
	st.mutex.Lock()
	defer st.mutex.Unlock()

	promo, ok := st.promos[id]
	if ok {
		return ResponsePromo(*promo), nil
	} else {
		return ResponsePromo{}, fmt.Errorf("There is no promo with id=%d", id)
	}
}

func (st *PromoStore) PutPromo(id int, name string, description string) error {
	st.mutex.Lock()
	defer st.mutex.Unlock()

	if _, ok := st.promos[id]; !ok {
		return fmt.Errorf("There is no promo with id=%d", id)
	}

	st.promos[id].Name = name
	st.promos[id].Description = description

	return nil
}

func (st *PromoStore) DeletePromo(id int) error {
	st.mutex.Lock()
	defer st.mutex.Unlock()

	if _, ok := st.promos[id]; !ok {
		return fmt.Errorf("There is no promo with id=%d", id)
	}
	delete(st.promos, id)
	return nil
}

func (st *PromoStore) CreateParticipant(id int,name string) (int,error) {
	st.mutex.Lock()
	defer st.mutex.Unlock()

	if _, ok := st.promos[id]; !ok {
		return -1,fmt.Errorf("There is no promo with id=%d", id)
	}
	part := &participant.Participant{
		Id:   st.promos[id].NextIdPart,
		Name: name,
	}
	st.promos[id].Part[st.promos[id].NextIdPart]= part
	st.promos[id].NextIdPart++
	return part.Id, nil
}

func (st *PromoStore) DeleteParticipant(id int,partId int) error {
	st.mutex.Lock()
	defer st.mutex.Unlock()

	if _, ok := st.promos[id]; !ok {
		return fmt.Errorf("There is no promo with id=%d", id)
	}
	if _, ok := st.promos[id].Part[partId]; !ok {
		return fmt.Errorf("There is no participant with id=%d", id)
	}
	delete(st.promos[id].Part, partId)
	return nil
}

func (st *PromoStore) CreatePrize(id int,desc string) (int,error) {
	st.mutex.Lock()
	defer st.mutex.Unlock()

	if _, ok := st.promos[id]; !ok {
		return -1,fmt.Errorf("There is no promo with id=%d", id)
	}
	prize := &prize.Prize{
		Id:   st.promos[id].NextIdPrize,
		Description: desc,
	}
	st.promos[id].Prizes[st.promos[id].NextIdPrize]= prize
	st.promos[id].NextIdPrize++
	return prize.Id, nil
}

func (st *PromoStore) DeletePrize(id int,prizeId int) error {
	st.mutex.Lock()
	defer st.mutex.Unlock()

	if _, ok := st.promos[id]; !ok {
		return fmt.Errorf("There is no promo with id=%d", id)
	}
	if _, ok := st.promos[id].Part[prizeId]; !ok {
		return fmt.Errorf("There is no prize with id=%d", id)
	}
	delete(st.promos[id].Prizes, prizeId)
	return nil
}

func (st *PromoStore) Raffle(id int) ([]RaffleResponse,error) {
	st.mutex.Lock()
	defer st.mutex.Unlock()

	if _, ok := st.promos[id]; !ok {
		return []RaffleResponse{},fmt.Errorf("There is no promo with id=%d", id)
	}
	if len(st.promos[id].Part) !=len(st.promos[id].Prizes) || len(st.promos[id].Prizes) == 0 || len(st.promos[id].Part)==0{
		return []RaffleResponse{},errors.New("Count of participants and prizes is not equal or empty")
	}
	prizes :=make([]prize.Prize, 0, len(st.promos[id].Prizes))
	for _, pr := range st.promos[id].Prizes {
		prizes = append(prizes, *pr)
	}
	parts := make([]participant.Participant, 0, len(st.promos[id].Part))
	for _, pr := range st.promos[id].Part {
		parts = append(parts, *pr)
	}
	rand.Shuffle(len(st.promos[id].Prizes),
		func(i, j int) { prizes[i], prizes[j] = prizes[j], prizes[i] })
	rand.Shuffle(len(st.promos[id].Part),
		func(i, j int) { parts[i], parts[j] = parts[j], parts[i] })
	raffle := make([]RaffleResponse, 0, len(st.promos[id].Prizes))
	for i:=0; i < len(st.promos[id].Prizes);i++{
		rr := RaffleResponse{
			Participant: parts[i],
			Prize:       prizes[i],
		}
		raffle = append(raffle,rr)
	}

	return raffle,nil
}




