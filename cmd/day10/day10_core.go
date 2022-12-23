package day10

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

type Sample struct {
	Cycle int
	Value int
}

func NewSample(cycle, value int) Sample {
	return Sample{Cycle: cycle, Value: value}
}

type CPU struct {
	X             int
	CRTPosition   int
	CRTScreen     Output
	Cycle         int
	XSampleCycles []int
	XSampleBuffer []Sample
}

func NewCPU() *CPU {
	return &CPU{X: 1, Cycle: 1, CRTPosition: 0, CRTScreen: NewNullOutput(), XSampleBuffer: make([]Sample, 0)}
}

func ShouldDrawSprite(spritePosition, beam int) bool {
	return beam >= spritePosition-1 && beam <= spritePosition+1
}

func (c *CPU) RunInstruction(i Instruction) {
	for cycle := i.CycleCount(); cycle > 0; cycle-- {
		// Check if we should draw the pixel
		if ShouldDrawSprite(c.X, c.CRTPosition) {
			c.CRTScreen.DrawLitPixel()
		} else {
			c.CRTScreen.DrawDarkPixel()
		}

		for _, sampleCycle := range c.XSampleCycles {
			if c.Cycle == sampleCycle {
				// Sample X now.
				c.XSampleBuffer = append(c.XSampleBuffer, NewSample(c.Cycle, c.X))
			}
		}

		c.Cycle++
		c.CRTPosition++
		if c.CRTPosition > 39 {
			c.CRTPosition = 0
			c.CRTScreen.NextLine()
		}
	}

	i.UpdateCPUState(c)
}

func (c *CPU) SetXSampleCycles(sampleCycles []int) {
	c.XSampleCycles = sampleCycles
	sort.Ints(c.XSampleCycles)
}

func (c *CPU) GetXSampleBuffer() []Sample {
	return c.XSampleBuffer
}

func (c *CPU) SetOutput(o Output) {
	c.CRTScreen = o

	o.SetScreenDimensions(40, 6)
	o.Reset()
}

type Output interface {
	SetScreenDimensions(width, height int)

	DrawLitPixel()
	DrawDarkPixel()

	Reset()
	NextLine()
}

type NullOutput struct {
}

func NewNullOutput() *NullOutput {
	return &NullOutput{}
}

func (o *NullOutput) SetScreenDimensions(width, height int) {
}

func (o *NullOutput) DrawLitPixel() {
}

func (o *NullOutput) DrawDarkPixel() {
}

func (o *NullOutput) Reset() {
}

func (o *NullOutput) NextLine() {
}

type StdoutOutput struct {
}

func NewStdoutOutput() *StdoutOutput {
	return &StdoutOutput{}
}

func (o *StdoutOutput) SetScreenDimensions(width, height int) {
}

func (o *StdoutOutput) DrawLitPixel() {
	fmt.Printf("#")
}

func (o *StdoutOutput) DrawDarkPixel() {
	fmt.Printf(".")
}

func (o *StdoutOutput) Reset() {
}

func (o *StdoutOutput) NextLine() {
	fmt.Printf("\n")
}

type BufferedOutput struct {
	Width  int
	Height int

	CurrentX int
	CurrentY int

	Screen [][]byte
}

func NewBufferedOutput() *BufferedOutput {
	o := &BufferedOutput{}
	o.SetScreenDimensions(40, 6)

	return o
}

func (o *BufferedOutput) SetScreenDimensions(width, height int) {
	o.Width = width
	o.Height = height

	o.Screen = make([][]byte, o.Height)

	for i := range o.Screen {
		o.Screen[i] = make([]byte, o.Width)
	}
}

func (o *BufferedOutput) DrawLitPixel() {
	currentRow := o.Screen[o.CurrentY]

	currentRow[o.CurrentX] = '#'
	o.CurrentX++
}

func (o *BufferedOutput) DrawDarkPixel() {
	currentRow := o.Screen[o.CurrentY]

	currentRow[o.CurrentX] = '.'
	o.CurrentX++
}

func (o *BufferedOutput) Reset() {
	o.CurrentX = 0
	o.CurrentY = 0
}

func (o *BufferedOutput) NextLine() {
	o.CurrentX = 0

	o.CurrentY++
	if o.CurrentY >= o.Height {
		o.CurrentY = 0
	}
}

type Instruction interface {
	Describe() string
	CycleCount() int

	UpdateCPUState(cpu *CPU)
}

type Noop struct {
}

func NewNoop() Noop {
	return Noop{}
}

func (i Noop) Describe() string {
	return "noop"
}

func (i Noop) CycleCount() int {
	return 1
}

func (i Noop) UpdateCPUState(cpu *CPU) {
}

type Addx struct {
	Value int
}

func NewAddx(v int) Addx {
	return Addx{Value: v}
}

func (i Addx) Describe() string {
	return fmt.Sprintf("addx %d", i.Value)
}

func (i Addx) CycleCount() int {
	return 2
}

func (i Addx) UpdateCPUState(cpu *CPU) {
	cpu.X += i.Value
}

func ParseInstruction(str string) (Instruction, error) {
	elements := strings.Split(str, " ")

	switch elements[0] {
	case "noop":
		return NewNoop(), nil
	case "addx":
		var value int
		c, err := fmt.Sscanf(str, "addx %d", &value)
		if err != nil {
			return nil, err
		}
		if c != 1 {
			return nil, errors.New("malformed instruction")
		}
		return NewAddx(value), nil
	}

	return nil, errors.New("unknown instruction")
}

func ParseInstructions(lines []string) ([]Instruction, error) {
	instructions := make([]Instruction, 0)
	for _, line := range lines {
		instruction, err := ParseInstruction(line)
		if err != nil {
			return []Instruction{}, err
		}

		instructions = append(instructions, instruction)
	}

	return instructions, nil
}

func day10(fileContents string) error {
	// Scan computer program in.
	instructions, err := ParseInstructions(strings.Split(string(fileContents), "\n"))
	if err != nil {
		return err
	}

	// Part 1: Find the signal strength during the 20th, 60th, 100th, 140th, 180th, and 220th cycles. What is the sum of these six signal strengths?
	c := NewCPU()

	cyclesToSampleX := []int{20, 60, 100, 140, 180, 220}
	c.SetXSampleCycles(cyclesToSampleX)

	for _, i := range instructions {
		c.RunInstruction(i)
	}

	totalSignalStrength := 0

	for _, sample := range c.GetXSampleBuffer() {
		totalSignalStrength += sample.Cycle * sample.Value
	}

	fmt.Printf("Cycle count: %d\n", c.Cycle)
	fmt.Printf("Total signal strength: %d\n", totalSignalStrength)

	// Part 2: Register X is the sprite location register.  If the CRT beam horizontal counter is +/-1 of X, then draw a lit pixel.  Otherwise, draw
	// a dark one.  Given a 40x6 "screen", what 8 capital letters are displayed?
	c2 := NewCPU()

	c2.SetOutput(NewStdoutOutput())

	for _, i := range instructions {
		c2.RunInstruction(i)
	}

	return nil
}
