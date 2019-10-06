package main

func main () {
	var opts Opts
	parseCmd(&opts)
	dispatch(&opts)
}
