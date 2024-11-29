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
type CellaCoords struct {
    x uint16;
    y uint16;
    cella Cella;
}
type CellaCoordsRef struct {
    x uint16;
    y uint16;
    cella *Cella;
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

func NewGame(width, height, n_bombe, tempo, first_x,first_y uint16) *Game {
    ret := Game {
        celle: make([][]Cella, height),
        state: Running,
        bomb_count: n_bombe,
        flag_count: 0,
        tempo: tempo,
    }

    r := rand.New(rand.NewSource(time.Now().UnixNano()));
    for {
        for y := range ret.celle {
            if ret.celle[y]==nil {ret.celle[y] = make([]Cella, width);}
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
                            if off_y==int(first_y) && off_x==int(first_x) {
                                goto out
                            }
                        }
                    }
                }
            }
        }
        out:
        if !ret.celle[first_y][first_x].is_bomb && ret.celle[first_y][first_x].label==0 {
            break;
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

func (g *Game) get_loosing_message() ([]CellaCoords) {
    changed := make([]CellaCoords, 0)
    for y, row := range g.celle {
        for x, cella := range row {
            if cella.is_bomb && !cella.is_flagged {
                changed = append(changed, CellaCoords{
                    x:     uint16(x),
                    y:     uint16(y),
                    cella: cella,
                })
            }
        }
    }
    return changed
}

func (g *Game) click(x uint16, y uint16) ([]CellaCoords, error) {
    if err:=g.check_bounds(int(x),int(y)); err!=nil {return nil, err;}
    if !g.celle[y][x].is_hidden {return nil, fmt.Errorf("Clicking uncovered cell");}

    changed := make([]CellaCoords, 0)

    if (g.celle[y][x].is_bomb) {
        g.state = Lost;
        return g.get_loosing_message(), nil;
    }

    stack := NewStack[CellaCoordsRef]();
    stack.Push(CellaCoordsRef{x, y, &g.celle[y][x]})

    for stack.Len()>0 {
        corrente := stack.Pop();

        if !corrente.cella.is_hidden {continue;}
        corrente.cella.is_hidden = false;

        changed = append(changed, CellaCoords{corrente.x, corrente.y, *corrente.cella})

        if corrente.cella.label > 0 {continue;}
        for i:=-1; i<=1; i++ {
            for j:=-1; j<=1; j++ {
                off_y := int(corrente.y)+i;
                off_x := int(corrente.x)+j;
                if g.is_inside(off_x, off_y) && g.celle[off_y][off_x].is_hidden {
                    stack.Push(CellaCoordsRef{uint16(off_x), uint16(off_y), &g.celle[off_y][off_x]})
                }
            }
        }
    }

    g.check_won();
    return changed, nil;
}

func (g *Game) check_won() {
    for _, row := range g.celle {
        for _, cel := range row {
            if cel.is_hidden && !cel.is_bomb {return;}
        }
    }
    g.state=Won
}

// Toggles the flagged state of a cell.
// Returns the new state.
// Errors if the coordinates lay outside the map
func (g *Game) flag(x uint16, y uint16) (bool, error) {
    if err:=g.check_bounds(int(x),int(y)); err!=nil {return false, err;}
    if !g.celle[y][x].is_hidden {return false, fmt.Errorf("flagging uncovered cell")}
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


