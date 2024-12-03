package main

import (
	"math"
	"math/rand/v2"
)

const (
	MAX_STARS   = 3000
	MAX_IDELTAT = 50
	DELTAT      = MAX_IDELTAT * 0.0001
	EPSILON     = 0.00000001

	GALAXYRANGESIZE = 0.1
	GALAXYMINSIZE   = 0.15
	QCONS           = 0.001

	COLORBASE = 12
)

type Star struct {
	Pos [3]float64
	Vel [3]float64
}

type Galaxy struct {
	Mass      int
	Nstars    int
	Stars     []Star
	Oldpoints [][2]int
	Newpoints [][2]int
	Pos       [3]float64
	Vel       [3]float64
	Galcol    int
}

type Universe struct {
	Mat             [3][3]float64
	Size            float64
	Diff            [3]float64
	Galaxies        []Galaxy
	Ngalaxies       int
	F_hititerations int
	Step            int
	Rot_y           float64
	Rot_x           float64
}

func (gp *Universe) startover() {
	var i, j int
	var w1, w2 float64
	var d, v, w, h float64

	gp.Step = 0
	gp.Rot_y = 0
	gp.Rot_x = 0

	gp.Ngalaxies = 3 + rand.IntN(3)
	gp.Galaxies = make([]Galaxy, gp.Ngalaxies)

	for i = 0; i < gp.Ngalaxies; i++ {
		gt := &gp.Galaxies[i]
		var sinw1, sinw2, cosw1, cosw2 float64

		gt.Galcol = rand.IntN(COLORBASE)
		gt.Nstars = MAX_STARS/2 + rand.IntN(MAX_STARS/2)
		gt.Stars = make([]Star, gt.Nstars)
		gt.Oldpoints = make([][2]int, gt.Nstars)
		gt.Newpoints = make([][2]int, gt.Nstars)

		w1 = 2 * math.Pi * rand.Float64()
		w2 = 2 * math.Pi * rand.Float64()
		sinw1 = math.Sin(w1)
		sinw2 = math.Sin(w2)
		cosw1 = math.Cos(w1)
		cosw2 = math.Cos(w2)

		gp.Mat[0][0] = cosw2
		gp.Mat[0][1] = -sinw1 * sinw2
		gp.Mat[0][2] = cosw1 * sinw2
		gp.Mat[1][0] = 0.0
		gp.Mat[1][1] = cosw1
		gp.Mat[1][2] = sinw1
		gp.Mat[2][0] = -sinw2
		gp.Mat[2][1] = -sinw1 * cosw2
		gp.Mat[2][2] = cosw1 * cosw2

		gt.Vel[0] = rand.Float64()*2.0 - 1.0
		gt.Vel[1] = rand.Float64()*2.0 - 1.0
		gt.Vel[2] = rand.Float64()*2.0 - 1.0
		gt.Pos[0] = -gt.Vel[0]*DELTAT*float64(gp.F_hititerations) + rand.Float64() - 0.5
		gt.Pos[1] = -gt.Vel[1]*DELTAT*float64(gp.F_hititerations) + rand.Float64() - 0.5
		gt.Pos[2] = -gt.Vel[2]*DELTAT*float64(gp.F_hititerations) + rand.Float64() - 0.5

		gt.Mass = int(rand.Float64()*1000.0) + 1
		gp.Size = GALAXYRANGESIZE*rand.Float64() + GALAXYMINSIZE

		for j = 0; j < gt.Nstars; j++ {
			st := &gt.Stars[j]
			oldp := &gt.Oldpoints[j]
			newp := &gt.Newpoints[j]

			var sinw, cosw float64

			w = 2.0 * math.Pi * rand.Float64()
			sinw = math.Sin(w)
			cosw = math.Cos(w)
			d = rand.Float64() * gp.Size
			h = rand.Float64() * math.Exp(-2.0*(d/gp.Size)) / 5.0 * gp.Size
			if rand.Float64() < 0.5 {
				h = -h
			}
			st.Pos[0] = gp.Mat[0][0]*d*cosw + gp.Mat[1][0]*d*sinw + gp.Mat[2][0]*h + gt.Pos[0]
			st.Pos[1] = gp.Mat[0][1]*d*cosw + gp.Mat[1][1]*d*sinw + gp.Mat[2][1]*h + gt.Pos[1]
			st.Pos[2] = gp.Mat[0][2]*d*cosw + gp.Mat[1][2]*d*sinw + gp.Mat[2][2]*h + gt.Pos[2]
			v = math.Sqrt(float64(gt.Mass) * QCONS / math.Sqrt(d*d+h*h))
			st.Vel[0] = -gp.Mat[0][0]*v*sinw + gp.Mat[1][0]*v*cosw + gt.Vel[0]
			st.Vel[1] = -gp.Mat[0][1]*v*sinw + gp.Mat[1][1]*v*cosw + gt.Vel[1]
			st.Vel[2] = -gp.Mat[0][2]*v*sinw + gp.Mat[1][2]*v*cosw + gt.Vel[2]

			st.Vel[0] *= DELTAT
			st.Vel[1] *= DELTAT
			st.Vel[2] *= DELTAT

			oldp[0] = 0
			oldp[1] = 0
			newp[0] = 0
			newp[1] = 0
		}
	}
}

