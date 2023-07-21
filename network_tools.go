package main

import (
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func New(layout fyne.Layout, objects ...fyne.CanvasObject) *fyne.Container {
	return fyne.NewContainerWithLayout(layout, objects...)
}

func NewAdaptiveGridWithRatios(ratios []float32, objects ...fyne.CanvasObject) *fyne.Container {
	return New(NewAdaptiveGridLayoutWithRatios(ratios), objects...)
}

// Declare conformity with Layout interface
var _ fyne.Layout = (*adaptiveGridLayoutWithRatios)(nil)

type adaptiveGridLayoutWithRatios struct {
	ratios          []float32
	adapt, vertical bool
}

func NewAdaptiveGridLayoutWithRatios(ratios []float32) fyne.Layout {
	return &adaptiveGridLayoutWithRatios{ratios: ratios, adapt: true}
}

func (g *adaptiveGridLayoutWithRatios) horizontal() bool {
	if g.adapt {
		return fyne.IsHorizontal(fyne.CurrentDevice().Orientation())
	}

	return !g.vertical
}

func (g *adaptiveGridLayoutWithRatios) countRows(objects []fyne.CanvasObject) int {
	count := 0
	for _, child := range objects {
		if child.Visible() {
			count++
		}
	}

	return int(math.Ceil(float64(count) / float64(len(g.ratios))))
}

// Layout is called to pack all child objects into a specified size.
// For a GridLayout this will pack objects into a table format with the number
// of columns specified in our constructor.
func (g *adaptiveGridLayoutWithRatios) Layout(objects []fyne.CanvasObject, size fyne.Size) {

	rows := g.countRows(objects)
	cols := len(g.ratios)

	padWidth := float32(cols-1) * float32(theme.Padding())
	padHeight := float32(rows-1) * float32(theme.Padding())
	tGap := float64(padWidth)
	tcellWidth := float64(size.Width) - tGap
	cellHeight := float64(size.Height-padHeight) / float64(rows)

	if !g.horizontal() {
		padWidth, padHeight = padHeight, padWidth
		tcellWidth = float64(size.Width-padWidth) - tGap
		cellHeight = float64(size.Height-padHeight) / float64(cols)
	}

	row, col := 0, 0
	i := 0
	var x1, x2, y1, y2 float32 = 0.0, 0.0, 0.0, 0.0
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		if i == 0 {
			x1 = 0
			y1 = 0
		} else {
			x1 = x2 + float32(theme.Padding())*float32(1)
			y1 = y2 - float32(cellHeight)
		}
		x2 = x1 + float32(tcellWidth*float64(g.ratios[i]))
		y2 = float32(cellHeight)

		child.Move(fyne.NewPos(x1, y1))
		child.Resize(fyne.NewSize((x2 - x1), y2-y1))

		if g.horizontal() {
			if (i+1)%cols == 0 {
				row++
				col = 0
			} else {
				col++
			}
		} else {
			if (i+1)%cols == 0 {
				col++
				row = 0
			} else {
				row++
			}
		}
		i++
	}
}

func (g *adaptiveGridLayoutWithRatios) MinSize(objects []fyne.CanvasObject) fyne.Size {
	rows := g.countRows(objects)
	minSize := fyne.NewSize(0, 0)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		minSize = minSize.Max(child.MinSize())
	}

	if g.horizontal() {
		minContentSize := fyne.NewSize(minSize.Width*float32(len(g.ratios)), minSize.Height*float32(rows))
		return minContentSize.Add(fyne.NewSize(float32(theme.Padding())*fyne.Max(float32(len(g.ratios)-1), 0), float32(theme.Padding())*fyne.Max(float32(rows-1), 0)))
	}

	minContentSize := fyne.NewSize(minSize.Width*float32(rows), minSize.Height*float32(len(g.ratios)))
	return minContentSize.Add(fyne.NewSize(float32(theme.Padding())*fyne.Max(float32(rows-1), 0), float32(theme.Padding())*fyne.Max(float32(len(g.ratios)-1), 0)))
}

