package v3

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/rudy3091/mooze/v3/event"
	"github.com/rudy3091/mooze/v3/ui"
	"golang.org/x/term"
)

func openTty() *os.File {
	in, err := os.OpenFile("/dev/tty", syscall.O_RDONLY, 0)
	if err != nil {
		log.Fatal("cannot find /dev/tty")
		panic(err)
	}
	return in
}

func getTtyFd() int {
	return int(openTty().Fd())
}

func makeRaw() (int, *term.State) {
	fd := getTtyFd()
	state, err := term.MakeRaw(fd)
	if err != nil {
		log.Fatal("cannot make terminal raw")
		panic(err)
	}
	return fd, state
}

func restoreRaw(fd int, s *term.State) {
	term.Restore(fd, s)
}

// returns (width, height)
func getSize() (int, int) {
	x, y, err := term.GetSize(getTtyFd())
	if err != nil {
		log.Fatal("cannot get terminal size")
		panic(err)
	}
	return x, y
}

func handleTermResize(sigchan chan os.Signal) {
	for {
		<-sigchan
		width, height := getSize()
		ui.WindowStore[0].Resize(15, width/3)
		ui.WindowStore[1].Resize(height-4, width-(width/3)).Relocate(1, width/3+1)
		ui.LensWindow.Resize(3, width).Relocate(height-3, 1)
		ui.ClearScreen()
		ui.ReloadAll()
		debug(width, width, "resized")
	}
}

func debug(w int, e interface{}, s interface{}) {
	ui.MoveCursorTo(w-1, 1)
	ui.Print(strings.Repeat(" ", w))
	ui.Print("\x1b[31mDEBUG: ", e, s, "\x1b[0m")
}

func Run() {
	ui.LoadAlternateScreen()
	ui.HideCursor()
	fd, state := makeRaw()
	width, height := getSize()
	req := NewRequest()

	defer restoreRaw(fd, state)
	defer ui.UnloadAlternateScreen()
	defer ui.ShowCursor()

	// key input event channel
	ev := make(chan event.Event, 1)
	// system event channel
	sigs := make(chan os.Signal, 1)
	// done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGWINCH)

	go event.EmitKeyEvent(getTtyFd(), ev)
	go handleTermResize(sigs)

	var main *ui.Window

	ui.NewWindow(1, 1, 15, width/3).
		Title("Ops").
		Content([]string{
			"info",
			"update target url",
			"change request method",
			"set request header",
			"set request body",
		}).
		/*
		 * registering event handlers
		 */
		OnSelect([]func(){
			func() {
				main.Clear()
				main.Content([]string{
					" __   __  _______  _______  _______  _______ ",
					"|  |_|  ||       ||       ||       ||       |",
					"|       ||   _   ||   _   ||____   ||    ___|",
					"|       ||  | |  ||  | |  | ____|  ||   |___ ",
					"|       ||  |_|  ||  |_|  || ______||    ___|",
					"| ||_|| ||       ||       || |_____ |   |___ ",
					"|_|   |_||_______||_______||_______||_______|",
					"",
					"Mooze 0.2.0",
					"https://github.com/rudy3091/mooze",
				}).Render()
			},
			func() {
				ui.OpenPopup("url", width, height, func(url string) {
					RequestInfo.url = url
				})
				<-ev
			},
			func() {
				ui.OpenPopup("method", width, height, func(method string) {
					RequestInfo.method = method
				})
				<-ev
			},
		}).
		Render()

	ui.NewWindow(16, 1, 15, width/3).
		Title("Request").
		Content([]string{
			"target:",
			"method:",
			"header:",
			"body:",
		}).
		Render()

	main = ui.NewWindow(1, width/3+1, height-4, width-(width/3)).
		Title("Main").
		Content([]string{
			" __   __  _______  _______  _______  _______ ",
			"|  |_|  ||       ||       ||       ||       |",
			"|       ||   _   ||   _   ||____   ||    ___|",
			"|       ||  | |  ||  | |  | ____|  ||   |___ ",
			"|       ||  |_|  ||  |_|  || ______||    ___|",
			"| ||_|| ||       ||       || |_____ |   |___ ",
			"|_|   |_||_______||_______||_______||_______|",
			"",
			"Mooze 0.2.0",
			"https://github.com/rudy3091/mooze",
		})
	main.Render()

	ui.LensWindow = ui.NewWindow(height-3, 1, 3, width).
		Disable().
		Title("Lens")
	ui.LensWindow.Render()

	ui.WindowStore[0].Focus()
	ui.WindowStore[0].Render()

	for {
		w, _ := getSize()
		e := <-ev
		debug(w, e, "")

		switch e.Buf[0] {
		case 113:
			// q
			return
		case 9:
			// tab
			ui.RotateFocus()
		case 106:
			// j
			ui.NextItem()
		case 107:
			// k
			ui.PrevItem()
		case 100:
			// d
			ui.ScrollHalfDown()
		case 117:
			// u
			ui.ScrollHalfUp()
		case 13:
			// enter
			ui.Select()
		case 115:
			// s
			res, _, err := RequestInfo.Send()
			if err != nil {
				ui.MoveCursorTo(w-1, 1)
				ui.Print(err)
			} else {
				main.Content(
					strings.Split(req.Json(res), "\n"),
				)
				main.Clear()
				main.Render()
			}
		case 104:
			// h
			ui.PrevPage()
		case 108:
			// l
			ui.NextPage()
		case 8:
			// ctrl l
			// scroll left
			ui.ScrollLeft()
		case 12:
			// ctrl h
			// scroll right
			ui.ScrollRight()
		case 112:
			// // p
			// input := ui.OpenPopup(w, h, func(url string) {
			// 	RequestInfo.url = url
			// })
			// debug(w, input, "")
			// <-ev // consume unwanted input
			// case 99:
			// 	ui.ClosePopup()
		}
	}
}
