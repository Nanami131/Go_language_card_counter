package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"strings"
)

type CardCounter struct {
	cards    [][]int
	nums1    []string
	nums2    []string
	m        map[rune]int
	s        string
	output   *widget.RichText
	input    *widget.Entry
	zeroCard *widget.Entry
}

func NewCardCounter() *CardCounter {
	cards := make([][]int, 13)
	for i := range cards {
		cards[i] = make([]int, 4)
		for j := range cards[i] {
			cards[i][j] = -1
		}
	}

	s := "3456789TJQKA2"
	m := make(map[rune]int)
	for i := 3; i <= 9; i++ {
		m[rune(i+'0')] = i - 3
	}
	m['T'] = 7
	m['J'] = 8
	m['Q'] = 9
	m['K'] = 10
	m['A'] = 11
	m['t'] = 7
	m['j'] = 8
	m['q'] = 9
	m['k'] = 10
	m['a'] = 11
	m['2'] = 12

	return &CardCounter{
		cards: cards,
		nums1: make([]string, 0),
		nums2: make([]string, 0),
		m:     m,
		s:     s,
	}
}

func (cc *CardCounter) updateCards(val string) string {
	flag := false
	for i := 1; i < len(val); i++ {
		if _, ok := cc.m[rune(val[i])]; !ok && val[i] != '-' {
			return "输入错误！"
		}
	}
	if flag {
		return ""
	}

	id := -1
	zero := false
	switch val[0] {
	case 'A', 'a':
		id = 1
	case 'D', 'd':
		id = 2
	}

	h := val[1:]
	if id == 1 {
		cc.nums1 = append(cc.nums1, h)
	} else if id == 2 {
		cc.nums2 = append(cc.nums2, h)
	}

	if strings.Contains(val, "-") {
		tem := "W"
		zero = true
		for k := 0; k < len(cc.s); k++ {
			if cc.s[k] == val[1] {
				for cc.s[k] != val[3] {
					tem += string(cc.s[k])
					k++
				}
				tem += string(cc.s[k])
				val = tem
				break
			}
		}
	}

	if zero {
		id = 0
	}

	for i := 1; i < len(val); i++ {
		num := cc.m[rune(val[i])]
		for j := 0; j < 4; j++ {
			if cc.cards[num][j] == -1 {
				cc.cards[num][j] = id
				break
			}
		}
	}

	return cc.display()
}

func (cc *CardCounter) display() string {
	var output strings.Builder
	output.WriteString("3 4 5 6 7 8 9 T J Q K A 2\n")
	output.WriteString("--------------------------\n")

	// 特殊处理的计数行
	counts := make([]int, 13)
	for l := 0; l < 13; l++ {
		count := 0
		for k := 0; k < 4; k++ {
			if cc.cards[l][k] == -1 {
				count++
			}
		}
		counts[l] = count
		output.WriteString(fmt.Sprintf("%d ", count))
	}
	output.WriteString("\n")

	output.WriteString("--------------------------\n")
	output.WriteString("左: ")
	for l := 0; l < 13; l++ {
		a0, a1, a2 := 0, 0, 0
		for k := 0; k < 4; k++ {
			switch cc.cards[l][k] {
			case 0:
				a0++
			case 1:
				a1++
			case 2:
				a2++
			}
		}
		if a1 == 0 && a2 != 0 {
			for m := 0; m < 4-a0-a2; m++ {
				output.WriteString(string(cc.s[l]))
			}
		}
	}

	output.WriteString("\n--------------------------\n")
	output.WriteString("右: ")
	for l := 0; l < 13; l++ {
		a0, a1, a2 := 0, 0, 0
		for k := 0; k < 4; k++ {
			switch cc.cards[l][k] {
			case 0:
				a0++
			case 1:
				a1++
			case 2:
				a2++
			}
		}
		if a2 == 0 && a1 != 0 {
			for m := 0; m < 4-a0-a1; m++ {
				output.WriteString(string(cc.s[l]))
			}
		}
	}

	output.WriteString("\n--------------------------\n")
	output.WriteString("左行牌: " + strings.Join(cc.nums1, " ") + "\n")
	output.WriteString("--------------------------\n")
	output.WriteString("右行牌: " + strings.Join(cc.nums2, " ") + "\n")
	output.WriteString("--------------------------\n")

	return output.String()
}

