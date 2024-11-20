package gamemodes

import (
	"fmt"
	"math/rand"
	"time"
)

type Cella struct {
    is_flagged bool;
    is_bomb bool;
    label uint8;
    is_hidden bool;
}

type GameState int;
const (
    Running GameState = iota
    Won GameState = iota
    Lost GameState = iota
);

type Game struct {
    celle [][]Cella;
    state GameState;
    bomb_count uint16;
    flag_count uint16;
    tempo uint16
}

func NewGame(width uint16, height uint16, n_bombe uint16, tempo uint16) *Game {
    ret := Game {
        celle: make([][]Cella, height),
        state: Running,
        bomb_count: n_bombe,
        flag_count: 0,
        tempo: tempo,
    }

    r := rand.New(rand.NewSource(time.Now().UnixNano()));
    for y := range ret.celle {
        ret.celle[y] = make([]Cella, width);
        for x := range ret.celle[y] {
            ret.celle[y][x] = Cella{
                is_flagged: false,
                is_bomb: false,
                label: 0,
                is_hidden: true,
            };
        }
    }
    for i:=0; i<int(n_bombe); i++ {
        x := r.Int() % (int(width))
        y := r.Int() % (int(height))
        if ret.celle[y][x].is_bomb {i--; continue}
        ret.celle[y][x].is_bomb=true;
    }

    for y := range ret.celle {
        for x := range ret.celle[y] {
            if !ret.celle[y][x].is_bomb {continue;}
            for i:=-1; i<=1; i++ {
                for j:=-1; j<=1; j++ {
                    off_y, off_x := y+i, x+j;
                    if ret.is_inside(off_x,off_y) && !ret.celle[off_y][off_x].is_bomb {
                        ret.celle[off_y][off_x].label++;
                    }
                }
            }
        }
    }

    return &ret;
}

func (g *Game) get_h() uint16 {
    return uint16(len(g.celle))
}
func (g *Game) get_w() uint16 {
    if len(g.celle) == 0 {
        return 0;
    }
    return uint16(len(g.celle[0]))
}

func (g *Game) click(x uint16, y uint16) error {
    if err:=g.check_bounds(int(x),int(y)); err!=nil {return err;}

    stack := NewStack[*Cella]();
    for stack.len()>0 {
        corrente := stack.pop();
        if !corrente.is_hidden {continue;}
        if corrente.is_bomb {
            g.state = Lost;
            return nil;
        }
        corrente.is_hidden = false;

        if corrente.label > 0 {continue;}
        for i:=-1; i<=1; i++ {
            for j:=-1; j<=1; j++ {
                off_y := int(y)+i;
                off_x := int(x)+i;
                if !g.is_inside(off_x, off_y) {continue;}
                stack.push(&g.celle[off_y][off_x]);
            }
        }
    }

    return nil;
}

// Toggles the flagged state of a cell.
// Returns the new state.
// Errors if the coordinates lay outside the map
func (g *Game) flag(x uint16, y uint16) (bool, error) {
    if err:=g.check_bounds(int(x),int(y)); err!=nil {return false, err;}
    g.celle[y][x].is_flagged = !g.celle[y][x].is_flagged;
    if g.celle[y][x].is_flagged {
        g.flag_count++;
    }else{
        g.flag_count--;
    }

    return g.celle[y][x].is_flagged, nil;
}

func (g *Game) is_inside(x int, y int) bool {
    return x>=0 && y>=0 && x<int(g.get_w()) && y<int(g.get_h())
}
func (g *Game) check_bounds(x int, y int) error {
    if g.is_inside(x, y) {return nil;}
    return fmt.Errorf("Coordinates (%d,%d) are outside the bounds for gama (%d,%d)", x,y, g.get_w(), g.get_h());
}


