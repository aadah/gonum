// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"log"

	"golang.org/x/exp/rand"

	"gonum.org/v1/exp/layout/barneshut"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

type mass struct {
	d barneshut.Vector2
	v barneshut.Vector2
	m float64
}

func (m *mass) Coord2() barneshut.Vector2 { return m.d }
func (m *mass) Mass() float64             { return m.m }
func (m *mass) move(f barneshut.Vector2) {
	f.X /= m.m
	f.Y /= m.m
	m.v = m.v.Add(f)
	m.d = m.d.Add(m.v)
}

func main() {
	rnd := rand.New(rand.NewSource(1))

	// Make 1000 stars in random locations.
	stars := make([]*mass, 40)
	p := make([]barneshut.Particle2, len(stars))
	for i := range stars {
		s := &mass{
			d: barneshut.Vector2{
				X: 100*rnd.Float64() - 50,
				Y: 100*rnd.Float64() - 50,
			},
			v: barneshut.Vector2{
				X: rnd.NormFloat64(),
				Y: rnd.NormFloat64(),
			},
			m: rnd.Float64(),
		}
		stars[i] = s
		p[i] = s
	}
	vectors := make([]barneshut.Vector2, len(stars))

	tracks := make([]plotter.XYs, len(stars))

	// Make a plane to calculate approximate forces
	plane := barneshut.Plane{Particles: p}

	// Run a simulation for 1000 updates.
	for i := 0; i < 10000; i++ {
		// Build the data structure, For small system
		// this step may be omitted and ForceOn will
		// perform the naive quadratic calculation
		// without building the data structure.
		plane.Reset()

		// Calculate the force vectors using the theta
		// parameter.
		const theta = 0.6
		for j, s := range stars {
			vectors[j] = plane.ForceOn(s, theta, barneshut.Gravity2).Scale(100)
		}

		// Update positions.
		for j, s := range stars {
			s.move(vectors[j])
			tracks[j] = append(tracks[j], plotter.XY{X: s.d.X, Y: s.d.Y})
		}
	}

	plt, err := plot.New()
	if err != nil {
		log.Fatalf("failed create plot:", err)
	}
	for i, t := range tracks {
		l, err := plotter.NewLine(t)
		if err != nil {
			log.Fatalf("failed create track:", err)
		}
		l.Color = plotutil.Color(i)
		l.Dashes = plotutil.Dashes(i)
		plt.Add(l)
	}
	plt.X.Min = -1000
	plt.X.Max = 1000
	plt.Y.Min = -1000
	plt.Y.Max = 1000
	err = plt.Save(20*vg.Centimeter, 20*vg.Centimeter, "galaxy.svg")
	if err != nil {
		log.Fatalf("failed to save file:", err)
	}
}