func refresh_interface() []string {
	var interfacesList []string
	var interfaceInValid bool

	interfaces, err := net.Interfaces()

	if err != nil {
		panic(err)
	}

	for _, i := range interfaces {

		byNameInterface, err := net.InterfaceByName(i.Name)
		interfaceInValid = true
		if err != nil {
			fmt.Println(err)
		}
		addresses, err := byNameInterface.Addrs()
		if err != nil {
			fmt.Println(err)
		}

		for k, addr := range addresses {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPAddr:
				interfaceInValid = false
				ip = v.IP
			case *net.IPNet:
				interfaceInValid = false
				ip = v.IP
			default:
				interfaceInValid = true
				continue
			}

			mac := byNameInterface.HardwareAddr.String()
			if mac == "" {
				interfaceInValid = true
			}

			// print the available ip addresses
			fmt.Printf("%v Interface Address #%v : %v invalid %v\n", i.Name, k, ip.String(), interfaceInValid)
		}
		if !interfaceInValid {
			interfacesList = append(interfacesList, i.Name)
			//fmt.Printf("%v \n", i.Name)
		}
	}
	return interfacesList
}

func generate_context(index int, objects ...fyne.CanvasObject) *fyne.Container {
	var content *fyne.Container
	switch index {
	case 1:
		content = container.NewVBox(
			NewAdaptiveGridWithRatios([]float32{0.25, 0.25, 0.25, 0.25}, objects[0], objects[1], objects[2], objects[3]),
			NewAdaptiveGridWithRatios([]float32{0.25, 0.25, 0.5}, objects[7], objects[8], objects[4]),
			NewAdaptiveGridWithRatios([]float32{1}, objects[9]),
			NewAdaptiveGridWithRatios([]float32{1}, objects[4]),
			//			NewAdaptiveGridWithRatios([]float32{0.1, 0.9}, objects[5], objects[6]),
		)
	case 2:
		content = container.NewVBox(
			NewAdaptiveGridWithRatios([]float32{0.25, 0.25, 0.25, 0.25}, objects[0], objects[1], objects[2], objects[3]),
			NewAdaptiveGridWithRatios([]float32{0.25, 0.5, 0.25}, objects[7], objects[8], objects[9]),
			NewAdaptiveGridWithRatios([]float32{1}, objects[10]),
			//NewAdaptiveGridWithRatios([]float32{1}, objects[4]),
			NewAdaptiveGridWithRatios([]float32{0.1, 0.9}, objects[5], objects[6]),
		)
	case 3:
		content = container.NewVBox(
			NewAdaptiveGridWithRatios([]float32{0.25, 0.25, 0.25, 0.25}, objects[0], objects[1], objects[2], objects[3]),
			NewAdaptiveGridWithRatios([]float32{0.25, 0.35, 0.15, 0.25}, objects[7], objects[8], objects[9], objects[10]),
			NewAdaptiveGridWithRatios([]float32{0.25, 0.35, 0.40}, objects[12], objects[13], objects[4]),
			NewAdaptiveGridWithRatios([]float32{1}, objects[11]),
			//NewAdaptiveGridWithRatios([]float32{1}, objects[4]),
			NewAdaptiveGridWithRatios([]float32{0.1, 0.9}, objects[5], objects[6]),
		)
	case 4:
		content = container.NewVBox(
			NewAdaptiveGridWithRatios([]float32{0.25, 0.25, 0.25, 0.25}, objects[0], objects[1], objects[2], objects[3]),
			NewAdaptiveGridWithRatios([]float32{1}, objects[4]),
			NewAdaptiveGridWithRatios([]float32{0.25, 0.5, 0.25}, objects[4], objects[5], objects[4]),
		)
	default:
		fmt.Println("unknow context")
		content = container.NewVBox(
			NewAdaptiveGridWithRatios([]float32{0.25, 0.25, 0.25, 0.25}, objects[0], objects[1], objects[2], objects[3]),
		)
	}

	return content
}

