package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var ops = []byte{'+', '-', '*', '/'}

type Question struct {
	A  int
	B  int
	op byte
}

type Answer struct {
	By    *Student
	Value int
}

func NewQuestion(A, B int, op byte) *Question {
	return &Question{
		A:  A,
		B:  B,
		op: op,
	}
}

func (q Question) Calculate() (int, bool) {
	switch q.op {
	case '+':
		return q.A + q.B, true
	case '-':
		return q.A - q.B, true
	case '*':
		return q.A * q.B, true
	case '/':
		if q.B == 0 {
			return 0, false
		}
		return q.A / q.B, true
	}
	return 0, false
}

func (q Question) ToString() string {
	return fmt.Sprintf("%d %c %d = ", q.A, q.op, q.B)
}

type Teacher struct{}

func (t *Teacher) GenerateQuestion() *Question {
	// A and B value range is [1,100)
	A, B := rand.Intn(99)+1, rand.Intn(99)+1
	return NewQuestion(A, B, ops[rand.Intn(4)])
}

func (t *Teacher) Say(s string) {
	fmt.Printf("Teacher: %s\n", s)
}

type Student struct {
	Name string
}

func (t *Student) Say(s string) {
	fmt.Printf("Student %s: %s\n", t.Name, s)
}

type Game struct {
	stop     chan any
	wg       sync.WaitGroup
	teacher  *Teacher
	students []*Student
}

func NewGame() *Game {
	students := make([]*Student, 5)
	for i := range students {
		students[i] = &Student{Name: string('A' + byte(i))}
	}
	return &Game{
		stop:     make(chan any, 1),
		teacher:  &Teacher{},
		students: students,
	}
}

func (g *Game) Start() {
	for {
		select {
		case <-g.stop:
			return
		default:
			g.teacher.Say("Guys, are you ready?")
			// count 3 secs
			time.Sleep(time.Second * 3)
			// teacher ask question
			q := g.teacher.GenerateQuestion()
			g.teacher.Say(q.ToString() + "?")

			var once sync.Once
			var winner *Student
			correctAnswer, valid := q.Calculate()
			if !valid {
				g.teacher.Say("Invalid question, skipping...")
				continue
			}

			// student answer
			for _, student := range g.students {
				g.wg.Add(1)
				go func(s *Student) {
					defer g.wg.Done()

					// simulate think randomly between 1 and 3 seconds
					time.Sleep(time.Duration(1000+rand.Intn(2000)) * time.Millisecond)

					once.Do(func() {
						winner = s
						s.Say(fmt.Sprintf("%s%d!", q.ToString(), correctAnswer))
						g.teacher.Say(fmt.Sprintf("%s, you are right!", s.Name))
					})
					if winner != nil && winner != s {
						s.Say(fmt.Sprintf("%s, you win.", winner.Name))
					}
				}(student)
			}

			g.wg.Wait()
			fmt.Println("---")
		}
	}
}

func (g *Game) Stop() {
	g.stop <- struct{}{}
}

func main() {
	g := NewGame()
	g.Start()
}
