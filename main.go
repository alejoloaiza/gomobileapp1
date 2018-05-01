package main

import (
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/ahmdrz/goinsta"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/gl"
)

func main() {
	app.Main(func(a app.App) {
		var glctx gl.Context

		sz := size.Event{}
		for {
			select {
			case e := <-a.Events():
				switch e := a.Filter(e).(type) {
				case lifecycle.Event:
					glctx, _ = e.DrawContext.(gl.Context)
				case size.Event:
					sz = e
				case paint.Event:
					if glctx == nil {
						continue
					}
					onDraw(glctx, sz)
					a.Publish()
				case touch.Event:
					if time.Now().After(t2) {
						t2 = time.Now().Add(time.Second * 5)
						ok = !ok
						a.Send(paint.Event{})
					}
					if counter == 0 {
						_ = GetConfig("config.json")
						Uploadlists()
						counter++
						go InstagramMain()
					}
				}
			}
		}
	})
}

var (
	ok           = false
	counter      = 0
	FemaleNames  = make(map[string]int)
	t2           = time.Now().Add(time.Second * 2)
	insta        *goinsta.Instagram
	myUsers      []FollowingUser
	myInboxUsers = make(map[string]int)
)

type FollowingUser struct {
	ID       int64
	Username string
	Fullname string
}

func InstagramMain() {
	rand.Seed(time.Now().UnixNano())

	insta = goinsta.New(Localconfig.InstaUser, Localconfig.InstaPass)
	err := insta.Login()
	if err != nil {
		panic(err)
	}
	users, err := insta.UserFollowing(insta.InstaType.LoggedInUser.ID, "")
	var response FollowingUser
	if err != nil {
		return
	}
	inbox, err := insta.GetV2Inbox()
	if err != nil {
		return
	}
	for _, thread := range inbox.Inbox.Threads {
		for _, userthreads := range thread.Users {
			myInboxUsers[userthreads.Username] = 1
		}
	}
	for _, user := range users.Users {
		if myInboxUsers[user.Username] != 1 {
			fullname := strings.Split(user.FullName, " ")
			firstname := strings.ToLower(fullname[0])
			response.Username = user.Username
			response.ID = user.ID
			response.Fullname = firstname
			myUsers = append(myUsers, response)
		}
	}

	for _, dmuser := range myUsers {
		DirectMessage(dmuser.Username, dmuser.Fullname, dmuser.ID)
		time.Sleep(2 * time.Minute)
	}

	defer insta.Logout()
}
func Uploadlists() {
	femalelistraw := Localconfig.FemaleNames
	for _, bname3 := range femalelistraw {
		FemaleNames[bname3] = 1
	}
}

func PrepareMessage(Message string, NameOfUser string) string {
	resp := ""
	if FemaleNames[strings.ToLower(NameOfUser)] == 1 {
		resp = strings.Replace(Message, "{name}", strings.ToLower(NameOfUser), 1)
	} else {
		resp = strings.Replace(Message, "{name}", "", 1)
	}
	return resp

}
func DirectMessage(To string, Name string, Id int64) {

	Message := Localconfig.Sentences[random(0, 9)]
	newMessage := PrepareMessage(Message, Name)

	_, err := insta.DirectMessage(strconv.FormatInt(Id, 10), newMessage)
	if err != nil {
		panic(err)
	}
}
func random(min int, max int) int {
	return rand.Intn(max-min) + min
}
func onDraw(glctx gl.Context, sz size.Event) {
	if ok {
		glctx.ClearColor(1, 1, 1, 1)
	} else {
		glctx.ClearColor(0, 0, 0, 1)
	}

	glctx.Clear(gl.COLOR_BUFFER_BIT)
}
