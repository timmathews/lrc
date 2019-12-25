/*
 * lrc is a Left, Right, Center simulator
 * Copyright (C) 2019 Tim Mathews <tim@signalk.org>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to
 *
 *  The Free Software Foundation, Inc.
 *  51 Franklin Street, Fifth Floor
 *  Boston, MA 02110-1301, USA
 *
 */

package main

import (
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"math"
	mrand "math/rand"
	"os"
)

var preamble = "Lay your money down\nLet's go!\n"
var players int

type cryptRandSrc struct{}

func (s cryptRandSrc) Seed(seed int64) {}

func (s cryptRandSrc) Int63() int64 {
	return int64(s.Uint64() & ^uint64(1<<63))
}

func (s cryptRandSrc) Uint64() (v uint64) {
	binary.Read(crand.Reader, binary.BigEndian, &v)
	return v
}

var src cryptRandSrc

type Player struct {
	dollars int
	left    *Player
	right   *Player
	name    string
}

func dice(num int) (s []int) {
	var rnd = mrand.New(src)
	s = make([]int, num)
	for i := 0; i < num; i++ {
		s[i] = rnd.Intn(6)
	}

	return s
}

func winner(table []Player) (s string) {
	for i := 0; i < len(table); i++ {
		if table[i].dollars > 0 {
			s = table[i].name
		}
	}

	return s
}

func playersInGame(table []Player) (s int) {
	s = 0
	for i := 0; i < len(table); i++ {
		if table[i].dollars > 0 {
			s++
		}
	}

	return s
}

func printRemainers(lastRoundPlayers int, table []Player) {
	fmt.Println()
	playersLeft := playersInGame(table)

	if playersLeft < lastRoundPlayers && playersLeft > 1 {
		lastRoundPlayers = playersLeft
		if playersLeft+1 == players {
			fmt.Println("First one out!")
		} else if playersLeft < players/2 {
			fmt.Printf("Just %v players left\n", playersLeft)
		} else if playersLeft == 2 {
			fmt.Println("And down to two")
		} else {
			fmt.Println("One more down!")
		}
	}
}

func main() {
	players = len(os.Args) - 1

	if players < 2 {
		fmt.Println("LRC needs at least two players")
		fmt.Println()
		fmt.Println("./lrc Joe Bob")
		os.Exit(1)
	}

	fmt.Println(preamble)
	table := make([]Player, players)
	kitty := 0

	for i := 0; i < len(os.Args)-1; i++ {
		table[i].dollars = 3
		table[i].name = os.Args[i+1]
		if i == 0 {
			table[i].left = &table[players-1]
			table[i].right = &table[i+1]
		} else if i == players-1 {
			table[i].left = &table[i-1]
			table[i].right = &table[0]
		} else {
			table[i].left = &table[i-1]
			table[i].right = &table[i+1]
		}
	}

	for {
		lastRoundPlayers := playersInGame(table)
		for i := 0; i < players; i++ {
			if lastRoundPlayers < 2 {
				fmt.Printf("And the winner is ... %s!!\n", winner(table))
				os.Exit(0)
			}
			if table[i].dollars > 0 {
				d := dice(int(math.Min(float64(table[i].dollars), 3)))
				for j := 0; j < len(d); j++ {
					if d[j] < 2 {
						table[i].left.dollars++
						table[i].dollars--
						fmt.Print("Left, ")
					} else if d[j] < 4 {
						table[i].right.dollars++
						table[i].dollars--
						fmt.Print("Right, ")
					} else if d[j] == 4 {
						kitty++
						table[i].dollars--
						fmt.Print("Center! ")
					}
				}

				printRemainers(lastRoundPlayers, table)
				lastRoundPlayers = playersInGame(table)
			}
		}
	}
}
