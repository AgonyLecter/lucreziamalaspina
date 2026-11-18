// Lucrezia Malaspina — versione desktop per Windows. Sviluppato da AgonyApp.
package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"

	webview2 "github.com/jchv/go-webview2"
)

//go:embed game
var gameFS embed.FS

var (
	user32            = syscall.NewLazyDLL("user32.dll")
	pGetWindowLongW   = user32.NewProc("GetWindowLongW")
	pSetWindowLongW   = user32.NewProc("SetWindowLongW")
	pSetWindowPos     = user32.NewProc("SetWindowPos")
	pMonitorFromWin   = user32.NewProc("MonitorFromWindow")
	pGetMonitorInfoW  = user32.NewProc("GetMonitorInfoW")
	pMessageBoxW      = user32.NewProc("MessageBoxW")
)

type rect struct{ L, T, R, B int32 }
type monitorInfo struct {
	Size    uint32
	Monitor rect
	Work    rect
	Flags   uint32
}

func fullscreen(hwnd uintptr) {
	const gwlStyle = ^uintptr(15) // -16
	style, _, _ := pGetWindowLongW.Call(hwnd, gwlStyle)
	style &^= 0x00CF0000 // WS_OVERLAPPEDWINDOW
	pSetWindowLongW.Call(hwnd, gwlStyle, style)
	mon, _, _ := pMonitorFromWin.Call(hwnd, 2)
	mi := monitorInfo{Size: uint32(unsafe.Sizeof(monitorInfo{}))}
	pGetMonitorInfoW.Call(mon, uintptr(unsafe.Pointer(&mi)))
	r := mi.Monitor
	pSetWindowPos.Call(hwnd, 0, uintptr(r.L), uintptr(r.T), uintptr(r.R-r.L), uintptr(r.B-r.T), 0x0020|0x0040)
}

func msg(text string) {
	t, _ := syscall.UTF16PtrFromString(text)
	c, _ := syscall.UTF16PtrFromString("Lucrezia Malaspina")
	pMessageBoxW.Call(0, uintptr(unsafe.Pointer(t)), uintptr(unsafe.Pointer(c)), 0x10)
}

func main() {
	sub, _ := fs.Sub(gameFS, "game")
	var ln net.Listener
	var err error
	for p := 47129; p < 47139; p++ { // porta fissa: così i salvataggi restano sempre gli stessi
		if ln, err = net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", p)); err == nil {
			break
		}
	}
	if err != nil {
		msg("Impossibile avviare il gioco / Unable to start the game.")
		return
	}
	go http.Serve(ln, http.FileServer(http.FS(sub)))

	cfg, _ := os.UserConfigDir()
	data := filepath.Join(cfg, "AgonyApp", "LucreziaMalaspina")
	os.MkdirAll(data, 0o755)

	w := webview2.NewWithOptions(webview2.WebViewOptions{
		AutoFocus: true,
		DataPath:  data,
		WindowOptions: webview2.WindowOptions{
			Title: "Lucrezia Malaspina", Width: 1280, Height: 720, IconId: 1, Center: true,
		},
	})
	if w == nil {
		msg("Per giocare serve Microsoft Edge WebView2 (già incluso in Windows 10 e 11).\nScaricalo da: https://go.microsoft.com/fwlink/p/?LinkId=2124703\n\nThis game needs Microsoft Edge WebView2 (included in Windows 10 and 11).")
		return
	}
	defer w.Destroy()
	w.Bind("lmQuit", func() { w.Dispatch(func() { w.Terminate() }) })
	fullscreen(uintptr(w.Window()))
	w.Navigate(fmt.Sprintf("http://%s/index.html", ln.Addr().String()))
	w.Run()
}
