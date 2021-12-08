Converts initial random prices to 1m-candles, 1m-candles to 2m-candles and 2m-candles to 10m candles.\
Saves all formed candler to the corresponding file. On SIGINT even not closed candles are saved to file as well.\
Using pipe pattern, async (goroutines, wait groups) and context.
