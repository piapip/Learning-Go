package main

import (
	"fmt"
	"time"
)

//Person ...
type Person struct {
	name        string
	age         int
	city, phone string
}

//People = Person acting
type People interface {
	SayHello()
	GetDetails()
}

//SayHello = jikkoshoukai
func (p Person) SayHello() {
	fmt.Printf("Hi, I am %s, from %s\n", p.name, p.city)
}

//GetDetails = identity card
func (p Person) GetDetails() {
	fmt.Printf("[Name: %s, Age: %d, City: %s, Phone: %s]\n", p.name, p.age, p.city, p.phone)
}

//Speaker function
type Speaker struct {
	Person
	speaksOn   []string
	pastEvents []string
}

//GetDetails overrides
func (s Speaker) GetDetails() {
	s.Person.GetDetails()
	fmt.Println("Speaker talks on following technologies:")
	for _, value := range s.speaksOn {
		fmt.Println(value)
	}
	fmt.Println("Presented on the following conferences:")
	for _, value := range s.pastEvents {
		fmt.Println(value)
	}
}

//Organizer a special person
type Organizer struct {
	Person
	meetups []string
}

//GetDetails overrides
func (o Organizer) GetDetails() {
	o.Person.GetDetails()
	fmt.Println("Organizer, conducting following Meetups:")
	for _, value := range o.meetups {
		fmt.Println(value)
	}
}

//Attendee just another special person
type Attendee struct {
	Person
	interest []string
}

//Athlete = sport activities
type Athlete interface {
	throw()
	swim()
}

func (p Person) throw() {
	fmt.Println("throwing...")
}

func (p Person) swim() {
	fmt.Println("Swimming...")
}

//Meetup struct
type Meetup struct {
	location string
	city     string
	date     time.Time
	people   []People
	//so this crapshoot define what people that pass into Meetup can do, even though by themselves they can throw and swim but they can't do
	//that in this Meetup, they'd better sit down or else
	//people   []Athlete
}

//MeetUpPeople not sure the outcome of this one
func (m Meetup) MeetUpPeople() {
	for _, v := range m.people {
		v.SayHello()
		v.GetDetails()
	}
}

func testMeetUp() {
	shiju := Speaker{Person{"Shiju", 35, "Kochi", "+91-94003372xx"},
		[]string{"Go", "Docker", "Azure", "AWS"},
		[]string{"FOSS", "JSFOO", "MS TechDays"}}

	satish := Organizer{Person{"Satish", 35, "Pune", "+91-94003372xx"},
		[]string{"Gophercon", "RubyConf"}}

	alex := Attendee{Person{"Alex", 22, "Bangalore", "+91-94003672xx"},
		[]string{"Go", "Ruby"}}

	meetup := Meetup{
		"Royal Orchid",
		"Bangalore",
		time.Date(2015, time.February, 19, 9, 0, 0, 0, time.UTC),
		[]People{shiju, satish, alex},
	}

	meetup.MeetUpPeople()

	// normalPerson := Person{"Shiju", 35, "Kochi", "+91-94003372xx"}
	// speaker := Speaker{normalPerson, []string{"Go", "Docker", "Azure", "AWS"}, []string{"FOSS", "JSFOO", "MS TechDays"}}
	// organizer := Organizer{normalPerson, []string{"Gophercon", "RubyConf"}}

	// attendees := []People{normalPerson, speaker, organizer}
	// for _, attendee := range attendees {
	// 	attendee.GetDetails()
	// }
	// speaker.GetDetails()
	// organizer.GetDetails()
}