func (cc *CardCounter) initCards(zero string) {
	for _, ch := range zero {
		num := cc.m[ch]
		for j := 0; j < 4; j++ {
			if cc.cards[num][j] == -1 {
				cc.cards[num][j] = 0
				break
			}
		}
	}
}

// CustomTextSegment 是一个支持自定义颜色的 TextSegment
type CustomTextSegment struct {
	Text      string
	Color     color.Color
	Alignment fyne.TextAlign
	TextStyle fyne.TextStyle
	IsInline  bool
}

// Inline 根据 IsInline 字段决定是否内联
func (c *CustomTextSegment) Inline() bool {
	return c.IsInline
}

func (c *CustomTextSegment) Textual() string {
	return c.Text
}

func (c *CustomTextSegment) Visual() fyne.CanvasObject {
	obj := canvas.NewText(c.Text, c.Color)
	obj.Alignment = c.Alignment
	obj.TextStyle = c.TextStyle
	return obj
}

func (c *CustomTextSegment) Update(o fyne.CanvasObject) {
	obj := o.(*canvas.Text)
	obj.Text = c.Text
	obj.Color = c.Color
	obj.Alignment = c.Alignment
	obj.TextStyle = c.TextStyle
	obj.Refresh()
}

// Select 未实现
func (c *CustomTextSegment) Select(begin, end fyne.Position) {}

// SelectedText 未实现
func (c *CustomTextSegment) SelectedText() string {
	return ""
}

// Unselect 未实现
func (c *CustomTextSegment) Unselect() {}

func (cc *CardCounter) updateOutput(text string) {
	segments := []widget.RichTextSegment{}

	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if line == "" {
			continue
		}
		if i == 2 {
			counts := strings.Fields(line)
			for _, count := range counts {
				seg := &CustomTextSegment{
					Text:     count + " ",
					IsInline: true,
				}
				switch count {
				case "4":
					seg.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255} // 红
				case "0":
					seg.Color = color.RGBA{R: 0, G: 0, B: 255, A: 255} // 蓝
				default:
					seg.Color = color.Black
				}
				segments = append(segments, seg)
			}

			segments = append(segments, &CustomTextSegment{
				Text:     "",
				IsInline: false,
			})
		} else {

			segments = append(segments, &CustomTextSegment{
				Text:     line,
				Color:    color.Black,
				IsInline: false,
			})
		}
	}

	cc.output.Segments = segments
	cc.output.Refresh()
}
func main() {
	a := app.New()
	w := a.NewWindow("Go语言记牌器")
	w.Resize(fyne.NewSize(1200, 800))

	cc := NewCardCounter()
	cc.output = widget.NewRichText()
	cc.output.Wrapping = fyne.TextWrapWord
	cc.updateOutput("欢迎使用记牌器！\n输入初始牌后点击初始化，然后输入当前牌点击提交。")

	cc.input = widget.NewEntry()
	cc.input.SetPlaceHolder("输入牌，如 A3 或 D5 或 3-5")
	cc.zeroCard = widget.NewEntry()
	cc.zeroCard.SetPlaceHolder("输入初始已出牌，如 345")

	submitBtn := widget.NewButton("提交", func() {
		if cc.input.Text != "" {
			result := cc.updateCards(cc.input.Text)
			cc.updateOutput(result)
			cc.input.SetText("")
		}
	})

	initBtn := widget.NewButton("初始化", func() {
		if cc.zeroCard.Text != "" {
			cc.initCards(cc.zeroCard.Text)
			cc.updateOutput(cc.display())
			cc.zeroCard.SetText("")
		}
	})

	inputArea := container.NewVBox(
		widget.NewLabel("初始已出牌："),
		cc.zeroCard,
		initBtn,
		widget.NewLabel("输入当前出牌："),
		cc.input,
		submitBtn,
	)

	content := container.NewBorder(
		inputArea, nil, nil, nil,
		container.NewScroll(cc.output),
	)

	w.SetContent(content)
	w.ShowAndRun()
}
