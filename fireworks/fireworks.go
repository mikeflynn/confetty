package fireworks

import (
	"math"
	"math/rand"
	"time"

	"github.com/mikeflynn/confetty/array"
	"github.com/mikeflynn/confetty/simulation"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/harmonica"
	"github.com/charmbracelet/lipgloss"
)

const (
	framesPerSecond = 30.0
	numParticles    = 50
)

var (
	colors     = []string{"#a864fd", "#29cdff", "#78ff44", "#ff718d", "#fdff6a"}
	characters = []string{"+", "*", "•"}
	head       = "▄"
	tail       = "│"
)

type frameMsg time.Time

func animate() tea.Cmd {
	return tea.Tick(time.Second/framesPerSecond, func(t time.Time) tea.Msg {
		return frameMsg(t)
	})
}

type Model struct {
	system *simulation.System
}

func SpawnShoot(width, height int) *simulation.Particle {
	color := lipgloss.Color(array.Sample(colors))
	v := float64(rand.Intn(15) + 15.0)
	x := rand.Float64() * float64(width)
	p := simulation.Particle{
		Physics: harmonica.NewProjectile(
			harmonica.FPS(framesPerSecond),
			harmonica.Point{X: x, Y: float64(height)},
			harmonica.Vector{X: 0, Y: -v},
			harmonica.TerminalGravity,
		),
		Char:          lipgloss.NewStyle().Foreground(color).Render(head),
		TailChar:      lipgloss.NewStyle().Foreground(color).Render(tail),
		Color:         color,
		Shooting:      true,
		ExplosionCall: SpawnExplosion,
	}
	return &p
}

func SpawnExplosion(color lipgloss.Color, x, y float64, width, height int) []*simulation.Particle {
	v := float64(rand.Intn(10) + 20.0)
	particles := []*simulation.Particle{}
	for i := 0; i < numParticles; i++ {
		p := simulation.Particle{
			Physics: harmonica.NewProjectile(
				harmonica.FPS(framesPerSecond),
				harmonica.Point{X: x, Y: y},
				harmonica.Vector{X: math.Cos(float64(i)) * v, Y: math.Sin(float64(i)) * v / 2},
				harmonica.TerminalGravity,
			),
			Char:     lipgloss.NewStyle().Foreground(color).Render(array.Sample(characters)),
			Shooting: false,
		}
		particles = append(particles, &p)
	}
	return particles
}

func InitialModel() Model {
	return Model{system: &simulation.System{
		Particles: []*simulation.Particle{},
		Frame:     simulation.Frame{},
	}}
}

// Init initializes the confetti after a small delay
func (m Model) Init() tea.Cmd {
	return animate()
}

// Update updates the model every frame, it handles the animation loop and
// updates the particle physics every frame
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
		m.system.Particles = append(m.system.Particles, SpawnShoot(m.system.Frame.Width, m.system.Frame.Height))

		return m, nil
	case frameMsg:
		m.system.Update()
		return m, animate()
	case tea.WindowSizeMsg:
		if m.system.Frame.Width == 0 && m.system.Frame.Height == 0 {
			// For the first frameMsg spawn a system of particles
			m.system.Particles = append(m.system.Particles, SpawnShoot(msg.Width, msg.Height))
		}
		m.system.Frame.Width = msg.Width
		m.system.Frame.Height = msg.Height
		return m, nil
	default:
		return m, nil
	}
}

// View displays all the particles on the screen
func (m Model) View() string {
	return m.system.Render()
}
