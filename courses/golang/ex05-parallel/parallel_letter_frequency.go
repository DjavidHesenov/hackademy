package letter

import "sync"

func Frequency(str string) map[rune]uint {
	frequencies := make(map[rune]uint)
	for _, ch := range str {
		frequencies[ch]++
	}
	return frequencies
}

func mergingMaps(maps ...map[rune]uint) map[rune]uint {
	res := make(map[rune]uint)

	for _, m := range maps {
		for key, val := range m {
			res[key] += val
		}
	}

	return res
}

func ConcurrentFrequency(strings []string) map[rune]uint {
	var wg sync.WaitGroup
	res := make([]map[rune]uint, len(strings))
	for index, str := range strings {
		wg.Add(1)
		go func(index int, str string) {
			res[index] = Frequency(str)
			wg.Done()
		}(index, str)
	}
	wg.Wait()
	return mergingMaps(res...)
}
