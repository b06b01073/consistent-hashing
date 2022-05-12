package main

import (
	cli "consistent-hashing/cli"
	hash "consistent-hashing/hash"
)

/*
 TODO:
	3. benchmarking
	4. testing
*/

func main() {
	hashRing := hash.GetRing()

	cli.Cli(hashRing)
}
