package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

func euclidean_distance(src [2]float64, dst [2]float64) float64 {
	dx := src[0] - dst[0]
	dy := src[1] - dst[1]
	return math.Sqrt(dx*dx + dy*dy)
}

func generate_random_passwords(word_list []string, n int, word_count int) ([]string, error) {
	passphrases := make([]string, n)
	for i := 0; i < n; i += 1 {
		passphrase := make([]string, word_count)
		for j := 0; j < word_count; j += 1 {
			k := rand.Int() % len(word_list)
			passphrase[j] = word_list[k]
		}
		passphrases[i] = strings.Join(passphrase, " ")
	}
	return passphrases, nil
}

func read_file_lines(filename string) ([]string, error) {
	file_byte_contents, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	// Line endings not made OS-agnostic by Go like in Python.
	file_no_carriages := strings.ReplaceAll(string(file_byte_contents), "\r", "")
	lines := strings.Split(file_no_carriages, "\n")
	return lines, nil
}

var finger_to_keys = map[int]string{
	2: "qaz",
	3: "xews",
	4: "tfgcvd",
	7: "bhnyjm",
	8: "uik",
	9: "lop",
}

var key_to_finger map[rune]int // Created at runtime

func qwerty_typing_distance(s string, key_distances map[rune][2]float64) float64 {
	// left  right (hand)
	// QWERT YUIOP
	// ASDFG HJKL
	// ZXCV  BNM

	// 2:qaz
	// 3:xews
	// 4:tfgcvd
	// 7:bhnyjm
	// 8:uik
	// 9:lop

	// PENALTY: moving between hands often to form a word (brain lag)
	//          also includes distance for each letter and the last letter that finger typed.
	//          should remember last position of finger
	// PENALTY, but not an option: shift and numericals

	finger_positions := make([][2]float64, 10)
	// assume middle row
	finger_positions[2] = key_distances['a']
	finger_positions[3] = key_distances['s']
	finger_positions[4] = key_distances['f']
	finger_positions[7] = key_distances['j']
	finger_positions[8] = key_distances['k']
	finger_positions[9] = key_distances['l']

	var cost_metric float64
	var last_finger int

	for _, c := range s {
		c_pos := key_distances[c]
		finger := key_to_finger[c]
		last_pos := finger_positions[finger]
		distance := euclidean_distance(c_pos, last_pos)
		cost_metric += distance
		fingers_same_side := last_finger <= 4 && finger <= 4 || last_finger >= 7 && finger >= 7
		if last_finger != 0 && !fingers_same_side {
			cost_metric += 0.4 //TODO, proper cost?
		}
		last_finger = finger
	}

	return cost_metric
}

func load_key_positions(filename string) (key_to_pos map[rune][2]float64, err error) {
	lines, err := read_file_lines(filename)
	if err != nil {
		return nil, err
	}
	//ignore first line, header
	key_to_pos = make(map[rune][2]float64)
	for i := 1; i < len(lines); i += 1 {
		fields := strings.Split(lines[i], ",")
		if len(fields[0]) == 0 {
			continue
		}
		key := rune(fields[0][0])
		x, err := strconv.ParseFloat(fields[1], 64)
		if err != nil {
			return nil, err
		}
		y, err := strconv.ParseFloat(fields[2], 64)
		if err != nil {
			return nil, err
		}
		key_to_pos[key] = [2]float64{x, y}
	}
	return key_to_pos, nil
}

func main() {
	key_to_finger = make(map[rune]int)
	for finger, keys := range finger_to_keys {
		for _, key := range keys {
			key_to_finger[key] = finger
		}
	}

	english_words, err := read_file_lines("google-10000-english.txt")
	if err != nil {
		log.Fatal(err)
	}
	english_words_filtered := make([]string, 0, len(english_words))
	for _, word := range english_words {
		if len(word) >= 6 {
			english_words_filtered = append(english_words_filtered, word)
		}
	}

	key_to_pos, err := load_key_positions("key_distances.csv")
	if err != nil {
		log.Fatal(err)
	}

	passphrases, err := generate_random_passwords(english_words_filtered, 100000, 3)
	if err != nil {
		log.Fatal(err)
	}

	lowest_cost := 1000.0
	var best_phrase string

	for _, passphrase := range passphrases {
		cost := qwerty_typing_distance(passphrase, key_to_pos)
		if cost < lowest_cost {
			lowest_cost = cost
			best_phrase = passphrase
		}
	}
	fmt.Println(best_phrase, lowest_cost)

	// create N random passphrases

	// score them all via a metric of how easy it is to type

	// output the top 3
}