func main() {

	var buttonSelect1 *widget.Button
	var buttonSelect2 *widget.Button
	var buttonSelect3 *widget.Button
	var buttonSelect4 *widget.Button
	var combo_interface *widget.Select

	a := app.New()
	w := a.NewWindow("fynetest")

	w.Resize(fyne.NewSize(600, 200))

	label_interface := widget.NewLabel("Interface : ")
	label_nslookup := widget.NewLabel("nslookup : ")
	label_wol := widget.NewLabel("Wol Mac : ")
	label_empty := widget.NewLabel("")
	data := widget.NewMultiLineEntry()
	data.Wrapping = fyne.TextTruncate
	data.SetMinRowsVisible(20)
	data2 := widget.NewMultiLineEntry()
	data2.Wrapping = fyne.TextTruncate
	data2.SetMinRowsVisible(20)
	data3 := widget.NewMultiLineEntry()
	data3.Wrapping = fyne.TextTruncate
	data3.SetMinRowsVisible(20)
	data4 := widget.NewLabel("\nStéphane BARTHELEMY\nJuillet 2023\n")
	data4.Alignment = fyne.TextAlignCenter
	entry_nsllokup := widget.NewEntry()
	entry_wolport := widget.NewEntry()
	entry_wolport.SetText("9")
	entry_wolmac := widget.NewEntry()
	entry_wolmac.SetPlaceHolder("11:22:33:44:55:66")
	entry_broadcastaddr := widget.NewEntry()
	entry_broadcastaddr.SetPlaceHolder("255.255.255.255")

	button_wol := widget.NewButton("Run", func() {

		var mp [102]byte
		log.Println(("Wol"))
		text := ""
		mac, err := net.ParseMAC(entry_wolmac.Text)
		if err != nil {
			log.Println("Error Parqing  mac")
			text = "Error Parsing Mac "
			data3.SetText(text)
			return
		}
		if len(mac) != 6 {
			log.Println("Error Len mac")
			text = "Error Len Mac "
			data3.SetText(text)
			return
		}

		addr := entry_broadcastaddr.Text + ":" + entry_wolport.Text
		udpAddr, err := net.ResolveUDPAddr("udp", addr)

		offset := 0
		copy(mp[0:], []byte{255, 255, 255, 255, 255, 255})
		offset += 6
		for i := 0; i < 16; i++ {
			copy(mp[offset:], mac)
			offset += 6
		}

		byNameInterface, err := net.InterfaceByName(combo_interface.Selected)
		if err != nil {
			fmt.Println(err)
		}

		addresses, err := byNameInterface.Addrs()
		if err != nil {
			fmt.Println(err)
		}

		var localUdpIp *net.UDPAddr

		for _, localaddr := range addresses {
			switch localip := localaddr.(type) {
			case *net.IPNet:
				if !localip.IP.IsLoopback() && localip.IP.To4() != nil {
					localUdpIp = &net.UDPAddr{IP: localip.IP}
				}
			}
		}

		conn, err := net.DialUDP("udp", localUdpIp, udpAddr)
		if err != nil {
			log.Println("Error set udp socket")
			text = "Error set udp socket " + err.Error()
			data3.SetText(text)
			return
		}

		//addr := "255.255.255.255" + ":" + entry_wolport.Text
		//conn, err := net.Dial("udp", addr)
		//if err != nil {
		//	log.Println("Error set udp socket")
		//	text = "Error set udp socket " + err.Error()
		//	data3.SetText(text)
		//	return
		//}
		defer conn.Close()

		_, err = conn.Write(mp[:])
		if err != nil {
			log.Println("Magic Packet error sent" + err.Error())
			text = "Magic Packet error sent" + err.Error()
			data3.SetText(text)
			return
		}

		conn.Close()

		text = "Magic Packet sent : \n"
		for i := 0; i < 102; i++ {
			text += fmt.Sprintf("%02x", mp[i])
			if ((i + 1) % 16) == 0 {
				text += "\n"
			} else {
				text += " "
			}
		}
		log.Println(text)
		data3.SetText(text)
	})
	button_nslookup := widget.NewButton("Run", func() {
		text := ""
		ips, err := net.LookupIP(entry_nsllokup.Text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not get IPs: %v\n", err)
			text += "Could not get IPs: " + entry_nsllokup.Text + "\n"
		}
		for _, ip := range ips {
			fmt.Printf("%v . IN A %s\n", entry_nsllokup.Text, ip.String())
			text += entry_nsllokup.Text + " IN A " + ip.String() + "\n"
		}
		data2.SetText(text)
	})

	pgBar := widget.NewLabelWithStyle("Progress :", fyne.TextAlignCenter, fyne.TextStyle{Italic: true})
	progressBar := widget.NewProgressBar()
	progressBar.SetValue(0)

	combo_interface = widget.NewSelect(refresh_interface(), func(value string) {
		log.Println("Select set to", value)
		text := ""

		byNameInterface, err := net.InterfaceByName(value)

		if err != nil {
			fmt.Println(err)
		}

		addresses, err := byNameInterface.Addrs()

		if err != nil {
			fmt.Println(err)
		}
		text += value + " : \n"
		text += "MAC : " + byNameInterface.HardwareAddr.String() + "\n"
		for k, addr := range addresses {
			fmt.Printf("Interface Address #%v : %v\n", k, addr.String())
			text += "Interface Address"
			text += strconv.Itoa(k)
			text += " : "
			text += addr.String()
			text += "\n"
		}
		fmt.Println("------------------------------------")
		text += "------------------------------------\n"
		data.SetText(text)

	})

	buttonSelect1 = widget.NewButton("Interface", func() {
		log.Println("Select set to Select 1")
		content := generate_context(1, buttonSelect1, buttonSelect2, buttonSelect3, buttonSelect4,
			label_empty, pgBar, progressBar,
			label_interface, combo_interface, data,
		)
		w.SetContent(content)
	})
	buttonSelect2 = widget.NewButton("NS Lookup", func() {
		log.Println("Select set to Select 2")
		content := generate_context(2, buttonSelect1, buttonSelect2, buttonSelect3, buttonSelect4,
			label_empty, pgBar, progressBar,
			label_nslookup, entry_nsllokup, button_nslookup, data2,
		)
		w.SetContent(content)
	})
	buttonSelect3 = widget.NewButton("Wake on lan", func() {
		log.Println("Select set to Select 3")
		content := generate_context(3, buttonSelect1, buttonSelect2, buttonSelect3, buttonSelect4,
			label_empty, pgBar, progressBar,
			label_wol, entry_wolmac, entry_wolport, button_wol, data3,
			label_interface, entry_broadcastaddr,
		)
		w.SetContent(content)
	})
	buttonSelect4 = widget.NewButton("Info", func() {
		log.Println("Select set to Select 4")
		content := generate_context(4, buttonSelect1, buttonSelect2, buttonSelect3, buttonSelect4,
			label_empty, data4,
		)
		w.SetContent(content)
	})

	//	content := container.NewVBox(
	//		NewAdaptiveGridWithRatios([]float32{0.25, 0.25, 0.25, 0.25}, buttonSelect1, buttonSelect2, buttonSelect3, buttonSelect4),
	//		NewAdaptiveGridWithRatios([]float32{0.25, 0.25, 0.5}, label_interface, combo_interface, label_empty),
	//		NewAdaptiveGridWithRatios([]float32{1}, data),
	//		NewAdaptiveGridWithRatios([]float32{0.1, 0.9}, pgBar, progressBar),
	//	)
	content := generate_context(1, buttonSelect1, buttonSelect2, buttonSelect3, buttonSelect4,
		label_empty, pgBar, progressBar,
		label_interface, combo_interface, data,
	)
	w.SetContent(content)

	w.ShowAndRun()
}