func InitGalaxy() Universe {
	universe := Universe{
		Mat:             [3][3]float64{},
		Size:            0,
		Diff:            [3]float64{},
		Galaxies:        nil,
		Ngalaxies:       0,
		F_hititerations: 250, // cycles, TODO make this an argument
		Step:            0,
		Rot_y:           0,
		Rot_x:           0,
	}
	universe.startover()
	return universe
}

func (gp *Universe) UpdateGalaxy(scale float64, midx, midy int) {
	var d, eps, cox, six, cor, sir float64
	var i, j, k int

	// Update rotation angles
	gp.Rot_y += 0.01
	gp.Rot_x += 0.004

	cox = math.Cos(gp.Rot_y)
	six = math.Sin(gp.Rot_y)
	cor = math.Cos(gp.Rot_x)
	sir = math.Sin(gp.Rot_x)

	eps = 1 / (EPSILON * math.Sqrt(EPSILON) * DELTAT * DELTAT * QCONS)

	for i = 0; i < gp.Ngalaxies; i++ {
		gt := &gp.Galaxies[i]

		// Swap old and new points
		gt.Oldpoints, gt.Newpoints = gt.Newpoints, gt.Oldpoints

		for j = 0; j < gt.Nstars; j++ {
			st := &gt.Stars[j]
			newp := &gt.Newpoints[j]
			v0 := st.Vel[0]
			v1 := st.Vel[1]
			v2 := st.Vel[2]

			for k = 0; k < gp.Ngalaxies; k++ {
				gtk := &gp.Galaxies[k]
				d0 := gtk.Pos[0] - st.Pos[0]
				d1 := gtk.Pos[1] - st.Pos[1]
				d2 := gtk.Pos[2] - st.Pos[2]

				d = d0*d0 + d1*d1 + d2*d2
				if d > EPSILON {
					d = float64(gtk.Mass) / (d * math.Sqrt(d)) * DELTAT * DELTAT * QCONS
				} else {
					d = float64(gtk.Mass) / (eps * math.Sqrt(eps))
				}
				v0 += d0 * d
				v1 += d1 * d
				v2 += d2 * d
			}

			st.Vel[0] = v0
			st.Vel[1] = v1
			st.Vel[2] = v2

			st.Pos[0] += v0
			st.Pos[1] += v1
			st.Pos[2] += v2

			newp[0] = int(((cox*st.Pos[0])-(six*st.Pos[2]))*scale) + midx
			newp[1] = int(((cor*st.Pos[1])-(sir*((six*st.Pos[0])+(cox*st.Pos[2]))))*scale) + midy
		}

		for k = i + 1; k < gp.Ngalaxies; k++ {
			gtk := &gp.Galaxies[k]
			d0 := gtk.Pos[0] - gt.Pos[0]
			d1 := gtk.Pos[1] - gt.Pos[1]
			d2 := gtk.Pos[2] - gt.Pos[2]

			d = d0*d0 + d1*d1 + d2*d2
			if d > EPSILON {
				d = 1 / (d * math.Sqrt(d)) * DELTAT * QCONS
			} else {
				d = 1 / (EPSILON * math.Sqrt(EPSILON)) * DELTAT * QCONS
			}

			d0 *= d
			d1 *= d
			d2 *= d
			gt.Vel[0] += d0 * float64(gtk.Mass)
			gt.Vel[1] += d1 * float64(gtk.Mass)
			gt.Vel[2] += d2 * float64(gtk.Mass)
			gtk.Vel[0] -= d0 * float64(gt.Mass)
			gtk.Vel[1] -= d1 * float64(gt.Mass)
			gtk.Vel[2] -= d2 * float64(gt.Mass)
		}

		gt.Pos[0] += gt.Vel[0] * DELTAT
		gt.Pos[1] += gt.Vel[1] * DELTAT
		gt.Pos[2] += gt.Vel[2] * DELTAT
	}

	gp.Step++
	if gp.Step > gp.F_hititerations*4 {
		gp.startover()
	}
}
